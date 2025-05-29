package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	oomodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/client"
	cmodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/client/model"
	"github.com/google/uuid"
	"github.com/mattermost/mattermost/server/public/model"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/language"
)

func BuildConvertHandler(plugin api.PluginAPI) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		_ = r.Context() // TODO: Proper timeouts
		var req oomodel.ConvertFile

		userID := r.Header.Get(plugin.Configuration.MMAuthHeader)
		if userID == "" {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not get user ID from request")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not decode request: " + err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := req.Validate(); err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "invalid request: " + err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var g errgroup.Group
		var file *model.FileInfo
		var user *model.User

		g.Go(func() error {
			var appErr *model.AppError
			file, appErr = plugin.API.GetFileInfo(req.FileID)
			if appErr != nil {
				return appErr
			}

			return nil
		})

		g.Go(func() error {
			var appErr *model.AppError
			user, appErr = plugin.API.GetUser(userID)
			if appErr != nil {
				return appErr
			}

			return nil
		})

		if err := g.Wait(); err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not get channel or user: " + err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if file == nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "file is nil after retrieval")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if user == nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "user is nil after retrieval")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if user.Id != file.CreatorId {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "user is not the owner of the file")
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		lang := "en"
		if user.Locale != "" {
			lang = user.Locale
		}

		tag, err := language.Parse(lang)
		if err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not parse language: " + err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		region, _ := tag.Region()

		serverURL := *plugin.API.GetConfig().ServiceSettings.SiteURL + "/" + _OnlyofficeAPIRootSuffix
		creq := cmodel.CommandConvertRequest{
			Async:      false,
			Key:        uuid.NewString(),
			Filetype:   file.Extension,
			Password:   req.Password,
			Outputtype: "ooxml",
			URL:        fmt.Sprintf("%s/download?id=%s", serverURL, file.Id),
			Region:     fmt.Sprintf("%s-%s", lang, region),
		}

		cresp, err := plugin.OnlyofficeCommandClient.SendConvert(
			plugin.Configuration.Address+client.OnlyofficeCommandConverterPath+"?shardkey="+uuid.New().String(),
			creq, 4*time.Second,
		)
		if err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not send convert api request: " + err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if cresp.Error != 0 {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not convert file: " + file.Id)
			resp := &oomodel.ConvertFileResponse{Error: cresp.Error}
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write(resp.ToJSON())
			return
		}

		fres, err := http.Get(cresp.FileURL)
		if err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not get converted file: " + err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		defer fres.Body.Close()

		fileData, err := io.ReadAll(fres.Body)
		if err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not read converted file: " + err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		uploadSession, sessionErr := plugin.API.CreateUploadSession(&model.UploadSession{
			Id:        model.NewId(),
			UserId:    userID,
			ChannelId: file.ChannelId,
			Filename:  strings.TrimSuffix(file.Name, "."+file.Extension) + "." + cresp.FileType,
			FileSize:  int64(len(fileData)),
			Type:      model.UploadTypeAttachment,
		})

		if sessionErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not create upload session: " + sessionErr.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if uploadSession == nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "upload session is nil after creation")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		fileInfo, uploadErr := plugin.API.UploadData(uploadSession, bytes.NewReader(fileData))
		if uploadErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not upload file data: " + uploadErr.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if fileInfo == nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "file info is nil after upload")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, postErr := plugin.API.CreatePost(&model.Post{
			ChannelId: file.ChannelId,
			RootId:    file.PostId,
			FileIds:   []string{fileInfo.Id},
			UserId:    userID,
		}); postErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not create post: " + postErr.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := &oomodel.ConvertFileResponse{Error: 0}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(resp.ToJSON())
	}
}

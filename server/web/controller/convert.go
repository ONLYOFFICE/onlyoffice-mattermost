/**
 *
 * (c) Copyright Ascensio System SIA 2025
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	oomodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	"github.com/google/uuid"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/language"
)

const (
	defaultOutputType = "ooxml"
	convertTimeout    = 4 * time.Second
	logPrefix         = onlyofficeLoggerPrefix
)

var (
	errNoValidCredentials = errors.New("no valid credentials and demo is not active")
	errFileNullPointer    = errors.New("file is nil after retrieval")
	errUserNullPointer    = errors.New("user is nil after retrieval")
	errUserNotOwner       = errors.New("user is not the owner of the file")
	errUnsupportedFormat  = errors.New("unsupported file format")
)

type ConvertHandler struct {
	api           plugin.API
	configuration *configuration.Configuration
	formatManager public.FormatManager
	jwtManager    crypto.JwtManager
	commandClient client.CommandClient
}

func NewConvertHandler(
	api plugin.API,
	configuration *configuration.Configuration,
	formatManager public.FormatManager,
	jwtManager crypto.JwtManager,
	commandClient client.CommandClient,
) ConvertHandler {
	return ConvertHandler{
		api:           api,
		configuration: configuration,
		formatManager: formatManager,
		jwtManager:    jwtManager,
		commandClient: commandClient,
	}
}

func (h *ConvertHandler) logErrorAndRespond(rw http.ResponseWriter, message string, statusCode int) {
	h.api.LogError(logPrefix + message)
	rw.WriteHeader(statusCode)
}

func (h *ConvertHandler) logErrorAndRespondWithJSON(rw http.ResponseWriter, message string, resp *oomodel.ConvertFileResponse) {
	h.api.LogError(logPrefix + message)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp.ToJSON())
}

func (h *ConvertHandler) validateCredentials(r *http.Request) (string, error) {
	userID := r.Header.Get(tools.MMAuthHeader)
	if userID == "" {
		return "", errNoValidCredentials
	}

	hasOwnCredentials := h.configuration.DESAddress != h.configuration.DemoAddress &&
		h.configuration.DESJwt != "" &&
		h.configuration.DESJwtHeader != "" &&
		h.configuration.DESJwtPrefix != ""

	demoActive := h.configuration.DemoEnabled &&
		h.configuration.DemoExpires >= time.Now().UnixMilli()

	if !demoActive && !hasOwnCredentials {
		return "", errNoValidCredentials
	}

	return userID, nil
}

func (h *ConvertHandler) fetchData(fileID string, userID string) (*model.FileInfo, *model.User, error) {
	var g errgroup.Group
	var file *model.FileInfo
	var user *model.User

	g.Go(func() error {
		var appErr *model.AppError
		file, appErr = h.api.GetFileInfo(fileID)
		if appErr != nil {
			return appErr
		}
		return nil
	})

	g.Go(func() error {
		var appErr *model.AppError
		user, appErr = h.api.GetUser(userID)
		if appErr != nil {
			return appErr
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	return file, user, nil
}

func (h *ConvertHandler) validateFileAndUser(file *model.FileInfo, user *model.User) error {
	if file == nil {
		return errFileNullPointer
	}

	if user == nil {
		return errUserNullPointer
	}

	if user.Id != file.CreatorId {
		return errUserNotOwner
	}

	return nil
}

func (h *ConvertHandler) prepareConvertRequest(file *model.FileInfo, user *model.User, req *oomodel.ConvertFileRequest) (*client.ConvertRequest, error) {
	tag, err := language.Parse(user.Locale)
	if err != nil {
		return nil, fmt.Errorf("could not parse language: %w", err)
	}

	region, _ := tag.Region()
	if _, supported := h.formatManager.GetFormatByName(file.Extension); !supported {
		return nil, errUnsupportedFormat
	}

	outputType := defaultOutputType
	if _, supported := h.formatManager.GetFormatByName(file.Extension); supported && req.OutputType != "" {
		outputType = req.OutputType
	}

	serverURL := *h.api.GetConfig().ServiceSettings.SiteURL + "/" + onlyofficeAPIRootSuffix
	convertReq := &client.ConvertRequest{
		Async:      false,
		Key:        uuid.NewString(),
		Filetype:   file.Extension,
		Password:   req.Password,
		Outputtype: outputType,
		URL:        fmt.Sprintf("%s/download?id=%s", serverURL, file.Id),
		Region:     fmt.Sprintf("%s-%s", user.Locale, region),
	}

	if h.configuration.DESJwt != "" {
		convertReq.Token, _ = h.jwtManager.Sign([]byte(h.configuration.DESJwt), convertReq)
	}

	return convertReq, nil
}

func (h *ConvertHandler) performConversion(convertReq *client.ConvertRequest, fileID string) (*client.ConvertResponse, error) {
	convertURL := h.configuration.DESAddress + client.OnlyofficeCommandConverterPath + "?shardkey=" + uuid.New().String()
	convertResp, err := h.commandClient.SendConvert(convertURL, *convertReq, convertTimeout)
	if err != nil {
		return nil, fmt.Errorf("could not send convert api request: %w", err)
	}

	if convertResp.Error != 0 {
		return &convertResp, fmt.Errorf("could not convert file: %s", fileID)
	}

	return &convertResp, nil
}

func (h *ConvertHandler) downloadConvertedFile(fileURL string) ([]byte, error) {
	response, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("could not get converted file: %w", err)
	}
	defer response.Body.Close()

	fileData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read converted file: %w", err)
	}

	return fileData, nil
}

func (h *ConvertHandler) uploadAndCreatePost(fileData []byte, file *model.FileInfo, userID string, fileType string) error {
	filename := strings.TrimSuffix(file.Name, "."+file.Extension) + "." + fileType
	uploadSession, sessionErr := h.api.CreateUploadSession(&model.UploadSession{
		Id:        model.NewId(),
		UserId:    userID,
		ChannelId: file.ChannelId,
		Filename:  filename,
		FileSize:  int64(len(fileData)),
		Type:      model.UploadTypeAttachment,
	})

	if sessionErr != nil {
		return fmt.Errorf("could not create upload session: %w", sessionErr)
	}

	if uploadSession == nil {
		return errors.New("upload session is nil after creation")
	}

	fileInfo, uploadErr := h.api.UploadData(uploadSession, bytes.NewReader(fileData))
	if uploadErr != nil {
		return fmt.Errorf("could not upload file data: %w", uploadErr)
	}

	if fileInfo == nil {
		return errors.New("file info is nil after upload")
	}

	if _, postErr := h.api.CreatePost(&model.Post{
		ChannelId: file.ChannelId,
		RootId:    file.PostId,
		FileIds:   []string{fileInfo.Id},
		UserId:    userID,
	}); postErr != nil {
		return fmt.Errorf("could not create post: %w", postErr)
	}

	return nil
}

func (h *ConvertHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	userID, err := h.validateCredentials(r)
	if err != nil {
		h.logErrorAndRespond(rw, "no valid credentials and demo is not active", http.StatusUnauthorized)
		return
	}

	var req oomodel.ConvertFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logErrorAndRespond(rw, "could not decode request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		h.logErrorAndRespond(rw, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, user, err := h.fetchData(req.FileID, userID)
	if err != nil {
		h.logErrorAndRespond(rw, "could not get channel or user: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validateFileAndUser(file, user); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user is not the owner of the file" {
			statusCode = http.StatusForbidden
		}
		h.logErrorAndRespond(rw, err.Error(), statusCode)
		return
	}

	convertReq, err := h.prepareConvertRequest(file, user, &req)
	if err != nil {
		h.logErrorAndRespond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	convertResp, err := h.performConversion(convertReq, file.Id)
	if err != nil {
		if convertResp != nil && convertResp.Error != 0 {
			// Conversion failed with specific error
			response := &oomodel.ConvertFileResponse{Error: convertResp.Error}
			h.logErrorAndRespondWithJSON(rw, err.Error(), response)
			return
		}
		h.logErrorAndRespond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	fileData, err := h.downloadConvertedFile(convertResp.FileURL)
	if err != nil {
		h.logErrorAndRespond(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.uploadAndCreatePost(fileData, file, userID, convertResp.FileType); err != nil {
		h.logErrorAndRespond(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &oomodel.ConvertFileResponse{Error: 0}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp.ToJSON())
}

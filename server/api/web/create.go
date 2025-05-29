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
package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	oomodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/onlyoffice"
	"github.com/mattermost/mattermost/server/public/model"
	"golang.org/x/sync/errgroup"
)

func BuildCreateHandler(plugin api.PluginAPI) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var req oomodel.NewFile

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
		var channel *model.Channel
		var user *model.User

		g.Go(func() error {
			var appErr *model.AppError
			channel, appErr = plugin.API.GetChannel(req.ChannelID)
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

		if channel == nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "channel is nil after retrieval")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if user == nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "user is nil after retrieval")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var language string
		if user.Locale != "" {
			language = onlyoffice.MapLanguageToTemplate(user.Locale)
		} else {
			language = "default"
		}

		templatePath := path.Join("template", language, "new."+req.FileType)
		fileData, readErr := public.Templates.ReadFile(templatePath)
		if readErr != nil {
			defaultTemplatePath := path.Join("template", "default", "new."+req.FileType)
			fileData, readErr = public.Templates.ReadFile(defaultTemplatePath)
			if readErr != nil {
				plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not read template file: " + readErr.Error())
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		uploadSession, sessionErr := plugin.API.CreateUploadSession(&model.UploadSession{
			Id:        model.NewId(),
			UserId:    userID,
			ChannelId: channel.Id,
			Filename:  req.FileName + "." + req.FileType,
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
			ChannelId: channel.Id,
			FileIds:   []string{fileInfo.Id},
			UserId:    userID,
		}); postErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not create post: " + postErr.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}

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
	"net/http"
	"path"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	oomodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"golang.org/x/sync/errgroup"
)

type CreateHandler struct {
	api           plugin.API
	configuration *configuration.Configuration
}

func NewCreateHandler(api plugin.API, configuration *configuration.Configuration) CreateHandler {
	return CreateHandler{
		api:           api,
		configuration: configuration,
	}
}

func (h *CreateHandler) fetchData(channelID, userID string) (*model.Channel, *model.User, error) {
	var g errgroup.Group
	var channel *model.Channel
	var user *model.User

	g.Go(func() error {
		var appErr *model.AppError
		channel, appErr = h.api.GetChannel(channelID)
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

	return channel, user, nil
}

func (h *CreateHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	var req oomodel.NewFileRequest

	userID := r.Header.Get(tools.MMAuthHeader)
	if userID == "" {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get user ID from request")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not decode request: " + err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "invalid request: " + err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	channel, user, err := h.fetchData(req.ChannelID, userID)
	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get channel or user: " + err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if channel == nil {
		h.api.LogError(onlyofficeLoggerPrefix + "channel is nil after retrieval")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if user == nil {
		h.api.LogError(onlyofficeLoggerPrefix + "user is nil after retrieval")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var language string
	if user.Locale != "" {
		language = tools.MapLanguageToTemplate(user.Locale)
	} else {
		language = "default"
	}

	templatePath := path.Join("template", language, "new."+req.FileType)
	fileData, readErr := public.Templates.ReadFile(templatePath)
	if readErr != nil {
		defaultTemplatePath := path.Join("template", "default", "new."+req.FileType)
		fileData, readErr = public.Templates.ReadFile(defaultTemplatePath)
		if readErr != nil {
			h.api.LogError(onlyofficeLoggerPrefix + "could not read template file: " + readErr.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	uploadSession, sessionErr := h.api.CreateUploadSession(&model.UploadSession{
		Id:        model.NewId(),
		UserId:    userID,
		ChannelId: channel.Id,
		Filename:  req.FileName + "." + req.FileType,
		FileSize:  int64(len(fileData)),
		Type:      model.UploadTypeAttachment,
	})

	if sessionErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not create upload session: " + sessionErr.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if uploadSession == nil {
		h.api.LogError(onlyofficeLoggerPrefix + "upload session is nil after creation")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileInfo, uploadErr := h.api.UploadData(uploadSession, bytes.NewReader(fileData))
	if uploadErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not upload file data: " + uploadErr.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if fileInfo == nil {
		h.api.LogError(onlyofficeLoggerPrefix + "file info is nil after upload")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, postErr := h.api.CreatePost(&model.Post{
		ChannelId: channel.Id,
		FileIds:   []string{fileInfo.Id},
		UserId:    userID,
	}); postErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not create post: " + postErr.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

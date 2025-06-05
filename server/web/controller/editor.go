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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	validator "github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	oomodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
)

type editorParameters struct {
	UserID   string `json:"userID" validate:"required"`
	Username string `json:"username" validate:"required"`
	FileID   string `json:"fileID" validate:"required"`
	Lang     string `json:"lang"`
}

func (c *editorParameters) Validate() error {
	return validator.New().Struct(c)
}

type EditorHandler struct {
	api            plugin.API
	configuration  *configuration.Configuration
	fileHelper     file.FileHelper
	encoder        crypto.Encoder
	jwtManager     crypto.JwtManager
	editorTemplate *template.Template
}

func NewEditorHandler(
	api plugin.API,
	configuration *configuration.Configuration,
	fileHelper file.FileHelper,
	encoder crypto.Encoder,
	jwtManager crypto.JwtManager,
	editorTemplate *template.Template,
) EditorHandler {
	return EditorHandler{
		api:            api,
		configuration:  configuration,
		fileHelper:     fileHelper,
		encoder:        encoder,
		jwtManager:     jwtManager,
		editorTemplate: editorTemplate,
	}
}

func (h *EditorHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	h.api.LogDebug(onlyofficeLoggerPrefix + "got an editor request")
	hasOwnCredentials := h.configuration.DESAddress != h.configuration.DemoAddress &&
		h.configuration.DESJwt != "" &&
		h.configuration.DESJwtHeader != "" &&
		h.configuration.DESJwtPrefix != ""

	demoActive := h.configuration.DemoEnabled &&
		h.configuration.DemoExpires >= time.Now().UnixMilli()

	if !demoActive && !hasOwnCredentials {
		h.api.LogError(onlyofficeLoggerPrefix + "no valid credentials and demo is not active")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	serverURL := *h.api.GetConfig().ServiceSettings.SiteURL + "/" + onlyofficeAPIRootSuffix
	user, err := h.api.GetUser(r.Header.Get(tools.MMAuthHeader))
	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get user info")
		return
	}

	query := r.URL.Query()
	payload := editorParameters{
		UserID:   user.Id,
		Username: user.Username,
		FileID:   query.Get("file"),
		Lang:     query.Get("lang"),
	}

	validationErr := payload.Validate()
	if validationErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "editor payload validation error: " + validationErr.Error())
		return
	}

	fileInfo, fileInfoErr := h.api.GetFileInfo(payload.FileID)
	if fileInfoErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not access file info " + payload.FileID + " Reason: " + fileInfoErr.Message)
		return
	}

	post, postErr := h.api.GetPost(fileInfo.PostId)
	if postErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not access post " + fileInfo.PostId + "Reason: " + postErr.Message)
		return
	}

	docType, typeErr := h.fileHelper.GetFileType(fileInfo.Extension)
	if typeErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get file type: " + typeErr.Error())
		return
	}

	docKey, keyErr := h.encoder.Encode(fileInfo.Id + strconv.FormatInt(post.UpdateAt, 10))
	if keyErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not encode document key: " + keyErr.Error())
		return
	}

	permissions := oomodel.OnlyofficeDefaultPermissions
	if h.fileHelper.IsExtensionEditable(fileInfo.Extension) {
		permissions = h.fileHelper.GetFilePermissionsByUserID(payload.UserID, payload.FileID, post)
	}

	code := h.fileHelper.GenerateKey()
	if err := h.api.KVSetWithExpiry(code, []byte(payload.UserID), 10); err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not set code: " + err.Error())
	}

	theme := "theme-classic-light"
	if strings.ToLower(query.Get("dark")) == "true" {
		theme = "theme-dark"
	}

	config := oomodel.Config{
		Document: oomodel.Document{
			FileType:    fileInfo.Extension,
			Key:         docKey,
			Title:       fileInfo.Name,
			URL:         fmt.Sprintf("%s/download?id=%s", serverURL, fileInfo.Id),
			Permissions: permissions,
		},
		DocumentType: docType,
		EditorConfig: oomodel.EditorConfig{
			User: oomodel.User{
				ID:    payload.UserID,
				Name:  payload.Username,
				Image: fmt.Sprintf("%s/image?code=%s", serverURL, code),
			},
			CallbackURL: serverURL + "/callback?file=" + payload.FileID,
			Customization: oomodel.Customization{
				Goback: oomodel.Goback{
					RequestClose: true,
				},
				UiTheme: theme,
			},
			Lang: payload.Lang,
		},
		Type: tools.IsMobile(r.Header.Get("User-Agent")),
	}

	config.IssuedAt, config.ExpiresAt = jwt.NewNumericDate(time.Now()),
		jwt.NewNumericDate(time.Now().Add(3*time.Minute))
	cToken, cTokenErr := h.jwtManager.Sign([]byte(h.configuration.DESJwt), config)
	if cTokenErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not sign config: " + cTokenErr.Error())
		return
	}

	config.Token = cToken
	encodedConfig, cerr := json.Marshal(config)
	if cerr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not marshal config: " + cerr.Error())
		return
	}

	data := map[string]interface{}{
		"apijs":  h.configuration.DESAddress + "/web-apps/apps/api/documents/api.js",
		"config": string(encodedConfig),
		"dark":   query.Get("dark"),
	}

	h.api.LogDebug(onlyofficeLoggerPrefix + "building an editor window")
	if err := h.editorTemplate.ExecuteTemplate(rw, "editor.html", data); err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not execute editor template: " + err.Error())
	}
}

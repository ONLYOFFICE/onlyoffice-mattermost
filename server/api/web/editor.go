/**
 *
 * (c) Copyright Ascensio System SIA 2022
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
	"fmt"
	"net/http"
	"strconv"
	"time"

	integration "github.com/ONLYOFFICE/onlyoffice-mattermost"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	oomodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
	oovalidator "github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/validator"
	"github.com/go-playground/validator/v10"
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

func BuildEditorHandler(plugin api.PluginAPI) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		serverURL := *plugin.API.GetConfig().ServiceSettings.SiteURL + "/" + _OnlyofficeApiRootSuffix

		plugin.API.LogDebug(_OnlyofficeLoggerPrefix + "got an editor request")

		user, err := plugin.API.GetUser(r.Header.Get(plugin.Configuration.MMAuthHeader))
		if err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not get user info")
			rw.WriteHeader(http.StatusBadRequest)
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
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "editor payload validation error: " + validationErr.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		post, fileInfo := GetPostInfo(plugin, payload.FileID, r)
		if post == nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		docType, typeErr := plugin.OnlyofficeHelper.GetFileType(fileInfo.Extension)
		if typeErr != nil {
			plugin.API.LogError(typeErr.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		docKey, keyErr := plugin.Encoder.Encode(fileInfo.Id + strconv.FormatInt(post.UpdateAt, 10))
		if keyErr != nil {
			plugin.API.LogError(keyErr.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		permissions := oomodel.OnlyofficeDefaultPermissions
		if plugin.OnlyofficeHelper.IsExtensionEditable(fileInfo.Extension) {
			permissions = plugin.OnlyofficeHelper.GetFilePermissionsByUserID(payload.UserID, payload.FileID, post)
		}

		downloadURL := fmt.Sprintf("%s/download?file=%s", serverURL, fileInfo.Id)
		if len(plugin.Manager.GetKey()) > 0 {
			var token oomodel.PlainToken
			token.Id = payload.FileID
			token.IssuedAt = time.Now().Unix()
			token.ExpiresAt = time.Now().Add(5 * time.Minute).Unix()
			token.Issuer = integration.Manifest.Id

			dToken, dTokenErr := plugin.Manager.Sign(token)
			if dTokenErr != nil {
				plugin.API.LogError(dTokenErr.Error())
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			downloadURL = fmt.Sprintf("%s/download?file=%s&token=%s", serverURL, fileInfo.Id, dToken)
		}

		config := oomodel.Config{
			Document: oomodel.Document{
				FileType:    fileInfo.Extension,
				Key:         docKey,
				Title:       fileInfo.Name,
				URL:         downloadURL,
				Permissions: permissions,
			},
			DocumentType: docType,
			EditorConfig: oomodel.EditorConfig{
				User: oomodel.User{
					ID:   payload.UserID,
					Name: payload.Username,
				},
				CallbackURL: serverURL + "/callback?file=" + payload.FileID,
				Customization: oomodel.Customization{
					Goback: oomodel.Goback{
						RequestClose: true,
					},
				},
				Lang: payload.Lang,
			},
			Type: oovalidator.IsMobile(r.Header.Get("User-Agent")),
		}

		if len(plugin.Manager.GetKey()) > 0 {
			config.IssuedAt = time.Now().Unix()
			config.ExpiresAt = time.Now().Add(5 * time.Minute).Unix()
			config.Issuer = integration.Manifest.Id
			token, err := plugin.Manager.Sign(config)
			if err != nil {
				plugin.API.LogError(err.Error())
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			config.Token = token
		}

		data := map[string]interface{}{
			"apijs":  plugin.Configuration.Address + "/web-apps/apps/api/documents/api.js",
			"config": config,
		}

		plugin.API.LogDebug(_OnlyofficeLoggerPrefix + "building an editor window")
		plugin.EditorTemplate.ExecuteTemplate(rw, "editor.html", data)
	}
}

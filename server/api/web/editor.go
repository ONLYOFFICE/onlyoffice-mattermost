/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
		plugin.API.LogDebug(_OnlyofficeLoggerPrefix + "got an editor request")
		serverURL := *plugin.API.GetConfig().ServiceSettings.SiteURL + "/" + _OnlyofficeApiRootSuffix

		user, err := plugin.API.GetUser(r.Header.Get(plugin.Configuration.MMAuthHeader))
		if err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not get user info")
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
			return
		}

		fileInfo, fileInfoErr := plugin.API.GetFileInfo(payload.FileID)
		if fileInfoErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not access file info " + payload.FileID + " Reason: " + fileInfoErr.Message)
			return
		}

		post, postErr := plugin.API.GetPost(fileInfo.PostId)
		if postErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not access post " + fileInfo.PostId + "Reason: " + postErr.Message)
			return
		}

		docType, typeErr := plugin.OnlyofficeHelper.GetFileType(fileInfo.Extension)
		if typeErr != nil {
			plugin.API.LogError(typeErr.Error())
			return
		}

		docKey, keyErr := plugin.Encoder.Encode(fileInfo.Id + strconv.FormatInt(post.UpdateAt, 10))
		if keyErr != nil {
			plugin.API.LogError(keyErr.Error())
			return
		}

		permissions := oomodel.OnlyofficeDefaultPermissions
		if plugin.OnlyofficeHelper.IsExtensionEditable(fileInfo.Extension) {
			permissions = plugin.OnlyofficeHelper.GetFilePermissionsByUserID(payload.UserID, payload.FileID, post)
		}

		dToken := &oomodel.DownloadToken{
			FileID: payload.FileID,
		}
		dToken.IssuedAt, dToken.ExpiresAt = time.Now().Unix(), time.Now().Add(3*time.Minute).Unix()
		dsignature, dTokenErr := plugin.Manager.Sign(dToken)
		if dTokenErr != nil {
			plugin.API.LogError(dTokenErr.Error())
			return
		}

		config := oomodel.Config{
			Document: oomodel.Document{
				FileType:    fileInfo.Extension,
				Key:         docKey,
				Title:       fileInfo.Name,
				URL:         fmt.Sprintf("%s/download?token=%s", serverURL, dsignature),
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

		config.IssuedAt, config.ExpiresAt = time.Now().Unix(), time.Now().Add(3*time.Minute).Unix()
		cToken, cTokenErr := plugin.Manager.Sign(config)
		if cTokenErr != nil {
			plugin.API.LogError(cTokenErr.Error())
			return
		}

		config.Token = cToken
		data := map[string]interface{}{
			"apijs":  plugin.Configuration.Address + "/web-apps/apps/api/documents/api.js",
			"config": config,
		}

		plugin.API.LogDebug(_OnlyofficeLoggerPrefix + "building an editor window")
		plugin.EditorTemplate.ExecuteTemplate(rw, "editor.html", data)
	}
}

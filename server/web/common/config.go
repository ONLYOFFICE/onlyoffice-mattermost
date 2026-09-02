/**
 *
 * (c) Copyright Ascensio System SIA 2026
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
package common

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

const (
	LoggerPrefix  = "[ONLYOFFICE Handler]: "
	apiRootSuffix = "plugins/com.onlyoffice.mattermost/api"
)

type editorParameters struct {
	UserID   string `json:"userID" validate:"required"`
	Username string `json:"username" validate:"required"`
	FileID   string `json:"fileID" validate:"required"`
	Lang     string `json:"lang"`
}

func (c *editorParameters) validate() error {
	return validator.New().Struct(c)
}

func BuildEditorConfig(
	r *http.Request,
	configuration *configuration.Configuration,
	api plugin.API,
	encoder crypto.Encoder,
	jwtManager crypto.JwtManager,
	fileHelper file.FileHelper,
) (config oomodel.Config, docKey string, isOwner bool, statusCode int, err error) {
	hasOwnCredentials := configuration.DESAddress != configuration.DemoAddress &&
		configuration.DESJwt != "" &&
		configuration.DESJwtHeader != "" &&
		configuration.DESJwtPrefix != ""

	demoActive := configuration.DemoEnabled &&
		configuration.DemoExpires >= time.Now().UnixMilli()

	if !demoActive && !hasOwnCredentials {
		return oomodel.Config{}, "", false, http.StatusUnauthorized, fmt.Errorf("no valid credentials and demo is not active")
	}

	serverURL := *api.GetConfig().ServiceSettings.SiteURL + "/" + apiRootSuffix
	user, userErr := api.GetUser(r.Header.Get(tools.MMAuthHeader))
	if userErr != nil {
		return oomodel.Config{}, "", false, http.StatusInternalServerError, fmt.Errorf("could not get user info")
	}

	query := r.URL.Query()
	payload := editorParameters{
		UserID:   user.Id,
		Username: user.Username,
		FileID:   query.Get("file"),
		Lang:     query.Get("lang"),
	}

	if validationErr := payload.validate(); validationErr != nil {
		return oomodel.Config{}, "", false, http.StatusBadRequest, fmt.Errorf("editor payload validation error: %w", validationErr)
	}

	fileInfo, fileInfoErr := api.GetFileInfo(payload.FileID)
	if fileInfoErr != nil {
		return oomodel.Config{}, "", false, http.StatusInternalServerError, fmt.Errorf("could not access file info %s: %s", payload.FileID, fileInfoErr.Message)
	}

	if !configuration.IsFormatAllowedForViewing(fileInfo.Extension) {
		return oomodel.Config{}, "", false, http.StatusForbidden, fmt.Errorf("format not allowed for viewing: %s", fileInfo.Extension)
	}

	post, postErr := api.GetPost(fileInfo.PostId)
	if postErr != nil {
		return oomodel.Config{}, "", false, http.StatusInternalServerError, fmt.Errorf("could not access post %s: %s", fileInfo.PostId, postErr.Message)
	}

	isOwner = payload.UserID == post.UserId

	docType, typeErr := fileHelper.GetFileType(fileInfo.Extension)
	if typeErr != nil {
		return oomodel.Config{}, "", false, http.StatusInternalServerError, fmt.Errorf("could not get file type: %w", typeErr)
	}

	docKey, keyErr := encoder.Encode(fileInfo.Id + strconv.FormatInt(post.UpdateAt, 10))
	if keyErr != nil {
		return oomodel.Config{}, "", false, http.StatusInternalServerError, fmt.Errorf("could not encode document key: %w", keyErr)
	}

	permissions := oomodel.OnlyofficeDefaultPermissions
	if fileHelper.IsExtensionEditable(fileInfo.Extension) {
		if configuration.IsFormatAllowedForEditing(fileInfo.Extension) {
			permissions = fileHelper.GetFilePermissionsByUserID(payload.UserID, payload.FileID, post)
		} else {
			api.LogDebug(LoggerPrefix + "format not allowed for editing, forcing read-only: " + fileInfo.Extension)
		}
	}

	if configuration.OwnerProtected {
		permissions.Protect = isOwner
	} else {
		permissions.Protect = true
	}

	code := fileHelper.GenerateKey()
	if err := api.KVSetWithExpiry(code, []byte(payload.UserID), 10); err != nil {
		api.LogError(LoggerPrefix + "could not set code: " + err.Error())
	}

	theme := "default-light"
	if strings.ToLower(query.Get("dark")) == "true" {
		theme = "default-dark"
	}

	editorConfig := oomodel.EditorConfig{
		User: oomodel.User{
			ID:    payload.UserID,
			Name:  payload.Username,
			Image: fmt.Sprintf("%s/image?code=%s", serverURL, code),
		},
		CallbackURL: serverURL + "/callback?file=" + payload.FileID,
		Customization: oomodel.Customization{
			UiTheme: theme,
			Close: oomodel.Close{
				Visible: true,
			},
			Plugins: configuration.PluginsEnabled,
			Macros:  configuration.MacrosEnabled,
		},
		Lang: payload.Lang,
	}

	editorConfig.CoEditing = oomodel.CoEditing{
		Mode:   "fast",
		Change: permissions.Edit,
	}

	if !permissions.Edit {
		editorConfig.Mode = "view"
	}

	config = oomodel.Config{
		Document: oomodel.Document{
			FileType:    fileInfo.Extension,
			Key:         docKey,
			Title:       fileInfo.Name,
			URL:         fmt.Sprintf("%s/download?id=%s", serverURL, fileInfo.Id),
			Permissions: permissions,
		},
		DocumentType: docType,
		EditorConfig: editorConfig,
		Type:         tools.IsMobile(r.Header.Get("User-Agent")),
	}

	config.IssuedAt, config.ExpiresAt = jwt.NewNumericDate(time.Now()),
		jwt.NewNumericDate(time.Now().Add(3*time.Minute))
	cToken, cTokenErr := jwtManager.Sign([]byte(configuration.DESJwt), config)
	if cTokenErr != nil {
		return oomodel.Config{}, "", false, http.StatusInternalServerError, fmt.Errorf("could not sign config: %w", cTokenErr)
	}

	config.Token = cToken
	return config, docKey, isOwner, http.StatusOK, nil
}

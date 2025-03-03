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
 *https://developers.mattermost.com/integrate/plugins/components/server/hello-world/
 */
package api

import (
	"html/template"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
	"github.com/golang-jwt/jwt/v5"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
)

type Encoder interface {
	Encode(text string) (string, error)
}

type JwtManager interface {
	Sign(payload jwt.Claims) (string, error)
	Verify(jwt string, body interface{}) error
	GetKey() []byte
}

type Bot interface {
	BotCreateDM(message string, userID string)
	BotCreatePost(message string, channelID string)
	BotCreateReply(message string, channelID string, parentID string)
}

type OnlyofficeHelper interface {
	IsExtensionSupported(fileExt string) bool
	IsExtensionEditable(fileExt string) bool
	GetFileType(fileExt string) (string, error)
	GenerateKey() string
	GetPostPermissionsByFileID(fileID string, post *mmModel.Post, getUser func(string) (*mmModel.User, *mmModel.AppError)) []model.UserInfoResponse
	GetFilePermissionsByUserID(userID string, fileID string, post *mmModel.Post) model.Permissions
	UserHasFilePermissions(userID string, fileID string, post *mmModel.Post) bool
	SetPostFilePermissions(post *mmModel.Post, permissions []model.PostPermission) []model.PostPermission
	GetWildcardUser() string
}

type OnlyofficeCovnerter interface {
	GetTimestamp() int64
	GetTime(timestamp int64) time.Time
}

type PluginAPI struct {
	API           plugin.API
	Configuration struct {
		Address      string
		Secret       string
		Header       string
		Prefix       string
		MMAuthHeader string
	}
	OnlyofficeHelper    OnlyofficeHelper
	OnlyofficeConverter OnlyofficeCovnerter
	Encoder             Encoder
	Manager             JwtManager
	Bot                 Bot
	EditorTemplate      *template.Template
	Filestore           filestore.FileBackend
}

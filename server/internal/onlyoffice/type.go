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
package onlyoffice

import (
	mmModel "github.com/mattermost/mattermost/server/public/model"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
)

var _ Helper = (*helper)(nil)

type Helper interface {
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

type helper struct {
	formatManager public.FormatManager
}

func NewHelper(formatManager public.FormatManager) Helper {
	return &helper{
		formatManager: formatManager,
	}
}

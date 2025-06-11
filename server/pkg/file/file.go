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
package file

import (
	mmModel "github.com/mattermost/mattermost/server/public/model"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
)

var _ FileHelper = (*fileHelperImpl)(nil)

type FileHelper interface {
	GetFileType(fileExt string) (string, error)

	UserHasFilePermissions(userID string, fileID string, post *mmModel.Post) bool
	GetPostPermissionsByFileID(fileID string, post *mmModel.Post, getUser func(string) (*mmModel.User, *mmModel.AppError)) []model.UserInfoResponse
	GetFilePermissionsByUserID(userID string, fileID string, post *mmModel.Post) model.Permissions
	SetPostFilePermissions(post *mmModel.Post, permissions []model.PostPermission) []model.PostPermission

	IsExtensionSupported(fileExt string) bool
	IsExtensionEditable(fileExt string) bool

	GenerateKey() string

	GetWildcardUser() string
	GetWordType() string
	GetCellType() string
	GetSlideType() string
}

type fileHelperImpl struct {
	formatManager public.FormatManager
}

func New(formatManager public.FormatManager) FileHelper {
	return &fileHelperImpl{
		formatManager: formatManager,
	}
}

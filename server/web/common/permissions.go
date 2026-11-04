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
	mmModel "github.com/mattermost/mattermost/server/public/model"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	oomodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
)

// EffectivePermissions returns what the user can do with the file in the
// editor. It is the single source of truth for both the editor config and the
// file details shown in the preview, so the two can never disagree.
func EffectivePermissions(
	configuration *configuration.Configuration,
	fileHelper file.FileHelper,
	userID string,
	fileInfo *mmModel.FileInfo,
	post *mmModel.Post,
) oomodel.Permissions {
	permissions := oomodel.OnlyofficeDefaultPermissions
	if fileHelper.IsExtensionEditable(fileInfo.Extension) && configuration.IsFormatAllowedForEditing(fileInfo.Extension) {
		permissions = fileHelper.GetFilePermissionsByUserID(userID, fileInfo.Id, post)
	}

	if configuration.OwnerProtected {
		permissions.Protect = userID == post.UserId
	} else {
		permissions.Protect = true
	}

	return permissions
}

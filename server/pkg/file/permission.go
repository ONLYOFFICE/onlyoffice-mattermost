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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	mmModel "github.com/mattermost/mattermost/server/public/model"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
)

func (h fileHelperImpl) GetPostPermissionsByFileID(
	fileID string,
	post *mmModel.Post,
	getUser func(string) (*mmModel.User, *mmModel.AppError),
) []model.UserInfoResponse {
	props := post.GetProps()
	response := []model.UserInfoResponse{}

	for key, value := range props {
		if strings.HasPrefix(key, createPermissionsKeyPrefix(fileID)) {
			permissions, err := toPermissions(value)

			if err != nil {
				continue
			}

			keyParts := strings.Split(key, onlyofficePermissionsPropSeparator)
			userID := keyParts[len(keyParts)-1]

			if userID != onlyofficePermissionsWildcardKey {
				user, err := getUser(userID)
				if err != nil {
					continue
				}

				response = append(response, model.UserInfoResponse{
					ID:          userID,
					Username:    user.Username,
					Email:       user.Email,
					Permissions: permissions,
				})
			} else {
				response = append(response, model.UserInfoResponse{
					ID:          userID,
					Username:    userID,
					Email:       userID,
					Permissions: permissions,
				})
			}
		}
	}

	return response
}

func (h fileHelperImpl) GetFilePermissionsByUserID(userID string, fileID string, post *mmModel.Post) model.Permissions {
	if userID == post.UserId {
		return model.OnlyofficeAuthorPermissions
	}

	permissionsProp := post.GetProp(createPermissionsPropKeyName(userID, fileID))
	wildcardProp := post.GetProp(createPermissionsPropKeyName(onlyofficePermissionsWildcardKey, fileID))

	if permissionsProp == nil && wildcardProp == nil {
		return model.OnlyofficeDefaultPermissions
	}

	var permissions model.Permissions
	var uerr, werr error

	if permissionsProp != nil {
		permissions, uerr = toPermissions(permissionsProp)
		if uerr == nil {
			return permissions
		}
	}

	if wildcardProp != nil {
		permissions, werr = toPermissions(wildcardProp)
		if werr == nil {
			return permissions
		}
	}

	return model.OnlyofficeDefaultPermissions
}

func (h fileHelperImpl) UserHasFilePermissions(userID string, fileID string, post *mmModel.Post) bool {
	prop := post.GetProp(createPermissionsPropKeyName(userID, fileID))
	return prop != nil
}

func (h fileHelperImpl) SetPostFilePermissions(post *mmModel.Post, permissions []model.PostPermission) []model.PostPermission {
	notifyPermissions := make([]model.PostPermission, 0, len(permissions))
	persistPermissions := make([]model.PostPermission, 0, len(permissions))

	for _, permission := range permissions {
		if post.UserId == permission.UserID {
			continue
		}

		hadPermissions := h.UserHasFilePermissions(permission.UserID, permission.FileID, post)
		equalsUserPrevious := reflect.DeepEqual(permission.Permissions, h.GetFilePermissionsByUserID(permission.UserID, permission.FileID, post))
		equalsPreviousRead := reflect.DeepEqual(permission.Permissions, model.OnlyofficeDefaultPermissions)
		if (!hadPermissions && !equalsPreviousRead) || !equalsUserPrevious {
			notifyPermissions = append(notifyPermissions, permission)
		}

		persistPermissions = append(persistPermissions, permission)
	}

	for _, permission := range permissions {
		purgeFilePermissions(permission.FileID, post)
	}

	for _, permission := range persistPermissions {
		if post.FileIds.Contains(permission.FileID) {
			key := createPermissionsPropKeyName(permission.UserID, permission.FileID)
			setFilePermission(key, permission.Permissions, post)
		}
	}

	return notifyPermissions
}

func (h fileHelperImpl) GetWildcardUser() string {
	return onlyofficePermissionsWildcardKey
}

func createPermissionsPropKeyName(userID string, fileID string) string {
	return onlyofficePermissionsProp + onlyofficePermissionsPropSeparator +
		fileID + onlyofficePermissionsPropSeparator + userID
}

func createPermissionsKeyPrefix(fileID string) string {
	return onlyofficePermissionsProp + onlyofficePermissionsPropSeparator + fileID + onlyofficePermissionsPropSeparator
}

func setFilePermission(propKey string, permissions model.Permissions, post *mmModel.Post) {
	pbytes, err := json.Marshal(permissions)

	if err != nil {
		return
	}

	post.DelProp(propKey)
	post.AddProp(propKey, pbytes)
}

func purgeFilePermissions(fileID string, post *mmModel.Post) {
	props := post.GetProps()

	for name := range props {
		if strings.HasPrefix(name, createPermissionsKeyPrefix(fileID)) {
			post.DelProp(name)
		}
	}
}

func toPermissions(prop interface{}) (model.Permissions, error) {
	pjson, err := base64.StdEncoding.DecodeString(fmt.Sprintf("%v", prop))
	if err != nil {
		if pbytes, ok := prop.([]uint8); ok {
			var permissions model.Permissions

			if uerr := json.Unmarshal(pbytes, &permissions); uerr != nil {
				return permissions, uerr
			}

			return permissions, nil
		}

		return model.Permissions{}, err
	}

	var permissions model.Permissions
	return permissions, json.Unmarshal(pjson, &permissions)
}

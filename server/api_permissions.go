/**
 *
 * (c) Copyright Ascensio System SIA 2021
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

package main

//TODO: Refactoring. Move some part to utils
import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/utils"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/models"

	"github.com/mattermost/mattermost-server/v5/model"
)

func GetPostPermissionsByFileId(fileId string, post model.Post, p *Plugin) []models.UserInfoResponse {
	postProps := post.GetProps()
	response := []models.UserInfoResponse{}

	for key, value := range postProps {
		if strings.HasPrefix(key, utils.CreateFilePermissionsPrefix(fileId)) {
			convertedPermissions, err := utils.ConvertInterfaceToPermissions(value)
			if err != nil {
				continue
			}
			//TODO: Refactor
			keyParts := strings.Split(key, utils.ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR)
			userId := keyParts[3]

			if userId != utils.ONLYOFFICE_PERMISSIONS_WILDCARD_KEY {
				user, _ := p.API.GetUser(userId)
				response = append(response, models.UserInfoResponse{
					Id:          userId,
					Username:    user.Username,
					Email:       user.Email,
					Permissions: convertedPermissions,
				})
			} else {
				response = append(response, models.UserInfoResponse{
					Id:          userId,
					Username:    userId,
					Email:       "",
					Permissions: convertedPermissions,
				})
			}
		}
	}

	return response
}

func GetFilePermissionsByUser(userId string, fileId string, post model.Post) (models.Permissions, error) {
	if userId == post.UserId {
		return models.ONLYOFFICE_AUTHOR_PERMISSIONS, nil
	}

	ONLYOFFICE_USER_PERMISSIONS_PROP := post.GetProp(utils.CreateUserPermissionsPropName(fileId, userId))
	ONLYOFFICE_WILDCARD_PERMISSIONS_PROP := post.GetProp(utils.CreateWildcardPermissionsPropName(fileId))

	//If no permissions set, we want to grant default rights
	if ONLYOFFICE_USER_PERMISSIONS_PROP == nil && ONLYOFFICE_WILDCARD_PERMISSIONS_PROP == nil {
		return models.ONLYOFFICE_DEFAULT_PERMISSIONS, nil
	}

	if ONLYOFFICE_USER_PERMISSIONS_PROP != nil {
		return utils.ConvertInterfaceToPermissions(ONLYOFFICE_USER_PERMISSIONS_PROP)
	}

	return utils.ConvertInterfaceToPermissions(ONLYOFFICE_WILDCARD_PERMISSIONS_PROP)
}

func SetFilePermissions(post *model.Post, propKey string, filePermissions models.Permissions) {
	permissionBytes, _ := json.Marshal(filePermissions)

	post.DelProp(propKey)
	post.AddProp(propKey, permissionBytes)
}

func UserHasFilePermissions(userId string, fileId string, post *model.Post) bool {
	ONLYOFFICE_USER_PERMISSIONS_PROP := post.GetProp(utils.CreateUserPermissionsPropName(fileId, userId))
	return ONLYOFFICE_USER_PERMISSIONS_PROP != nil
}

func PurgeFilePermissions(post *model.Post, fileId string) {
	postProps := post.GetProps()

	for propName := range postProps {
		if strings.HasPrefix(propName, utils.CreateFilePermissionsPrefix(fileId)) {
			delete(postProps, propName)
		}
	}
}

func (p *Plugin) SetPostFilesPermissions(postPermissions []models.PostPermission, postId string) error {
	post, postErr := p.API.GetPost(postId)

	if postErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Invalid post id")
	}

	_, wildcardFiles := utils.ExtractUsernames(postPermissions)

	for fileId := range wildcardFiles {
		PurgeFilePermissions(post, fileId)
	}

	for _, postPermission := range postPermissions {
		if post.FileIds.Contains(postPermission.FileId) {
			if postPermission.Id == post.UserId {
				continue
			}
			propKey := utils.CreateUserPermissionsPropName(postPermission.FileId, postPermission.Id)
			SetFilePermissions(post, propKey, postPermission.Permissions)
			if utils.CompareUserAndWildcard(postPermission.Id) {
				propKey := utils.CreateWildcardPermissionsPropName(postPermission.FileId)
				SetFilePermissions(post, propKey, postPermission.Permissions)
			}
		}
	}

	p.API.UpdatePost(post)

	return nil
}

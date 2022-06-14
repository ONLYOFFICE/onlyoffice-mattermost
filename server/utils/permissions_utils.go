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

package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/models"
)

func GetPermissionsName(permissions models.Permissions) string {
	if permissions.Edit {
		return "Edit"
	}
	return "Read only"
}

func ConvertBase64ToPermissions(base64permissions string) (models.Permissions, error) {
	jsonPermissions, jsonErr := base64.StdEncoding.DecodeString(base64permissions)

	if jsonErr != nil {
		return models.Permissions{}, jsonErr
	}

	var permissions models.Permissions

	unmarshallingErr := json.Unmarshal(jsonPermissions, &permissions)

	if unmarshallingErr != nil {
		return models.Permissions{}, unmarshallingErr
	}

	return permissions, nil
}

func ConvertInterfaceToPermissions(ONLYOFFICE_PERMISSIONS_PROP interface{}) (models.Permissions, error) {
	base64permissions := fmt.Sprintf("%v", ONLYOFFICE_PERMISSIONS_PROP)

	permissions, permissionsErr := ConvertBase64ToPermissions(base64permissions)

	if permissionsErr != nil {
		return models.Permissions{}, permissionsErr
	}

	return permissions, nil
}

func ExtractUsernames(postPermissions []models.PostPermission) ([]string, map[string]bool) {
	var usernames []string = []string{}
	var wildcardFiles map[string]bool = make(map[string]bool)

	for _, postPermission := range postPermissions {
		if !CompareUserAndWildcard(postPermission.Username) {
			usernames = append(usernames, postPermission.Username)
		} else {
			wildcardFiles[postPermission.FileId] = true
		}
	}
	return usernames, wildcardFiles
}

func CompareUserAndWildcard(username string) bool {
	return username == ONLYOFFICE_PERMISSIONS_WILDCARD_KEY
}

func CreateUserPermissionsPropName(fileId string, userId string) string {
	return ONLYOFFICE_PERMISSIONS_PROP + ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR +
		fileId + ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR + userId
}

func CreateWildcardPermissionsPropName(fileId string) string {
	return ONLYOFFICE_PERMISSIONS_PROP + ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR + fileId +
		ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR + ONLYOFFICE_PERMISSIONS_WILDCARD_KEY
}

func CreateFilePermissionsPrefix(fileId string) string {
	return ONLYOFFICE_PERMISSIONS_PROP + ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR + fileId + ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR
}

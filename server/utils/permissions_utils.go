package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"models"
)

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

func CreateUserPermissionsPropName(fileId string, userId string, username string) string {
	return ONLYOFFICE_PERMISSIONS_PROP + ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR +
		fileId + ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR + userId + ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR + username
}

func CreateWildcardPermissionsPropName(fileId string) string {
	return ONLYOFFICE_PERMISSIONS_PROP + ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR + fileId +
		ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR + ONLYOFFICE_PERMISSIONS_WILDCARD_KEY +
		ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR + ONLYOFFICE_PERMISSIONS_WILDCARD_KEY
}

func CreateFilePermissionsPrefix(fileId string) string {
	return ONLYOFFICE_PERMISSIONS_PROP + ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR + fileId
}

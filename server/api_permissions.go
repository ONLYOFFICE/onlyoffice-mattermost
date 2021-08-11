package main

import (
	"encoding/json"
	"errors"
	"models"
	"strings"
	"utils"

	"github.com/mattermost/mattermost-server/v5/model"
)

func GetFilePermissionsByFileId(fileId string, post model.Post) []models.UserInfoResponse {
	postProps := post.GetProps()
	response := []models.UserInfoResponse{}

	for key, value := range postProps {
		if strings.HasPrefix(key, utils.CreateFilePermissionsPrefix(fileId)) {
			convertedPermissions, err := utils.ConvertInterfaceToPermissions(value)
			if err != nil {
				continue
			}
			//TODO: Refactor later
			keyParts := strings.Split(key, utils.ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR)
			response = append(response, models.UserInfoResponse{
				Id:          keyParts[3],
				Username:    keyParts[4],
				Permissions: convertedPermissions,
			})
		}
	}

	return response
}

func GetFilePermissionsByUser(userId string, username string, fileId string, post model.Post) (models.Permissions, error) {
	if userId == post.UserId {
		return models.ONLYOFFICE_AUTHOR_PERMISSIONS, nil
	}

	ONLYOFFICE_USER_PERMISSIONS_PROP := post.GetProp(utils.CreateUserPermissionsPropName(fileId, userId, username))
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

	usernames, wildcardFiles := utils.ExtractUsernames(postPermissions)

	users, usersErr := p.API.GetUsersByUsernames(usernames)

	if usersErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Invalid users while setting file permissions")
	}

	for fileId := range wildcardFiles {
		PurgeFilePermissions(post, fileId)
	}

	for _, postPermission := range postPermissions {
		if post.FileIds.Contains(postPermission.FileId) {
			for _, user := range users {
				if user.Id == post.UserId {
					continue
				}

				if user.Username == postPermission.Username {
					propKey := utils.CreateUserPermissionsPropName(postPermission.FileId, user.Id, user.Username)
					SetFilePermissions(post, propKey, postPermission.Permissions)
				}
			}
			if utils.CompareUserAndWildcard(postPermission.Username) {
				propKey := utils.CreateWildcardPermissionsPropName(postPermission.FileId)
				SetFilePermissions(post, propKey, postPermission.Permissions)
			}
		}
	}

	p.API.UpdatePost(post)

	return nil
}

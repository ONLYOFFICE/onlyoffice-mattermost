package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"models"

	"github.com/mattermost/mattermost-server/v5/model"
)

type PostPermission struct {
	FileId      string
	Username    string
	Permissions models.Permissions
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

func getFilePermissionsByUserId(userId string, fileId string, post model.Post) (models.Permissions, error) {
	if userId == post.UserId {
		return models.ONLYOFFICE_AUTHOR_PERMISSIONS, nil
	}

	ONLYOFFICE_USER_PERMISSIONS_PROP := post.GetProp(ONLYOFFICE_PERMISSIONS_PROP + "_" + userId + "_" + fileId)
	ONLYOFFICE_WILDCARD_PERMISSIONS_PROP := post.GetProp(ONLYOFFICE_PERMISSIONS_PROP + "_*_" + fileId)

	//If no permissions set, we want to grant default rights
	if ONLYOFFICE_USER_PERMISSIONS_PROP == nil && ONLYOFFICE_WILDCARD_PERMISSIONS_PROP == nil {
		return models.ONLYOFFICE_DEFAULT_PERMISSIONS, nil
	}

	if ONLYOFFICE_USER_PERMISSIONS_PROP != nil {
		return ConvertInterfaceToPermissions(ONLYOFFICE_USER_PERMISSIONS_PROP)
	}

	return ConvertInterfaceToPermissions(ONLYOFFICE_WILDCARD_PERMISSIONS_PROP)
}

func extractUsernames(postPermissions []PostPermission) ([]string, string) {
	var usernames []string = []string{}
	var wildcardKey string = ""
	for _, postPermission := range postPermissions {
		if postPermission.Username != ONLYOFFICE_PERMISSIONS_WILDCARD_KEY {
			usernames = append(usernames, postPermission.Username)
		} else {
			wildcardKey = postPermission.Username
		}
	}
	return usernames, wildcardKey
}

func setFilePermissions(post *model.Post, propKey string, filePermissions models.Permissions) {
	permissionBytes, _ := json.Marshal(filePermissions)

	post.DelProp(propKey)
	post.AddProp(propKey, permissionBytes)
}

func (p *Plugin) SetPostFilesPermissions(postPermissions []PostPermission, postId string) error {
	post, postErr := p.API.GetPost(postId)

	if postErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Invalid post id")
	}

	usernames, wildcardKey := extractUsernames(postPermissions)
	users, usersErr := p.API.GetUsersByUsernames(usernames)

	if usersErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Invalid users while setting file permissions")
	}

	for _, postPermission := range postPermissions {
		if post.FileIds.Contains(postPermission.FileId) {
			for _, user := range users {
				if user.Id == post.UserId {
					continue
				}

				if user.Username == postPermission.Username {
					propKey := ONLYOFFICE_PERMISSIONS_PROP + "_" + user.Id + "_" + postPermission.FileId

					setFilePermissions(post, propKey, postPermission.Permissions)
				}
			}
			if postPermission.Username == wildcardKey && wildcardKey != "" {
				propKey := ONLYOFFICE_PERMISSIONS_PROP + "_*_" + postPermission.FileId

				setFilePermissions(post, propKey, postPermission.Permissions)
			}
		}
	}

	p.API.UpdatePost(post)

	return nil
}

func (p *Plugin) PurgePostFilePermissions(postId string) error {
	post, postErr := p.API.GetPost(postId)

	if postErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Invalid post id")
	}

	post.DelProp(ONLYOFFICE_PERMISSIONS_PROP)
	return nil
}

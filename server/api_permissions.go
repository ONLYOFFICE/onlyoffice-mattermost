package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"models"

	"github.com/mattermost/mattermost-server/v5/model"
)

func ConvertBase64ToPermissionSet(base64permissions string) (map[string]map[string]models.Permissions, error) {
	jsonPermissions, jsonErr := base64.StdEncoding.DecodeString(base64permissions)

	if jsonErr != nil {
		return nil, jsonErr
	}

	var permissionsSet map[string]map[string]models.Permissions

	unmarshallingErr := json.Unmarshal(jsonPermissions, &permissionsSet)

	if unmarshallingErr != nil {
		return nil, unmarshallingErr
	}

	return permissionsSet, nil
}

func getFilePermissionsByUserId(userId string, fileId string, post model.Post) (models.Permissions, error) {
	if userId == post.UserId {
		return models.ONLYOFFICE_AUTHOR_PERMISSIONS, nil
	}

	ONLYOFFICE_PERMISSIONS_PROP := post.GetProp(ONLYOFFICE_PERMISSIONS_PROP)

	//If no permissions set, we want to grant default rights
	if ONLYOFFICE_PERMISSIONS_PROP == nil {
		return models.ONLYOFFICE_DEFAULT_PERMISSIONS, nil
	}

	base64permissions := fmt.Sprintf("%v", ONLYOFFICE_PERMISSIONS_PROP)

	permissionsSet, unmarshallingErr := ConvertBase64ToPermissionSet(base64permissions)

	if unmarshallingErr != nil {
		return models.Permissions{}, unmarshallingErr
	}

	_, fileExists := permissionsSet[fileId]

	if !fileExists {
		return models.Permissions{}, errors.New("No file with id: " + fileId + " in this post")
	}

	userPermissions, userPermissionsExists := permissionsSet[fileId][userId]

	if !userPermissionsExists {
		wildcardPermissions, wildcardPermissionsExist := permissionsSet[fileId][ONLYOFFICE_PERMISSIONS_WILDCARD_KEY]

		if !wildcardPermissionsExist {
			return models.ONLYOFFICE_DEFAULT_PERMISSIONS, nil
		}

		return wildcardPermissions, nil
	}

	return userPermissions, nil
}

type PostPermission struct {
	FileId      string
	Username    string
	Permissions models.Permissions
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

	ONLYOFFICE_PERMISSIONS := post.GetProp(ONLYOFFICE_PERMISSIONS_PROP)

	if ONLYOFFICE_PERMISSIONS != nil {
		base64permissions := fmt.Sprintf("%v", ONLYOFFICE_PERMISSIONS)

		filesPermissions, unmarshallingErr := ConvertBase64ToPermissionSet(base64permissions)

		if unmarshallingErr != nil {
			return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Conversion error while setting file permissions")
		}

		for _, postPermission := range postPermissions {
			if post.FileIds.Contains(postPermission.FileId) {
				if filesPermissions[postPermission.FileId] == nil {
					filesPermissions[postPermission.FileId] = make(map[string]models.Permissions)
				}
				for _, user := range users {
					if user.Id == post.UserId {
						continue
					}

					if user.Username == postPermission.Username {
						filesPermissions[postPermission.FileId][user.Id] = postPermission.Permissions
					}
				}
				if postPermission.Username == wildcardKey && wildcardKey != "" {
					filesPermissions[postPermission.FileId][wildcardKey] = postPermission.Permissions
				}
			}
		}

		jsonPermissionsBytes, bytesErr := json.Marshal(filesPermissions)

		if bytesErr != nil {
			return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Conversion error while setting file permissions")
		}

		post.DelProp(ONLYOFFICE_PERMISSIONS_PROP)
		post.AddProp(ONLYOFFICE_PERMISSIONS_PROP, jsonPermissionsBytes)

		p.API.UpdatePost(post)
		return nil
	}

	var filesPermissions map[string]map[string]models.Permissions = make(map[string]map[string]models.Permissions)

	for _, postPermission := range postPermissions {
		filesPermissions[postPermission.FileId] = make(map[string]models.Permissions)
		for _, user := range users {
			if user.Id == post.UserId {
				continue
			}

			if user.Username == postPermission.Username {
				filesPermissions[postPermission.FileId][user.Id] = postPermission.Permissions
			}
		}
		if postPermission.Username == wildcardKey && wildcardKey != "" {
			filesPermissions[postPermission.FileId][wildcardKey] = postPermission.Permissions
		}
	}

	jsonPermissionsBytes, bytesErr := json.Marshal(filesPermissions)

	if bytesErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Conversion error while setting file permissions")
	}

	post.AddProp(ONLYOFFICE_PERMISSIONS_PROP, jsonPermissionsBytes)
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

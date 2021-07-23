package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"models"
	"utils"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
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
		return utils.ONLYOFFICE_AUTHOR_PERMISSIONS, nil
	}

	ONLYOFFICE_PERMISSIONS_PROP := post.GetProp(utils.ONLYOFFICE_PERMISSIONS_PROP)

	//If no permissions set, we want to grant default rights
	if ONLYOFFICE_PERMISSIONS_PROP == nil {
		return utils.ONLYOFFICE_ALL_USERS_PERMISSIONS, nil
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
		wildcardPermissions, wildcardPermissionsExist := permissionsSet[fileId][utils.ONLYOFFICE_PERMISSIONS_WILDCARD_KEY]

		if !wildcardPermissionsExist {
			return utils.ONLYOFFICE_ALL_USERS_PERMISSIONS, nil
		}

		return wildcardPermissions, nil
	}

	return userPermissions, nil
}

//TODO: Refactor
func (p *Plugin) SetFilePermissionsByUsername(username string, fileId string, permissions models.Permissions) error {
	fileInfo, fileInfoErr := p.API.GetFileInfo(fileId)

	if fileInfoErr != nil {
		return fileInfoErr
	}

	user, userErr := p.API.GetUserByUsername(username)

	if userErr != nil {
		return userErr
	}

	post, _ := p.API.GetPost(fileInfo.PostId)

	//Do not change Author's permissions
	if post.UserId == user.Id {
		return nil
	}

	ONLYOFFICE_PERMISSIONS_PROP := post.GetProp(utils.ONLYOFFICE_PERMISSIONS_PROP)

	if ONLYOFFICE_PERMISSIONS_PROP != nil {
		base64permissions := fmt.Sprintf("%v", ONLYOFFICE_PERMISSIONS_PROP)
		filePermissions, unmarshallingErr := ConvertBase64ToPermissionSet(base64permissions)

		if unmarshallingErr != nil {
			return unmarshallingErr
		}

		filePermissions[fileId][user.Id] = permissions
		jsonPermissions, _ := json.Marshal(filePermissions)

		post.DelProp(utils.ONLYOFFICE_PERMISSIONS_PROP)
		post.AddProp(utils.ONLYOFFICE_PERMISSIONS_PROP, jsonPermissions)

		p.API.UpdatePost(post)
		return nil
	}

	var filePermissions map[string]map[string]models.Permissions = map[string]map[string]models.Permissions{
		fileId: {
			user.Id: permissions,
		},
	}

	jsonPermissionsBytes, marshallingErr := json.Marshal(filePermissions)

	if marshallingErr != nil {
		return marshallingErr
	}

	post.AddProp(utils.ONLYOFFICE_PERMISSIONS_PROP, jsonPermissionsBytes)
	p.API.UpdatePost(post)

	return nil
}

func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	return post, ""
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	//TODO: Custom file permissions logic on upload (with a file wizzard)
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Plugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {

	return info, ""
}

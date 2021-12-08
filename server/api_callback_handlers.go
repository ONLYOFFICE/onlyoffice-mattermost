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

import (
	"errors"
	"strings"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/utils"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/security"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/models"

	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
)

//TODO: Generalize the function
func ConvertJwtToBody(body *models.CallbackBody, jwtKey []byte, jwtString string) error {
	var decodedCallback jwt.MapClaims
	var jwtDecodingErr error

	if jwtString == "" {
		decodedCallback, jwtDecodingErr = security.JwtDecode(body.Token, jwtKey)
	} else {
		decodedCallback, jwtDecodingErr = security.JwtDecode(strings.Split(jwtString, " ")[1], jwtKey)
	}

	if jwtDecodingErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Could not process JWT in callback body")
	}

	err := mapstructure.Decode(decodedCallback, &body)

	if err != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Could not populate callback body with decoded JWT")
	}

	return nil
}

//Status 2 and 6
func handleSave(body *models.CallbackBody, p *Plugin) error {
	url := body.Url

	if url == "" {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Invalid download URL")
	}

	response, getErr := p.GetHTTPClient().GetRequest(url)
	if getErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Could not download the file requested")
	}

	file := response.Body
	defer file.Close()

	fileID := body.FileId

	if fileID == "" {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Invalid file id")
	}

	fileInfo, fileInfoErr := p.API.GetFileInfo(fileID)
	if fileInfoErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Could not find given file's FileInfo")
	}

	_, filestoreErr := p.WriteFile(file, fileInfo.Path)

	if filestoreErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Filestore error when writing changes to the file " + fileInfo.Name)
	}

	if body.Status == 2 {
		post, _ := p.API.GetPost(fileInfo.PostId)
		post.UpdateAt = utils.GetTimestamp()
		p.API.UpdatePost(post)

		last := body.Users[0]

		if last == "" {
			return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Invalid callback user")
		}

		user, err := p.API.GetUser(last)

		if err != nil {
			return err
		}

		var newReplyMessage string = "File " + fileInfo.Name + " was updated" + " by @" + user.Username

		p.onlyoffice_bot.BOT_CREATE_REPLY(newReplyMessage, post.ChannelId, post.Id)
	}

	return nil
}

//Status 4
func handleNoChanges(body *models.CallbackBody, p *Plugin) error {
	return nil
}

//Status 1
func handleIsBeingEdited(body *models.CallbackBody, p *Plugin) error {
	return nil
}

//Status 3
func handleSavingError(body *models.CallbackBody, p *Plugin) error {
	return nil
}

//Status 7
func handleForcesavingError(body *models.CallbackBody, p *Plugin) error {
	return nil
}

func (p *Plugin) getCallbackHandler(callbackBody *models.CallbackBody) (func(body *models.CallbackBody, plugin *Plugin) error, bool) {
	docServerStatus := map[int]func(body *models.CallbackBody, plugin *Plugin) error{
		1: handleIsBeingEdited,
		2: handleSave,
		3: handleSavingError,
		4: handleNoChanges,
		6: handleSave,
		7: handleForcesavingError,
	}

	handler, exists := docServerStatus[callbackBody.Status]

	return handler, exists
}

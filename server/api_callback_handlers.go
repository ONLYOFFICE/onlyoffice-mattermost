package main

import (
	"errors"
	"models"
	"security"
	"utils"

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
		decodedCallback, jwtDecodingErr = security.JwtDecode(jwtString, jwtKey)
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
	var url string = body.Url

	response, getErr := p.GetHTTPClient().GetRequest(url)
	if getErr != nil {
		return errors.New(ONLYOFFICE_LOGGER_PREFIX + "Could not download the file requested")
	}

	file := response.Body
	defer file.Close()

	fileInfo, fileInfoErr := p.API.GetFileInfo(body.FileId)
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

		decryptedUser, _ := security.EncryptorAES{}.Decrypt(body.Users[0], p.internalKey)
		user, _ := p.API.GetUser(decryptedUser)
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

package main

import (
	"models"
	"utils"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mitchellh/mapstructure"
)

func processJwtBody(body *models.CallbackBody, jwtKey []byte) error {
	decodedCallback, jwtDecodingErr := utils.JwtDecode(body.Token, jwtKey)

	if jwtDecodingErr != nil {
		return jwtDecodingErr
	}

	err := mapstructure.Decode(decodedCallback, &body)

	if err != nil {
		return err
	}

	return nil
}

//Status 2 and 6
func handleSave(body *models.CallbackBody, p *Plugin) {
	var url string = body.Url
	response, getErr := p.GetHTTPClient().GetRequest(url)

	if getErr != nil {
		p.API.LogError("[ONLYOFFICE] Couldn't fetch the file requested")
		return
	}

	file := response.Body
	defer file.Close()

	fileInfo, err := p.API.GetFileInfo(body.FileId)

	if err != nil {
		p.API.LogError("[ONLYOFFICE]: Fileinfo error - ", err.Error())
		return
	}

	_, exception := p.WriteFile(file, fileInfo.Path)

	if exception != nil {
		p.API.LogError("[ONLYOFFICE]: Filestore error - ", exception.Error())
		return
	}

	if body.Status == 2 {
		post, _ := p.API.GetPost(fileInfo.PostId)
		post.UpdateAt = utils.GetTimestamp()
		p.API.UpdatePost(post)
		//TODO: Move to a separate function
		user, _ := p.API.GetUser(body.Users[0])
		var newPostMessage string = "File " + fileInfo.Name + " was updated" + " by @" + user.Username
		newPost := model.Post{
			Message:   newPostMessage,
			ParentId:  post.Id,
			RootId:    post.Id,
			ChannelId: post.ChannelId,
			UserId:    p.onlyoffice_bot_id,
		}
		_, creationErr := p.API.CreatePost(&newPost)
		if creationErr != nil {
			p.API.LogError("[ONLYOFFICE] Post creation error: ", creationErr.Error())
			return
		}
	}
}

//Status 4
func handleNoChanges(body *models.CallbackBody, p *Plugin) {
}

//Status 1
func handleIsBeingEdited(body *models.CallbackBody, p *Plugin) {
}

//Status 3
func handleSavingError(body *models.CallbackBody, p *Plugin) {

}

//Status 7
func handleForcesavingError(body *models.CallbackBody, p *Plugin) {

}

func (p *Plugin) getCallbackHandler(callbackBody *models.CallbackBody) (func(body *models.CallbackBody, plugin *Plugin), bool) {
	docServerStatus := map[int]func(body *models.CallbackBody, plugin *Plugin){
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

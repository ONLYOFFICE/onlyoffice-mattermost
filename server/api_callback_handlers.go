package main

import (
	"io"
	"models"
	"utils"

	"github.com/mattermost/mattermost-server/v5/shared/filestore"
)

//Status 2 and 6
func (p *Plugin) handleSave(body *models.CallbackBody) {
	var url string = body.Url
	var file io.ReadCloser = p.GetHTTPClient().GetRequest(url)

	defer file.Close()

	serverConfig := p.API.GetUnsanitizedConfig()
	filestore, _ := filestore.NewFileBackend(serverConfig.FileSettings.ToFileBackendSettings(false))

	fileInfo, err := p.API.GetFileInfo(body.FileId)

	if err != nil {
		p.API.LogError("[ONLYOFFICE]: Fileinfo error - ", err.Error())
	}

	_, exception := filestore.WriteFile(file, fileInfo.Path)

	if exception != nil {
		p.API.LogError("[ONLYOFFICE]: Filestore error - ", exception.Error())
		return
	}

	//TODO: To a separate function
	if body.Status == 2 {
		post, _ := p.API.GetPost(fileInfo.PostId)
		post.UpdateAt = utils.GetTimestamp()
		p.API.UpdatePost(post)
	}
}

//Status 4
func (p *Plugin) handleNoChanges(body *models.CallbackBody) {
}

//Status 1
func (p *Plugin) handleIsBeingEdited(body *models.CallbackBody) {
}

//Status 3
func (p *Plugin) handleSavingError(body *models.CallbackBody) {

}

//Status 7
func (p *Plugin) handleForcesavingError(body *models.CallbackBody) {

}

func (p *Plugin) getCallbackHandler(callbackBody *models.CallbackBody) (func(body *models.CallbackBody), bool) {
	docServerStatus := map[int]func(body *models.CallbackBody){
		1: p.handleIsBeingEdited,
		2: p.handleSave,
		3: p.handleSavingError,
		4: p.handleNoChanges,
		6: p.handleSave,
		7: p.handleForcesavingError,
	}

	handler, exists := docServerStatus[callbackBody.Status]

	return handler, exists
}

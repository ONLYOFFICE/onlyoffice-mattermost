package main

import (
	"io"
	"models"
	"utils"

	"github.com/mattermost/mattermost-server/v5/shared/filestore"
)

//Status 2 and 6
func handleSave(body *models.CallbackBody, p *Plugin) {
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

	if body.Status == 2 {
		post, _ := p.API.GetPost(fileInfo.PostId)
		post.UpdateAt = utils.GetTimestamp()
		p.API.UpdatePost(post)
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

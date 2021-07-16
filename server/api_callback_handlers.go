package main

import (
	"dto"
	"io"

	"github.com/mattermost/mattermost-server/v5/shared/filestore"
)

//Status 2 and 6
func (p *Plugin) handleSave(body *dto.CallbackBody) {
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
	}

	if body.Status == 2 {
		p.globalCache.Delete("ONLYOFFICE_" + body.FileId)
	}
}

//Status 4
func (p *Plugin) handleNoChanges(body *dto.CallbackBody) {
	_, found := p.globalCache.Get("ONLYOFFICE_" + body.FileId)

	if found {
		p.globalCache.Delete("ONLYOFFICE_" + body.FileId)
	}
}

//Status 1
func (p *Plugin) handleIsBeingEdited(body *dto.CallbackBody) {

}

//Status 3
func (p *Plugin) handleSavingError(body *dto.CallbackBody) {

}

//Status 7
func (p *Plugin) handleForcesavingError(body *dto.CallbackBody) {

}

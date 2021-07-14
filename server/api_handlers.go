package main

import (
	"dto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"text/template"
	"utils"

	"github.com/mattermost/mattermost-server/v5/shared/filestore"
	"github.com/patrickmn/go-cache"
)

func (p *Plugin) editor(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		p.API.LogError("[ONLYOFFICE]: Editor error ", err.Error())
		return
	}

	var fileId string = request.PostForm.Get("fileid")

	docKey, found := p.globalCache.Get("ONLYOFFICE_" + fileId)
	userId, _ := request.Cookie("MMUSERID")

	if !found {
		docKey = utils.GenerateKey()
		p.globalCache.Set("ONLYOFFICE_"+fileId, docKey, cache.NoExpiration)
	}

	fileInfo, _ := p.API.GetFileInfo(fileId)
	user, _ := p.API.GetUser(userId.Value)
	//TODO: Use id from manifest
	var serverURL string = *p.API.GetConfig().ServiceSettings.SiteURL

	bundlePath, _ := p.API.GetBundlePath()

	temp := template.New("onlyoffice")
	temp, _ = temp.ParseFiles(filepath.Join(bundlePath, "public/editor.html"))

	fileId, _ = p.encryptAES(fileId, p.internalKey)

	data := map[string]interface{}{
		"apijs":        p.configuration.DESAddress + utils.DESApijs,
		"key":          docKey,
		"title":        fileInfo.Name,
		"fileType":     fileInfo.Extension,
		"fileId":       fileId,
		"documentType": utils.GetFileType(fileInfo.Extension),
		"serverURL":    serverURL,
		"userId":       userId.Value,
		"username":     user.Username,
	}
	temp.ExecuteTemplate(writer, "editor.html", data)
}

func (p *Plugin) saveFile(body *dto.CallbackBody) {
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

func noImplementation(body *dto.CallbackBody) {
	fmt.Println("[ONLYOFFICE]: No implementation")
}

func (p *Plugin) noChangesClosed(body *dto.CallbackBody) {
	_, found := p.globalCache.Get("ONLYOFFICE_" + body.FileId)

	if found {
		p.globalCache.Delete("ONLYOFFICE_" + body.FileId)
	}
}

func (p *Plugin) callback(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	response := "{\"error\": 0}"

	body := dto.CallbackBody{}
	json.NewDecoder(request.Body).Decode(&body)

	fileId, _ := p.decryptAES(query.Get("fileId"), p.internalKey)

	body.FileId = fileId

	docServerStatus := map[int]func(body *dto.CallbackBody){
		1: noImplementation,
		2: p.saveFile,
		4: p.noChangesClosed,
		6: p.saveFile,
	}

	handler, exists := docServerStatus[body.Status]

	if !exists {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(500)
		writer.Write([]byte(response))
	}

	handler(&body)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write([]byte(response))
}

func (p *Plugin) downloadFile(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	fileId, _ := p.decryptAES(query.Get("fileId"), p.internalKey)
	fileContent, _ := p.API.GetFile(fileId)

	writer.Write(fileContent)
}

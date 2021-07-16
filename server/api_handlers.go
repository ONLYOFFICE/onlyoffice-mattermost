package main

import (
	"dto"
	"encoders"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
	"utils"

	"github.com/patrickmn/go-cache"
)

func (p *Plugin) editor(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		p.API.LogError("[ONLYOFFICE]: Editor error ", err.Error())
		return
	}

	var fileId string = request.PostForm.Get("fileid")
	var docKey string = p.getDocKey(fileId)
	fileInfo, _ := p.API.GetFileInfo(fileId)

	userId, _ := request.Cookie("MMUSERID")
	user, _ := p.API.GetUser(userId.Value)

	var serverURL string = *p.API.GetConfig().ServiceSettings.SiteURL + "/" + utils.MMPluginApi

	temp := template.New("onlyoffice")
	bundlePath, _ := p.API.GetBundlePath()
	temp, _ = temp.ParseFiles(filepath.Join(bundlePath, "public/editor.html"))

	p.encoder = encoders.EncoderAES{}
	fileId, _ = p.encoder.Encode(fileId, p.internalKey)

	var config dto.Config = dto.Config{
		Document: dto.Document{
			FileType: fileInfo.Extension,
			Key:      docKey,
			Title:    fileInfo.Name,
			Url:      serverURL + "/download?fileId=" + fileId,
		},
		DocumentType: utils.GetFileType(fileInfo.Extension),
		EditorConfig: dto.EditorConfig{
			User: dto.User{
				Id:   userId.Value,
				Name: user.Username,
			},
			CallbackUrl: serverURL + "/callback?fileId=" + fileId,
		},
	}

	jwtString, _ := utils.JwtSign(config, []byte(p.configuration.DESJwt))

	config.Token = jwtString

	data := map[string]interface{}{
		"apijs":  p.configuration.DESAddress + utils.DESApijs,
		"config": config,
	}

	temp.ExecuteTemplate(writer, "editor.html", data)
}

func (p *Plugin) callback(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	response := "{\"error\": 0}"

	body := dto.CallbackBody{}
	json.NewDecoder(request.Body).Decode(&body)

	handler, exists := p.getCallbackHandler(&body)

	p.encoder = encoders.EncoderAES{}
	fileId, _ := p.encoder.Decode(query.Get("fileId"), p.internalKey)
	body.FileId = fileId

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

func (p *Plugin) download(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	p.encoder = encoders.EncoderAES{}
	fileId, _ := p.encoder.Decode(query.Get("fileId"), p.internalKey)
	fileContent, _ := p.API.GetFile(fileId)

	writer.Write(fileContent)
}

func (p *Plugin) getDocKey(fileId string) string {
	docKey, found := p.globalCache.Get("ONLYOFFICE_" + fileId)
	if !found {
		docKey = utils.GenerateKey()
		p.globalCache.Set("ONLYOFFICE_"+fileId, docKey, cache.NoExpiration)
	}
	return fmt.Sprintf("%v", docKey)
}

func (p *Plugin) getCallbackHandler(callbackBody *dto.CallbackBody) (func(body *dto.CallbackBody), bool) {
	docServerStatus := map[int]func(body *dto.CallbackBody){
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

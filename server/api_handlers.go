package main

import (
	"encoding/json"
	"encryptors"
	"models"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"
	"utils"
)

func (p *Plugin) editor(writer http.ResponseWriter, request *http.Request) {
	var fileId string = request.PostForm.Get("fileid")

	var docKey string = p.generateDocKey(fileId)

	fileInfo, _ := p.API.GetFileInfo(fileId)
	docType, _ := utils.GetFileType(fileInfo.Extension)

	userId, _ := request.Cookie(utils.MMUserCookie)
	user, _ := p.API.GetUser(userId.Value)

	var serverURL string = *p.API.GetConfig().ServiceSettings.SiteURL + "/" + utils.MMPluginApi

	temp := template.New("onlyoffice")
	bundlePath, _ := p.API.GetBundlePath()
	temp, _ = temp.ParseFiles(filepath.Join(bundlePath, "public/editor.html"))

	p.encryptor = encryptors.EncryptorAES{}
	fileId, _ = p.encryptor.Encrypt(fileId, p.internalKey)

	var config models.Config = models.Config{
		Document: models.Document{
			FileType: fileInfo.Extension,
			Key:      docKey,
			Title:    fileInfo.Name,
			Url:      serverURL + "/download?fileId=" + fileId,
		},
		DocumentType: docType,
		EditorConfig: models.EditorConfig{
			User: models.User{
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

	body := models.CallbackBody{}
	json.NewDecoder(request.Body).Decode(&body)

	handler, exists := p.getCallbackHandler(&body)

	p.encryptor = encryptors.EncryptorAES{}
	fileId, _ := p.encryptor.Decrypt(query.Get("fileId"), p.internalKey)
	body.FileId = fileId

	if !exists {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(500)
		writer.Write([]byte(response))
	}

	handler(&body, p)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write([]byte(response))
}

func (p *Plugin) download(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	p.encryptor = encryptors.EncryptorAES{}
	fileId, _ := p.encryptor.Decrypt(query.Get("fileId"), p.internalKey)
	fileContent, _ := p.API.GetFile(fileId)

	writer.Write(fileContent)
}

func (p *Plugin) generateDocKey(fileId string) string {
	fileInfo, err := p.API.GetFileInfo(fileId)
	if err != nil {
		return ""
	}

	post, _ := p.API.GetPost(fileInfo.PostId)

	var postUpdatedAt string = strconv.FormatInt(post.UpdateAt, 10)

	p.encryptor = encryptors.EncryptorRC4{}
	docKey, encodeErr := p.encryptor.Encrypt(fileId+postUpdatedAt, []byte(utils.RC4Key))

	if encodeErr != nil {
		p.API.LogError("[ONLYOFFICE] Key generation error: ", encodeErr.Error())
		return ""
	}
	return docKey
}

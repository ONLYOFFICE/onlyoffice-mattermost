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

	post, _ := p.API.GetPost(fileInfo.PostId)

	userPermissions, _ := getFilePermissionsByUserId(userId.Value, fileInfo.Id, *post)

	var config models.Config = models.Config{
		Document: models.Document{
			FileType: fileInfo.Extension,
			Key:      docKey,
			Title:    fileInfo.Name,
			Url:      serverURL + "/download?fileId=" + fileId,
			P:        userPermissions,
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

	if p.configuration.DESJwt != "" {
		jwtString, _ := utils.JwtSign(config, []byte(p.configuration.DESJwt))
		config.Token = jwtString
	}

	jsonBytes, _ := json.Marshal(config)
	jsonConfig := string(jsonBytes)

	data := map[string]interface{}{
		"apijs":  p.configuration.DESAddress + utils.DESApijs,
		"config": jsonConfig,
	}

	temp.ExecuteTemplate(writer, "editor.html", data)
}

func (p *Plugin) callback(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	response := "{\"error\": 0}"

	body := models.CallbackBody{}

	//TODO: Refactor
	decodingErr := json.NewDecoder(request.Body).Decode(&body)

	if decodingErr != nil {
		p.API.LogError("[ONLYOFFICE] Callback body decoding error - ", decodingErr.Error())
		http.Error(writer, response, http.StatusInternalServerError)
		return
	}

	if p.configuration.DESJwt != "" {
		jwtBodyHandlerErr := processJwtBody(&body, []byte(p.configuration.DESJwt))

		if jwtBodyHandlerErr != nil {
			p.API.LogError("[ONLYOFFICE] JWT Body processing error - ", jwtBodyHandlerErr.Error())
			http.Error(writer, response, http.StatusInternalServerError)
			return
		}
	}

	handler, exists := p.getCallbackHandler(&body)

	p.encryptor = encryptors.EncryptorAES{}
	fileId, decryptionErr := p.encryptor.Decrypt(query.Get("fileId"), p.internalKey)
	body.FileId = fileId

	if !exists || decryptionErr != nil {
		http.Error(writer, response, http.StatusInternalServerError)
		return
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

func (p *Plugin) permissions(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	fileInfo, fileInfoErr := p.API.GetFileInfo(query.Get("fileId"))
	user, userErr := p.API.GetUserByUsername(query.Get("username"))

	if fileInfoErr != nil || userErr != nil {
		writer.WriteHeader(400)
		return
	}

	userId, _ := request.Cookie(utils.MMUserCookie)

	if fileInfo.CreatorId != userId.Value {
		writer.WriteHeader(403)
		return
	}

	body := models.Permissions{}

	//TODO: Refactor
	decodingErr := json.NewDecoder(request.Body).Decode(&body)

	if decodingErr != nil {
		writer.WriteHeader(500)
		p.API.LogError("[ONLYOFFICE] Permissions endpoint error: Decoding error")
		return
	}

	setPermissionsErr := p.SetFilePermissionsByUsername(user.Username, fileInfo.Id, body)

	if setPermissionsErr != nil {
		writer.WriteHeader(500)
		p.API.LogError("[ONLYOFFICE] Permissions endpoint error: Permissions update error")
		return
	}

	writer.WriteHeader(200)
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

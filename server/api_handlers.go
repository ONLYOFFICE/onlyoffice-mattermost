package main

import (
	"encoding/json"
	"models"
	"net/http"
	"path/filepath"
	"security"
	"strconv"
	"text/template"
	"utils"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (p *Plugin) editor(writer http.ResponseWriter, request *http.Request) {
	var serverURL string = *p.API.GetConfig().ServiceSettings.SiteURL + "/" + ONLYOFFICE_API_PATH
	var fileId string = request.PostForm.Get("fileid")
	fileInfo, _ := p.API.GetFileInfo(fileId)
	docType, _ := utils.GetFileType(fileInfo.Extension)

	//We expect only authorized by middlewares users
	var userId string = request.Header.Get(ONLYOFFICE_AUTHORIZATION_USERID_HEADER)
	var username string = request.Header.Get(ONLYOFFICE_AUTHORIZATION_USERNAME_HEADER)

	htmlTemplate := template.New("onlyoffice")
	bundlePath, _ := p.API.GetBundlePath()
	htmlTemplate, _ = htmlTemplate.ParseFiles(filepath.Join(bundlePath, "public/editor.html"))

	var encryptor security.Encryptor = security.EncryptorAES{}
	fileId, _ = encryptor.Encrypt(fileId, p.internalKey)

	post, _ := p.API.GetPost(fileInfo.PostId)

	var docKey string = generateDocKey(*fileInfo, *post)

	userPermissions, _ := getFilePermissionsByUserId(userId, fileInfo.Id, *post)

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
				Id:   userId,
				Name: username,
			},
			CallbackUrl: serverURL + "/callback?fileId=" + fileId,
		},
	}

	//TODO: Think up a better JWT logic
	if p.configuration.DESJwt != "" {
		jwtString, _ := security.JwtSign(config, []byte(p.configuration.DESJwt))
		config.Token = jwtString
	}

	jsonBytes, _ := json.Marshal(config)
	jsonConfig := string(jsonBytes)

	data := map[string]interface{}{
		"apijs":  p.configuration.DESAddress + ONLYOFFICE_API_JS,
		"config": jsonConfig,
	}

	htmlTemplate.ExecuteTemplate(writer, "editor.html", data)
}

func (p *Plugin) callback(writer http.ResponseWriter, request *http.Request) {
	body := models.CallbackBody{}
	decodingErr := json.NewDecoder(request.Body).Decode(&body)

	if decodingErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Callback body decoding error")
		writer.WriteHeader(500)
		return
	}

	if p.configuration.DESJwt != "" {
		jwtBodyProcessingErr := processJwtBody(&body, []byte(p.configuration.DESJwt))

		if jwtBodyProcessingErr != nil {
			p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "JWT Body processing error")
			writer.WriteHeader(500)
			return
		}
	}

	handler, exists := p.getCallbackHandler(&body)

	if !exists {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Could not find a proper callback handler")
		writer.WriteHeader(500)
		return
	}

	var encryptor security.Encryptor = security.EncryptorAES{}
	fileId, _ := encryptor.Decrypt(request.URL.Query().Get("fileId"), p.internalKey)
	body.FileId = fileId

	handler(&body, p)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write([]byte("{\"error\": 0}"))
}

func (p *Plugin) download(writer http.ResponseWriter, request *http.Request) {
	var encryptor security.Encryptor = security.EncryptorAES{}
	fileId, _ := encryptor.Decrypt(request.URL.Query().Get("fileId"), p.internalKey)
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

	var userId string = request.Header.Get(ONLYOFFICE_AUTHORIZATION_USERID_HEADER)

	if fileInfo.CreatorId != userId {
		writer.WriteHeader(403)
		return
	}

	body := models.Permissions{}
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

func generateDocKey(fileInfo model.FileInfo, post model.Post) string {
	var postUpdatedAt string = strconv.FormatInt(post.UpdateAt, 10)

	var encryptor security.Encryptor = security.EncryptorRC4{}
	docKey, _ := encryptor.Encrypt(fileInfo.Id+postUpdatedAt, []byte(ONLYOFFICE_RC4_KEY))

	return docKey
}

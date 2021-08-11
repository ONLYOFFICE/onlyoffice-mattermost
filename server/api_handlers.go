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

	var userId string = request.Header.Get(ONLYOFFICE_AUTHORIZATION_USERID_HEADER)
	var username string = request.Header.Get(ONLYOFFICE_AUTHORIZATION_USERNAME_HEADER)

	htmlTemplate := template.New("onlyoffice")
	bundlePath, _ := p.API.GetBundlePath()
	htmlTemplate, _ = htmlTemplate.ParseFiles(filepath.Join(bundlePath, "public/editor.html"))

	encryptor := security.EncryptorAES{}
	fileId, _ = encryptor.Encrypt(fileId, p.internalKey)
	userIdEnc, _ := encryptor.Encrypt(userId, p.internalKey)

	post, _ := p.API.GetPost(fileInfo.PostId)

	var docKey string = GenerateDocKey(*fileInfo, *post)

	userPermissions, _ := GetFilePermissionsByUser(userId, username, fileInfo.Id, *post)

	var config models.Config = models.Config{
		Document: models.Document{
			FileType: fileInfo.Extension,
			Key:      docKey,
			Title:    fileInfo.Name,
			Url:      serverURL + ONLYOFFICE_ROUTE_DOWNLOAD + "?fileId=" + fileId,
			P:        userPermissions,
		},
		DocumentType: docType,
		EditorConfig: models.EditorConfig{
			User: models.User{
				Id:   userIdEnc,
				Name: username,
			},
			CallbackUrl: serverURL + ONLYOFFICE_ROUTE_CALLBACK + "?fileId=" + fileId,
			Customization: models.Customization{
				Goback: models.Goback{
					RequestClose: true,
				},
			},
		},
	}

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
		jwtBodyProcessingErr := ConvertJwtToBody(&body, []byte(p.configuration.DESJwt), request.Header.Get(p.configuration.DESJwtHeader))

		if jwtBodyProcessingErr != nil {
			p.API.LogError(jwtBodyProcessingErr.Error())
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

	fileId, _ := security.EncryptorAES{}.Decrypt(request.URL.Query().Get("fileId"), p.internalKey)
	body.FileId = fileId

	handler(&body, p)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write([]byte("{\"error\": 0}"))
}

func (p *Plugin) download(writer http.ResponseWriter, request *http.Request) {
	fileId, _ := security.EncryptorAES{}.Decrypt(request.URL.Query().Get("fileId"), p.internalKey)
	fileContent, _ := p.API.GetFile(fileId)

	writer.Write(fileContent)
}

func (p *Plugin) setFilePermissions(writer http.ResponseWriter, request *http.Request) {
	var postPermissionsBody []models.PostPermission = []models.PostPermission{}

	decodingErr := json.NewDecoder(request.Body).Decode(&postPermissionsBody)
	if decodingErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Permissions body decoding error")
		writer.WriteHeader(400)
		return
	}

	if len(postPermissionsBody) == 0 {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Invalid permissions body length")
		writer.WriteHeader(400)
		return
	}

	fileInfo, fileInfoErr := p.API.GetFileInfo(postPermissionsBody[0].FileId)

	if fileInfoErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Invalid file id in permissions body")
		writer.WriteHeader(400)
		return
	}

	post, postErr := p.API.GetPost(fileInfo.PostId)

	if postErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Invalid post id")
		writer.WriteHeader(400)
		return
	}

	var userId string = request.Header.Get(ONLYOFFICE_AUTHORIZATION_USERID_HEADER)

	if post.UserId != userId {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Only post's author can change file permissions")
		writer.WriteHeader(403)
		return
	}

	setPermissionsErr := p.SetPostFilesPermissions(postPermissionsBody, post.Id)

	if setPermissionsErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Permissions update error")
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(200)
}

func (p *Plugin) getFilePermissions(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	fileId := query.Get("fileId")

	fileInfo, fileInfoErr := p.API.GetFileInfo(fileId)

	if fileInfoErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Invalid file id")
		writer.WriteHeader(400)
		return
	}

	var userId string = request.Header.Get(ONLYOFFICE_AUTHORIZATION_USERID_HEADER)
	post, postErr := p.API.GetPost(fileInfo.PostId)

	if post.UserId != userId {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Unauthorized request")
		writer.WriteHeader(403)
		return
	}

	if postErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Invalid post id")
		writer.WriteHeader(400)
		return
	}

	filePermissions := GetFilePermissionsByFileId(fileId, *post)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(filePermissions)
}

func (p *Plugin) channelUsers(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	page, parsePageErr := strconv.Atoi(query.Get("page"))
	limit, parseLimitErr := strconv.Atoi(query.Get("limit"))

	if parsePageErr != nil || parseLimitErr != nil {
		writer.WriteHeader(400)
		return
	}

	authorName := request.Header.Get(ONLYOFFICE_AUTHORIZATION_USERNAME_HEADER)

	var channelId string = request.Header.Get(ONYLOFFICE_CHANNELVALIDATION_CHANNELID_HEADER)

	usersInChannel, usersErr := p.API.GetUsersInChannel(channelId, "username", page, limit)

	if usersErr != nil {
		writer.WriteHeader(400)
		return
	}

	fileId := request.Header.Get(ONLYOFFICE_FILEVALIDATION_FILEID_HEADER)
	postId := request.Header.Get(ONLYOFFICE_FILEVALIDATION_POSTID_HEADER)
	post, _ := p.API.GetPost(postId)

	var response []models.UserInfoResponse

	for _, user := range usersInChannel {
		if user.Username == authorName {
			continue
		}
		userPermissions, _ := GetFilePermissionsByUser(user.Id, user.Username, fileId, *post)
		response = append(response, models.UserInfoResponse{
			Id:          user.Id,
			Username:    user.Username,
			Permissions: userPermissions,
		})
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(response)
}

func (p *Plugin) channelUser(writer http.ResponseWriter, request *http.Request) {
	channelId := request.Header.Get(ONYLOFFICE_CHANNELVALIDATION_CHANNELID_HEADER)

	query := request.URL.Query()
	username := query.Get("username")

	response := models.UserInfoResponse{}

	users, usersErr := p.API.GetUsersByUsernames([]string{username})

	if usersErr != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)
		json.NewEncoder(writer).Encode(response)
		return
	}

	_, membershipErr := p.API.GetChannelMember(channelId, users[0].Id)

	if membershipErr != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)
		json.NewEncoder(writer).Encode(response)
		return
	}

	fileId := request.Header.Get(ONLYOFFICE_FILEVALIDATION_FILEID_HEADER)
	postId := request.Header.Get(ONLYOFFICE_FILEVALIDATION_POSTID_HEADER)
	post, _ := p.API.GetPost(postId)

	if users[0].Id == post.UserId {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)
		json.NewEncoder(writer).Encode(response)
		return
	}

	userPermissions, _ := GetFilePermissionsByUser(users[0].Id, users[0].Username, fileId, *post)
	response.Id = users[0].Id
	response.Username = users[0].Username
	response.Permissions = userPermissions

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(response)
}

func GenerateDocKey(fileInfo model.FileInfo, post model.Post) string {
	var postUpdatedAt string = strconv.FormatInt(post.UpdateAt, 10)

	docKey, _ := security.EncryptorMD5{}.Encrypt(fileInfo.Id+postUpdatedAt, nil)

	return docKey
}

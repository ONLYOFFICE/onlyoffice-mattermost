/**
 *
 * (c) Copyright Ascensio System SIA 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/utils"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/security"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/models"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (p *Plugin) editor(writer http.ResponseWriter, request *http.Request) {
	var serverURL string = *p.API.GetConfig().ServiceSettings.SiteURL + "/" + ONLYOFFICE_API_PATH
	var fileId string = request.PostForm.Get("fileid")
	var lang string = request.PostForm.Get("lang")
	p.API.LogDebug(ONLYOFFICE_LOGGER_PREFIX + "Got an editor request")
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

	var userPermissions models.Permissions = models.ONLYOFFICE_DEFAULT_PERMISSIONS

	if utils.IsExtensionEditable(fileInfo.Extension) {
		userPermissions, _ = GetFilePermissionsByUser(userId, fileInfo.Id, *post)
	}

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
			Lang: lang,
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

func sendDocumentServerResponse(writer http.ResponseWriter, isError bool) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	if isError {
		writer.Write([]byte("{\"error\": 1}"))
	} else {
		writer.Write([]byte("{\"error\": 0}"))
	}
}

func (p *Plugin) callback(writer http.ResponseWriter, request *http.Request) {
	body := models.CallbackBody{}
	decodingErr := json.NewDecoder(request.Body).Decode(&body)

	if decodingErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Callback body decoding error")
		sendDocumentServerResponse(writer, true)
		return
	}

	p.API.LogDebug(ONLYOFFICE_LOGGER_PREFIX + "Got a valid callback payload")

	if p.configuration.DESJwt != "" {
		jwtBodyProcessingErr := ConvertJwtToBody(&body, []byte(p.configuration.DESJwt), request.Header.Get(p.configuration.DESJwtHeader))

		if jwtBodyProcessingErr != nil {
			p.API.LogError(jwtBodyProcessingErr.Error())
			sendDocumentServerResponse(writer, true)
			return
		}
	}

	handler, exists := p.getCallbackHandler(&body)

	if !exists {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Could not find a proper callback handler")
		sendDocumentServerResponse(writer, true)
		return
	}

	fileId, _ := security.EncryptorAES{}.Decrypt(request.URL.Query().Get("fileId"), p.internalKey)
	body.FileId = fileId

	handlingErr := handler(&body, p)

	if handlingErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX+"A callback handling error has occured: ", handlingErr.Error())
		sendDocumentServerResponse(writer, true)
		return
	}

	p.API.LogDebug(ONLYOFFICE_LOGGER_PREFIX + "The callback request had no errors")
	sendDocumentServerResponse(writer, false)
}

func (p *Plugin) download(writer http.ResponseWriter, request *http.Request) {
	fileId, _ := security.EncryptorAES{}.Decrypt(request.URL.Query().Get("fileId"), p.internalKey)
	fileContent, fileErr := p.API.GetFile(fileId)

	if fileErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Invalid file id when trying to download")
		return
	}

	p.API.LogDebug(ONLYOFFICE_LOGGER_PREFIX + "Downloading a file")

	writer.Write(fileContent)
}

//TODO: Refactoring
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

	prevPermissions := GetPostPermissionsByFileId(fileInfo.Id, *post, p)

	channel, _ := p.API.GetChannel(post.ChannelId)
	team, _ := p.API.GetTeam(channel.TeamId)

	//TODO: Wrong order
	for _, permissionBody := range postPermissionsBody {
		if !UserHasFilePermissions(permissionBody.Id, fileInfo.Id, post) && permissionBody.Id != utils.ONLYOFFICE_PERMISSIONS_WILDCARD_KEY {
			permissionsName := utils.GetPermissionsName(permissionBody.Permissions)
			p.onlyoffice_bot.BOT_CREATE_DM("Your "+fileInfo.Name+" file permissions have been changed to "+permissionsName+": "+*p.API.GetConfig().ServiceSettings.SiteURL+"/"+team.Name+MATTERMOST_COPY_POST_LINK_SEPARATOR+post.Id, permissionBody.Id)
		}
	}

	setPermissionsErr := p.SetPostFilesPermissions(postPermissionsBody, post.Id)

	if setPermissionsErr != nil {
		p.API.LogError(ONLYOFFICE_LOGGER_PREFIX + "Permissions update error")
		writer.WriteHeader(500)
		return
	}

	//TODO: Replace with maps and write proper util functions (it is just a workaround)
	if len(prevPermissions) > 0 {
		for _, prevPermission := range prevPermissions {
			for _, permissionBody := range postPermissionsBody {
				if prevPermission.Id == permissionBody.Id &&
					prevPermission.Permissions.Edit != permissionBody.Permissions.Edit {
					permissionsName := utils.GetPermissionsName(permissionBody.Permissions)
					if prevPermission.Username == utils.ONLYOFFICE_PERMISSIONS_WILDCARD_KEY {
						p.onlyoffice_bot.BOT_CREATE_REPLY(fileInfo.Name+" permissions have been changed to "+permissionsName, post.ChannelId, post.Id)
					} else {
						p.onlyoffice_bot.BOT_CREATE_DM("Your "+fileInfo.Name+" file permissions have been changed to "+permissionsName+": "+*p.API.GetConfig().ServiceSettings.SiteURL+"/"+team.Name+MATTERMOST_COPY_POST_LINK_SEPARATOR+post.Id, prevPermission.Id)
					}
				}
			}
		}
	} else {
		if postPermissionsBody[0].Permissions.Edit {
			permissionsName := utils.GetPermissionsName(postPermissionsBody[0].Permissions)
			p.onlyoffice_bot.BOT_CREATE_REPLY(fileInfo.Name+" permissions have been changed to "+permissionsName, post.ChannelId, post.Id)
		}
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

	channel, _ := p.API.GetChannel(post.ChannelId)

	filePermissions := GetPostPermissionsByFileId(fileId, *post, p)

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Add("Channel-Type", channel.Type)
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
		userPermissions, _ := GetFilePermissionsByUser(user.Id, fileId, *post)
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

func GenerateDocKey(fileInfo model.FileInfo, post model.Post) string {
	var postUpdatedAt string = strconv.FormatInt(post.UpdateAt, 10)

	docKey, _ := security.EncryptorMD5{}.Encrypt(fileInfo.Id+postUpdatedAt, nil)

	return docKey
}

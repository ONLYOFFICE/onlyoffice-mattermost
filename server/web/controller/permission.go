/**
 *
 * (c) Copyright Ascensio System SIA 2025
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
package controller

import (
	"encoding/json"
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/common"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type PermissionsHandler struct {
	api           plugin.API
	configuration *configuration.Configuration
	fileHelper    file.FileHelper
	bot           bot.Bot
}

func NewPermissionsHandler(
	api plugin.API,
	configuration *configuration.Configuration,
	fileHelper file.FileHelper,
	bot bot.Bot,
) PermissionsHandler {
	return PermissionsHandler{
		api:           api,
		configuration: configuration,
		fileHelper:    fileHelper,
		bot:           bot,
	}
}

func (h *PermissionsHandler) SetPermissions(rw http.ResponseWriter, r *http.Request) {
	h.api.LogDebug(onlyofficeLoggerPrefix + "got a new set file permissions request")

	var postPermissions []model.PostPermission
	if json.NewDecoder(r.Body).Decode(&postPermissions) != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not decode permissions body")
		common.WriteJSON(rw, callbackErr, http.StatusBadRequest)
		return
	}

	if len(postPermissions) < 1 {
		h.api.LogError(onlyofficeLoggerPrefix + "invalid permissions body length")
		common.WriteJSON(rw, callbackErr, http.StatusBadRequest)
		return
	}

	for _, permission := range postPermissions {
		if permission.FileID != postPermissions[0].FileID {
			h.api.LogWarn(onlyofficeLoggerPrefix + "an unauthorized attempt to change file permissions: " + permission.FileID)
			common.WriteJSON(rw, callbackErr, http.StatusForbidden)
			return
		}
	}

	fileInfo, fileInfoErr := h.api.GetFileInfo(postPermissions[0].FileID)
	if fileInfoErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not access file info " + postPermissions[0].FileID + " Reason: " + fileInfoErr.Message)
		common.WriteJSON(rw, callbackErr, http.StatusInternalServerError)
		return
	}

	post, postErr := h.api.GetPost(fileInfo.PostId)
	if postErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not access post " + fileInfo.PostId + "Reason: " + postErr.Message)
		common.WriteJSON(rw, callbackErr, http.StatusInternalServerError)
		return
	}

	if post.UserId != r.Header.Get(tools.MMAuthHeader) {
		h.api.LogWarn(onlyofficeLoggerPrefix + "only author can set file permissions")
		common.WriteJSON(rw, callbackErr, http.StatusForbidden)
		return
	}

	channel, channelErr := h.api.GetChannel(post.ChannelId)
	if channelErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get channel with id " + post.ChannelId + " Reason: " + channelErr.Message)
		common.WriteJSON(rw, callbackErr, http.StatusBadRequest)
		return
	}

	team, teamErr := h.api.GetTeam(channel.TeamId)
	if teamErr != nil && len(postPermissions) > 1 && postPermissions[0].UserID != h.fileHelper.GetWildcardUser() {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get team with id " + channel.TeamId + " Reason: " + teamErr.Message)
		common.WriteJSON(rw, callbackErr, http.StatusBadRequest)
		return
	}

	newPermissions := h.fileHelper.SetPostFilePermissions(post, postPermissions)
	_, err := h.api.UpdatePost(post)

	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not update post. Reason: " + err.Message)
		common.WriteJSON(rw, callbackErr, http.StatusInternalServerError)
		return
	}

	postID := post.Id
	if post.RootId != "" {
		postID = post.RootId
	}

	for _, permission := range newPermissions {
		if permission.UserID == h.fileHelper.GetWildcardUser() {
			h.bot.BotCreateReply(fileInfo.Name+" permissions have been changed to "+common.GetPermissionsName(permission.Permissions), post.ChannelId, postID)
		} else if team != nil {
			h.bot.BotCreateDM("Your "+fileInfo.Name+" file permissions have been changed to "+common.GetPermissionsName(permission.Permissions)+": "+*h.api.GetConfig().ServiceSettings.SiteURL+"/"+team.Name+"/pl/"+postID, permission.UserID)
		}
	}

	common.WriteJSON(rw, callbackErr)
}

func (h *PermissionsHandler) GetPermissions(rw http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("file")

	fileInfo, fileInfoErr := h.api.GetFileInfo(fileID)
	if fileInfoErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not access file info " + fileID + " Reason: " + fileInfoErr.Message)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	post, postErr := h.api.GetPost(fileInfo.PostId)
	if postErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not access post " + fileInfo.PostId + "Reason: " + postErr.Message)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if post.UserId != r.Header.Get(tools.MMAuthHeader) {
		h.api.LogWarn(onlyofficeLoggerPrefix + "only author can get file permissions")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	permissions := h.fileHelper.GetPostPermissionsByFileID(fileID, post, h.api.GetUser)
	common.WriteJSON(rw, permissions)
}

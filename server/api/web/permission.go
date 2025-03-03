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
package web

import (
	"encoding/json"
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
)

func BuildSetFilePermissionsHandler(plugin api.PluginAPI) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		plugin.API.LogDebug(_OnlyofficeLoggerPrefix + "got a new set file permissions request")

		var postPermissions []model.PostPermission
		if json.NewDecoder(r.Body).Decode(&postPermissions) != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not decode permissions body")
			api.WriteJSON(rw, _CallbackErr, http.StatusBadRequest)
			return
		}

		if len(postPermissions) < 1 {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "invalid permissions body length")
			api.WriteJSON(rw, _CallbackErr, http.StatusBadRequest)
			return
		}

		for _, permission := range postPermissions {
			if permission.FileID != postPermissions[0].FileID {
				plugin.API.LogWarn(_OnlyofficeLoggerPrefix + "an unauthorized attempt to change file permissions: " + permission.FileID)
				api.WriteJSON(rw, _CallbackErr, http.StatusForbidden)
				return
			}
		}

		fileInfo, fileInfoErr := plugin.API.GetFileInfo(postPermissions[0].FileID)
		if fileInfoErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not access file info " + postPermissions[0].FileID + " Reason: " + fileInfoErr.Message)
			api.WriteJSON(rw, _CallbackErr, http.StatusInternalServerError)
			return
		}

		post, postErr := plugin.API.GetPost(fileInfo.PostId)
		if postErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not access post " + fileInfo.PostId + "Reason: " + postErr.Message)
			api.WriteJSON(rw, _CallbackErr, http.StatusInternalServerError)
			return
		}

		if post.UserId != r.Header.Get(plugin.Configuration.MMAuthHeader) {
			plugin.API.LogWarn(_OnlyofficeLoggerPrefix + "only author can set file permissions")
			api.WriteJSON(rw, _CallbackErr, http.StatusForbidden)
			return
		}

		channel, channelErr := plugin.API.GetChannel(post.ChannelId)
		if channelErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not get channel with id " + post.ChannelId + " Reason: " + channelErr.Message)
			api.WriteJSON(rw, _CallbackErr, http.StatusBadRequest)
			return
		}

		team, teamErr := plugin.API.GetTeam(channel.TeamId)
		if teamErr != nil && len(postPermissions) > 1 && postPermissions[0].UserID != plugin.OnlyofficeHelper.GetWildcardUser() {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not get team with id " + channel.TeamId + " Reason: " + teamErr.Message)
			api.WriteJSON(rw, _CallbackErr, http.StatusBadRequest)
			return
		}

		newPermissions := plugin.OnlyofficeHelper.SetPostFilePermissions(post, postPermissions)
		_, err := plugin.API.UpdatePost(post)

		if err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not update post. Reason: " + err.Message)
			api.WriteJSON(rw, _CallbackErr, http.StatusInternalServerError)
			return
		}

		for _, permission := range newPermissions {
			if permission.UserID == plugin.OnlyofficeHelper.GetWildcardUser() {
				plugin.Bot.BotCreateReply(fileInfo.Name+" permissions have been changed to "+api.GetPermissionsName(permission.Permissions), post.ChannelId, post.Id)
			} else if team != nil {
				plugin.Bot.BotCreateDM("Your "+fileInfo.Name+" file permissions have been changed to "+api.GetPermissionsName(permission.Permissions)+": "+*plugin.API.GetConfig().ServiceSettings.SiteURL+"/"+team.Name+"/pl/"+post.Id, permission.UserID)
			}
		}

		api.WriteJSON(rw, _CallbackErr)
	}
}

func BuildGetFilePermissionsHandler(plugin api.PluginAPI) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		fileID := r.URL.Query().Get("file")

		fileInfo, fileInfoErr := plugin.API.GetFileInfo(fileID)
		if fileInfoErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not access file info " + fileID + " Reason: " + fileInfoErr.Message)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		post, postErr := plugin.API.GetPost(fileInfo.PostId)
		if postErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not access post " + fileInfo.PostId + "Reason: " + postErr.Message)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if post.UserId != r.Header.Get(plugin.Configuration.MMAuthHeader) {
			plugin.API.LogWarn(_OnlyofficeLoggerPrefix + "only author can get file permissions")
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		permissions := plugin.OnlyofficeHelper.GetPostPermissionsByFileID(fileID, post, plugin.API.GetUser)

		api.WriteJSON(rw, permissions)
	}
}

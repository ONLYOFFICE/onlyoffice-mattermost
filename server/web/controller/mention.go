/**
 *
 * (c) Copyright Ascensio System SIA 2026
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
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type MentionsHandler struct {
	api plugin.API
	bot bot.Bot
}

func NewMentionsHandler(
	api plugin.API,
	bot bot.Bot,
) MentionsHandler {
	return MentionsHandler{
		api: api,
		bot: bot,
	}
}

func (h *MentionsHandler) writeErrorResponse(rw http.ResponseWriter, errorMsg string, statusCode int) {
	response := &model.MentionNotifyResponse{
		Status: "error",
		Error:  errorMsg,
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(response.ToJSON())
}

func (h *MentionsHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write(model.MentionUsersResponse{}.ToJSON())
		}
	}()

	fileID := r.URL.Query().Get("file")
	if fileID == "" {
		h.api.LogError(onlyofficeLoggerPrefix + "missing file parameter")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	fileInfo, fileInfoErr := h.api.GetFileInfo(fileID)
	if fileInfoErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not access file info " + fileID + " Reason: " + fileInfoErr.Message)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	post, postErr := h.api.GetPost(fileInfo.PostId)
	if postErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not access post " + fileInfo.PostId + " Reason: " + postErr.Message)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	channel, channelErr := h.api.GetChannel(post.ChannelId)
	if channelErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get channel with id " + post.ChannelId + " Reason: " + channelErr.Message)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	channelMembers, membersErr := h.api.GetUsersInChannel(channel.Id, "username", 0, 200)
	if membersErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get channel members: " + membersErr.Message)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(model.MentionUsersResponse{}.ToJSON())
		return
	}

	if channelMembers == nil {
		channelMembers = []*mmModel.User{}
	}

	mentionUsers := make(model.MentionUsersResponse, 0, len(channelMembers))
	currentUserID := r.Header.Get(tools.MMAuthHeader)
	for _, user := range channelMembers {
		if user == nil {
			continue
		}

		if user.IsBot || user.Id == currentUserID {
			continue
		}

		mentionUsers = append(mentionUsers, model.MentionUser{
			Email: user.Email,
			Name:  user.Username,
			ID:    user.Id,
		})
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(mentionUsers.ToJSON()); err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not write response: " + err.Error())
	}
}

func (h *MentionsHandler) SendNotifications(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			h.writeErrorResponse(rw, "internal server error", http.StatusInternalServerError)
		}
	}()

	var payload model.MentionNotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not decode mention notification body: " + err.Error())
		h.writeErrorResponse(rw, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if payload.FileID == "" {
		h.api.LogError(onlyofficeLoggerPrefix + "missing file ID in mention notification")
		h.writeErrorResponse(rw, "missing file ID", http.StatusBadRequest)
		return
	}

	fileInfo, fileInfoErr := h.api.GetFileInfo(payload.FileID)
	if fileInfoErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not access file info " + payload.FileID + " Reason: " + fileInfoErr.Message)
		h.writeErrorResponse(rw, "file not found", http.StatusBadRequest)
		return
	}

	post, postErr := h.api.GetPost(fileInfo.PostId)
	if postErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not access post " + fileInfo.PostId + " Reason: " + postErr.Message)
		h.writeErrorResponse(rw, "post not found", http.StatusBadRequest)
		return
	}

	channel, channelErr := h.api.GetChannel(post.ChannelId)
	if channelErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get channel with id " + post.ChannelId + " Reason: " + channelErr.Message)
		h.writeErrorResponse(rw, "channel not found", http.StatusBadRequest)
		return
	}

	team, teamErr := h.api.GetTeam(channel.TeamId)
	if teamErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get team with id " + channel.TeamId + " Reason: " + teamErr.Message)
		h.writeErrorResponse(rw, "team not found", http.StatusBadRequest)
		return
	}

	postID := post.Id
	if post.RootId != "" {
		postID = post.RootId
	}

	documentLink := *h.api.GetConfig().ServiceSettings.SiteURL + "/" + team.Name + "/pl/" + postID
	currentUser, _ := h.api.GetUser(r.Header.Get(tools.MMAuthHeader))
	mentionedByMention := "@Someone"
	if currentUser != nil {
		mentionedByMention = "@" + currentUser.Username
	}

	notificationsSent := 0
	for _, email := range payload.Emails {
		user, userErr := h.api.GetUserByEmail(email)
		if userErr != nil {
			h.api.LogDebug(onlyofficeLoggerPrefix + "could not find user with email " + email + ": " + userErr.Message)
			continue
		}

		message := mentionedByMention + " mentioned you in the document **" + fileInfo.Name + "**: " + documentLink
		h.bot.BotCreateDM(message, user.Id)
		notificationsSent++
	}

	response := &model.MentionNotifyResponse{
		Status: "ok",
	}

	if notificationsSent == 0 && len(payload.Emails) > 0 {
		response.Error = "no valid users found"
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response.ToJSON())
}

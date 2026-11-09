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
package common

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

func CanReadFile(api plugin.API, userID string, fileInfo *model.FileInfo, post *model.Post) bool {
	if userID == "" || fileInfo == nil {
		return false
	}

	if fileInfo.CreatorId == userID {
		return true
	}

	channelID := fileInfo.ChannelId
	if channelID == "" && post != nil {
		channelID = post.ChannelId
	}

	return CanReadChannel(api, userID, channelID)
}

func CanReadChannel(api plugin.API, userID, channelID string) bool {
	if userID == "" {
		return false
	}

	if api.HasPermissionToChannel(userID, channelID, model.PermissionReadChannelContent) {
		return true
	}

	channel, appErr := api.GetChannel(channelID)
	if appErr != nil || channel == nil || channel.Type != model.ChannelTypeOpen {
		return false
	}

	if config := api.GetConfig(); config != nil && model.SafeDereference(config.ComplianceSettings.Enable) {
		return false
	}

	return api.HasPermissionToTeam(userID, channel.TeamId, model.PermissionReadPublicChannel)
}

func PostLink(api plugin.API, channel *model.Channel, postID string) string {
	siteURL := ""
	if config := api.GetConfig(); config != nil {
		siteURL = model.SafeDereference(config.ServiceSettings.SiteURL)
	}

	if channel != nil && channel.TeamId != "" {
		if team, appErr := api.GetTeam(channel.TeamId); appErr == nil && team != nil {
			return siteURL + "/" + team.Name + "/pl/" + postID
		}
	}

	return siteURL + "/_redirect/pl/" + postID
}

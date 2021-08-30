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
	"github.com/mattermost/mattermost-server/v5/model"
)

type ONLYOFFICE_BOT struct {
	Id           string
	LoggerPrefix string
	P            *Plugin
}

func (bot *ONLYOFFICE_BOT) BOT_CREATE_POST(message string, channelId string) {
	ONLYOFFICE_BOT_POST := model.Post{
		Message:   message,
		ChannelId: channelId,
		UserId:    bot.Id,
	}

	_, creationErr := bot.P.API.CreatePost(&ONLYOFFICE_BOT_POST)
	if creationErr != nil {
		bot.P.API.LogError(ONLYOFFICE_BOT_LOGGER_PREFIX + "Post creation error")
		return
	}

	bot.P.API.LogInfo(bot.LoggerPrefix + "Created a new post")
}

func (bot *ONLYOFFICE_BOT) BOT_CREATE_REPLY(message string, channelId string, parentId string) {
	ONLYOFFICE_BOT_POST := model.Post{
		Message:   message,
		ParentId:  parentId,
		RootId:    parentId,
		ChannelId: channelId,
		UserId:    bot.Id,
	}

	_, creationErr := bot.P.API.CreatePost(&ONLYOFFICE_BOT_POST)
	if creationErr != nil {
		bot.P.API.LogError(ONLYOFFICE_BOT_LOGGER_PREFIX + "Post creation error")
		return
	}
}

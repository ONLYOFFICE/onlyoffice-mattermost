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
package bot

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type Bot interface {
	BotCreateDM(message string, userID string)
	BotCreatePost(message string, channelID string)
	BotCreateReply(message string, channelID string, parentID string)
}

type onlyofficeBot struct {
	ID  string
	API plugin.API
}

func New(botID string, api plugin.API) Bot {
	return &onlyofficeBot{
		ID:  botID,
		API: api,
	}
}

func (b *onlyofficeBot) BotCreateDM(message string, userID string) {
	direct, err := b.API.GetDirectChannel(b.ID, userID)

	if err != nil {
		b.API.LogError(onlyofficeLoggerBotPrefix + "could not get direct channel")
		return
	}

	post := model.Post{
		Message:   message,
		ChannelId: direct.Id,
		UserId:    b.ID,
	}

	_, err = b.API.CreatePost(&post)
	if err != nil {
		b.API.LogError(onlyofficeLoggerBotPrefix + "post creation error")
		return
	}

	b.API.LogDebug(onlyofficeLoggerBotPrefix + "created a new DM post")
}

func (b *onlyofficeBot) BotCreatePost(message string, channelID string) {
	post := model.Post{
		Message:   message,
		ChannelId: channelID,
		UserId:    b.ID,
	}

	_, err := b.API.CreatePost(&post)
	if err != nil {
		b.API.LogError(onlyofficeLoggerBotPrefix + "bot post creation error")
		return
	}

	b.API.LogDebug(onlyofficeLoggerBotPrefix + "created a new post")
}

func (b *onlyofficeBot) BotCreateReply(message string, channelID string, parentID string) {
	post := model.Post{
		Message:   message,
		RootId:    parentID,
		ChannelId: channelID,
		UserId:    b.ID,
	}

	_, err := b.API.CreatePost(&post)
	if err != nil {
		b.API.LogError(onlyofficeLoggerBotPrefix + "reply creation error")
		return
	}

	b.API.LogDebug(onlyofficeLoggerBotPrefix + "created a new reply")
}

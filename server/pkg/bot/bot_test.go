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
package bot

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/mock"
)

func TestBotCreateDM(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetDirectChannel", "bot-1", "user-1").Return(&model.Channel{Id: "dm-1"}, nil)
	api.On("CreatePost", mock.MatchedBy(func(p *model.Post) bool {
		return p.ChannelId == "dm-1" && p.UserId == "bot-1" && p.Message == "hello"
	})).Return(&model.Post{Id: "post-1"}, nil)
	api.On("LogDebug", mock.Anything).Return().Maybe()

	bot := New("bot-1", api)
	bot.BotCreateDM("hello", "user-1")
	api.AssertExpectations(t)
}

func TestBotCreatePost(t *testing.T) {
	api := &plugintest.API{}
	api.On("CreatePost", mock.MatchedBy(func(p *model.Post) bool {
		return p.ChannelId == "channel-1" && p.Message == "post"
	})).Return(&model.Post{Id: "post-1"}, nil)
	api.On("LogDebug", mock.Anything).Return().Maybe()

	bot := New("bot-1", api)
	bot.BotCreatePost("post", "channel-1")
	api.AssertExpectations(t)
}

func TestBotCreateReply(t *testing.T) {
	api := &plugintest.API{}
	api.On("CreatePost", mock.MatchedBy(func(p *model.Post) bool {
		return p.RootId == "parent-1" && p.ChannelId == "channel-1"
	})).Return(&model.Post{Id: "reply-1"}, nil)
	api.On("LogDebug", mock.Anything).Return().Maybe()

	bot := New("bot-1", api)
	bot.BotCreateReply("reply", "channel-1", "parent-1")
	api.AssertExpectations(t)
}

func TestBotCreateDMChannelError(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetDirectChannel", "bot-1", "user-1").Return((*model.Channel)(nil), model.NewAppError("GetDirectChannel", "err", nil, "fail", 500))
	api.On("LogError", mock.Anything).Return().Maybe()

	bot := New("bot-1", api)
	bot.BotCreateDM("hello", "user-1")
	api.AssertExpectations(t)
}

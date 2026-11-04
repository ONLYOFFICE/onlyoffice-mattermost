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
	"net/http"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func accessAPI() *plugintest.API {
	api := &plugintest.API{}
	member := func(user, channel string) bool {
		return user == "alice" && (channel == "private" || channel == "public")
	}
	api.On("HasPermissionToChannel", mock.Anything, mock.Anything, mock.Anything).Return(func(user, channel string, p *model.Permission) bool {
		switch p.Id {
		case model.PermissionReadChannelContent.Id, model.PermissionReadChannel.Id, model.PermissionCreatePost.Id, model.PermissionUploadFile.Id:
			return member(user, channel)
		}
		return false
	})
	api.On("HasPermissionToTeam", mock.Anything, "team1", model.PermissionReadPublicChannel).Return(func(user, _ string, _ *model.Permission) bool {
		return user == "alice" || user == "eve"
	})
	api.On("HasPermissionToTeam", mock.Anything, mock.Anything, mock.Anything).Return(false)
	return api
}

func TestCanReadChannel(t *testing.T) {
	api := accessAPI()
	private := &model.Channel{Id: "private", TeamId: "team1", Type: model.ChannelTypePrivate}
	public := &model.Channel{Id: "public", TeamId: "team1", Type: model.ChannelTypeOpen}

	assert.True(t, CanReadChannel(api, "alice", private), "member reads private channel")
	assert.False(t, CanReadChannel(api, "eve", private), "team member outside a private channel")
	assert.True(t, CanReadChannel(api, "eve", public), "team member reads public channel")
	assert.False(t, CanReadChannel(api, "mallory", public), "user outside the team")
	assert.False(t, CanReadChannel(api, "", private))
	assert.False(t, CanReadChannel(api, "alice", nil))
}

func TestCanPostFileToChannel(t *testing.T) {
	api := accessAPI()
	assert.True(t, CanPostFileToChannel(api, "alice", &model.Channel{Id: "private"}))
	assert.False(t, CanPostFileToChannel(api, "eve", &model.Channel{Id: "private"}))
	assert.False(t, CanPostFileToChannel(api, "alice", &model.Channel{Id: "private", DeleteAt: 1}), "archived channel")
}

func TestResolveFileForUser(t *testing.T) {
	api := accessAPI()
	api.On("GetFileInfo", "f1").Return(&model.FileInfo{Id: "f1", PostId: "p1"}, nil)
	api.On("GetFileInfo", "deleted").Return(&model.FileInfo{Id: "deleted", PostId: "p1", DeleteAt: 1}, nil)
	api.On("GetFileInfo", "nope").Return(nil, model.NewAppError("t", "id", nil, "", http.StatusNotFound))
	api.On("GetPost", "p1").Return(&model.Post{Id: "p1", ChannelId: "private"}, nil)
	api.On("GetChannel", "private").Return(&model.Channel{Id: "private", TeamId: "team1", Type: model.ChannelTypePrivate}, nil)

	ctx, status, err := ResolveFileForUser(api, "alice", "f1")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "private", ctx.Channel.Id)

	_, status, err = ResolveFileForUser(api, "eve", "f1")
	assert.Error(t, err)
	assert.Equal(t, http.StatusForbidden, status)

	_, status, _ = ResolveFileForUser(api, "alice", "deleted")
	assert.Equal(t, http.StatusNotFound, status)
	_, status, _ = ResolveFileForUser(api, "alice", "nope")
	assert.Equal(t, http.StatusNotFound, status)
	_, status, _ = ResolveFileForUser(api, "alice", "")
	assert.Equal(t, http.StatusBadRequest, status)
}

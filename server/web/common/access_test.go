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
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
)

func complianceConfig(enabled bool) *model.Config {
	config := &model.Config{}
	config.SetDefaults()
	config.ComplianceSettings.Enable = model.NewPointer(enabled)
	return config
}

func TestCanReadFile(t *testing.T) {
	file := &model.FileInfo{Id: "file-1", CreatorId: "owner", ChannelId: "ch-1"}

	t.Run("no user", func(t *testing.T) {
		assert.False(t, CanReadFile(&plugintest.API{}, "", file, nil))
	})

	t.Run("uploader without channel access", func(t *testing.T) {
		assert.True(t, CanReadFile(&plugintest.API{}, "owner", file, nil))
	})

	t.Run("channel member", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("HasPermissionToChannel", "member", "ch-1", model.PermissionReadChannelContent).Return(true)
		assert.True(t, CanReadFile(api, "member", file, nil))
	})

	t.Run("outsider of a private channel", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("HasPermissionToChannel", "outsider", "ch-1", model.PermissionReadChannelContent).Return(false)
		api.On("GetChannel", "ch-1").Return(&model.Channel{Id: "ch-1", Type: model.ChannelTypePrivate, TeamId: "team-1"}, nil)
		assert.False(t, CanReadFile(api, "outsider", file, nil))
		api.AssertNotCalled(t, "HasPermissionToTeam", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("outsider of a direct message", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("HasPermissionToChannel", "outsider", "ch-1", model.PermissionReadChannelContent).Return(false)
		api.On("GetChannel", "ch-1").Return(&model.Channel{Id: "ch-1", Type: model.ChannelTypeDirect}, nil)
		assert.False(t, CanReadFile(api, "outsider", file, nil))
	})

	t.Run("team member reading a public channel", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("HasPermissionToChannel", "reader", "ch-1", model.PermissionReadChannelContent).Return(false)
		api.On("GetChannel", "ch-1").Return(&model.Channel{Id: "ch-1", Type: model.ChannelTypeOpen, TeamId: "team-1"}, nil)
		api.On("GetConfig").Return(complianceConfig(false))
		api.On("HasPermissionToTeam", "reader", "team-1", model.PermissionReadPublicChannel).Return(true)
		assert.True(t, CanReadFile(api, "reader", file, nil))
	})

	t.Run("public channel under compliance mode", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("HasPermissionToChannel", "reader", "ch-1", model.PermissionReadChannelContent).Return(false)
		api.On("GetChannel", "ch-1").Return(&model.Channel{Id: "ch-1", Type: model.ChannelTypeOpen, TeamId: "team-1"}, nil)
		api.On("GetConfig").Return(complianceConfig(true))
		assert.False(t, CanReadFile(api, "reader", file, nil))
	})

	t.Run("public channel outside of the user's teams", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("HasPermissionToChannel", "stranger", "ch-1", model.PermissionReadChannelContent).Return(false)
		api.On("GetChannel", "ch-1").Return(&model.Channel{Id: "ch-1", Type: model.ChannelTypeOpen, TeamId: "team-1"}, nil)
		api.On("GetConfig").Return(complianceConfig(false))
		api.On("HasPermissionToTeam", "stranger", "team-1", model.PermissionReadPublicChannel).Return(false)
		assert.False(t, CanReadFile(api, "stranger", file, nil))
	})

	t.Run("falls back to the post channel", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("HasPermissionToChannel", "member", "post-ch", model.PermissionReadChannelContent).Return(true)
		legacy := &model.FileInfo{Id: "file-1", CreatorId: "owner"}
		assert.True(t, CanReadFile(api, "member", legacy, &model.Post{ChannelId: "post-ch"}))
	})
}

func TestBuildEditorConfigDeniesUserWithoutAccess(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "outsider").Return(&model.User{Id: "outsider", Username: "mallory"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Extension: "docx", PostId: "post-1", ChannelId: "private-ch", CreatorId: "owner",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", ChannelId: "private-ch", UserId: "owner"}, nil)
	api.On("HasPermissionToChannel", "outsider", "private-ch", model.PermissionReadChannelContent).Return(false)
	api.On("GetChannel", "private-ch").Return(&model.Channel{Id: "private-ch", Type: model.ChannelTypePrivate}, nil)

	req := httptest.NewRequest(http.MethodGet, "/editor?file=file-1", nil)
	req.Header.Set(tools.MMAuthHeader, "outsider")

	_, _, status, err := BuildEditorConfig(req, validEditorConfiguration(), api, crypto.NewMD5Encoder(), crypto.NewJwtManager(), newFileHelper(t))
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, status)
	api.AssertNotCalled(t, "KVSetWithExpiry", mock.Anything, mock.Anything, mock.Anything)
}

func TestPostLink(t *testing.T) {
	t.Run("team channel", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("GetConfig").Return(siteURLConfig())
		api.On("GetTeam", "team-1").Return(&model.Team{Id: "team-1", Name: "acme"}, nil)
		assert.Equal(t, "https://mm.example.com/acme/pl/post-1", PostLink(api, &model.Channel{TeamId: "team-1"}, "post-1"))
	})

	t.Run("direct message has no team", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("GetConfig").Return(siteURLConfig())
		assert.Equal(t, "https://mm.example.com/_redirect/pl/post-1", PostLink(api, &model.Channel{Type: model.ChannelTypeDirect}, "post-1"))
		api.AssertNotCalled(t, "GetTeam", mock.Anything)
	})
}

func TestBuildEditorConfigCarriesActionLink(t *testing.T) {
	newRequest := func(action string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/editor?file=file-1&lang=en&action="+url.QueryEscape(action), nil)
		req.Header.Set(tools.MMAuthHeader, "user-1")
		return req
	}
	newAPI := func() *plugintest.API {
		api := &plugintest.API{}
		allowAccess(api)
		api.On("GetConfig").Return(siteURLConfig())
		api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
		api.On("GetFileInfo", "file-1").Return(&model.FileInfo{Id: "file-1", Name: "notes.docx", Extension: "docx", PostId: "post-1"}, nil)
		api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", UserId: "user-1", UpdateAt: 1}, nil)
		api.On("KVSetWithExpiry", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		return api
	}

	config, _, _, err := BuildEditorConfig(newRequest(`{ "action": {"type":"comment","data":"c1"} }`), validEditorConfiguration(), newAPI(), crypto.NewMD5Encoder(), crypto.NewJwtManager(), newFileHelper(t))
	require.NoError(t, err)
	assert.JSONEq(t, `{"action":{"type":"comment","data":"c1"}}`, string(config.EditorConfig.ActionLink))

	var claims map[string]any
	require.NoError(t, crypto.NewJwtManager().Verify([]byte(validEditorConfiguration().DESJwt), config.Token, &claims))
	assert.Equal(t, map[string]any{"action": map[string]any{"type": "comment", "data": "c1"}}, claims["editorConfig"].(map[string]any)["actionLink"], "the action link must be inside the signed config")

	config, _, _, err = BuildEditorConfig(newRequest(`not json`), validEditorConfiguration(), newAPI(), crypto.NewMD5Encoder(), crypto.NewJwtManager(), newFileHelper(t))
	require.NoError(t, err)
	assert.Nil(t, config.EditorConfig.ActionLink)
}

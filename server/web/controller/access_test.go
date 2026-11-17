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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
)

func privateFileAPI() *plugintest.API {
	api := quietAPI()
	api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{
		Id: "file-1", Name: "salaries.docx", PostId: "post-1", ChannelId: "private-ch", CreatorId: "owner",
	}, nil)
	api.On("GetPost", "post-1").Return(&mmModel.Post{Id: "post-1", ChannelId: "private-ch", UserId: "owner"}, nil)
	api.On("GetChannel", "private-ch").Return(&mmModel.Channel{Id: "private-ch", Type: mmModel.ChannelTypePrivate, TeamId: "team-1"}, nil)
	api.On("GetTeam", "team-1").Return(&mmModel.Team{Id: "team-1", Name: "acme"}, nil)
	api.On("GetConfig").Return(siteURLConfig())
	api.On("HasPermissionToChannel", "member", "private-ch", mmModel.PermissionReadChannelContent).Return(true)
	api.On("HasPermissionToChannel", mock.Anything, "private-ch", mmModel.PermissionReadChannelContent).Return(false)
	return api
}

func TestMentionsGetUsersDeniesOutsider(t *testing.T) {
	api := privateFileAPI()
	handler := NewMentionsHandler(api, &stubBot{})
	req := httptest.NewRequest(http.MethodGet, "/mentions/users?file=file-1", nil)
	req.Header.Set(tools.MMAuthHeader, "outsider")
	recorder := httptest.NewRecorder()
	handler.GetUsers(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.NotContains(t, recorder.Body.String(), "@")
	api.AssertNotCalled(t, "GetUsersInChannel", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestMentionsSendNotificationsDeniesOutsider(t *testing.T) {
	api := privateFileAPI()
	bot := &stubBot{}
	handler := NewMentionsHandler(api, bot)
	body, _ := json.Marshal(model.MentionNotifyRequest{FileID: "file-1", Emails: []string{"member@ex.com"}})
	req := httptest.NewRequest(http.MethodPost, "/mentions/notify", bytes.NewReader(body))
	req.Header.Set(tools.MMAuthHeader, "outsider")
	recorder := httptest.NewRecorder()
	handler.SendNotifications(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.Empty(t, bot.dms)
}

func TestMentionsSendNotificationsSkipsRecipientsWithoutAccess(t *testing.T) {
	api := privateFileAPI()
	api.On("GetUser", "member").Return(&mmModel.User{Id: "member", Username: "alice"}, nil)
	api.On("GetUserByEmail", "member@ex.com").Return(&mmModel.User{Id: "member"}, nil)
	api.On("GetUserByEmail", "outsider@ex.com").Return(&mmModel.User{Id: "outsider"}, nil)
	bot := &stubBot{}
	handler := NewMentionsHandler(api, bot)
	body, _ := json.Marshal(model.MentionNotifyRequest{FileID: "file-1", Emails: []string{"member@ex.com", "outsider@ex.com"}})
	req := httptest.NewRequest(http.MethodPost, "/mentions/notify", bytes.NewReader(body))
	req.Header.Set(tools.MMAuthHeader, "member")
	recorder := httptest.NewRecorder()
	handler.SendNotifications(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	if assert.Len(t, bot.dms, 1) {
		assert.Contains(t, bot.dms[0], "member:")
	}
}

func directFileAPI() *plugintest.API {
	api := quietAPI()
	api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{
		Id: "file-1", Name: "direct.docx", PostId: "post-1", ChannelId: "dm-ch", CreatorId: "owner",
	}, nil)
	api.On("GetPost", "post-1").Return(&mmModel.Post{Id: "post-1", ChannelId: "dm-ch", UserId: "owner"}, nil)
	api.On("GetChannel", "dm-ch").Return(&mmModel.Channel{Id: "dm-ch", Type: mmModel.ChannelTypeDirect}, nil)
	api.On("GetConfig").Return(siteURLConfig())
	api.On("HasPermissionToChannel", "member", "dm-ch", mmModel.PermissionReadChannelContent).Return(true)
	api.On("HasPermissionToChannel", "owner", "dm-ch", mmModel.PermissionReadChannelContent).Return(true)
	api.On("HasPermissionToChannel", mock.Anything, "dm-ch", mmModel.PermissionReadChannelContent).Return(false)
	return api
}

func TestMentionsSendNotificationsInDirectMessage(t *testing.T) {
	api := directFileAPI()
	api.On("GetUser", "owner").Return(&mmModel.User{Id: "owner", Username: "alice"}, nil)
	api.On("GetUserByEmail", "member@ex.com").Return(&mmModel.User{Id: "member"}, nil)
	bot := &stubBot{}
	handler := NewMentionsHandler(api, bot)
	body, _ := json.Marshal(model.MentionNotifyRequest{FileID: "file-1", Emails: []string{"member@ex.com"}})
	req := httptest.NewRequest(http.MethodPost, "/mentions/notify", bytes.NewReader(body))
	req.Header.Set(tools.MMAuthHeader, "owner")
	recorder := httptest.NewRecorder()
	handler.SendNotifications(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	if assert.Len(t, bot.dms, 1) {
		assert.Contains(t, bot.dms[0], "https://mm.example.com/_redirect/pl/post-1")
	}

	api.AssertNotCalled(t, "GetTeam", mock.Anything)
}

func TestPermissionsSetNotifiesUserInDirectMessage(t *testing.T) {
	api := directFileAPI()
	api.On("GetPost", "post-1").Unset()
	api.On("GetPost", "post-1").Return(&mmModel.Post{Id: "post-1", ChannelId: "dm-ch", UserId: "owner", FileIds: mmModel.StringArray{"file-1"}}, nil)
	api.On("UpdatePost", mock.AnythingOfType("*model.Post")).Return(&mmModel.Post{Id: "post-1"}, nil)
	api.On("GetUser", mock.Anything).Return(&mmModel.User{Id: "member", Username: "bob"}, nil).Maybe()
	bot := &stubBot{}
	handler := NewPermissionsHandler(api, validDESConfig(), newFileHelper(t), bot)
	body, _ := json.Marshal([]model.PostPermission{
		{FileID: "file-1", UserID: "member", Username: "bob", Permissions: model.Permissions{Edit: true}},
	})
	req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewReader(body))
	req.Header.Set(tools.MMAuthHeader, "owner")
	recorder := httptest.NewRecorder()
	handler.SetPermissions(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	if assert.Len(t, bot.dms, 1) {
		assert.Contains(t, bot.dms[0], "member:")
		assert.Contains(t, bot.dms[0], "/_redirect/pl/post-1")
	}
}

func TestMentionsSendNotificationsQuotesCommentAndLinksToIt(t *testing.T) {
	api := directFileAPI()
	api.On("GetTeamsForUser", "member").Return([]*mmModel.Team{{Id: "team-1", Name: "acme"}}, nil)
	api.On("GetUser", "owner").Return(&mmModel.User{Id: "owner", Username: "alice"}, nil)
	api.On("GetUserByEmail", "member@ex.com").Return(&mmModel.User{Id: "member"}, nil)
	bot := &stubBot{}
	handler := NewMentionsHandler(api, bot)
	body, _ := json.Marshal(model.MentionNotifyRequest{
		FileID:     "file-1",
		Emails:     []string{"member@ex.com"},
		Message:    "please check\nthe totals",
		ActionLink: json.RawMessage(`{"action": {"type": "comment", "data": "c1"}}`),
	})
	req := httptest.NewRequest(http.MethodPost, "/mentions/notify", bytes.NewReader(body))
	req.Header.Set(tools.MMAuthHeader, "owner")
	recorder := httptest.NewRecorder()
	handler.SendNotifications(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	if assert.Len(t, bot.dms, 1) {
		dm := bot.dms[0]
		assert.Contains(t, dm, "> please check\n> the totals")
		assert.Contains(t, dm, "[Open the comment](https://mm.example.com/acme/com.onlyoffice.mattermost/comment?file=file-1&action=%7B%22action%22%3A%7B%22type%22%3A%22comment%22%2C%22data%22%3A%22c1%22%7D%7D)")
		assert.Contains(t, dm, "[Go to the post](https://mm.example.com/_redirect/pl/post-1)")
	}
}

func TestMentionMessageWithoutCommentOrActionLink(t *testing.T) {
	message := mentionMessage("@alice", "doc.docx", "  ", "https://mm/pl/p1", "")
	assert.Equal(t, "@alice mentioned you in the document **doc.docx**\n[Go to the post](https://mm/pl/p1)", message)
}

func TestMentionMessageTruncatesLongComment(t *testing.T) {
	message := mentionMessage("@alice", "doc.docx", strings.Repeat("ж", 600), "https://mm/pl/p1", "")
	assert.Contains(t, message, strings.Repeat("ж", 500)+"…")
	assert.NotContains(t, message, strings.Repeat("ж", 501))
}

func TestCommentLinkIgnoresInvalidActionLink(t *testing.T) {
	api := quietAPI()
	api.On("GetConfig").Return(siteURLConfig()).Maybe()
	for _, raw := range []string{``, `null`, `"text"`, `[1]`, `{broken`, `{"a":"` + strings.Repeat("x", 3000) + `"}`} {
		assert.Empty(t, commentLink(api, nil, "member", "file-1", json.RawMessage(raw)), raw)
	}
}

func TestCommentLinkTeam(t *testing.T) {
	action := json.RawMessage(`{"action":{"type":"comment","data":"c1"}}`)

	t.Run("channel team", func(t *testing.T) {
		api := quietAPI()
		api.On("GetConfig").Return(siteURLConfig())
		api.On("GetTeam", "team-1").Return(&mmModel.Team{Id: "team-1", Name: "acme"}, nil)
		link := commentLink(api, &mmModel.Channel{TeamId: "team-1"}, "member", "file-1", action)
		assert.True(t, strings.HasPrefix(link, "https://mm.example.com/acme/com.onlyoffice.mattermost/comment?file=file-1&action="), link)
		api.AssertNotCalled(t, "GetTeamsForUser", mock.Anything)
	})

	t.Run("no team anywhere opens the editor page", func(t *testing.T) {
		api := quietAPI()
		api.On("GetConfig").Return(siteURLConfig())
		api.On("GetTeamsForUser", "member").Return([]*mmModel.Team{}, nil)
		link := commentLink(api, &mmModel.Channel{Type: mmModel.ChannelTypeDirect}, "member", "file-1", action)
		assert.True(t, strings.HasPrefix(link, "https://mm.example.com/plugins/com.onlyoffice.mattermost/api/editor?file=file-1&action="), link)
	})
}

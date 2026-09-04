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
	"testing"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMentionsGetUsers(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		api := loggingAPI()
		handler := NewMentionsHandler(api, &stubBot{})
		recorder := httptest.NewRecorder()
		handler.GetUsers(recorder, httptest.NewRequest(http.MethodGet, "/mentions/users", nil))

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("success with search and pagination", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{Id: "file-1", PostId: "post-1"}, nil)
		api.On("GetPost", "post-1").Return(&mmModel.Post{Id: "post-1", ChannelId: "ch-1"}, nil)
		api.On("GetChannel", "ch-1").Return(&mmModel.Channel{Id: "ch-1"}, nil)
		api.On("GetUsersInChannel", "ch-1", "username", 0, 200).Return([]*mmModel.User{
			{Id: "self", Username: "me", Email: "me@ex.com"},
			{Id: "u2", Username: "alice", Email: "alice@ex.com"},
			{Id: "u3", Username: "bob", Email: "bob@ex.com", IsBot: true},
			{Id: "u4", Username: "alex", Email: "alex@ex.com"},
		}, nil)

		handler := NewMentionsHandler(api, &stubBot{})
		req := httptest.NewRequest(http.MethodGet, "/mentions/users?file=file-1&search=al&from=0&count=1", nil)
		req.Header.Set(tools.MMAuthHeader, "self")
		recorder := httptest.NewRecorder()
		handler.GetUsers(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response model.MentionUsersResponse

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
		require.Len(t, response.Users, 1)

		assert.Equal(t, "alice", response.Users[0].Name)
		assert.True(t, response.HasMore)
	})
}

func TestMentionsSendNotifications(t *testing.T) {
	t.Run("bad body", func(t *testing.T) {
		api := loggingAPI()
		handler := NewMentionsHandler(api, &stubBot{})
		recorder := httptest.NewRecorder()
		handler.SendNotifications(recorder, httptest.NewRequest(http.MethodPost, "/mentions/notify", bytes.NewBufferString("{")))

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("missing file id", func(t *testing.T) {
		api := loggingAPI()
		handler := NewMentionsHandler(api, &stubBot{})
		body, _ := json.Marshal(model.MentionNotifyRequest{Emails: []string{"a@ex.com"}})
		recorder := httptest.NewRecorder()
		handler.SendNotifications(recorder, httptest.NewRequest(http.MethodPost, "/mentions/notify", bytes.NewReader(body)))

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("success", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{Id: "file-1", Name: "doc.docx", PostId: "post-1"}, nil)
		api.On("GetPost", "post-1").Return(&mmModel.Post{Id: "post-1", ChannelId: "ch-1"}, nil)
		api.On("GetChannel", "ch-1").Return(&mmModel.Channel{Id: "ch-1", TeamId: "team-1"}, nil)
		api.On("GetTeam", "team-1").Return(&mmModel.Team{Id: "team-1", Name: "town"}, nil)
		api.On("GetConfig").Return(siteURLConfig())
		api.On("GetUser", "author").Return(&mmModel.User{Id: "author", Username: "alice"}, nil)
		api.On("GetUserByEmail", "bob@ex.com").Return(&mmModel.User{Id: "bob"}, nil)

		bot := &stubBot{}
		handler := NewMentionsHandler(api, bot)
		body, _ := json.Marshal(model.MentionNotifyRequest{
			FileID: "file-1",
			Emails: []string{"bob@ex.com"},
		})
		req := httptest.NewRequest(http.MethodPost, "/mentions/notify", bytes.NewReader(body))
		req.Header.Set(tools.MMAuthHeader, "author")
		recorder := httptest.NewRecorder()
		handler.SendNotifications(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		require.Len(t, bot.dms, 1)

		assert.Contains(t, bot.dms[0], "@alice")
	})
}

func TestConvertHandler(t *testing.T) {
	fm := newFormatManager(t)

	t.Run("unauthorized", func(t *testing.T) {
		api := loggingAPI()
		handler := NewConvertHandler(api, validDESConfig(), fm, crypto.NewJwtManager(), &stubCommandClient{})
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, httptest.NewRequest(http.MethodPost, "/convert", bytes.NewBufferString("{}")))

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("bad body", func(t *testing.T) {
		api := loggingAPI()
		handler := NewConvertHandler(api, validDESConfig(), fm, crypto.NewJwtManager(), &stubCommandClient{})
		req := httptest.NewRequest(http.MethodPost, "/convert", bytes.NewBufferString("{"))
		req.Header.Set(tools.MMAuthHeader, "user-1")
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("not owner", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{Id: "file-1", CreatorId: "owner", Extension: "doc"}, nil)
		api.On("GetUser", "user-1").Return(&mmModel.User{Id: "user-1", Locale: "en-US"}, nil)
		handler := NewConvertHandler(api, validDESConfig(), fm, crypto.NewJwtManager(), &stubCommandClient{})
		body, _ := json.Marshal(model.ConvertFileRequest{FileID: "file-1"})
		req := httptest.NewRequest(http.MethodPost, "/convert", bytes.NewReader(body))
		req.Header.Set(tools.MMAuthHeader, "user-1")
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, req)

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("conversion api error json", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{
			Id: "file-1", CreatorId: "user-1", Extension: "doc", Name: "a.doc", ChannelId: "ch-1", PostId: "post-1",
		}, nil)
		api.On("GetUser", "user-1").Return(&mmModel.User{Id: "user-1", Locale: "en-US"}, nil)
		api.On("GetConfig").Return(siteURLConfig())

		commandClient := &stubCommandClient{convertResp: client.ConvertResponse{Error: 5}}
		handler := NewConvertHandler(api, validDESConfig(), fm, crypto.NewJwtManager(), commandClient)
		body, _ := json.Marshal(model.ConvertFileRequest{FileID: "file-1", OutputType: "docx"})
		req := httptest.NewRequest(http.MethodPost, "/convert", bytes.NewReader(body))
		req.Header.Set(tools.MMAuthHeader, "user-1")
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response model.ConvertFileResponse

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))

		assert.Equal(t, 5, response.Error)
	})

	t.Run("success", func(t *testing.T) {
		download := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("converted-bytes"))
		}))

		t.Cleanup(download.Close)

		api := loggingAPI()
		api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{
			Id: "file-1", CreatorId: "user-1", Extension: "doc", Name: "a.doc", ChannelId: "ch-1", PostId: "post-1",
		}, nil)
		api.On("GetUser", "user-1").Return(&mmModel.User{Id: "user-1", Locale: "en-US"}, nil)
		api.On("GetConfig").Return(siteURLConfig())
		api.On("CreateUploadSession", mock.AnythingOfType("*model.UploadSession")).Return(&mmModel.UploadSession{Id: "up-1"}, nil)
		api.On("UploadData", mock.AnythingOfType("*model.UploadSession"), mock.Anything).Return(&mmModel.FileInfo{Id: "file-2"}, nil)
		api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&mmModel.Post{Id: "post-2"}, nil)

		commandClient := &stubCommandClient{convertResp: client.ConvertResponse{
			Error: 0, FileURL: download.URL, FileType: "docx",
		}}

		handler := NewConvertHandler(api, validDESConfig(), fm, crypto.NewJwtManager(), commandClient)
		body, _ := json.Marshal(model.ConvertFileRequest{FileID: "file-1", OutputType: "docx"})
		req := httptest.NewRequest(http.MethodPost, "/convert", bytes.NewReader(body))
		req.Header.Set(tools.MMAuthHeader, "user-1")
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response model.ConvertFileResponse

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))

		assert.Equal(t, 0, response.Error)
	})
}

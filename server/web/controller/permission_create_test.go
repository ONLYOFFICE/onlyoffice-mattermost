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

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPermissionsGet(t *testing.T) {
	helper := newFileHelper(t)

	t.Run("file error", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetFileInfo", "missing").Return((*mmModel.FileInfo)(nil), mmModel.NewAppError("GetFileInfo", "fail", nil, "fail", 404))
		handler := NewPermissionsHandler(api, validDESConfig(), helper, &stubBot{})
		recorder := httptest.NewRecorder()
		handler.GetPermissions(recorder, httptest.NewRequest(http.MethodGet, "/permissions?file=missing", nil))

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("forbidden for non-author", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{Id: "file-1", PostId: "post-1"}, nil)
		api.On("GetPost", "post-1").Return(&mmModel.Post{Id: "post-1", UserId: "author"}, nil)
		handler := NewPermissionsHandler(api, validDESConfig(), helper, &stubBot{})
		req := httptest.NewRequest(http.MethodGet, "/permissions?file=file-1", nil)
		req.Header.Set(tools.MMAuthHeader, "other")
		recorder := httptest.NewRecorder()
		handler.GetPermissions(recorder, req)

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("success", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{Id: "file-1", PostId: "post-1"}, nil)
		api.On("GetPost", "post-1").Return(&mmModel.Post{Id: "post-1", UserId: "author"}, nil)
		handler := NewPermissionsHandler(api, validDESConfig(), helper, &stubBot{})
		req := httptest.NewRequest(http.MethodGet, "/permissions?file=file-1", nil)
		req.Header.Set(tools.MMAuthHeader, "author")
		recorder := httptest.NewRecorder()
		handler.GetPermissions(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}

func TestPermissionsSet(t *testing.T) {
	helper := newFileHelper(t)

	t.Run("bad json", func(t *testing.T) {
		api := loggingAPI()
		handler := NewPermissionsHandler(api, validDESConfig(), helper, &stubBot{})
		recorder := httptest.NewRecorder()
		handler.SetPermissions(recorder, httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBufferString("{")))

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("empty list", func(t *testing.T) {
		api := loggingAPI()
		handler := NewPermissionsHandler(api, validDESConfig(), helper, &stubBot{})
		body, _ := json.Marshal([]model.PostPermission{})
		recorder := httptest.NewRecorder()
		handler.SetPermissions(recorder, httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewReader(body)))

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("mismatched file ids", func(t *testing.T) {
		api := loggingAPI()
		handler := NewPermissionsHandler(api, validDESConfig(), helper, &stubBot{})
		body, _ := json.Marshal([]model.PostPermission{
			{FileID: "a", UserID: "u1", Permissions: model.Permissions{Edit: true}},
			{FileID: "b", UserID: "u2", Permissions: model.Permissions{Edit: false}},
		})

		recorder := httptest.NewRecorder()
		handler.SetPermissions(recorder, httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewReader(body)))

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("success notifies user", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{
			Id: "file-1", Name: "doc.docx", PostId: "post-1",
		}, nil)
		api.On("GetPost", "post-1").Return(&mmModel.Post{
			Id: "post-1", UserId: "author", ChannelId: "ch-1",
			FileIds: mmModel.StringArray{"file-1"},
		}, nil)
		api.On("GetChannel", "ch-1").Return(&mmModel.Channel{Id: "ch-1", TeamId: "team-1"}, nil)
		api.On("GetTeam", "team-1").Return(&mmModel.Team{Id: "team-1", Name: "town"}, nil)
		api.On("UpdatePost", mock.AnythingOfType("*model.Post")).Return(&mmModel.Post{Id: "post-1"}, nil)
		api.On("GetConfig").Return(siteURLConfig())

		bot := &stubBot{}
		handler := NewPermissionsHandler(api, validDESConfig(), helper, bot)
		body, _ := json.Marshal([]model.PostPermission{
			{FileID: "file-1", UserID: "editor", Permissions: model.Permissions{Edit: true}},
		})

		req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewReader(body))
		req.Header.Set(tools.MMAuthHeader, "author")
		recorder := httptest.NewRecorder()
		handler.SetPermissions(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		require.NotEmpty(t, bot.dms)
	})
}

func TestCreateHandler(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		api := loggingAPI()
		handler := NewCreateHandler(api, validDESConfig())
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString("{}")))

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("bad json", func(t *testing.T) {
		api := loggingAPI()
		handler := NewCreateHandler(api, validDESConfig())
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString("{"))
		req.Header.Set(tools.MMAuthHeader, "user-1")
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		api := loggingAPI()
		handler := NewCreateHandler(api, validDESConfig())
		body, _ := json.Marshal(map[string]string{"channel_id": "ch"})
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req.Header.Set(tools.MMAuthHeader, "user-1")
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("success", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetChannel", "ch-1").Return(&mmModel.Channel{Id: "ch-1"}, nil)
		api.On("GetUser", "user-1").Return(&mmModel.User{Id: "user-1", Locale: "en"}, nil)
		api.On("CreateUploadSession", mock.AnythingOfType("*model.UploadSession")).Return(&mmModel.UploadSession{Id: "up-1"}, nil)
		api.On("UploadData", mock.AnythingOfType("*model.UploadSession"), mock.Anything).Return(&mmModel.FileInfo{Id: "file-1"}, nil)
		api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&mmModel.Post{Id: "post-1"}, nil)

		handler := NewCreateHandler(api, validDESConfig())
		body, _ := json.Marshal(model.NewFileRequest{
			ChannelID: "ch-1",
			FileName:  "Notes",
			FileType:  "docx",
		})

		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req.Header.Set(tools.MMAuthHeader, "user-1")
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}

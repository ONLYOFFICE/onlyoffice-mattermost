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
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	oomodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEditorHandlerError(t *testing.T) {
	api := loggingAPI()
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return((*model.User)(nil), model.NewAppError("GetUser", "fail", nil, "fail", 500))

	tmpl := template.Must(template.New("editor.html").Parse("ok"))
	handler := NewEditorHandler(api, validDESConfig(), newFileHelper(t), crypto.NewMD5Encoder(), crypto.NewJwtManager(), tmpl)

	req := httptest.NewRequest(http.MethodGet, "/editor?file=file-1", nil)
	req.Header.Set(tools.MMAuthHeader, "user-1")
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestEditorHandlerSuccess(t *testing.T) {
	api := loggingAPI()
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Name: "doc.docx", Extension: "docx", PostId: "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", UserId: "user-1", UpdateAt: 1}, nil)
	api.On("KVSetWithExpiry", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return((*model.AppError)(nil))
	api.On("KVGet", mock.AnythingOfType("string")).Return([]byte(""), (*model.AppError)(nil))

	tmpl := template.Must(template.New("editor.html").Parse("apijs={{.apijs}};mentions={{.mentionscode}};configlen={{len .config}}"))
	handler := NewEditorHandler(api, validDESConfig(), newFileHelper(t), crypto.NewMD5Encoder(), crypto.NewJwtManager(), tmpl)

	req := httptest.NewRequest(http.MethodGet, "/editor?file=file-1", nil)
	req.Header.Set(tools.MMAuthHeader, "user-1")
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	body := recorder.Body.String()

	assert.Contains(t, body, "apijs=https://docs.example.com/web-apps/apps/api/documents/api.js")
	assert.Contains(t, body, "mentions=")
}

func TestRefreshHandler(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetConfig").Return(siteURLConfig())
		api.On("GetUser", "user-1").Return((*model.User)(nil), model.NewAppError("GetUser", "fail", nil, "fail", 500))

		handler := NewRefreshHandler(api, validDESConfig(), newFileHelper(t), crypto.NewMD5Encoder(), crypto.NewJwtManager())
		req := httptest.NewRequest(http.MethodGet, "/editor/config?file=file-1", nil)
		req.Header.Set(tools.MMAuthHeader, "user-1")
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("success", func(t *testing.T) {
		api := loggingAPI()
		api.On("GetConfig").Return(siteURLConfig())
		api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
		api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
			Id: "file-1", Name: "doc.docx", Extension: "docx", PostId: "post-1",
		}, nil)
		api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", UserId: "user-1", UpdateAt: 1}, nil)
		api.On("KVSetWithExpiry", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("int64")).Return(nil)

		handler := NewRefreshHandler(api, validDESConfig(), newFileHelper(t), crypto.NewMD5Encoder(), crypto.NewJwtManager())
		req := httptest.NewRequest(http.MethodGet, "/editor/config?file=file-1", nil)
		req.Header.Set(tools.MMAuthHeader, "user-1")
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var config oomodel.Config

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&config))

		assert.Equal(t, "docx", config.Document.FileType)
		assert.NotEmpty(t, config.Token)
	})
}

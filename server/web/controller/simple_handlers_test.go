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

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNotFoundHandler(t *testing.T) {
	handler := NewNotFoundHandler()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusMovedPermanently, recorder.Code)
	assert.Equal(t, "https://onlyoffice.com", recorder.Header().Get("Location"))
}

func TestCodeHandler(t *testing.T) {
	api := loggingAPI()
	api.On("KVSetWithExpiry", mock.AnythingOfType("string"), []byte("user-1"), int64(120)).Return(nil)

	handler := NewCodeHandler(api, newFileHelper(t))
	req := httptest.NewRequest(http.MethodGet, "/code", nil)
	req.Header.Set(tools.MMAuthHeader, "user-1")
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var code string

	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&code))
	assert.NotEmpty(t, code)
}

func TestCodeHandlerKVError(t *testing.T) {
	api := loggingAPI()
	api.On("KVSetWithExpiry", mock.AnythingOfType("string"), mock.Anything, int64(120)).
		Return(mmModel.NewAppError("KVSetWithExpiry", "fail", nil, "fail", 500))

	handler := NewCodeHandler(api, newFileHelper(t))
	req := httptest.NewRequest(http.MethodGet, "/code", nil)
	req.Header.Set(tools.MMAuthHeader, "user-1")
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestImageHandler(t *testing.T) {
	t.Run("missing code", func(t *testing.T) {
		api := loggingAPI()
		handler := NewImageHandler(api)
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, httptest.NewRequest(http.MethodGet, "/image", nil))

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("kv error", func(t *testing.T) {
		api := loggingAPI()
		api.On("KVGet", "bad").Return([]byte(nil), mmModel.NewAppError("KVGet", "fail", nil, "fail", 500))
		handler := NewImageHandler(api)
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, httptest.NewRequest(http.MethodGet, "/image?code=bad", nil))

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("profile error", func(t *testing.T) {
		api := loggingAPI()
		api.On("KVGet", "ok").Return([]byte("user-1"), nil)
		api.On("GetProfileImage", "user-1").Return([]byte(nil), mmModel.NewAppError("GetProfileImage", "fail", nil, "fail", 500))
		handler := NewImageHandler(api)
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, httptest.NewRequest(http.MethodGet, "/image?code=ok", nil))

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("success", func(t *testing.T) {
		api := loggingAPI()
		api.On("KVGet", "ok").Return([]byte("user-1"), nil)
		api.On("GetProfileImage", "user-1").Return([]byte("\x89PNG\r\n\x1a\n"), nil)
		handler := NewImageHandler(api)
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, httptest.NewRequest(http.MethodGet, "/image?code=ok", nil))

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.NotEmpty(t, recorder.Body.Bytes())
	})
}

func TestConfigHandler(t *testing.T) {
	api := loggingAPI()
	formatManager := newFormatManager(t)

	t.Run("all formats", func(t *testing.T) {
		handler := NewConfigHandler(api, &configuration.Configuration{
			OwnerProtected: true,
			PluginsEnabled: true,
			MacrosEnabled:  true,
		}, formatManager)
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, httptest.NewRequest(http.MethodGet, "/config", nil))

		assert.Equal(t, http.StatusOK, recorder.Code)

		var resp model.FormatResponse

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp))

		assert.NotEmpty(t, resp.Formats)
		assert.True(t, resp.OwnerProtected)
		assert.True(t, resp.PluginsEnabled)
		assert.True(t, resp.MacrosEnabled)
	})

	t.Run("none formats", func(t *testing.T) {
		handler := NewConfigHandler(api, &configuration.Configuration{Formats: configuration.EmptyFormats}, formatManager)
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, httptest.NewRequest(http.MethodGet, "/config", nil))

		var resp model.FormatResponse

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp))

		assert.Empty(t, resp.Formats)
	})

	t.Run("explicit formats", func(t *testing.T) {
		handler := NewConfigHandler(api, &configuration.Configuration{Formats: "docx, xlsx"}, formatManager)
		recorder := httptest.NewRecorder()
		handler.Handle(recorder, httptest.NewRequest(http.MethodGet, "/config", nil))

		var resp model.FormatResponse

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp))

		assert.Equal(t, []string{"docx", "xlsx"}, resp.Formats)
	})
}

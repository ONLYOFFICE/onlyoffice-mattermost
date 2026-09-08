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
package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	model "github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
)

func TestPluginHTTP(t *testing.T) {
	api := containerAPI(t)
	plugin := &Plugin{}
	plugin.API = api
	plugin.configuration = validPluginConfig()

	require.NoError(t, plugin.reinitializeContainer(plugin.configuration))

	t.Cleanup(func() {
		require.NoError(t, plugin.OnDeactivate())
	})

	t.Run("health", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		plugin.ServeHTTP(nil, recorder, httptest.NewRequest(http.MethodGet, "/api/health", nil))

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("config", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		plugin.ServeHTTP(nil, recorder, httptest.NewRequest(http.MethodGet, "/api/config", nil))

		require.Equal(t, http.StatusOK, recorder.Code)

		var body map[string]any

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
		assert.Contains(t, body, "formats")
	})

	t.Run("callback rejects unsigned payload", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/callback?file=file-1", bytes.NewBufferString(`{"key":"k","status":1}`))
		plugin.ServeHTTP(nil, recorder, req)

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("callback accepts signed editing status", func(t *testing.T) {
		token, err := crypto.NewJwtManager().Sign([]byte(plugin.configuration.DESJwt), jwt.MapClaims{
			"key":    "doc-key",
			"status": float64(1),
		})

		require.NoError(t, err)

		body, _ := json.Marshal(map[string]any{
			"key":    "doc-key",
			"status": 1,
			"token":  token,
		})

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/callback?file=file-1", bytes.NewReader(body))
		plugin.ServeHTTP(nil, recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)

		var response model.CallbackResponse

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
		assert.Equal(t, int8(0), response.Error)
	})

	t.Run("download requires jwt", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		plugin.ServeHTTP(nil, recorder, httptest.NewRequest(http.MethodGet, "/api/download", nil))

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("code requires auth", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		plugin.ServeHTTP(nil, recorder, httptest.NewRequest(http.MethodGet, "/api/code", nil))

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("code with user header", func(t *testing.T) {
		api.On("KVSetWithExpiry", mock.AnythingOfType("string"), []byte("user-1"), int64(120)).Return(nil).Once()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/code", nil)
		req.Header.Set(tools.MMAuthHeader, "user-1")
		plugin.ServeHTTP(nil, recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("unknown route is not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		plugin.ServeHTTP(nil, recorder, httptest.NewRequest(http.MethodGet, "/api/does-not-exist", nil))

		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
}

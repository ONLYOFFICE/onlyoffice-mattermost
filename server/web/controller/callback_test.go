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
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/callback"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type stubCallbackHandler struct {
	last callback.Callback
	err  error
}

func (s *stubCallbackHandler) Handle(ctx context.Context, callback callback.Callback) error {
	s.last = callback
	return s.err
}

func callbackConfig() *configuration.Configuration {
	return &configuration.Configuration{
		DESJwt:       "callback-secret",
		DESJwtHeader: "AuthorizationJWT",
		DESJwtPrefix: "Bearer ",
	}
}

func signCallbackClaims(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token, err := crypto.NewJwtManager().Sign([]byte(secret), claims)
	require.NoError(t, err)
	return token
}

func TestCallbackHandlerSuccessWithBodyToken(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogError", mock.Anything).Return().Maybe()

	config := callbackConfig()
	token := signCallbackClaims(t, config.DESJwt, jwt.MapClaims{
		"key":    "doc-key",
		"status": float64(2),
		"url":    "https://docs.example.com/file",
	})

	stub := &stubCallbackHandler{}
	handler := NewCallbackHandler(api, config, crypto.NewJwtManager(), stub)

	body, _ := json.Marshal(map[string]any{
		"key":    "doc-key",
		"status": 2,
		"url":    "https://docs.example.com/file",
		"token":  token,
	})

	req := httptest.NewRequest(http.MethodPost, "/callback?file=file-1", bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var resp model.CallbackResponse

	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp))

	assert.Equal(t, int8(0), resp.Error)
	assert.Equal(t, "file-1", stub.last.FileID)
	assert.Equal(t, 2, stub.last.Status)
}

func TestCallbackHandlerSuccessWithHeaderToken(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything).Return().Maybe()

	config := callbackConfig()
	token := signCallbackClaims(t, config.DESJwt, jwt.MapClaims{
		"key":    "doc-key",
		"status": float64(1),
	})

	stub := &stubCallbackHandler{}
	handler := NewCallbackHandler(api, config, crypto.NewJwtManager(), stub)

	body, _ := json.Marshal(map[string]any{
		"key":    "doc-key",
		"status": 1,
	})

	req := httptest.NewRequest(http.MethodPost, "/callback?file=file-1", bytes.NewReader(body))
	req.Header.Set(config.DESJwtHeader, config.DESJwtPrefix+token)
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response model.CallbackResponse

	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))

	assert.Equal(t, int8(0), response.Error)
}

func TestCallbackHandlerMissingToken(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogError", mock.Anything).Return().Maybe()

	handler := NewCallbackHandler(api, callbackConfig(), crypto.NewJwtManager(), &stubCallbackHandler{})
	body, _ := json.Marshal(map[string]any{"key": "doc-key", "status": 1})
	req := httptest.NewRequest(http.MethodPost, "/callback?file=file-1", bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)

	var response model.CallbackResponse

	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))

	assert.Equal(t, int8(1), response.Error)
}

func TestCallbackHandlerInvalidBody(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogError", mock.Anything).Return().Maybe()

	handler := NewCallbackHandler(api, callbackConfig(), crypto.NewJwtManager(), &stubCallbackHandler{})
	req := httptest.NewRequest(http.MethodPost, "/callback?file=file-1", bytes.NewReader([]byte("{")))
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response model.CallbackResponse

	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))

	assert.Equal(t, int8(1), response.Error)
}

func TestCallbackHandlerValidationError(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogError", mock.Anything).Return().Maybe()

	handler := NewCallbackHandler(api, callbackConfig(), crypto.NewJwtManager(), &stubCallbackHandler{})
	body, _ := json.Marshal(map[string]any{"status": 1})
	req := httptest.NewRequest(http.MethodPost, "/callback?file=file-1", bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	var response model.CallbackResponse

	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))

	assert.Equal(t, int8(1), response.Error)
}

func TestCallbackHandlerInvalidJWT(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogError", mock.Anything).Return().Maybe()

	handler := NewCallbackHandler(api, callbackConfig(), crypto.NewJwtManager(), &stubCallbackHandler{})
	body, _ := json.Marshal(map[string]any{
		"key":    "doc-key",
		"status": 1,
		"token":  "not.a.jwt",
	})

	req := httptest.NewRequest(http.MethodPost, "/callback?file=file-1", bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

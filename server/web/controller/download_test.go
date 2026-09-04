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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	ctrlmodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	jwt "github.com/golang-jwt/jwt/v5"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func downloadConfig() *configuration.Configuration {
	return &configuration.Configuration{
		DESJwt:       "download-secret",
		DESJwtHeader: "AuthorizationJWT",
	}
}

func TestDownloadHandlerSuccess(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("GetFile", "file").Return([]byte("file-bytes"), nil)

	config := downloadConfig()
	token, err := crypto.NewJwtManager().Sign([]byte(config.DESJwt), ctrlmodel.DownloadTokenRequest{
		Payload: ctrlmodel.DownloadTokenPayload{URL: "https://mm.example.com/download?id=file"},
	})

	require.NoError(t, err)

	handler := NewDownloadHandler(api, config, crypto.NewJwtManager())
	req := httptest.NewRequest(http.MethodGet, "/download", nil)
	req.Header.Set(config.DESJwtHeader, "Bearer "+token)
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "file-bytes", recorder.Body.String())

	api.AssertExpectations(t)
}

func TestDownloadHandlerMissingToken(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogError", mock.Anything).Return().Maybe()

	handler := NewDownloadHandler(api, downloadConfig(), crypto.NewJwtManager())
	req := httptest.NewRequest(http.MethodGet, "/download", nil)
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestDownloadHandlerInvalidToken(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogError", mock.Anything).Return().Maybe()

	config := downloadConfig()
	handler := NewDownloadHandler(api, config, crypto.NewJwtManager())
	req := httptest.NewRequest(http.MethodGet, "/download", nil)
	req.Header.Set(config.DESJwtHeader, "Bearer bad.token")
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestDownloadHandlerFileNotFound(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("GetFile", "missing").Return([]byte(nil), mmModel.NewAppError("GetFile", "file.missing", nil, "missing", http.StatusNotFound))

	config := downloadConfig()
	token, err := crypto.NewJwtManager().Sign([]byte(config.DESJwt), jwt.MapClaims{
		"payload": map[string]string{"url": "https://mm.example.com/download?id=missing"},
	})

	require.NoError(t, err)

	handler := NewDownloadHandler(api, config, crypto.NewJwtManager())
	req := httptest.NewRequest(http.MethodGet, "/download", nil)
	req.Header.Set(config.DESJwtHeader, "Bearer "+token)
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

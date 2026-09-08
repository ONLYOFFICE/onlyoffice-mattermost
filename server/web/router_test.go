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
package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/callback"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller"
)

type routerHealthStub struct{}

func (routerHealthStub) Start() {
}

func (routerHealthStub) Stop() {
}

func (routerHealthStub) IsHealthy() bool {
	return true
}

func newTestRouter(t *testing.T, api *plugintest.API) http.Handler {
	t.Helper()
	formatManager, err := public.NewMapFormatManager()

	require.NoError(t, err)

	return NewRouter(
		api,
		&configuration.Configuration{
			DESAddress:   "https://docs.example.com",
			DESJwt:       "secret",
			DESJwtHeader: "AuthorizationJWT",
			DESJwtPrefix: "Bearer ",
		},
		crypto.NewJwtManager(),
		callback.New(&callback.Config{
			PluginAPI: api,
			Converter: converter.New(),
			Filestore: &fxStubFileBackend{},
			Bot:       &fxStubBot{},
		}),
		controller.NewMentionsHandler(api, &fxStubBot{}),
		controller.NewHealthHandler(api, routerHealthStub{}),
		file.New(formatManager),
		crypto.NewMD5Encoder(),
		formatManager,
		&integrationCommandClient{},
		&fxStubBot{},
	)
}

func TestRecoverRoutesCatchesPanic(t *testing.T) {
	handler := recoverRoutes(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	}))

	recorder := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/panic", nil))
	})
}

func TestTimeoutRoutesCompletesBeforeDeadline(t *testing.T) {
	handler := timeoutRoutes(200 * time.Millisecond)(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/fast", nil))

	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestTimeoutRoutesReturnsGatewayTimeout(t *testing.T) {
	handler := timeoutRoutes(20 * time.Millisecond)(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(200 * time.Millisecond):
			w.WriteHeader(http.StatusOK)
		}
	})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/slow", nil))

	assert.Equal(t, http.StatusGatewayTimeout, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "request timeout")
}

func TestNewRouterBundlePathError(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetBundlePath").Return("", errors.New("no bundle"))
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogWarn", mock.Anything).Return().Maybe()

	router := newTestRouter(t, api)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/health", nil))

	assert.Equal(t, http.StatusOK, recorder.Code)
	api.AssertCalled(t, "LogError", mock.Anything)
}

func TestNewRouterTemplateParseError(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetBundlePath").Return(t.TempDir(), nil)
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("LogDebug", mock.Anything).Return().Maybe()

	router := newTestRouter(t, api)

	require.NotNil(t, router)
	api.AssertCalled(t, "LogError", mock.Anything)
}

func TestNewRouterEditorRoutesUseTimeoutAndAuth(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetBundlePath").Return(integrationRoot(t), nil)
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogWarn", mock.Anything).Return().Maybe()
	api.On("KVGet", mock.Anything).Return([]byte(""), nil).Maybe()

	router := newTestRouter(t, api)

	for _, path := range []string{"/api/editor?file=file-1", "/api/editor/config?file=file-1"} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))

		assert.Equal(t, http.StatusForbidden, recorder.Code, path)
	}
}

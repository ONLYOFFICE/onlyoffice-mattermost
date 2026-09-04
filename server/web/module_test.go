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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/callback"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/health"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type fxStubBot struct{}

func (b *fxStubBot) BotCreateDM(message string, userID string) {}

func (b *fxStubBot) BotCreatePost(message string, channelID string) {}

func (b *fxStubBot) BotCreateReply(message string, channelID string, parentID string) {}

type fxStubFileBackend struct{}

func (s *fxStubFileBackend) DriverName() string { return "local" }

func (s *fxStubFileBackend) TestConnection() error { return nil }

func (s *fxStubFileBackend) Reader(path string) (filestore.ReadCloseSeeker, error) {
	return nil, nil
}

func (s *fxStubFileBackend) ReadFile(path string) ([]byte, error) { return nil, nil }

func (s *fxStubFileBackend) FileExists(path string) (bool, error) { return false, nil }

func (s *fxStubFileBackend) FileSize(path string) (int64, error) { return 0, nil }

func (s *fxStubFileBackend) FileModTime(path string) (time.Time, error) { return time.Time{}, nil }

func (s *fxStubFileBackend) CopyFile(oldPath, newPath string) error { return nil }

func (s *fxStubFileBackend) MoveFile(oldPath, newPath string) error { return nil }

func (s *fxStubFileBackend) WriteFile(fr io.Reader, path string) (int64, error) {
	return 0, nil
}

func (s *fxStubFileBackend) AppendFile(fr io.Reader, path string) (int64, error) {
	return 0, nil
}
func (s *fxStubFileBackend) RemoveFile(path string) error { return nil }
func (s *fxStubFileBackend) ListDirectory(path string) ([]string, error) {
	return nil, nil
}
func (s *fxStubFileBackend) ListDirectoryRecursively(path string) ([]string, error) {
	return nil, nil
}

func (s *fxStubFileBackend) RemoveDirectory(path string) error { return nil }

func (s *fxStubFileBackend) ZipReader(path string, deflate bool) (io.ReadCloser, error) {
	return nil, nil
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}

func TestWebModuleProvidesRouter(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetBundlePath").Return(moduleRoot(t), nil)
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("LogDebug", mock.Anything).Return().Maybe()

	formatManager, err := public.NewMapFormatManager()
	require.NoError(t, err)

	config := &configuration.Configuration{
		DESAddress:   "https://docs.example.com",
		DESJwt:       "secret",
		DESJwtHeader: "AuthorizationJWT",
		DESJwtPrefix: "Bearer ",
	}

	var router *mux.Router
	app := fxtest.New(t,
		fx.NopLogger,
		fx.Provide(
			func() plugin.API { return api },
			func() *configuration.Configuration { return config },
			func() public.FormatManager { return formatManager },
			func() bot.Bot { return &fxStubBot{} },
			func() filestore.FileBackend { return &fxStubFileBackend{} },
			controller.NewMentionsHandler,
			controller.NewHealthHandler,
		),
		crypto.Module,
		client.Module,
		converter.Module,
		file.Module,
		callback.Module,
		health.Module,
		Module,
		fx.Populate(&router),
	)

	app.RequireStart().RequireStop()

	require.NotNil(t, router)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var body map[string]any

	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))

	assert.Equal(t, true, body["healthy"])
}

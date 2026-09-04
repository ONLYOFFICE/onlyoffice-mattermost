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
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/common"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}

func localFileConfig(t *testing.T) *model.Config {
	t.Helper()
	config := &model.Config{}
	config.SetDefaults()
	driver := model.ImageDriverLocal
	dir := t.TempDir()
	config.FileSettings.DriverName = &driver
	config.FileSettings.Directory = &dir
	return config
}

func validPluginConfig() *configuration.Configuration {
	return &configuration.Configuration{
		DESAddress:   "https://docs.example.com",
		DESJwt:       "secret",
		DESJwtHeader: "AuthorizationJWT",
		DESJwtPrefix: "Bearer ",
	}
}

func baseAPI(t *testing.T) *plugintest.API {
	t.Helper()
	api := &plugintest.API{}
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogInfo", mock.Anything).Return().Maybe()
	api.On("LogWarn", mock.Anything).Return().Maybe()
	api.On("DisablePlugin", mock.Anything).Return(nil).Maybe()
	api.On("PublishWebSocketEvent", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	api.On("GetUsers", mock.Anything).Return([]*model.User{}, nil).Maybe()
	api.On("KVGet", mock.Anything).Return([]byte(""), nil).Maybe()
	api.On("KVSet", mock.Anything, mock.Anything).Return(nil).Maybe()
	return api
}

func containerAPI(t *testing.T) *plugintest.API {
	t.Helper()
	api := baseAPI(t)
	api.On("EnsureBotUser", mock.AnythingOfType("*model.Bot")).Return("bot-id", nil)
	api.On("GetBundlePath").Return(repoRoot(t), nil)
	api.On("SetProfileImage", "bot-id", mock.AnythingOfType("[]uint8")).Return(nil)
	api.On("GetLicense").Return((*model.License)(nil))
	api.On("GetUnsanitizedConfig").Return(localFileConfig(t))
	return api
}

type stubHealthChecker struct {
	started int
	stopped int
}

func (s *stubHealthChecker) Start()          { s.started++ }
func (s *stubHealthChecker) Stop()           { s.stopped++ }
func (s *stubHealthChecker) IsHealthy() bool { return true }

type stubCommandClient struct {
	response client.VersionResponse
	err      error
	url      string
}

func (c *stubCommandClient) SendVersion(commandURL string, request client.VersionRequest, timeout time.Duration) (client.VersionResponse, error) {
	c.url = commandURL
	return c.response, c.err
}

func (c *stubCommandClient) SendConvert(commandURL string, request client.ConvertRequest, timeout time.Duration) (client.ConvertResponse, error) {
	return client.ConvertResponse{}, nil
}

func TestInitializeContainerWiresDependencies(t *testing.T) {
	api := containerAPI(t)
	p := &Plugin{}
	p.API = api
	p.configuration = validPluginConfig()

	p.app = p.initializeContainer()
	require.NotNil(t, p.app)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, p.app.Start(ctx))
	t.Cleanup(func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer stopCancel()
		_ = p.app.Stop(stopCtx)
	})

	require.NotNil(t, p.router)
	require.NotNil(t, p.commandClient)
	require.NotNil(t, p.jwtManager)
	require.NotNil(t, p.healthChecker)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()
	p.ServeHTTP(nil, rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestInitializeContainerStartStopCycle(t *testing.T) {
	api := containerAPI(t)
	plugin := &Plugin{}
	plugin.API = api
	plugin.configuration = validPluginConfig()

	require.NoError(t, plugin.reinitializeContainer(plugin.configuration))
	assert.True(t, plugin.ready)
	require.NotNil(t, plugin.router)
	require.NotNil(t, plugin.healthChecker)
	require.NoError(t, plugin.OnDeactivate())
}

func TestOnActivate(t *testing.T) {
	t.Run("nil configuration", func(t *testing.T) {
		plugin := &Plugin{}
		err := plugin.OnActivate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("configuration error", func(t *testing.T) {
		configErr := errors.New("bad config")
		plugin := &Plugin{configuration: &configuration.Configuration{Error: configErr}}

		assert.Equal(t, configErr, plugin.OnActivate())
	})

	t.Run("starts health checker", func(t *testing.T) {
		checker := &stubHealthChecker{}
		plugin := &Plugin{
			configuration: validPluginConfig(),
			healthChecker: checker,
		}

		require.NoError(t, plugin.OnActivate())
		assert.Equal(t, 1, checker.started)
	})
}

func TestOnDeactivateNilDeps(t *testing.T) {
	plugin := &Plugin{}

	require.NoError(t, plugin.OnDeactivate())
}

func TestServeHTTPNilRouter(t *testing.T) {
	api := baseAPI(t)
	plugin := &Plugin{}
	plugin.API = api

	recorder := httptest.NewRecorder()
	plugin.ServeHTTP(nil, recorder, httptest.NewRequest(http.MethodGet, "/api/health", nil))

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetConfiguration(t *testing.T) {
	plugin := &Plugin{}
	assert.NotNil(t, plugin.getConfiguration())

	config := validPluginConfig()
	plugin.configuration = config

	assert.Equal(t, config, plugin.getConfiguration())
}

func TestSetConfiguration(t *testing.T) {
	api := baseAPI(t)
	plugin := &Plugin{}
	plugin.API = api

	config := validPluginConfig()
	plugin.setConfiguration(config)
	assert.Equal(t, "https://docs.example.com", plugin.getConfiguration().DESAddress)

	assert.Panics(t, func() {
		plugin.setConfiguration(config)
	})
}

func TestPublishConfigChange(t *testing.T) {
	t.Run("nil api", func(t *testing.T) {
		plugin := &Plugin{}
		plugin.publishConfigChange()
	})

	t.Run("publishes event", func(t *testing.T) {
		api := baseAPI(t)
		plugin := &Plugin{}
		plugin.API = api
		plugin.publishConfigChange()
		api.AssertCalled(t, "PublishWebSocketEvent", "config_changed", mock.Anything, mock.Anything)
	})
}

func TestLogErrorWithoutAPI(t *testing.T) {
	plugin := &Plugin{}
	plugin.logError("test message without api")
}

func TestHandleConfigError(t *testing.T) {
	api := baseAPI(t)
	plugin := &Plugin{}
	plugin.API = api

	cfg := validPluginConfig()
	err := errors.New("boom")
	got := plugin.handleConfigError(cfg, err, "failed")
	assert.Equal(t, err, got)
	assert.Equal(t, err, cfg.Error)
	assert.Equal(t, err, plugin.getConfiguration().Error)

	time.Sleep(150 * time.Millisecond)
	api.AssertCalled(t, "DisablePlugin", PluginID)
}

func TestValidateServerVersion(t *testing.T) {
	api := baseAPI(t)
	plugin := &Plugin{}
	plugin.API = api

	t.Run("empty", func(t *testing.T) {
		err := plugin.validateServerVersion(validPluginConfig(), "")
		assert.ErrorIs(t, err, common.ErrParseDocumentServerVersion)
		time.Sleep(150 * time.Millisecond)
	})

	t.Run("unparseable", func(t *testing.T) {
		err := plugin.validateServerVersion(validPluginConfig(), "x.0.0")
		assert.ErrorIs(t, err, common.ErrParseDocumentServerVersion)
		time.Sleep(150 * time.Millisecond)
	})

	t.Run("deprecated", func(t *testing.T) {
		err := plugin.validateServerVersion(validPluginConfig(), "7.5.0")
		assert.ErrorIs(t, err, common.ErrDeprecatedDocumentServerVersion)
		time.Sleep(150 * time.Millisecond)
	})

	t.Run("supported", func(t *testing.T) {
		require.NoError(t, plugin.validateServerVersion(validPluginConfig(), "8.2.0"))
		require.NoError(t, plugin.validateServerVersion(validPluginConfig(), "9.0.1"))
	})
}

func TestValidateDependencies(t *testing.T) {
	api := baseAPI(t)

	t.Run("nil command client", func(t *testing.T) {
		plugin := &Plugin{}
		plugin.API = api
		err := plugin.validateDependencies(validPluginConfig())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "command client is nil")
		time.Sleep(150 * time.Millisecond)
	})

	t.Run("nil jwt manager", func(t *testing.T) {
		plugin := &Plugin{commandClient: &stubCommandClient{}}
		plugin.API = api
		err := plugin.validateDependencies(validPluginConfig())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "jwt manager is nil")
		time.Sleep(150 * time.Millisecond)
	})

	t.Run("ok", func(t *testing.T) {
		plugin := &Plugin{
			commandClient: &stubCommandClient{},
			jwtManager:    crypto.NewJwtManager(),
		}
		plugin.API = api
		require.NoError(t, plugin.validateDependencies(validPluginConfig()))
	})
}

func TestCreateVersionTokenAndValidateDocumentServer(t *testing.T) {
	api := baseAPI(t)
	plugin := &Plugin{}
	plugin.API = api
	plugin.jwtManager = crypto.NewJwtManager()

	token, err := plugin.createVersionToken(validPluginConfig())
	require.NoError(t, err)
	require.NotEmpty(t, token)

	t.Run("command error", func(t *testing.T) {
		plugin.commandClient = &stubCommandClient{err: errors.New("unreachable")}
		err := plugin.validateDocumentServer(validPluginConfig())
		require.Error(t, err)
		time.Sleep(150 * time.Millisecond)
	})

	t.Run("response error code", func(t *testing.T) {
		plugin.commandClient = &stubCommandClient{response: client.VersionResponse{Error: 1}}
		err := plugin.validateDocumentServer(validPluginConfig())
		var dsErr *common.DocumentServerCommandResponseError
		require.ErrorAs(t, err, &dsErr)
		time.Sleep(150 * time.Millisecond)
	})

	t.Run("success", func(t *testing.T) {
		stub := &stubCommandClient{response: client.VersionResponse{Error: 0, Version: "8.2.1"}}
		plugin.commandClient = stub
		require.NoError(t, plugin.validateDocumentServer(validPluginConfig()))
		assert.Contains(t, stub.url, client.OnlyofficeCommandServicePath)
	})
}

func TestEnsureBot(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		api := baseAPI(t)
		api.On("EnsureBotUser", mock.AnythingOfType("*model.Bot")).Return("bot-1", nil)
		api.On("GetBundlePath").Return(repoRoot(t), nil)
		api.On("SetProfileImage", "bot-1", mock.AnythingOfType("[]uint8")).Return(nil)

		plugin := &Plugin{}
		plugin.API = api
		id, err := plugin.EnsureBot()
		require.NoError(t, err)
		assert.Equal(t, "bot-1", id)
	})

	t.Run("ensure bot user fails", func(t *testing.T) {
		api := baseAPI(t)
		api.On("EnsureBotUser", mock.AnythingOfType("*model.Bot")).
			Return("", model.NewAppError("EnsureBotUser", "fail", nil, "fail", http.StatusInternalServerError))

		plugin := &Plugin{}
		plugin.API = api
		_, err := plugin.EnsureBot()
		assert.ErrorIs(t, err, common.ErrCreateBotProfile)
	})

	t.Run("bundle path fails", func(t *testing.T) {
		api := baseAPI(t)
		api.On("EnsureBotUser", mock.AnythingOfType("*model.Bot")).Return("bot-1", nil)
		api.On("GetBundlePath").Return("", errors.New("no bundle"))

		plugin := &Plugin{}
		plugin.API = api
		_, err := plugin.EnsureBot()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no bundle")
	})

	t.Run("missing logo", func(t *testing.T) {
		api := baseAPI(t)
		api.On("EnsureBotUser", mock.AnythingOfType("*model.Bot")).Return("bot-1", nil)
		api.On("GetBundlePath").Return(t.TempDir(), nil)

		plugin := &Plugin{}
		plugin.API = api
		_, err := plugin.EnsureBot()
		assert.ErrorIs(t, err, common.ErrLoadBotProfileImage)
	})

	t.Run("set profile image fails", func(t *testing.T) {
		api := baseAPI(t)
		api.On("EnsureBotUser", mock.AnythingOfType("*model.Bot")).Return("bot-1", nil)
		api.On("GetBundlePath").Return(repoRoot(t), nil)
		api.On("SetProfileImage", "bot-1", mock.AnythingOfType("[]uint8")).
			Return(model.NewAppError("SetProfileImage", "fail", nil, "fail", http.StatusInternalServerError))

		plugin := &Plugin{}
		plugin.API = api
		_, err := plugin.EnsureBot()
		assert.ErrorIs(t, err, common.ErrSetBotProfileImage)
	})
}

func TestPrepareConfiguration(t *testing.T) {
	t.Run("load error", func(t *testing.T) {
		api := baseAPI(t)
		api.On("LoadPluginConfiguration", mock.Anything).Return(errors.New("load failed"))

		plugin := &Plugin{}
		plugin.API = api
		_, err := plugin.prepareConfiguration()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load plugin configuration")
	})

	t.Run("success", func(t *testing.T) {
		api := baseAPI(t)
		api.On("LoadPluginConfiguration", mock.AnythingOfType("*configuration.Configuration")).
			Run(func(args mock.Arguments) {
				cfg := args.Get(0).(*configuration.Configuration)
				*cfg = *validPluginConfig()
			}).
			Return(nil)

		plugin := &Plugin{}
		plugin.API = api
		cfg, err := plugin.prepareConfiguration()
		require.NoError(t, err)
		assert.Equal(t, "https://docs.example.com", cfg.DESAddress)
	})
}

func TestValidateConfiguration(t *testing.T) {
	t.Run("invalid config", func(t *testing.T) {
		api := baseAPI(t)
		api.On("LoadPluginConfiguration", mock.AnythingOfType("*configuration.Configuration")).
			Run(func(args mock.Arguments) {
				cfg := args.Get(0).(*configuration.Configuration)
				*cfg = configuration.Configuration{}
			}).
			Return(nil)

		plugin := &Plugin{}
		plugin.API = api
		err := plugin.validateConfiguration()
		require.Error(t, err)
		time.Sleep(150 * time.Millisecond)
	})

	t.Run("valid with document server", func(t *testing.T) {
		api := baseAPI(t)
		api.On("LoadPluginConfiguration", mock.AnythingOfType("*configuration.Configuration")).
			Run(func(args mock.Arguments) {
				cfg := args.Get(0).(*configuration.Configuration)
				*cfg = *validPluginConfig()
			}).
			Return(nil)

		plugin := &Plugin{}
		plugin.API = api
		plugin.commandClient = &stubCommandClient{response: client.VersionResponse{Error: 0, Version: "8.3.0"}}
		plugin.jwtManager = crypto.NewJwtManager()

		require.NoError(t, plugin.validateConfiguration())
		api.AssertCalled(t, "LogInfo", mock.MatchedBy(func(msg string) bool {
			return msg != "" && msg[0] != 0
		}))
	})
}

func TestProvideFormatManagerAndFileBackend(t *testing.T) {
	api := containerAPI(t)
	plugin := &Plugin{}
	plugin.API = api

	fm := plugin.provideFormatManager()
	require.NotNil(t, fm)

	backend := plugin.provideFileBackend()
	require.NotNil(t, backend)
}

func TestReinitializeContainerReplacesApp(t *testing.T) {
	api := containerAPI(t)
	plugin := &Plugin{}
	plugin.API = api
	plugin.configuration = validPluginConfig()

	require.NoError(t, plugin.reinitializeContainer(plugin.configuration))
	firstRouter := plugin.router
	require.NotNil(t, firstRouter)

	require.NoError(t, plugin.reinitializeContainer(plugin.configuration))
	require.NotNil(t, plugin.router)
	assert.True(t, plugin.ready)
	require.NoError(t, plugin.OnDeactivate())
}

func TestOnConfigurationChangeSuccess(t *testing.T) {
	docs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(client.VersionResponse{Error: 0, Version: "8.2.0"})
	}))
	t.Cleanup(docs.Close)

	api := containerAPI(t)
	api.On("LoadPluginConfiguration", mock.AnythingOfType("*configuration.Configuration")).
		Run(func(args mock.Arguments) {
			cfg := args.Get(0).(*configuration.Configuration)
			*cfg = configuration.Configuration{
				DESAddress:     docs.URL,
				DESJwt:         "secret",
				DESJwtHeader:   "AuthorizationJWT",
				DESJwtPrefix:   "Bearer ",
				Formats:        "docx",
				OwnerProtected: true,
			}
		}).
		Return(nil)

	plugin := &Plugin{}
	plugin.API = api
	plugin.configuration = &configuration.Configuration{
		DESAddress:   docs.URL,
		DESJwt:       "secret",
		DESJwtHeader: "AuthorizationJWT",
		DESJwtPrefix: "Bearer ",
		Formats:      "",
	}

	require.NoError(t, plugin.OnConfigurationChange())
	assert.True(t, plugin.ready)
	require.NotNil(t, plugin.router)
	api.AssertCalled(t, "PublishWebSocketEvent", "config_changed", mock.Anything, mock.Anything)

	require.NoError(t, plugin.OnDeactivate())
}

func TestOnConfigurationChangePrepareFails(t *testing.T) {
	api := baseAPI(t)
	api.On("LoadPluginConfiguration", mock.Anything).Return(errors.New("load failed"))

	plugin := &Plugin{}
	plugin.API = api
	err := plugin.OnConfigurationChange()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load plugin configuration")
}

func TestSendVersionRequestUsesStubClient(t *testing.T) {
	plugin := &Plugin{
		commandClient: &stubCommandClient{response: client.VersionResponse{Version: "8.1.0"}},
	}
	cfg := validPluginConfig()
	resp, err := plugin.sendVersionRequest(cfg, "token")
	require.NoError(t, err)
	assert.Equal(t, "8.1.0", resp.Version)
}

func TestServeHTTPWithRouter(t *testing.T) {
	api := baseAPI(t)
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	plugin := &Plugin{router: router}
	plugin.API = api

	recorder := httptest.NewRecorder()
	plugin.ServeHTTP(nil, recorder, httptest.NewRequest(http.MethodGet, "/ping", nil))
	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestEnsureBotReadsLogoBytes(t *testing.T) {
	logo, err := os.ReadFile(filepath.Join(repoRoot(t), "assets", "logo.png"))
	require.NoError(t, err)
	require.NotEmpty(t, logo)
}

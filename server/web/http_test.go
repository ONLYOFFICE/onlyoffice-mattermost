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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/callback"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/health"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller"
	model "github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
)

type integrationBot struct {
	mutex   sync.Mutex
	replies []string
}

func (b *integrationBot) BotCreateDM(message string, userID string) {
}

func (b *integrationBot) BotCreatePost(message string, channelID string) {
}

func (b *integrationBot) BotCreateReply(message string, channelID string, parentID string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.replies = append(b.replies, message)
}

type integrationCommandClient struct {
	convertResp client.ConvertResponse
	convertErr  error
	lastURL     string
	lastReq     client.ConvertRequest
}

func (c *integrationCommandClient) SendVersion(commandURL string, request client.VersionRequest, timeout time.Duration) (client.VersionResponse, error) {
	return client.VersionResponse{Version: "8.2.0"}, nil
}

func (c *integrationCommandClient) SendConvert(commandURL string, request client.ConvertRequest, timeout time.Duration) (client.ConvertResponse, error) {
	c.lastURL = commandURL
	c.lastReq = request
	return c.convertResp, c.convertErr
}

type integrationStore struct {
	mutex   sync.Mutex
	written map[string][]byte
	err     error
}

func newIntegrationStore() *integrationStore {
	return &integrationStore{written: map[string][]byte{}}
}

func (s *integrationStore) DriverName() string {
	return "local"
}

func (s *integrationStore) TestConnection() error {
	return nil
}

func (s *integrationStore) Reader(path string) (filestore.ReadCloseSeeker, error) {
	return nil, nil
}

func (s *integrationStore) ReadFile(path string) ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return append([]byte(nil), s.written[path]...), nil
}

func (s *integrationStore) FileExists(path string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.written[path]
	return ok, nil
}

func (s *integrationStore) FileSize(path string) (int64, error) {
	return 0, nil
}

func (s *integrationStore) FileModTime(path string) (time.Time, error) {
	return time.Time{}, nil
}

func (s *integrationStore) CopyFile(oldPath, newPath string) error {
	return nil
}

func (s *integrationStore) MoveFile(oldPath, newPath string) error {
	return nil
}

func (s *integrationStore) AppendFile(fr io.Reader, path string) (int64, error) {
	return 0, nil
}

func (s *integrationStore) RemoveFile(path string) error {
	return nil
}

func (s *integrationStore) ListDirectory(path string) ([]string, error) {
	return nil, nil
}

func (s *integrationStore) ListDirectoryRecursively(path string) ([]string, error) {
	return nil, nil
}

func (s *integrationStore) RemoveDirectory(path string) error {
	return nil
}

func (s *integrationStore) ZipReader(path string, deflate bool) (io.ReadCloser, error) {
	return nil, nil
}

func (s *integrationStore) WriteFile(fr io.Reader, path string) (int64, error) {
	if s.err != nil {
		return 0, s.err
	}

	data, err := io.ReadAll(fr)
	if err != nil {
		return 0, err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.written[path] = data
	return int64(len(data)), nil
}

type applicationContext struct {
	t        *testing.T
	router   *mux.Router
	api      *plugintest.API
	config   *configuration.Configuration
	store    *integrationStore
	bot      *integrationBot
	commands *integrationCommandClient
	jwt      crypto.JwtManager
}

func integrationRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}

func newApplicationContext(t *testing.T) *applicationContext {
	t.Helper()

	api := &plugintest.API{}
	api.On("GetBundlePath").Return(integrationRoot(t), nil)
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogWarn", mock.Anything).Return().Maybe()
	api.On("LogInfo", mock.Anything).Return().Maybe()
	api.On("PublishWebSocketEvent", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	formatManager, err := public.NewMapFormatManager()

	require.NoError(t, err)

	config := &configuration.Configuration{
		DESAddress:                 "https://docs.example.com",
		DESJwt:                     "integration-secret",
		DESJwtHeader:               "AuthorizationJWT",
		DESJwtPrefix:               "Bearer ",
		AllowPrivateDocumentServer: true,
		PluginsEnabled:             true,
		MacrosEnabled:              false,
	}

	store := newIntegrationStore()
	botInstance := &integrationBot{}
	commands := &integrationCommandClient{
		convertResp: client.ConvertResponse{FileURL: "https://docs.example.com/out.docx", FileType: "docx"},
	}

	var router *mux.Router
	app := fxtest.New(t,
		fx.NopLogger,
		fx.Provide(
			func() plugin.API { return api },
			func() *configuration.Configuration { return config },
			func() public.FormatManager { return formatManager },
			func() bot.Bot { return botInstance },
			func() filestore.FileBackend { return store },
			func() client.CommandClient { return commands },
			func() *http.Client {
				return client.NewHTTPClient(config)
			},
			controller.NewMentionsHandler,
			controller.NewHealthHandler,
		),
		crypto.Module,
		converter.Module,
		file.Module,
		callback.Module,
		health.Module,
		Module,
		fx.Populate(&router),
	)

	app.RequireStart()

	t.Cleanup(func() { app.RequireStop() })

	return &applicationContext{
		t:        t,
		router:   router,
		api:      api,
		config:   config,
		store:    store,
		bot:      botInstance,
		commands: commands,
		jwt:      crypto.NewJwtManager(),
	}
}

func (e *applicationContext) do(method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	e.t.Helper()
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req := httptest.NewRequest(method, path, reader)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	recorder := httptest.NewRecorder()
	e.router.ServeHTTP(recorder, req)
	return recorder
}

func (e *applicationContext) sign(claims jwt.Claims) string {
	e.t.Helper()
	token, err := e.jwt.Sign([]byte(e.config.DESJwt), claims)
	require.NoError(e.t, err)
	return token
}

var _ filestore.FileBackend = (*integrationStore)(nil)

func TestHealthAndConfig(t *testing.T) {
	context := newApplicationContext(t)

	healthResponse := context.do(http.MethodGet, "/api/health", nil, nil)

	require.Equal(t, http.StatusOK, healthResponse.Code)

	var healthBody map[string]any
	require.NoError(t, json.NewDecoder(healthResponse.Body).Decode(&healthBody))

	assert.Equal(t, true, healthBody["healthy"])

	configResponse := context.do(http.MethodGet, "/api/config", nil, nil)

	require.Equal(t, http.StatusOK, configResponse.Code)

	var configBody map[string]any

	require.NoError(t, json.NewDecoder(configResponse.Body).Decode(&configBody))
	assert.Contains(t, configBody, "formats")
	assert.Equal(t, true, configBody["pluginsEnabled"])
}

func TestCallbackEditingAndUnauthorized(t *testing.T) {
	context := newApplicationContext(t)

	unauthorized := context.do(http.MethodPost, "/api/callback?file=file-1", []byte(`{"key":"k","status":1}`), nil)

	require.Equal(t, http.StatusForbidden, unauthorized.Code)

	var errResp model.CallbackResponse

	require.NoError(t, json.NewDecoder(unauthorized.Body).Decode(&errResp))
	assert.Equal(t, int8(1), errResp.Error)

	token := context.sign(jwt.MapClaims{"key": "doc-key", "status": float64(1)})
	body, _ := json.Marshal(map[string]any{
		"key":    "doc-key",
		"status": 1,
		"token":  token,
	})

	ok := context.do(http.MethodPost, "/api/callback?file=file-1", body, nil)

	require.Equal(t, http.StatusOK, ok.Code)

	var okResp model.CallbackResponse

	require.NoError(t, json.NewDecoder(ok.Body).Decode(&okResp))
	assert.Equal(t, int8(0), okResp.Error)
}

func TestCallbackSavePersistsFile(t *testing.T) {
	context := newApplicationContext(t)

	docsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("saved-from-docs"))
	}))

	t.Cleanup(docsServer.Close)

	context.api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{
		Id:     "file-1",
		PostId: "post-1",
		Path:   "files/file-1.docx",
		Name:   "file-1.docx",
	}, nil)
	context.api.On("GetPost", "post-1").Return(&mmModel.Post{
		Id:        "post-1",
		ChannelId: "channel-1",
	}, nil)
	context.api.On("UpdatePost", mock.AnythingOfType("*model.Post")).Return(&mmModel.Post{Id: "post-1"}, nil)
	context.api.On("GetUser", "user-1").Return(&mmModel.User{Id: "user-1", Username: "alice"}, nil)

	for _, status := range []int{2, 6} {
		t.Run("status_"+strconv.Itoa(status), func(t *testing.T) {
			context.store.mutex.Lock()
			context.store.written = map[string][]byte{}
			context.store.mutex.Unlock()
			context.bot.mutex.Lock()
			context.bot.replies = nil
			context.bot.mutex.Unlock()

			token := context.sign(jwt.MapClaims{
				"key":    "doc-key",
				"status": float64(status),
				"url":    docsServer.URL,
			})

			body, _ := json.Marshal(map[string]any{
				"key":    "doc-key",
				"status": status,
				"url":    docsServer.URL,
				"users":  []string{"user-1"},
				"token":  token,
			})

			callbackResp := context.do(http.MethodPost, "/api/callback?file=file-1", body, nil)

			require.Equal(t, http.StatusOK, callbackResp.Code)

			var resp model.CallbackResponse

			require.NoError(t, json.NewDecoder(callbackResp.Body).Decode(&resp))
			assert.Equal(t, int8(0), resp.Error)

			context.store.mutex.Lock()
			result := append([]byte(nil), context.store.written["files/file-1.docx"]...)
			context.store.mutex.Unlock()

			assert.Equal(t, []byte("saved-from-docs"), result)
		})
	}
}

func TestDownloadRoundTrip(t *testing.T) {
	context := newApplicationContext(t)
	context.api.On("GetFile", "file-1").Return([]byte("binary-payload"), nil)

	missingResp := context.do(http.MethodGet, "/api/download", nil, nil)

	assert.Equal(t, http.StatusForbidden, missingResp.Code)

	badResp := context.do(http.MethodGet, "/api/download", nil, map[string]string{
		context.config.DESJwtHeader: "Bearer not.a.jwt",
	})

	assert.Equal(t, http.StatusForbidden, badResp.Code)

	okResp := context.do(http.MethodGet, "/api/download", nil, map[string]string{
		context.config.DESJwtHeader: "Bearer " + context.sign(model.DownloadTokenRequest{
			Payload: model.DownloadTokenPayload{URL: "https://mm.example.com/plugins/x/api/download?id=file-1"},
		}),
	})

	require.Equal(t, http.StatusOK, okResp.Code)
	assert.Equal(t, "binary-payload", okResp.Body.String())
}

func TestAuthMiddlewareAndCode(t *testing.T) {
	context := newApplicationContext(t)
	context.api.On("KVGet", mock.AnythingOfType("string")).Return([]byte(""), nil).Maybe()
	context.api.On("KVSetWithExpiry", mock.AnythingOfType("string"), mock.Anything, int64(120)).Return(nil).Maybe()

	denied := context.do(http.MethodGet, "/api/code", nil, nil)
	assert.Equal(t, http.StatusForbidden, denied.Code)

	okResp := context.do(http.MethodGet, "/api/code", nil, map[string]string{
		tools.MMAuthHeader: "user-1",
	})

	require.Equal(t, http.StatusOK, okResp.Code)

	var code string

	require.NoError(t, json.NewDecoder(okResp.Body).Decode(&code))
	assert.NotEmpty(t, code)

	context.api.ExpectedCalls = nil
	context.api.Calls = nil
	context.api.On("LogError", mock.Anything).Return().Maybe()
	context.api.On("LogDebug", mock.Anything).Return().Maybe()
	context.api.On("LogWarn", mock.Anything).Return().Maybe()
	context.api.On("LogInfo", mock.Anything).Return().Maybe()
	context.api.On("GetBundlePath").Return(integrationRoot(t), nil).Maybe()
	context.api.On("KVGet", "fallback-code").Return([]byte("user-2"), nil)
	context.api.On("KVSetWithExpiry", mock.AnythingOfType("string"), []byte("user-2"), int64(120)).Return(nil)

	codeResp := context.do(http.MethodGet, "/api/code?code=fallback-code", nil, nil)

	require.Equal(t, http.StatusOK, codeResp.Code)
}

func TestPermissionsRequiresAuthor(t *testing.T) {
	context := newApplicationContext(t)

	context.api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{
		Id:     "file-1",
		PostId: "post-1",
		Name:   "doc.docx",
	}, nil)
	context.api.On("GetPost", "post-1").Return(&mmModel.Post{
		Id:     "post-1",
		UserId: "author",
	}, nil)

	deniedResp := context.do(http.MethodGet, "/api/permissions?file=file-1", nil, map[string]string{
		tools.MMAuthHeader: "other-user",
	})

	assert.Equal(t, http.StatusForbidden, deniedResp.Code)

	context.api.On("GetUser", "author").Return(&mmModel.User{
		Id:       "author",
		Username: "author",
		Email:    "a@ex.com",
	}, nil).Maybe()

	okResp := context.do(http.MethodGet, "/api/permissions?file=file-1", nil, map[string]string{
		tools.MMAuthHeader: "author",
	})

	require.Equal(t, http.StatusOK, okResp.Code)
}

func TestConvertOwnerFlow(t *testing.T) {
	context := newApplicationContext(t)

	convertedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("docx-bytes"))
	}))

	t.Cleanup(convertedServer.Close)

	context.commands.convertResp = client.ConvertResponse{
		Error:    0,
		FileURL:  convertedServer.URL,
		FileType: "docx",
	}

	siteURL := "https://mm.example.com"
	config := &mmModel.Config{}
	config.ServiceSettings.SiteURL = &siteURL
	context.api.On("GetConfig").Return(config)
	context.api.On("GetFileInfo", "file-doc").Return(&mmModel.FileInfo{
		Id:        "file-doc",
		CreatorId: "user-1",
		PostId:    "post-1",
		Name:      "legacy.doc",
		Extension: "doc",
		ChannelId: "channel-1",
	}, nil)
	context.api.On("GetUser", "user-1").Return(&mmModel.User{Id: "user-1", Username: "alice", Locale: "en-US"}, nil)
	context.api.On("CreateUploadSession", mock.AnythingOfType("*model.UploadSession")).Return(&mmModel.UploadSession{Id: "up-1"}, nil)
	context.api.On("UploadData", mock.AnythingOfType("*model.UploadSession"), mock.Anything).Return(&mmModel.FileInfo{Id: "file-2"}, nil)
	context.api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&mmModel.Post{Id: "post-2"}, nil)

	body, _ := json.Marshal(map[string]any{
		"file_id":     "file-doc",
		"output_type": "docx",
	})

	convertResp := context.do(http.MethodPost, "/api/convert", body, map[string]string{
		tools.MMAuthHeader: "user-1",
	})

	require.Equal(t, http.StatusOK, convertResp.Code)

	var resp map[string]any

	require.NoError(t, json.NewDecoder(convertResp.Body).Decode(&resp))
	assert.Equal(t, float64(0), resp["error"])
	assert.NotEmpty(t, context.commands.lastURL)
}

func TestImageEndpoint(t *testing.T) {
	context := newApplicationContext(t)

	imageResp := context.do(http.MethodGet, "/api/image", nil, nil)

	assert.Equal(t, http.StatusForbidden, imageResp.Code)

	context.api.On("KVGet", "img-code").Return([]byte("user-1"), nil)
	context.api.On("GetProfileImage", "user-1").Return([]byte("avatar"), nil)

	okResp := context.do(http.MethodGet, "/api/image?code=img-code", nil, nil)

	require.Equal(t, http.StatusOK, okResp.Code)
	assert.NotEmpty(t, okResp.Body.Bytes())
}

func TestCreateFile(t *testing.T) {
	context := newApplicationContext(t)
	context.api.On("GetChannel", "channel-1").Return(&mmModel.Channel{Id: "channel-1"}, nil)
	context.api.On("GetUser", "user-1").Return(&mmModel.User{Id: "user-1", Locale: "en"}, nil)
	context.api.On("CreateUploadSession", mock.AnythingOfType("*model.UploadSession")).Return(&mmModel.UploadSession{Id: "up-1"}, nil)
	context.api.On("UploadData", mock.AnythingOfType("*model.UploadSession"), mock.Anything).Return(&mmModel.FileInfo{Id: "file-new"}, nil)
	context.api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&mmModel.Post{Id: "post-new"}, nil)

	body, _ := json.Marshal(map[string]any{
		"channel_id": "channel-1",
		"file_name":  "Notes",
		"file_type":  "docx",
	})

	createResp := context.do(http.MethodPost, "/api/create", body, map[string]string{
		tools.MMAuthHeader: "user-1",
	})

	assert.Equal(t, http.StatusOK, createResp.Code)
}

func TestMentionsUsers(t *testing.T) {
	context := newApplicationContext(t)
	context.api.On("GetFileInfo", "file-m").Return(&mmModel.FileInfo{Id: "file-m", PostId: "post-m"}, nil)
	context.api.On("GetPost", "post-m").Return(&mmModel.Post{Id: "post-m", ChannelId: "channel-m"}, nil)
	context.api.On("GetChannel", "channel-m").Return(&mmModel.Channel{Id: "channel-m"}, nil)
	context.api.On("GetUsersInChannel", "channel-m", "username", 0, 200).Return([]*mmModel.User{
		{Id: "user-1", Username: "self"},
		{Id: "user-2", Username: "bob", Email: "b@ex.com"},
	}, nil)

	mentionsResp := context.do(http.MethodGet, "/api/mentions/users?file=file-m", nil, map[string]string{
		tools.MMAuthHeader: "user-1",
	})

	require.Equal(t, http.StatusOK, mentionsResp.Code)
	assert.Contains(t, mentionsResp.Body.String(), "bob")
}

func TestSetPermissions(t *testing.T) {
	context := newApplicationContext(t)
	siteURL := "https://mm.example.com"
	config := &mmModel.Config{}
	config.ServiceSettings.SiteURL = &siteURL
	context.api.On("GetConfig").Return(config).Maybe()
	context.api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{
		Id: "file-1", Name: "doc.docx", PostId: "post-1",
	}, nil)
	context.api.On("GetPost", "post-1").Return(&mmModel.Post{
		Id: "post-1", UserId: "author", ChannelId: "channel-1", FileIds: mmModel.StringArray{"file-1"},
	}, nil)
	context.api.On("GetChannel", "channel-1").Return(&mmModel.Channel{Id: "channel-1", TeamId: "team-1"}, nil)
	context.api.On("GetTeam", "team-1").Return(&mmModel.Team{Id: "team-1", Name: "team"}, nil).Maybe()
	context.api.On("UpdatePost", mock.AnythingOfType("*model.Post")).Return(&mmModel.Post{Id: "post-1"}, nil)

	body, _ := json.Marshal([]map[string]any{
		{
			"fileID":   "file-1",
			"userID":   "*",
			"username": "*",
			"permissions": map[string]any{
				"edit": false,
			},
		},
	})

	permissionsResp := context.do(http.MethodPost, "/api/permissions", body, map[string]string{
		tools.MMAuthHeader: "author",
	})

	assert.Equal(t, http.StatusOK, permissionsResp.Code)
}

func TestCallbackHeaderAndEditorAuth(t *testing.T) {
	context := newApplicationContext(t)
	context.api.On("KVGet", mock.AnythingOfType("string")).Return([]byte(""), nil).Maybe()

	body, _ := json.Marshal(map[string]any{"key": "doc-key", "status": 4})
	callbackResp := context.do(http.MethodPost, "/api/callback?file=file-1", body, map[string]string{
		context.config.DESJwtHeader: context.config.DESJwtPrefix + context.sign(jwt.MapClaims{"key": "doc-key", "status": float64(4)}),
	})

	require.Equal(t, http.StatusOK, callbackResp.Code)
	assert.Equal(t, http.StatusForbidden, context.do(http.MethodGet, "/api/editor?file=file-1", nil, nil).Code)
	assert.Equal(t, http.StatusForbidden, context.do(http.MethodGet, "/api/editor/config?file=file-1", nil, nil).Code)
}

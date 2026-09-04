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
package common

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type failingEncoder struct{}

func (failingEncoder) Encode(text string) (string, error) {
	return "", errors.New("encode failed")
}

type viewOnlyHelper struct {
	file.FileHelper
}

func (viewOnlyHelper) IsExtensionEditable(fileExt string) bool { return false }

func (viewOnlyHelper) GetFileType(fileExt string) (string, error) {
	return "word", nil
}

func (viewOnlyHelper) GenerateKey() string { return "fixed-code" }

func validEditorConfiguration() *configuration.Configuration {
	return &configuration.Configuration{
		DESAddress:     "https://docs.example.com",
		DESJwt:         "editor-secret",
		DESJwtHeader:   "AuthorizationJWT",
		DESJwtPrefix:   "Bearer ",
		DemoAddress:    "https://onlinedocs.docs.onlyoffice.com",
		PluginsEnabled: true,
		MacrosEnabled:  true,
	}
}

func newFileHelper(t *testing.T) file.FileHelper {
	t.Helper()
	formatManager, err := public.NewMapFormatManager()
	require.NoError(t, err)
	return file.New(formatManager)
}

func siteURLConfig() *model.Config {
	config := &model.Config{}
	siteURL := "https://mm.example.com"
	config.ServiceSettings.SiteURL = &siteURL
	return config
}

func editorRequest(userID, fileID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/editor?file="+fileID+"&lang=en", nil)
	req.Header.Set(tools.MMAuthHeader, userID)
	return req
}

func TestBuildEditorConfigHappyPath(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id:        "file-1",
		Name:      "notes.docx",
		Extension: "docx",
		PostId:    "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{
		Id:       "post-1",
		UserId:   "user-1",
		UpdateAt: 1700000000000,
	}, nil)
	api.On("KVSetWithExpiry", mock.AnythingOfType("string"), []byte("user-1"), int64(10)).Return(nil)

	config, docKey, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		validEditorConfiguration(),
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.NotEmpty(t, docKey)
	assert.Equal(t, docKey, config.Document.Key)
	assert.Equal(t, "docx", config.Document.FileType)
	assert.Equal(t, "notes.docx", config.Document.Title)
	assert.Equal(t, "word", config.DocumentType)
	assert.Equal(t, "user-1", config.EditorConfig.User.ID)
	assert.Equal(t, "alice", config.EditorConfig.User.Name)
	assert.Contains(t, config.EditorConfig.CallbackURL, "/callback?file=file-1")
	assert.Contains(t, config.Document.URL, "/download?id=file-1")
	assert.Contains(t, config.EditorConfig.User.Image, "/image?code=")
	assert.Equal(t, "en", config.EditorConfig.Lang)
	assert.Equal(t, "default-light", config.EditorConfig.Customization.UiTheme)
	assert.True(t, config.EditorConfig.Customization.Plugins)
	assert.True(t, config.EditorConfig.Customization.Macros)
	assert.True(t, config.Document.Permissions.Protect)
	assert.NotEmpty(t, config.Token)
	assert.Equal(t, "desktop", config.Type)
	assert.Empty(t, config.EditorConfig.Mode)

	api.AssertExpectations(t)
}

func TestBuildEditorConfigDarkThemeAndViewMode(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-2").Return(&model.User{Id: "user-2", Username: "bob"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id:        "file-1",
		Name:      "notes.docx",
		Extension: "docx",
		PostId:    "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{
		Id:       "post-1",
		UserId:   "author",
		UpdateAt: 1,
	}, nil)
	api.On("KVSetWithExpiry", mock.AnythingOfType("string"), []byte("user-2"), int64(10)).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/editor?file=file-1&dark=true", nil)
	req.Header.Set(tools.MMAuthHeader, "user-2")

	config, _, status, err := BuildEditorConfig(
		req,
		validEditorConfiguration(),
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "default-dark", config.EditorConfig.Customization.UiTheme)
	assert.Equal(t, "view", config.EditorConfig.Mode)
	assert.False(t, config.Document.Permissions.Edit)
}

func TestBuildEditorConfigNoCredentials(t *testing.T) {
	_, _, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		&configuration.Configuration{DemoAddress: "https://onlinedocs.docs.onlyoffice.com"},
		&plugintest.API{},
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.Error(t, err)

	assert.Equal(t, http.StatusUnauthorized, status)
	assert.Contains(t, err.Error(), "no valid credentials")
}

func TestBuildEditorConfigDemoActive(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Name: "a.docx", Extension: "docx", PostId: "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", UserId: "user-1", UpdateAt: 1}, nil)
	api.On("KVSetWithExpiry", mock.AnythingOfType("string"), mock.Anything, int64(10)).Return(nil)

	config := &configuration.Configuration{
		DemoEnabled:  true,
		DemoExpires:  time.Now().UnixMilli() + 60_000,
		DemoAddress:  "https://onlinedocs.docs.onlyoffice.com",
		DESAddress:   "https://onlinedocs.docs.onlyoffice.com",
		DESJwt:       "demo-secret",
		DESJwtHeader: "AuthorizationJWT",
		DESJwtPrefix: "Bearer ",
	}

	_, _, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		config,
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, status)
}

func TestBuildEditorConfigUserError(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return((*model.User)(nil), model.NewAppError("GetUser", "fail", nil, "missing", http.StatusNotFound))

	_, _, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		validEditorConfiguration(),
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.Error(t, err)

	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestBuildEditorConfigMissingFileQuery(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/editor", nil)
	req.Header.Set(tools.MMAuthHeader, "user-1")

	_, _, status, err := BuildEditorConfig(
		req,
		validEditorConfiguration(),
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
}

func TestBuildEditorConfigFileInfoError(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
	api.On("GetFileInfo", "file-1").Return((*model.FileInfo)(nil), model.NewAppError("GetFileInfo", "fail", nil, "missing", http.StatusNotFound))

	_, _, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		validEditorConfiguration(),
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.Error(t, err)

	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestBuildEditorConfigFormatNotAllowed(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Name: "a.docx", Extension: "docx", PostId: "post-1",
	}, nil)

	config := validEditorConfiguration()
	config.Formats = "xlsx"

	_, _, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		config,
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.Error(t, err)

	assert.Equal(t, http.StatusForbidden, status)
}

func TestBuildEditorConfigPostError(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Name: "a.docx", Extension: "docx", PostId: "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return((*model.Post)(nil), model.NewAppError("GetPost", "fail", nil, "missing", http.StatusNotFound))

	_, _, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		validEditorConfiguration(),
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.Error(t, err)

	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestBuildEditorConfigUnsupportedExtension(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Name: "a.bin", Extension: "bin", PostId: "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", UserId: "user-1", UpdateAt: 1}, nil)

	_, _, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		validEditorConfiguration(),
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.Error(t, err)

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Contains(t, err.Error(), "could not get file type")
}

func TestBuildEditorConfigOwnerProtectedNonOwner(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-2").Return(&model.User{Id: "user-2", Username: "bob"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Name: "a.docx", Extension: "docx", PostId: "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", UserId: "author", UpdateAt: 1}, nil)
	api.On("KVSetWithExpiry", mock.AnythingOfType("string"), []byte("user-2"), int64(10)).Return(nil)

	config := validEditorConfiguration()
	config.OwnerProtected = true

	built, _, status, err := BuildEditorConfig(
		editorRequest("user-2", "file-1"),
		config,
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.False(t, built.Document.Permissions.Protect)
}

func TestBuildEditorConfigOwnerProtectedAuthor(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "author").Return(&model.User{Id: "author", Username: "owner"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Name: "a.docx", Extension: "docx", PostId: "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", UserId: "author", UpdateAt: 1}, nil)
	api.On("KVSetWithExpiry", mock.AnythingOfType("string"), []byte("author"), int64(10)).Return(nil)

	config := validEditorConfiguration()
	config.OwnerProtected = true

	built, _, status, err := BuildEditorConfig(
		editorRequest("author", "file-1"),
		config,
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.True(t, built.Document.Permissions.Protect)
}

func TestBuildEditorConfigKVSetErrorIsLogged(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Name: "a.docx", Extension: "docx", PostId: "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", UserId: "user-1", UpdateAt: 1}, nil)
	api.On("KVSetWithExpiry", mock.AnythingOfType("string"), mock.Anything, int64(10)).
		Return(model.NewAppError("KVSetWithExpiry", "fail", nil, "kv", http.StatusInternalServerError))
	api.On("LogError", mock.Anything).Return()

	_, _, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		validEditorConfiguration(),
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, status)

	api.AssertCalled(t, "LogError", mock.Anything)
}
func TestBuildEditorConfigEncodeError(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Name: "a.docx", Extension: "docx", PostId: "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", UserId: "user-1", UpdateAt: 1}, nil)

	_, _, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		validEditorConfiguration(),
		api,
		failingEncoder{},
		crypto.NewJwtManager(),
		newFileHelper(t),
	)

	require.Error(t, err)

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Contains(t, err.Error(), "could not encode document key")
}

func TestBuildEditorConfigNonEditableExtension(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetConfig").Return(siteURLConfig())
	api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Username: "alice"}, nil)
	api.On("GetFileInfo", "file-1").Return(&model.FileInfo{
		Id: "file-1", Name: "a.pdf", Extension: "pdf", PostId: "post-1",
	}, nil)
	api.On("GetPost", "post-1").Return(&model.Post{Id: "post-1", UserId: "user-1", UpdateAt: 1}, nil)
	api.On("KVSetWithExpiry", "fixed-code", []byte("user-1"), int64(10)).Return(nil)

	helper := viewOnlyHelper{FileHelper: newFileHelper(t)}
	config, _, status, err := BuildEditorConfig(
		editorRequest("user-1", "file-1"),
		validEditorConfiguration(),
		api,
		crypto.NewMD5Encoder(),
		crypto.NewJwtManager(),
		helper,
	)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "view", config.EditorConfig.Mode)
	assert.False(t, config.Document.Permissions.Edit)
}

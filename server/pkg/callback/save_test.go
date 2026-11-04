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
package callback

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type fakeBot struct {
	mu      sync.Mutex
	replies []string
}

func (b *fakeBot) BotCreateDM(string, string)   {}
func (b *fakeBot) BotCreatePost(string, string) {}
func (b *fakeBot) BotCreateReply(message, _, _ string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.replies = append(b.replies, message)
}

type fixedClock struct{}

func (fixedClock) GetTimestamp() int64     { return 4242 }
func (fixedClock) GetTime(int64) time.Time { return time.UnixMilli(4242) }

const (
	fileID   = "file1"
	postID   = "post1"
	filePath = "data/Report.docx"
	original = "ORIGINAL-DOCUMENT"
)

type env struct {
	api     *plugintest.API
	store   filestore.FileBackend
	bot     *fakeBot
	dir     string
	updated bool
}

func newEnv(t *testing.T) *env {
	t.Helper()
	dir := t.TempDir()
	store, err := filestore.NewFileBackend(filestore.FileBackendSettings{DriverName: model.ImageDriverLocal, Directory: dir})
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "data"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, filePath), []byte(original), 0o600))

	e := &env{api: &plugintest.API{}, store: store, bot: &fakeBot{}, dir: dir}
	cfg := &model.Config{}
	cfg.SetDefaults()

	e.api.On("GetFileInfo", fileID).Return(&model.FileInfo{Id: fileID, PostId: postID, Path: filePath, Name: "Report.docx"}, nil)
	e.api.On("GetPost", postID).Return(&model.Post{Id: postID, ChannelId: "ch1", UpdateAt: 1}, nil)
	e.api.On("UpdatePost", mock.Anything).Run(func(mock.Arguments) { e.updated = true }).Return(&model.Post{}, nil).Maybe()
	e.api.On("GetConfig").Return(cfg).Maybe()
	e.api.On("GetUser", "bob").Return(&model.User{Id: "bob", Username: "bob"}, nil).Maybe()
	for _, level := range []string{"LogDebug", "LogWarn", "LogError", "LogInfo"} {
		e.api.On(level, mock.Anything).Maybe()
	}
	return e
}

func (e *env) stored(t *testing.T) string {
	b, err := os.ReadFile(filepath.Join(e.dir, filePath))
	require.NoError(t, err)
	return string(b)
}

func (e *env) save(cb Callback) error {
	cb.FileID = fileID
	if cb.Key == "" {
		cb.Key = "k"
	}
	return _saveFile(context.Background(), cb, e.api, fixedClock{}, e.store, e.bot)
}

func docServer(t *testing.T) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			_, _ = w.Write([]byte("NEW-DOCUMENT"))
		case "/empty":
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("<!DOCTYPE html><h1>404</h1>"))
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestSaveStoresDocumentAndNotifies(t *testing.T) {
	e, ds := newEnv(t), docServer(t)
	err := e.save(Callback{Status: 2, Users: []string{"bob"}, URL: ds.URL + "/ok", DocumentServer: ds.URL})
	require.NoError(t, err)
	assert.Equal(t, "NEW-DOCUMENT", e.stored(t))
	assert.True(t, e.updated, "document key (post.UpdateAt) is bumped after saving")
	assert.Equal(t, []string{"File Report.docx was updated by @bob"}, e.bot.replies)
}

func TestSaveFailureNeverCorruptsTheFile(t *testing.T) {
	for _, path := range []string{"/missing", "/empty"} {
		e, ds := newEnv(t), docServer(t)
		err := e.save(Callback{Status: 2, Users: []string{"bob"}, URL: ds.URL + path, DocumentServer: ds.URL})
		assert.Error(t, err, path)
		assert.Equal(t, original, e.stored(t), path+": stored file must stay intact")
		assert.False(t, e.updated, path+": document key must not change on failure")
		assert.Empty(t, e.bot.replies, path)
	}
}

func TestSaveRefusesForeignHosts(t *testing.T) {
	hit := false
	foreign := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { hit = true }))
	defer foreign.Close()

	e := newEnv(t)
	err := e.save(Callback{Status: 2, Users: []string{"bob"}, URL: foreign.URL + "/ok", DocumentServer: "https://docs.example.com"})
	assert.Error(t, err)
	assert.False(t, hit, "no request must be sent to a host other than the document server")
	assert.Equal(t, original, e.stored(t))
}

func TestSaveWithoutUsersDoesNotPanic(t *testing.T) {
	e, ds := newEnv(t), docServer(t)
	assert.NotPanics(t, func() {
		assert.NoError(t, e.save(Callback{Status: 2, URL: ds.URL + "/ok", DocumentServer: ds.URL}))
	})
	assert.Equal(t, "NEW-DOCUMENT", e.stored(t))
	assert.Empty(t, e.bot.replies)
}

func TestForceSaveDoesNotNotify(t *testing.T) {
	e, ds := newEnv(t), docServer(t)
	require.NoError(t, e.save(Callback{Status: 6, Users: []string{"bob"}, URL: ds.URL + "/ok", DocumentServer: ds.URL}))
	assert.Equal(t, "NEW-DOCUMENT", e.stored(t))
	assert.Empty(t, e.bot.replies)
}

func TestHandlePropagatesErrorsAndIgnoresUnknownStatuses(t *testing.T) {
	e, ds := newEnv(t), docServer(t)
	h := newHandler(e.api, fixedClock{}, e.store, e.bot)

	err := h.Handle(context.Background(), Callback{FileID: fileID, Key: "k", Status: 2, Users: []string{"bob"}, URL: ds.URL + "/missing", DocumentServer: ds.URL})
	assert.Error(t, err, "save errors reach the controller so Document Server gets error:1 and retries")

	assert.NoError(t, h.Handle(context.Background(), Callback{FileID: fileID, Key: "k", Status: 99}))
}

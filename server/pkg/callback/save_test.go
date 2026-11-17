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
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"testing"

	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
)

const original = "original-content"

type saveFixture struct {
	api   *plugintest.API
	store *stubFileBackend
	bot   *stubBot
	ops   *[]string
}

func newSaveFixture(t *testing.T, maxFileSize int64) saveFixture {
	t.Helper()
	ops := &[]string{}
	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogWarn", mock.Anything).Return().Maybe()
	api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{Id: "file-1", PostId: "post-1", Path: "files/file-1.docx", Name: "file-1.docx"}, nil)
	api.On("GetPost", "post-1").Return(&mmModel.Post{Id: "post-1", ChannelId: "channel-1"}, nil)
	api.On("GetChannel", "channel-1").Return(&mmModel.Channel{Id: "channel-1"}, nil).Maybe()
	api.On("GetUser", "user-1").Return(&mmModel.User{Id: "user-1", Username: "alice"}, nil).Maybe()
	api.On("UpdatePost", mock.AnythingOfType("*model.Post")).
		Run(func(mock.Arguments) { *ops = append(*ops, "update post") }).
		Return(&mmModel.Post{Id: "post-1"}, nil).Maybe()
	config := &mmModel.Config{}
	if maxFileSize > 0 {
		config.FileSettings.MaxFileSize = mmModel.NewPointer(maxFileSize)
	}
	api.On("GetConfig").Return(config).Maybe()

	store := &stubFileBackend{files: map[string][]byte{"files/file-1.docx": []byte(original)}, ops: ops}
	return saveFixture{api: api, store: store, bot: &stubBot{}, ops: ops}
}

func (f saveFixture) save(url string, users ...string) error {
	return registryContainer.RunHandler(
		context.Background(), 2,
		Callback{FileID: "file-1", Status: 2, URL: url, Users: users},
		f.api, client.NewHTTPClient(&configuration.Configuration{DESAllowPrivate: true}), converter.New(), f.store, f.bot,
	)
}

func (f saveFixture) assertUntouched(t *testing.T) {
	t.Helper()
	assert.Equal(t, original, string(f.store.files["files/file-1.docx"]), "the stored file must not change")
	assert.Len(t, f.store.files, 1, "no temporary file may be left behind")
	f.api.AssertNotCalled(t, "UpdatePost", mock.Anything)
	assert.Empty(t, f.bot.replies)
}

func serve(t *testing.T, handler http.HandlerFunc) string {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server.URL
}

func TestSaveReplacesFileOnlyAfterCompleteDownloadAndThenRotatesKey(t *testing.T) {
	f := newSaveFixture(t, 0)
	url := serve(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("new-content")) })

	require.NoError(t, f.save(url, "user-1"))
	assert.Equal(t, "new-content", string(f.store.files["files/file-1.docx"]))
	assert.Len(t, f.store.files, 1)
	require.Len(t, *f.ops, 3)
	assert.True(t, strings.HasPrefix((*f.ops)[0], "write files/file-1.docx.onlyoffice-"), "the download goes to a temporary file first")
	assert.True(t, strings.HasPrefix((*f.ops)[1], "move files/file-1.docx.onlyoffice-") && strings.HasSuffix((*f.ops)[1], "-> files/file-1.docx"))
	assert.Equal(t, "update post", (*f.ops)[2], "the document key changes only after the file is replaced")
	require.Len(t, f.bot.replies, 1)
}

func TestSaveRejectsErrorResponse(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusInternalServerError, http.StatusNoContent} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			f := newSaveFixture(t, 0)
			url := serve(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				_, _ = w.Write([]byte("<html>error page</html>"))
			})

			var downloadErr *FileDownloadError
			require.ErrorAs(t, f.save(url, "user-1"), &downloadErr)
			assert.Contains(t, downloadErr.Reason, strconv.Itoa(status))
			f.assertUntouched(t)
			assert.Empty(t, *f.ops, "nothing is written for an error response")
		})
	}
}

func TestSaveRejectsFileOverTheLimit(t *testing.T) {
	t.Run("declared size", func(t *testing.T) {
		f := newSaveFixture(t, 10)
		url := serve(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("0123456789ABC")) })

		var downloadErr *FileDownloadError
		require.ErrorAs(t, f.save(url, "user-1"), &downloadErr)
		assert.Contains(t, downloadErr.Reason, "exceeds the limit")
		f.assertUntouched(t)
		assert.Empty(t, *f.ops, "an oversized download is refused before writing")
	})

	t.Run("chunked without a declared size", func(t *testing.T) {
		f := newSaveFixture(t, 10)
		url := serve(t, func(w http.ResponseWriter, r *http.Request) {
			for i := 0; i < 5; i++ {
				_, _ = w.Write([]byte("0123"))
				w.(http.Flusher).Flush()
			}
		})

		var downloadErr *FileDownloadError
		require.ErrorAs(t, f.save(url, "user-1"), &downloadErr)
		assert.Contains(t, downloadErr.Reason, "exceeds the limit")
		f.assertUntouched(t)
	})

	t.Run("exactly at the limit is saved", func(t *testing.T) {
		f := newSaveFixture(t, 10)
		url := serve(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("0123456789")) })

		require.NoError(t, f.save(url, "user-1"))
		assert.Equal(t, "0123456789", string(f.store.files["files/file-1.docx"]))
	})
}

func TestSaveRejectsIncompleteDownload(t *testing.T) {
	f := newSaveFixture(t, 0)
	url := serve(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("0123456789"))
		conn, _, err := w.(http.Hijacker).Hijack()
		if err == nil {
			_ = conn.(*net.TCPConn).SetLinger(0)
			_ = conn.Close()
		}
	})

	require.Error(t, f.save(url, "user-1"))
	f.assertUntouched(t)
}

func TestSaveRejectsEmptyDownload(t *testing.T) {
	f := newSaveFixture(t, 0)
	url := serve(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	var downloadErr *FileDownloadError
	require.ErrorAs(t, f.save(url, "user-1"), &downloadErr)
	assert.Equal(t, "empty download", downloadErr.Reason)
	f.assertUntouched(t)
}

func TestSaveKeepsFileWhenReplacingFails(t *testing.T) {
	f := newSaveFixture(t, 0)
	f.store.moveErr = errors.New("storage unavailable")
	url := serve(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("new-content")) })

	var persistErr *FilePersistenceError
	require.ErrorAs(t, f.save(url, "user-1"), &persistErr)
	f.assertUntouched(t)
}

func TestSaveWithoutUsersChangesNothing(t *testing.T) {
	for name, users := range map[string][]string{"no users": nil, "empty user": {""}} {
		t.Run(name, func(t *testing.T) {
			f := newSaveFixture(t, 0)
			requested := false
			url := serve(t, func(w http.ResponseWriter, r *http.Request) {
				requested = true
				_, _ = w.Write([]byte("new-content"))
			})

			require.NotPanics(t, func() {
				assert.ErrorIs(t, f.save(url, users...), ErrInvalidUserID)
			})
			assert.False(t, requested, "the file is not even downloaded")
			f.assertUntouched(t)
		})
	}
}

func TestForceSaveDoesNotNeedUsers(t *testing.T) {
	f := newSaveFixture(t, 0)
	url := serve(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("forced")) })

	err := registryContainer.RunHandler(
		context.Background(), 6,
		Callback{FileID: "file-1", Status: 6, URL: url},
		f.api, client.NewHTTPClient(&configuration.Configuration{DESAllowPrivate: true}), converter.New(), f.store, f.bot,
	)
	require.NoError(t, err)
	assert.Equal(t, "forced", string(f.store.files["files/file-1.docx"]))
	assert.Empty(t, f.bot.replies)
}

func TestSaveRefusesArchivedChannel(t *testing.T) {
	f := newSaveFixture(t, 0)
	f.api.ExpectedCalls = slices.DeleteFunc(f.api.ExpectedCalls, func(c *mock.Call) bool { return c.Method == "GetChannel" })
	f.api.On("GetChannel", "channel-1").Return(&mmModel.Channel{Id: "channel-1", DeleteAt: 1700000000000}, nil)
	downloaded := false
	url := serve(t, func(w http.ResponseWriter, r *http.Request) {
		downloaded = true
		_, _ = w.Write([]byte("new-content"))
	})

	var persistErr *FilePersistenceError
	require.ErrorAs(t, f.save(url, "user-1"), &persistErr)
	assert.Contains(t, persistErr.Reason, "archived")
	assert.False(t, downloaded, "nothing is downloaded for an archived channel")
	f.assertUntouched(t)
	assert.Empty(t, *f.ops)
}

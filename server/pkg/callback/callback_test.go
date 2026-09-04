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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type stubBot struct {
	replies []string
}

func (b *stubBot) BotCreateDM(message string, userID string) {}

func (b *stubBot) BotCreatePost(message string, channelID string) {}

func (b *stubBot) BotCreateReply(message string, channelID string, parentID string) {
	b.replies = append(b.replies, message)
}

type stubFileBackend struct {
	written []byte
	path    string
	err     error
}

func (s *stubFileBackend) DriverName() string { return "local" }

func (s *stubFileBackend) TestConnection() error { return nil }

func (s *stubFileBackend) Reader(path string) (filestore.ReadCloseSeeker, error) {
	return nil, nil
}

func (s *stubFileBackend) ReadFile(path string) ([]byte, error) { return nil, nil }

func (s *stubFileBackend) FileExists(path string) (bool, error) { return false, nil }

func (s *stubFileBackend) FileSize(path string) (int64, error) { return 0, nil }

func (s *stubFileBackend) FileModTime(path string) (time.Time, error) {
	return time.Time{}, nil
}

func (s *stubFileBackend) CopyFile(oldPath, newPath string) error { return nil }

func (s *stubFileBackend) MoveFile(oldPath, newPath string) error { return nil }

func (s *stubFileBackend) AppendFile(fr io.Reader, path string) (int64, error) {
	return 0, nil
}

func (s *stubFileBackend) RemoveFile(path string) error { return nil }

func (s *stubFileBackend) ListDirectory(path string) ([]string, error) {
	return nil, nil
}

func (s *stubFileBackend) ListDirectoryRecursively(path string) ([]string, error) {
	return nil, nil
}

func (s *stubFileBackend) RemoveDirectory(path string) error { return nil }

func (s *stubFileBackend) ZipReader(path string, deflate bool) (io.ReadCloser, error) {
	return nil, nil
}

func (s *stubFileBackend) WriteFile(fr io.Reader, path string) (int64, error) {
	if s.err != nil {
		return 0, s.err
	}

	data, err := io.ReadAll(fr)
	if err != nil {
		return 0, err
	}

	s.written = data
	s.path = path
	return int64(len(data)), nil
}

var _ filestore.FileBackend = (*stubFileBackend)(nil)

func TestCallbackValidate(t *testing.T) {
	callback := &Callback{Key: "k", Status: 1, FileID: "f"}
	require.NoError(t, callback.Validate())

	assert.Error(t, (&Callback{Status: 1, FileID: "f"}).Validate())
}

func TestRegistryRunHandlerKnownStatuses(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogWarn", mock.Anything).Return().Maybe()

	bot := &stubBot{}
	store := &stubFileBackend{}
	converter := converter.New()
	ctx := context.Background()

	for _, status := range []int{1, 3, 4, 7} {
		err := registryContainer.RunHandler(ctx, status, Callback{FileID: "f1", Status: status}, api, converter, store, bot)
		assert.NoError(t, err, "status %d", status)
	}
}

func TestRegistryRunHandlerUnknownStatus(t *testing.T) {
	api := &plugintest.API{}
	err := registryContainer.RunHandler(
		context.Background(),
		99,
		Callback{FileID: "f1"},
		api,
		converter.New(),
		&stubFileBackend{},
		&stubBot{},
	)

	var missing *CallbackHandlerDoesNotExistError
	require.ErrorAs(t, err, &missing)
	assert.Equal(t, 99, missing.Code)
}

func TestSaveHandlerStatus2(t *testing.T) {
	download := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("updated-content"))
	}))

	t.Cleanup(download.Close)

	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("GetFileInfo", "file-1").Return(&mmModel.FileInfo{
		Id:     "file-1",
		PostId: "post-1",
		Path:   "files/file-1.docx",
		Name:   "file-1.docx",
	}, nil)
	api.On("GetPost", "post-1").Return(&mmModel.Post{
		Id:        "post-1",
		ChannelId: "channel-1",
	}, nil)
	api.On("UpdatePost", mock.AnythingOfType("*model.Post")).Return(&mmModel.Post{Id: "post-1"}, nil)
	api.On("GetUser", "user-1").Return(&mmModel.User{Id: "user-1", Username: "alice"}, nil)

	bot := &stubBot{}
	store := &stubFileBackend{}

	err := registryContainer.RunHandler(
		context.Background(),
		2,
		Callback{
			FileID: "file-1",
			Status: 2,
			URL:    download.URL,
			Users:  []string{"user-1"},
		},
		api,
		converter.New(),
		store,
		bot,
	)

	require.NoError(t, err)
	assert.Equal(t, []byte("updated-content"), store.written)
	assert.Equal(t, "files/file-1.docx", store.path)
	require.Len(t, bot.replies, 1)
	assert.Contains(t, bot.replies[0], "@alice")
	api.AssertExpectations(t)
}

func TestSaveHandlerEmptyURL(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything).Return().Maybe()

	err := registryContainer.RunHandler(
		context.Background(),
		6,
		Callback{FileID: "file-1", Status: 6, URL: ""},
		api,
		converter.New(),
		&stubFileBackend{},
		&stubBot{},
	)

	var invalidErr *InvalidFileDownloadURLError
	require.ErrorAs(t, err, &invalidErr)
}

func TestHandlerHandleDiscardsErrors(t *testing.T) {
	api := &plugintest.API{}
	handler := newHandler(api, converter.New(), &stubFileBackend{}, &stubBot{})
	err := handler.Handle(context.Background(), Callback{Status: 99, FileID: "f", Key: "k"})
	assert.NoError(t, err)
}

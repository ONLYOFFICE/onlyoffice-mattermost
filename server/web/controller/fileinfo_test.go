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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
)

func TestFileDetails(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "doc.docx"), make([]byte, 12345), 0o600))
	store, err := filestore.NewFileBackend(filestore.FileBackendSettings{DriverName: model.ImageDriverLocal, Directory: dir})
	require.NoError(t, err)
	fm, err := public.NewMapFormatManager()
	require.NoError(t, err)

	api := &plugintest.API{}
	api.On("GetFileInfo", "f1").Return(&model.FileInfo{Id: "f1", PostId: "p1", Path: "doc.docx", Name: "doc.docx", Extension: "docx", Size: 100, CreatorId: "alice", CreateAt: 1000}, nil)
	api.On("GetPost", "p1").Return(&model.Post{Id: "p1", ChannelId: "c1", UserId: "alice"}, nil)
	api.On("GetChannel", "c1").Return(&model.Channel{Id: "c1", TeamId: "t1", Type: model.ChannelTypePrivate}, nil)
	api.On("HasPermissionToChannel", mock.Anything, "c1", mock.Anything).Return(func(u, _ string, _ *model.Permission) bool { return u == "alice" || u == "bob" })
	api.On("HasPermissionToTeam", mock.Anything, mock.Anything, mock.Anything).Return(false)
	api.On("GetUser", "alice").Return(&model.User{Id: "alice", Username: "alice", FirstName: "Alice", LastName: "A"}, nil)
	api.On("GetUser", "bob").Return(&model.User{Id: "bob", Username: "bob"}, nil)
	api.On("KVGet", "onlyoffice_lastedit_f1").Return([]byte(`{"user_id":"bob","at":2000}`), nil)
	api.On("LogWarn", mock.Anything).Maybe()

	h := NewFileDetailsHandler(api, &configuration.Configuration{}, file.New(fm), store)
	get := func(user string) (*httptest.ResponseRecorder, FileDetailsResponse) {
		req := httptest.NewRequest(http.MethodGet, "/api/fileinfo?file=f1", nil)
		req.Header.Set(tools.MMAuthHeader, user)
		rw := httptest.NewRecorder()
		h.Handle(rw, req)
		var resp FileDetailsResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &resp)
		return rw, resp
	}

	rw, owner := get("alice")
	require.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, int64(12345), owner.Size, "size comes from the file store, not the stale database value")
	assert.Equal(t, "alice", owner.UploadedBy.Username)
	assert.Equal(t, int64(1000), owner.UploadedAt)
	assert.Equal(t, "bob", owner.EditedBy.Username)
	assert.Equal(t, int64(2000), owner.EditedAt)
	assert.True(t, owner.CanEdit)
	assert.True(t, owner.IsOwner)

	_, reader := get("bob")
	assert.False(t, reader.CanEdit, "channel members are read-only by default")
	assert.False(t, reader.IsOwner)

	rw, _ = get("eve")
	assert.Equal(t, http.StatusForbidden, rw.Code, "no details for users without channel access")
}

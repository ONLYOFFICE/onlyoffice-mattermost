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

	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/callback"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/common"
)

// FileUser is a user reference shown in the file details.
type FileUser struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

// FileDetailsResponse is the extra information shown in the file preview.
type FileDetailsResponse struct {
	Name       string    `json:"name"`
	Extension  string    `json:"extension"`
	Size       int64     `json:"size"`
	UploadedBy *FileUser `json:"uploaded_by"`
	UploadedAt int64     `json:"uploaded_at"`
	EditedBy   *FileUser `json:"edited_by,omitempty"`
	EditedAt   int64     `json:"edited_at,omitempty"`
	CanEdit    bool      `json:"can_edit"`
	IsOwner    bool      `json:"is_owner"`
}

type FileDetailsHandler struct {
	api           plugin.API
	configuration *configuration.Configuration
	fileHelper    file.FileHelper
	filestore     filestore.FileBackend
}

func NewFileDetailsHandler(
	api plugin.API,
	configuration *configuration.Configuration,
	fileHelper file.FileHelper,
	filestore filestore.FileBackend,
) FileDetailsHandler {
	return FileDetailsHandler{api: api, configuration: configuration, fileHelper: fileHelper, filestore: filestore}
}

func (h *FileDetailsHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(tools.MMAuthHeader)
	fileCtx, status, err := common.ResolveFileForUser(h.api, userID, r.URL.Query().Get("file"))
	if err != nil {
		h.api.LogWarn(onlyofficeLoggerPrefix + "file details: " + err.Error())
		rw.WriteHeader(status)
		return
	}

	info, post := fileCtx.File, fileCtx.Post
	resp := FileDetailsResponse{
		Name:       info.Name,
		Extension:  info.Extension,
		Size:       h.actualSize(info),
		UploadedBy: h.user(firstNonEmpty(info.CreatorId, post.UserId)),
		UploadedAt: info.CreateAt,
		CanEdit:    common.EffectivePermissions(h.configuration, h.fileHelper, userID, info, post).Edit,
		IsOwner:    userID == post.UserId,
	}

	if last := callback.GetLastEdit(h.api, info.Id); last != nil {
		resp.EditedAt = last.At
		resp.EditedBy = h.user(last.UserID)
	}

	common.WriteJSON(rw, resp)
}

// actualSize reads the size from the file store: FileInfo.Size in the database
// is not updated when ONLYOFFICE saves a new version of the file.
func (h *FileDetailsHandler) actualSize(info *mmModel.FileInfo) int64 {
	if h.filestore != nil {
		if size, err := h.filestore.FileSize(info.Path); err == nil {
			return size
		}
	}

	return info.Size
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func (h *FileDetailsHandler) user(id string) *FileUser {
	if id == "" {
		return nil
	}

	u, appErr := h.api.GetUser(id)
	if appErr != nil {
		return &FileUser{ID: id}
	}

	return &FileUser{ID: u.Id, Username: u.Username, DisplayName: u.GetDisplayName(mmModel.ShowNicknameFullName)}
}

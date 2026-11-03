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
	"fmt"
	"net/http"

	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// The Mattermost plugin API (GetFileInfo, GetPost, CreatePost, ...) does not
// perform any permission checks on behalf of the calling user, so every
// handler acting on behalf of a user must verify access explicitly.

// CanReadChannel reports whether the user may read the channel content,
// mirroring Mattermost's own rules (members, or any team member with
// read_public_channel for public channels; system admins always pass).
func CanReadChannel(api plugin.API, userID string, channel *mmModel.Channel) bool {
	if userID == "" || channel == nil {
		return false
	}

	if api.HasPermissionToChannel(userID, channel.Id, mmModel.PermissionReadChannelContent) ||
		api.HasPermissionToChannel(userID, channel.Id, mmModel.PermissionReadChannel) {
		return true
	}

	return channel.Type == mmModel.ChannelTypeOpen && channel.TeamId != "" &&
		api.HasPermissionToTeam(userID, channel.TeamId, mmModel.PermissionReadPublicChannel)
}

// CanPostFileToChannel reports whether the user may create a post with an
// attachment in the channel.
func CanPostFileToChannel(api plugin.API, userID string, channel *mmModel.Channel) bool {
	if userID == "" || channel == nil || channel.DeleteAt != 0 {
		return false
	}

	return api.HasPermissionToChannel(userID, channel.Id, mmModel.PermissionCreatePost) &&
		api.HasPermissionToChannel(userID, channel.Id, mmModel.PermissionUploadFile)
}

// FileContext bundles a file with the post and channel it belongs to.
type FileContext struct {
	File    *mmModel.FileInfo
	Post    *mmModel.Post
	Channel *mmModel.Channel
}

// ResolveFileForUser loads the file, its post and channel and verifies that the
// user can read the channel. On failure it returns a suitable HTTP status.
// Missing files and missing access both yield 403/404 without revealing which.
func ResolveFileForUser(api plugin.API, userID, fileID string) (*FileContext, int, error) {
	if fileID == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("missing file id")
	}

	fileInfo, appErr := api.GetFileInfo(fileID)
	if appErr != nil {
		return nil, http.StatusNotFound, fmt.Errorf("could not access file info %s: %s", fileID, appErr.Message)
	}

	if fileInfo.DeleteAt != 0 || fileInfo.PostId == "" {
		return nil, http.StatusNotFound, fmt.Errorf("file %s is deleted or not attached to a post", fileID)
	}

	post, appErr := api.GetPost(fileInfo.PostId)
	if appErr != nil {
		return nil, http.StatusNotFound, fmt.Errorf("could not access post %s: %s", fileInfo.PostId, appErr.Message)
	}

	channel, appErr := api.GetChannel(post.ChannelId)
	if appErr != nil {
		return nil, http.StatusNotFound, fmt.Errorf("could not access channel %s: %s", post.ChannelId, appErr.Message)
	}

	if !CanReadChannel(api, userID, channel) {
		return nil, http.StatusForbidden, fmt.Errorf("user %s has no access to channel %s (file %s)", userID, channel.Id, fileID)
	}

	return &FileContext{File: fileInfo, Post: post, Channel: channel}, http.StatusOK, nil
}

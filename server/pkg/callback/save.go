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
	"fmt"
	"io"
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"github.com/pkg/errors"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
)

const defaultMaxDownloadSize int64 = 100 * 1024 * 1024

var _ = registryContainer.RegisterHandler(2, _saveFile)
var _ = registryContainer.RegisterHandler(6, _saveFile)

func maxDownloadSize(api plugin.API) int64 {
	if config := api.GetConfig(); config != nil {
		if size := model.SafeDereference(config.FileSettings.MaxFileSize); size > 0 {
			return size
		}
	}

	return defaultMaxDownloadSize
}

func replaceFile(
	ctx context.Context,
	api plugin.API,
	httpClient *http.Client,
	store filestore.FileBackend,
	callback Callback,
	fileInfo *model.FileInfo,
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, callback.URL, nil)
	if err != nil {
		return errors.Wrap(err, onlyofficeLoggerCallbackPrefix+"failed to create request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, onlyofficeLoggerCallbackPrefix)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &FileDownloadError{FileID: callback.FileID, Reason: fmt.Sprintf("unexpected response status %d", resp.StatusCode)}
	}

	limit := maxDownloadSize(api)
	if resp.ContentLength > limit {
		return &FileDownloadError{FileID: callback.FileID, Reason: fmt.Sprintf("file size %d exceeds the limit of %d bytes", resp.ContentLength, limit)}
	}

	tmpPath := fileInfo.Path + ".onlyoffice-" + model.NewId() + ".tmp"
	cleanup := func() {
		if rmErr := store.RemoveFile(tmpPath); rmErr != nil {
			api.LogWarn(onlyofficeLoggerCallbackPrefix + "could not remove temporary file " + tmpPath + ": " + rmErr.Error())
		}
	}

	written, err := filestore.TryWriteFileContext(ctx, store, io.LimitReader(resp.Body, limit+1), tmpPath)
	if err != nil {
		cleanup()
		return &FilePersistenceError{FileID: callback.FileID, Reason: err.Error()}
	}

	var reason string
	switch {
	case written > limit:
		reason = fmt.Sprintf("file exceeds the limit of %d bytes", limit)
	case resp.ContentLength >= 0 && written != resp.ContentLength:
		reason = fmt.Sprintf("incomplete download: got %d of %d bytes", written, resp.ContentLength)
	case written == 0:
		reason = "empty download"
	}

	if reason != "" {
		cleanup()
		return &FileDownloadError{FileID: callback.FileID, Reason: reason}
	}

	if err := store.MoveFile(tmpPath, fileInfo.Path); err != nil {
		cleanup()
		return &FilePersistenceError{FileID: callback.FileID, Reason: err.Error()}
	}

	return nil
}

func _saveFile(
	ctx context.Context,
	callback Callback,
	api plugin.API,
	httpClient *http.Client,
	converter converter.TimeConverter,
	filestore filestore.FileBackend,
	bot bot.Bot,
) error {
	api.LogDebug(onlyofficeLoggerCallbackPrefix + "file " + callback.FileID + " save call")

	if callback.URL == "" {
		return &InvalidFileDownloadURLError{
			FileID: callback.FileID,
		}
	}

	if httpClient == nil {
		return ErrHTTPClientRequired
	}

	fileInfo, fileErr := api.GetFileInfo(callback.FileID)
	if fileErr != nil {
		return &FileNotFoundError{
			FileID: callback.FileID,
			Reason: fileErr.Error(),
		}
	}

	post, postErr := api.GetPost(fileInfo.PostId)
	if postErr != nil {
		return &FilePersistenceError{
			FileID: callback.FileID,
			Reason: postErr.Error(),
		}
	}

	var editor *model.User
	if callback.Status == 2 {
		if len(callback.Users) == 0 || callback.Users[0] == "" {
			return ErrInvalidUserID
		}

		user, userErr := api.GetUser(callback.Users[0])
		if userErr != nil {
			return &FilePersistenceError{
				FileID: callback.FileID,
				Reason: userErr.Error(),
			}
		}

		editor = user
	}

	if err := replaceFile(ctx, api, httpClient, filestore, callback, fileInfo); err != nil {
		return err
	}

	post.UpdateAt = converter.GetTimestamp()
	if _, uErr := api.UpdatePost(post); uErr != nil {
		return &FilePersistenceError{
			FileID: callback.FileID,
			Reason: uErr.Error(),
		}
	}

	if editor != nil {
		replyMsg := fmt.Sprintf("File %s was updated by @%s", fileInfo.Name, editor.Username)
		postID := post.Id
		if post.RootId != "" {
			postID = post.RootId
		}

		bot.BotCreateReply(replyMsg, post.ChannelId, postID)
	}

	return nil
}

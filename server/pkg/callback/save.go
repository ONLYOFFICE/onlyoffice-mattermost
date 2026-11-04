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
	"os"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"github.com/pkg/errors"
)

const (
	downloadTimeout     = 5 * time.Minute
	defaultMaxFileBytes = int64(100 << 20)
)

var _ = registryContainer.RegisterHandler(2, _saveFile)
var _ = registryContainer.RegisterHandler(6, _saveFile)

// _saveFile persists a document produced by ONLYOFFICE Docs (status 2 - ready
// for saving, status 6 - force save). The stored file is replaced only after
// the new version has been downloaded completely and successfully, so a failed
// download never corrupts the original. Any error is returned to the caller
// which reports it to Document Server (error: 1) so that it can retry.
func _saveFile(
	ctx context.Context,
	callback Callback,
	api plugin.API,
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

	if filestore == nil {
		return &FilePersistenceError{FileID: callback.FileID, Reason: "file backend is not initialized"}
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

	tmp, size, err := downloadToTempFile(ctx, callback.URL, callback.DocumentServer, maxFileSize(api))
	if err != nil {
		return &FilePersistenceError{
			FileID: callback.FileID,
			Reason: err.Error(),
		}
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	if _, storeErr := filestore.WriteFile(tmp, fileInfo.Path); storeErr != nil {
		return &FilePersistenceError{
			FileID: callback.FileID,
			Reason: storeErr.Error(),
		}
	}

	api.LogDebug(onlyofficeLoggerCallbackPrefix + fmt.Sprintf("file %s saved (%d bytes)", callback.FileID, size))

	// The document key is derived from post.UpdateAt, so bump it only after the
	// new content has been stored successfully.
	post.UpdateAt = converter.GetTimestamp()
	if _, uErr := api.UpdatePost(post); uErr != nil {
		return &FilePersistenceError{
			FileID: callback.FileID,
			Reason: uErr.Error(),
		}
	}

	if callback.Status == 2 {
		notifyFileUpdated(api, bot, callback, fileInfo.Name, post.Id, post.RootId, post.ChannelId)
	}

	return nil
}

func notifyFileUpdated(api plugin.API, bot bot.Bot, callback Callback, fileName, postID, rootID, channelID string) {
	if len(callback.Users) == 0 || callback.Users[0] == "" {
		api.LogWarn(onlyofficeLoggerCallbackPrefix + "file " + callback.FileID + " saved, but callback has no users to report")
		return
	}

	user, userErr := api.GetUser(callback.Users[0])
	if userErr != nil {
		api.LogWarn(onlyofficeLoggerCallbackPrefix + "file " + callback.FileID + " saved, but could not get user: " + userErr.Error())
		return
	}

	threadID := postID
	if rootID != "" {
		threadID = rootID
	}

	bot.BotCreateReply(fmt.Sprintf("File %s was updated by @%s", fileName, user.Username), channelID, threadID)
}

func maxFileSize(api plugin.API) int64 {
	if cfg := api.GetConfig(); cfg != nil && cfg.FileSettings.MaxFileSize != nil && *cfg.FileSettings.MaxFileSize > 0 {
		return *cfg.FileSettings.MaxFileSize
	}

	return defaultMaxFileBytes
}

// downloadToTempFile fetches url into a temporary file and returns it rewound
// to the beginning. The caller is responsible for closing and removing it.
func downloadToTempFile(ctx context.Context, url, documentServer string, limit int64) (*os.File, int64, error) {
	if err := tools.CheckDocumentServerURL(url, documentServer); err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to create download request")
	}

	resp, err := tools.NewDocumentServerClient(documentServer, downloadTimeout).Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to download file")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("document server returned HTTP %d for the file download", resp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "onlyoffice-save-*")
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to create temporary file")
	}

	cleanup := func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}

	size, err := io.Copy(tmp, io.LimitReader(resp.Body, limit+1))
	if err != nil {
		cleanup()
		return nil, 0, errors.Wrap(err, "failed to read downloaded file")
	}

	if size > limit {
		cleanup()
		return nil, 0, fmt.Errorf("downloaded file exceeds the maximum file size of %d bytes", limit)
	}

	if size == 0 {
		cleanup()
		return nil, 0, errors.New("downloaded file is empty")
	}

	if resp.ContentLength > 0 && size != resp.ContentLength {
		cleanup()
		return nil, 0, fmt.Errorf("incomplete download: got %d of %d bytes", size, resp.ContentLength)
	}

	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		cleanup()
		return nil, 0, errors.Wrap(err, "failed to rewind temporary file")
	}

	return tmp, size, nil
}

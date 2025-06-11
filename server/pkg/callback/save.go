/**
 *
 * (c) Copyright Ascensio System SIA 2025
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
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"github.com/pkg/errors"
)

var _ = registryContainer.RegisterHandler(2, _saveFile)
var _ = registryContainer.RegisterHandler(6, _saveFile)

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

	fileInfo, fileErr := api.GetFileInfo(callback.FileID)
	if fileErr != nil {
		return &FileNotFoundError{
			FileID: callback.FileID,
			Reason: fileErr.Error(),
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", callback.URL, nil)
	if err != nil {
		return errors.Wrap(err, onlyofficeLoggerCallbackPrefix+"failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return errors.Wrap(err, onlyofficeLoggerCallbackPrefix)
	}

	post, postErr := api.GetPost(fileInfo.PostId)
	if postErr != nil {
		return &FilePersistenceError{
			FileID: callback.FileID,
			Reason: postErr.Error(),
		}
	}

	post.UpdateAt = converter.GetTimestamp()
	_, uErr := api.UpdatePost(post)
	if uErr != nil {
		return &FilePersistenceError{
			FileID: callback.FileID,
			Reason: uErr.Error(),
		}
	}

	_, storeErr := filestore.WriteFile(resp.Body, fileInfo.Path)
	if storeErr != nil {
		return &FilePersistenceError{
			FileID: callback.FileID,
			Reason: storeErr.Error(),
		}
	}

	if callback.Status == 2 {
		last := callback.Users[0]
		if last == "" {
			return ErrInvalidUserID
		}

		user, userErr := api.GetUser(last)
		if userErr != nil {
			return &FilePersistenceError{
				FileID: callback.FileID,
				Reason: userErr.Error(),
			}
		}

		replyMsg := fmt.Sprintf("File %s was updated by @%s", fileInfo.Name, user.Username)
		postID := post.Id
		if post.RootId != "" {
			postID = post.RootId
		}

		bot.BotCreateReply(replyMsg, post.ChannelId, postID)
	}

	return nil
}

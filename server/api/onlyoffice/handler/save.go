/**
 *
 * (c) Copyright Ascensio System SIA 2022
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
package handler

import (
	"fmt"
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
	"github.com/pkg/errors"
)

var _ = Registry.RegisterHandler(2, _saveFile)
var _ = Registry.RegisterHandler(6, _saveFile)

func _saveFile(c model.Callback, a api.PluginAPI) error {
	a.API.LogDebug(_OnlyofficeLoggerPrefix + "file " + c.FileID + " save call")

	if c.URL == "" {
		return &InvalidFileDownloadUrlError{
			FileID: c.FileID,
		}
	}

	fileInfo, fileErr := a.API.GetFileInfo(c.FileID)
	if fileErr != nil {
		return &FileNotFoundError{
			FileID: c.FileID,
			Reason: fileErr.Error(),
		}
	}

	resp, err := http.Get(c.URL)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return errors.Wrap(err, _OnlyofficeLoggerPrefix)
	}

	post, postErr := a.API.GetPost(fileInfo.PostId)
	if postErr != nil {
		return &FilePersistenceError{
			FileID: c.FileID,
			Reason: postErr.Error(),
		}
	}

	post.UpdateAt = a.OnlyofficeConverter.GetTimestamp()
	_, uErr := a.API.UpdatePost(post)
	if uErr != nil {
		return &FilePersistenceError{
			FileID: c.FileID,
			Reason: uErr.Error(),
		}
	}

	_, storeErr := a.Filestore.WriteFile(resp.Body, fileInfo.Path)
	if storeErr != nil {
		return &FilePersistenceError{
			FileID: c.FileID,
			Reason: storeErr.Error(),
		}
	}

	if c.Status == 2 {
		last := c.Users[0]
		if last == "" {
			return ErrInvalidUserID
		}

		user, userErr := a.API.GetUser(last)
		if userErr != nil {
			return &FilePersistenceError{
				FileID: c.FileID,
				Reason: userErr.Error(),
			}
		}

		replyMsg := fmt.Sprintf("File %s was updated by @%s", fileInfo.Name, user.Username)
		a.Bot.BotCreateReply(replyMsg, post.ChannelId, post.Id)
	}

	return nil
}

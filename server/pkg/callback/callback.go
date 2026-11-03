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

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
)

type Handler interface {
	Handle(ctx context.Context, callback Callback) error
}

type handlerImpl struct {
	api       plugin.API
	converter converter.TimeConverter
	filestore filestore.FileBackend
	bot       bot.Bot
}

func newHandler(
	api plugin.API,
	converter converter.TimeConverter,
	filestore filestore.FileBackend,
	bot bot.Bot,
) Handler {
	return &handlerImpl{
		api:       api,
		converter: converter,
		filestore: filestore,
		bot:       bot,
	}
}

// Handle dispatches a Document Server callback to the handler registered for
// its status. Errors are propagated so that Document Server receives
// {"error": 1} and retries instead of assuming the document was saved.
// Unknown statuses are logged and acknowledged.
func (h *handlerImpl) Handle(ctx context.Context, callback Callback) error {
	err := registryContainer.RunHandler(ctx, callback.Status, callback, h.api, h.converter, h.filestore, h.bot)

	var unknown *CallbackHandlerDoesNotExistError
	if errors.As(err, &unknown) {
		h.api.LogWarn(err.Error())
		return nil
	}

	return err
}

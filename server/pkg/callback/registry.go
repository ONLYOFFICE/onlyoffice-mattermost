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
	"sync"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
)

var registryContainer = registry{
	handlers: make(map[int]func(context.Context, Callback, plugin.API, converter.TimeConverter, filestore.FileBackend, bot.Bot) error),
}

type registry struct {
	handlers map[int]func(context.Context, Callback, plugin.API, converter.TimeConverter, filestore.FileBackend, bot.Bot) error
	locker   sync.Mutex
}

func (r *registry) RegisterHandler(
	code int,
	processor func(
		ctx context.Context,
		callback Callback,
		api plugin.API,
		converter converter.TimeConverter,
		filestore filestore.FileBackend,
		bot bot.Bot,
	) error,
) error {
	r.locker.Lock()
	defer r.locker.Unlock()
	if _, exists := r.handlers[code]; exists {
		return ErrHandlerAlreadyRegistered
	}

	r.handlers[code] = processor
	return nil
}

func (r *registry) RunHandler(
	ctx context.Context,
	code int,
	callback Callback,
	api plugin.API,
	converter converter.TimeConverter,
	filestore filestore.FileBackend,
	bot bot.Bot,
) error {
	if handler, exists := r.handlers[code]; exists {
		return handler(ctx, callback, api, converter, filestore, bot)
	}

	return &CallbackHandlerDoesNotExistError{
		Code: code,
	}
}

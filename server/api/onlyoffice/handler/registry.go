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
	"sync"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	oomodel "github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
)

var Registry registry = registry{
	handlers: make(map[int]func(oomodel.Callback, api.PluginAPI) error),
}

type registry struct {
	handlers map[int]func(oomodel.Callback, api.PluginAPI) error
	locker   sync.Mutex
}

func (r *registry) RegisterHandler(code int, processor func(oomodel.Callback, api.PluginAPI) error) error {
	r.locker.Lock()
	defer r.locker.Unlock()
	if _, exists := r.handlers[code]; exists {
		return ErrHandlerAlreadyRegistered
	}
	r.handlers[code] = processor
	return nil
}

func (r *registry) RunHandler(code int, callback oomodel.Callback, api api.PluginAPI) error {
	if handler, exists := r.handlers[code]; exists {
		return handler(callback, api)
	} else {
		return &CallbackHandlerDoesNotExistError{
			Code: code,
		}
	}
}

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
package controller

import (
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/common"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type CodeHandler struct {
	api        plugin.API
	fileHelper file.FileHelper
}

func NewCodeHandler(api plugin.API, fileHelper file.FileHelper) CodeHandler {
	return CodeHandler{
		api:        api,
		fileHelper: fileHelper,
	}
}

func (h *CodeHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	code := h.fileHelper.GenerateKey()
	if err := h.api.KVSetWithExpiry(code, []byte(r.Header.Get(tools.MMAuthHeader)), 120); err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not set code: " + err.Error())
	}

	common.WriteJSON(rw, code)
}

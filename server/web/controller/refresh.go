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
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	common "github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/common"
)

type RefreshHandler struct {
	api           plugin.API
	configuration *configuration.Configuration
	fileHelper    file.FileHelper
	encoder       crypto.Encoder
	jwtManager    crypto.JwtManager
}

func NewRefreshHandler(
	api plugin.API,
	configuration *configuration.Configuration,
	fileHelper file.FileHelper,
	encoder crypto.Encoder,
	jwtManager crypto.JwtManager,
) RefreshHandler {
	return RefreshHandler{
		api:           api,
		configuration: configuration,
		fileHelper:    fileHelper,
		encoder:       encoder,
		jwtManager:    jwtManager,
	}
}

func (h *RefreshHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	h.api.LogDebug(common.LoggerPrefix + "got an editor config request")

	config, _, _, statusCode, err := common.BuildEditorConfig(r, h.configuration, h.api, h.encoder, h.jwtManager, h.fileHelper)
	if err != nil {
		h.api.LogError(common.LoggerPrefix + err.Error())
		rw.WriteHeader(statusCode)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(config); err != nil {
		h.api.LogError(common.LoggerPrefix + "could not encode config response: " + err.Error())
	}
}

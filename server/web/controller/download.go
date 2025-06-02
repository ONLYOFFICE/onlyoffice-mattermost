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
	"net/url"
	"strings"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type DownloadHandler struct {
	api           plugin.API
	configuration *configuration.Configuration
	jwtManager    crypto.JwtManager
}

func NewDownloadHandler(
	api plugin.API,
	configuration *configuration.Configuration,
	jwtManager crypto.JwtManager,
) DownloadHandler {
	return DownloadHandler{
		api:           api,
		configuration: configuration,
		jwtManager:    jwtManager,
	}
}

func (h *DownloadHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	var jwt model.DownloadTokenRequest

	token := strings.ReplaceAll(r.Header.Get(h.configuration.DESJwtHeader), "Bearer ", "")
	if token == "" {
		h.api.LogError(onlyofficeLoggerPrefix + "could not extract jwt with the header specified. Please validate your JWT Header settings")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	err := h.jwtManager.Verify([]byte(h.configuration.DESJwt), token, &jwt)
	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not verify jwt: " + err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	u, err := url.Parse(jwt.Payload.URL)
	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not parse url: " + err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	file, fileErr := h.api.GetFile(u.Query().Get("id"))
	if fileErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not download file. Reason: " + fileErr.Message)
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	h.api.LogDebug(onlyofficeLoggerPrefix + "downloading file " + u.Query().Get("id"))
	if _, err := rw.Write(file); err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "downloading file error: " + err.Error())
	}
}

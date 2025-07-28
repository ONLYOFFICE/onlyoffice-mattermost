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
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/callback"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/common"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type CallbackHandler struct {
	api             plugin.API
	configuration   *configuration.Configuration
	jwtManager      crypto.JwtManager
	callbackHandler callback.Handler
}

func NewCallbackHandler(
	api plugin.API,
	configuration *configuration.Configuration,
	jwtManager crypto.JwtManager,
	callbackHandler callback.Handler,
) CallbackHandler {
	return CallbackHandler{
		api:             api,
		configuration:   configuration,
		jwtManager:      jwtManager,
		callbackHandler: callbackHandler,
	}
}

func (h *CallbackHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	var body model.CallbackRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "callback body decoding error: " + err.Error())
		common.WriteJSON(rw, callbackErr)
		return
	}

	body.FileID = r.URL.Query().Get("file")
	validationErr := body.Validate()
	if validationErr != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "callback body validation error: " + validationErr.Error())
		common.WriteJSON(rw, callbackErr)
		return
	}

	h.api.LogDebug(onlyofficeLoggerPrefix + "valid callback payload")

	headerToken := strings.TrimSpace(strings.ReplaceAll(r.Header.Get(h.configuration.DESJwtHeader), h.configuration.DESJwtPrefix, ""))
	if body.Token == "" && headerToken == "" {
		h.api.LogError(onlyofficeLoggerPrefix + "expected to get a callback jwt. Got null")
		common.WriteJSON(rw, callbackErr, http.StatusForbidden)
		return
	}

	if body.Token != "" {
		err := h.jwtManager.Verify([]byte(h.configuration.DESJwt), body.Token, &body)
		if err != nil {
			h.api.LogError(err.Error())
			common.WriteJSON(rw, callbackErr, http.StatusForbidden)
			return
		}
	}

	if body.Token == "" && headerToken != "" {
		err := h.jwtManager.Verify([]byte(h.configuration.DESJwt), headerToken, &body)
		if err != nil {
			h.api.LogError(err.Error())
			common.WriteJSON(rw, callbackErr, http.StatusForbidden)
			return
		}

	}

	if handlerErr := h.callbackHandler.Handle(r.Context(), callback.Callback{
		Actions: body.Actions,
		Key:     body.Key,
		Status:  body.Status,
		Users:   body.Users,
		URL:     body.URL,
		FileID:  body.FileID,
		Token:   body.Token,
	}); handlerErr != nil {
		h.api.LogError(handlerErr.Error())
		common.WriteJSON(rw, callbackErr)
		return
	}

	h.api.LogDebug(onlyofficeLoggerPrefix + "callback request had no errors")
	common.WriteJSON(rw, callbackOK)
}

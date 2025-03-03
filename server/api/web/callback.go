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
package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/handler"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
)

func BuildCallbackHandler(plugin api.PluginAPI) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var body model.Callback

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "callback body decoding error: " + err.Error())
			api.WriteJSON(rw, _CallbackErr)
			return
		}

		body.FileID = r.URL.Query().Get("file")
		validationErr := body.Validate()
		if validationErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "callback body validation error: " + validationErr.Error())
			api.WriteJSON(rw, _CallbackErr)
			return
		}

		plugin.API.LogDebug(_OnlyofficeLoggerPrefix + "valid callback payload")

		headerToken := strings.TrimSpace(strings.ReplaceAll(r.Header.Get(plugin.Configuration.Header), plugin.Configuration.Prefix, ""))
		if body.Token == "" && headerToken == "" {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "expected to get a callback jwt. Got null")
			api.WriteJSON(rw, _CallbackErr, http.StatusForbidden)
			return
		}

		if body.Token != "" {
			err := plugin.Manager.Verify(body.Token, &body)
			if err != nil {
				plugin.API.LogError(err.Error())
				api.WriteJSON(rw, _CallbackErr, http.StatusForbidden)
				return
			}
		}

		if body.Token == "" && headerToken != "" {
			err := plugin.Manager.Verify(headerToken, &body)
			if err != nil {
				plugin.API.LogError(err.Error())
				api.WriteJSON(rw, _CallbackErr, http.StatusForbidden)
				return
			}
		}

		handlerErr := handler.Registry.RunHandler(body.Status, body, plugin)
		if handlerErr != nil {
			plugin.API.LogError(handlerErr.Error())
			api.WriteJSON(rw, _CallbackErr)
			return
		}

		plugin.API.LogDebug(_OnlyofficeLoggerPrefix + "callback request had no errors")
		api.WriteJSON(rw, _CallbackOK)
	}
}

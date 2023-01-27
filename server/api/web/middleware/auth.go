/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
package middleware

import (
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
)

func MattermostAuthorizationMiddleware(plugin api.PluginAPI) func(next func(plugin api.PluginAPI) func(rw http.ResponseWriter, r *http.Request)) func(rw http.ResponseWriter, r *http.Request) {
	return func(next func(plugin api.PluginAPI) func(rw http.ResponseWriter, r *http.Request)) func(rw http.ResponseWriter, r *http.Request) {
		return func(rw http.ResponseWriter, r *http.Request) {
			uid, err := plugin.API.KVGet(r.URL.Query().Get("code"))

			if r.Header.Get(plugin.Configuration.MMAuthHeader) == "" && err != nil {
				plugin.API.LogWarn("[ONLYOFFICE Mattermost Authorization]: could not find uid")
				rw.WriteHeader(http.StatusForbidden)
				return
			} else if r.Header.Get(plugin.Configuration.MMAuthHeader) == "" && err == nil {
				r.Header.Set(plugin.Configuration.MMAuthHeader, string(uid))
			}

			next(plugin)(rw, r)
		}
	}
}

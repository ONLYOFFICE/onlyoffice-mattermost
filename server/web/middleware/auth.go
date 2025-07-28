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
package middleware

import (
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type AuthorizationMiddleware struct {
	api plugin.API
}

func NewAuthorizationMiddleware(api plugin.API) AuthorizationMiddleware {
	return AuthorizationMiddleware{api: api}
}

func (h *AuthorizationMiddleware) Handle(
	next func(plugin plugin.API) func(rw http.ResponseWriter, r *http.Request),
) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(tools.MMAuthHeader)
		if userID == "" {
			code := r.URL.Query().Get("code")
			uid, err := h.api.KVGet(code)

			if err != nil || len(uid) == 0 {
				h.api.LogWarn("[ONLYOFFICE Mattermost Authorization]: could not find uid")
				rw.WriteHeader(http.StatusForbidden)
				return
			}

			userID = string(uid)
			r.Header.Set(tools.MMAuthHeader, userID)
		}

		next(h.api)(rw, r)
	}
}

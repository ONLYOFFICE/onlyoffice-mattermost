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

package main

import (
	"crypto/tls"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

func (p *Plugin) GetHTTPClient() *HTTPClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := HTTPClient{client: http.Client{Transport: tr}}

	return &client
}

func (p *Plugin) DebugRoutes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				p.API.LogError(
					"url", r.URL.String(),
					"error", x,
					"stack", string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (p *Plugin) forkRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(p.DebugRoutes)

	subrouter := router.PathPrefix("/onlyofficeapi").Subrouter()
	subrouter.HandleFunc(ONLYOFFICE_ROUTE_DOWNLOAD, p.callbackMiddleware(p.download)).Methods(http.MethodGet)
	subrouter.HandleFunc(ONLYOFFICE_ROUTE_EDITOR, p.userAccessMiddleware(p.editor)).Methods(http.MethodPost)
	subrouter.HandleFunc(ONLYOFFICE_ROUTE_CALLBACK, p.callbackMiddleware(p.callback)).Methods(http.MethodPost)
	subrouter.HandleFunc(ONLYOFFICE_ROUTE_SET_PERMISSIONS, p.permissionsMiddleware(p.setFilePermissions)).Methods(http.MethodPost)
	subrouter.HandleFunc(ONLYOFFICE_ROUTE_GET_PERMISSIONS, p.permissionsMiddleware(p.getFilePermissions)).Methods(http.MethodGet)
	subrouter.HandleFunc(ONLYOFFICE_ROUTE_GET_CHANNEL_USERS, p.userAccessMiddleware(p.channelUsers)).Methods(http.MethodGet)

	return router
}

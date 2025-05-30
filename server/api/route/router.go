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
package route

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/web"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/web/middleware"
)

func recoverRoutes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("recovering from %v with url: %s", err, r.URL.String())
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func timeoutRoutes(timeout time.Duration) func(next func(rw http.ResponseWriter, r *http.Request)) http.Handler {
	return func(next func(rw http.ResponseWriter, r *http.Request)) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			done := make(chan struct{})
			go func() {
				next(w, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-ctx.Done():
				w.WriteHeader(http.StatusGatewayTimeout)
				w.Write([]byte("request timeout"))
				return
			case <-done:
				return
			}
		})
	}
}

func NewRouter(api api.PluginAPI) *mux.Router {
	router := mux.NewRouter()
	router.Use(recoverRoutes)

	subrouter := router.PathPrefix("/api").Subrouter()
	subrouter.Handle("/callback", timeoutRoutes(5*time.Second)(web.BuildCallbackHandler(api))).Methods(http.MethodPost)
	subrouter.Handle("/download", timeoutRoutes(5*time.Second)(web.BuildDownloadHandler(api))).Methods(http.MethodGet)
	subrouter.Handle("/image", timeoutRoutes(5*time.Second)(web.BuildImageHandler(api))).Methods(http.MethodGet)

	authMiddleware := middleware.MattermostAuthorizationMiddleware(api)

	subrouter.Handle("/permissions", timeoutRoutes(2*time.Second)(authMiddleware(web.BuildSetFilePermissionsHandler))).Methods(http.MethodPost)
	subrouter.HandleFunc("/editor", authMiddleware(web.BuildEditorHandler)).Methods(http.MethodGet)
	subrouter.Handle("/permissions", timeoutRoutes(2*time.Second)(authMiddleware(web.BuildGetFilePermissionsHandler))).Methods(http.MethodGet)
	subrouter.Handle("/create", timeoutRoutes(2*time.Second)(authMiddleware(web.BuildCreateHandler))).Methods(http.MethodPost)
	subrouter.Handle("/convert", timeoutRoutes(10*time.Second)(authMiddleware(web.BuildConvertHandler))).Methods(http.MethodPost)
	subrouter.Handle("/code", timeoutRoutes(2*time.Second)(authMiddleware(web.BuildCodeHandler))).Methods(http.MethodGet)

	return router
}

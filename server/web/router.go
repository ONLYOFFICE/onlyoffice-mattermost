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
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/callback"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/middleware"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/plugin"
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

func timeoutRoutes(timeout time.Duration) func(next http.HandlerFunc) http.Handler {
	return func(next http.HandlerFunc) http.Handler {
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

func NewRouter(
	api plugin.API,
	configuration *configuration.Configuration,
	jwtManager crypto.JwtManager,
	callbackHandler callback.Handler,
	fileHelper file.FileHelper,
	encoder crypto.Encoder,
	formatManager public.FormatManager,
	commandClient client.CommandClient,
	bot bot.Bot,
) *mux.Router {
	router := mux.NewRouter()
	router.Use(recoverRoutes)
	bpath, err := api.GetBundlePath()
	var editorTemplate *template.Template
	if err != nil {
		api.LogError("[ONLYOFFICE] Failed to get bundle path: " + err.Error())
	} else {
		tmpl, parseErr := template.New("onlyoffice").ParseFiles(filepath.Join(bpath, "public/editor.html"))
		if parseErr != nil {
			api.LogError("[ONLYOFFICE] Failed to parse editor template: " + parseErr.Error())
		} else {
			editorTemplate = tmpl
		}
	}

	ch := controller.NewCallbackHandler(api, configuration, jwtManager, callbackHandler)
	dh := controller.NewDownloadHandler(api, configuration, jwtManager)
	ih := controller.NewImageHandler(api)
	eh := controller.NewEditorHandler(api, configuration, fileHelper, encoder, jwtManager, editorTemplate)
	ph := controller.NewPermissionsHandler(api, configuration, fileHelper, bot)
	crh := controller.NewCreateHandler(api, configuration)
	cvh := controller.NewConvertHandler(api, configuration, formatManager, jwtManager, commandClient)
	cdh := controller.NewCodeHandler(api, fileHelper)

	subrouter := router.PathPrefix("/api").Subrouter()
	subrouter.HandleFunc("/callback", ch.Handle).Methods(http.MethodPost)
	subrouter.HandleFunc("/download", dh.Handle).Methods(http.MethodGet)
	subrouter.HandleFunc("/image", ih.Handle).Methods(http.MethodGet)

	authMiddleware := middleware.NewAuthorizationMiddleware(api)

	subrouter.HandleFunc("/permissions", authMiddleware.Handle(func(api plugin.API) func(rw http.ResponseWriter, r *http.Request) {
		return ph.SetPermissions
	})).Methods(http.MethodPost)
	subrouter.Handle("/editor", timeoutRoutes(5*time.Second)(authMiddleware.Handle(func(api plugin.API) func(rw http.ResponseWriter, r *http.Request) {
		return eh.Handle
	}))).Methods(http.MethodGet)
	subrouter.HandleFunc("/permissions", authMiddleware.Handle(func(api plugin.API) func(rw http.ResponseWriter, r *http.Request) {
		return ph.GetPermissions
	})).Methods(http.MethodGet)
	subrouter.HandleFunc("/create", authMiddleware.Handle(func(api plugin.API) func(rw http.ResponseWriter, r *http.Request) {
		return crh.Handle
	})).Methods(http.MethodPost)
	subrouter.HandleFunc("/convert", authMiddleware.Handle(func(api plugin.API) func(rw http.ResponseWriter, r *http.Request) {
		return cvh.Handle
	})).Methods(http.MethodPost)
	subrouter.HandleFunc("/code", authMiddleware.Handle(func(api plugin.API) func(rw http.ResponseWriter, r *http.Request) {
		return cdh.Handle
	})).Methods(http.MethodGet)

	return router
}

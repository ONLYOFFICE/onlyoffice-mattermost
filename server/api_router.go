package main

import (
	"crypto/tls"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

func (p *Plugin) GetHTTPClient() *HTTPClient {
	config := p.getConfiguration()

	var client HTTPClient = HTTPClient{client: http.Client{}}

	if !config.DESEnableTLS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = HTTPClient{client: http.Client{Transport: tr}}
	}

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
	subrouter.HandleFunc(ONLYOFFICE_ROUTE_GET_CHANNEL_USER, p.userAccessMiddleware(p.channelUser)).Methods(http.MethodGet)

	return router
}

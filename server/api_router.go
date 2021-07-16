package main

import (
	"crypto/tls"
	"net/http"

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

func (p *Plugin) forkRouter() *mux.Router {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/onlyofficeapi").Subrouter()
	subrouter.HandleFunc("/download", p.docServerOnlyMiddleware(p.download)).Methods(http.MethodGet)
	subrouter.HandleFunc("/editor", p.authenticationMiddleware((p.editor))).Methods(http.MethodPost)
	subrouter.HandleFunc("/callback", p.docServerOnlyMiddleware(p.callback)).Methods(http.MethodPost)

	return router
}

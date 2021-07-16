package main

import (
	"encoders"
	"net/http"
	"sync"
	"utils"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
	router            *mux.Router
	internalKey       []byte
	encoder           encoders.Encoder
	configurationLock sync.RWMutex
	configuration     *configuration
}

func (p *Plugin) OnActivate() error {
	p.router = p.forkRouter()
	p.internalKey = []byte(utils.GenerateKey())
	return nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

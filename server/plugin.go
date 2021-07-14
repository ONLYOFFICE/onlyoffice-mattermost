package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/patrickmn/go-cache"
)

type Plugin struct {
	plugin.MattermostPlugin
	router            *mux.Router
	globalCache       *cache.Cache
	signingKey        []byte
	configurationLock sync.RWMutex
	configuration     *configuration
}

func (p *Plugin) OnActivate() error {
	p.router = p.forkRouter()
	p.signingKey = []byte(p.configuration.DESSecret)
	p.globalCache = cache.New(5*time.Minute, 5*time.Minute)
	return nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

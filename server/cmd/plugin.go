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
package cmd

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"github.com/pkg/errors"

	integration "github.com/ONLYOFFICE/onlyoffice-mattermost"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/route"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/converter"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/onlyoffice"
)

var (
	PluginID      = integration.Manifest.Id
	PluginVersion = integration.Manifest.Version
)

type Plugin struct {
	plugin.MattermostPlugin
	configurationLock       sync.RWMutex
	configuration           *configuration
	Bot                     bot.Bot
	OnlyofficeHelper        onlyoffice.Helper
	OnlyofficeConverter     converter.Converter
	Encoder                 crypto.Encoder
	Manager                 crypto.JwtManager
	EditorTemplate          *template.Template
	Filestore               filestore.FileBackend
	OnlyofficeCommandClient client.OnlyofficeCommandClient
}

func (p *Plugin) OnActivate() error {
	if p.configuration.Error != nil {
		return p.configuration.Error
	}

	bot, err := p.EnsureBot()
	if err != nil {
		return err
	}

	p.Bot = bot

	return nil
}

func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	if err := p.MattermostPlugin.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	p.setConfiguration(configuration)

	configuration.Error = configuration.IsValid()
	if configuration.Error != nil {
		time.AfterFunc(100*time.Millisecond, func() {
			if err := p.MattermostPlugin.API.DisablePlugin(PluginID); err != nil {
				p.MattermostPlugin.API.LogInfo(_OnlyofficeLoggerPrefix+"Could not disable the plugin via Mattermost API: ", err.Message)
			}
		})

		return nil
	}

	p.Encoder = crypto.NewMD5Encoder()
	p.Manager = crypto.NewJwtManager([]byte(p.configuration.DESJwt))
	p.OnlyofficeHelper = onlyoffice.NewHelper()
	p.OnlyofficeConverter = converter.NewConverter()
	p.OnlyofficeCommandClient = client.NewOnlyofficeCommandClient(p.Manager)
	bpath, _ := p.MattermostPlugin.API.GetBundlePath()
	p.EditorTemplate, configuration.Error = template.New("onlyoffice").ParseFiles(filepath.Join(bpath, "public/editor.html"))
	if configuration.Error != nil {
		time.AfterFunc(100*time.Millisecond, func() {
			if err := p.MattermostPlugin.API.DisablePlugin(PluginID); err != nil {
				p.MattermostPlugin.API.LogInfo(_OnlyofficeLoggerPrefix + "Could not disable the plugin via Mattermost API: " + err.Message)
			}
		})
		return nil
	}

	license := p.MattermostPlugin.API.GetLicense()
	serverConfig := p.MattermostPlugin.API.GetUnsanitizedConfig()
	serverConfig.FileSettings.SetDefaults(true)
	p.Filestore, configuration.Error = filestore.NewFileBackend(filestore.NewFileBackendSettingsFromConfig(&serverConfig.FileSettings, (license != nil && *license.Features.Compliance), true))
	if configuration.Error != nil {
		time.AfterFunc(100*time.Millisecond, func() {
			if err := p.MattermostPlugin.API.DisablePlugin(PluginID); err != nil {
				p.MattermostPlugin.API.LogInfo(_OnlyofficeLoggerPrefix + "Could not disable the plugin via Mattermost API: " + err.Message)
			}
		})
		return nil
	}

	p.MattermostPlugin.API.LogInfo(_OnlyofficeLoggerPrefix + "The server responded without errors")
	return nil
}

func (p *Plugin) EnsureBot() (bot.Bot, error) {
	botID, err := p.MattermostPlugin.API.EnsureBotUser(&model.Bot{
		Username:    "onlyoffice",
		DisplayName: "ONLYOFFICE",
		Description: "ONLYOFFICE Helper",
	})
	if err != nil {
		return nil, ErrCreateBotProfile
	}

	bundlePath, err := p.MattermostPlugin.API.GetBundlePath()
	if err != nil {
		return nil, err
	}

	profileImage, err := os.ReadFile(filepath.Join(bundlePath, "assets", "logo.png"))
	if err != nil {
		return nil, ErrLoadBotProfileImage
	}

	if appErr := p.MattermostPlugin.API.SetProfileImage(botID, profileImage); appErr != nil {
		return nil, ErrSetBotProfileImage
	}

	return bot.NewBot(botID, p.MattermostPlugin.API), nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	route.NewRouter(api.PluginAPI{
		API: p.MattermostPlugin.API,
		Configuration: struct {
			Address      string
			Secret       string
			Header       string
			Prefix       string
			MMAuthHeader string
		}{
			Address:      p.configuration.DESAddress,
			Secret:       p.configuration.DESJwt,
			Header:       p.configuration.DESJwtHeader,
			Prefix:       p.configuration.DESJwtPrefix,
			MMAuthHeader: "Mattermost-User-Id",
		},
		OnlyofficeHelper:        p.OnlyofficeHelper,
		OnlyofficeConverter:     p.OnlyofficeConverter,
		Encoder:                 p.Encoder,
		Manager:                 p.Manager,
		Bot:                     p.Bot,
		EditorTemplate:          p.EditorTemplate,
		Filestore:               p.Filestore,
		OnlyofficeCommandClient: p.OnlyofficeCommandClient,
	}).ServeHTTP(w, r)
}

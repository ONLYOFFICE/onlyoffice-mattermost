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
	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
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
	FormatManager           public.FormatManager
}

func (p *Plugin) OnActivate() error {
	if p.configuration == nil {
		return errors.New("plugin configuration is not initialized")
	}

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

	configuration.sanitizeConfiguration()
	configuration.handleDemoConfiguration(p.API)
	if err := configuration.IsValid(); err != nil {
		configuration.Error = err
		p.MattermostPlugin.API.LogError(_OnlyofficeLoggerPrefix + "Configuration validation failed: " + err.Error())
		time.AfterFunc(100*time.Millisecond, func() {
			if disableErr := p.MattermostPlugin.API.DisablePlugin(PluginID); disableErr != nil {
				p.MattermostPlugin.API.LogError(_OnlyofficeLoggerPrefix + "Could not disable the plugin via Mattermost API: " + disableErr.Message)
			}
		})
		p.setConfiguration(configuration)
		return err
	}

	p.setConfiguration(configuration)

	p.Encoder = crypto.NewMD5Encoder()
	p.Manager = crypto.NewJwtManager([]byte(p.configuration.DESJwt))

	var err error
	p.FormatManager, err = public.NewMapFormatManager()
	if err != nil {
		p.MattermostPlugin.API.LogError(_OnlyofficeLoggerPrefix + "Failed to initialize format manager: " + err.Error())
		time.AfterFunc(100*time.Millisecond, func() {
			if disableErr := p.MattermostPlugin.API.DisablePlugin(PluginID); disableErr != nil {
				p.MattermostPlugin.API.LogError(_OnlyofficeLoggerPrefix + "Could not disable the plugin via Mattermost API: " + disableErr.Message)
			}
		})
		return err
	}

	p.OnlyofficeHelper = onlyoffice.NewHelper(p.FormatManager)
	p.OnlyofficeConverter = converter.NewConverter()
	p.OnlyofficeCommandClient = client.NewOnlyofficeCommandClient(p.Manager)

	bpath, _ := p.MattermostPlugin.API.GetBundlePath()
	p.EditorTemplate, err = template.New("onlyoffice").ParseFiles(filepath.Join(bpath, "public/editor.html"))
	if err != nil {
		p.MattermostPlugin.API.LogError(_OnlyofficeLoggerPrefix + "Failed to parse editor template: " + err.Error())
		time.AfterFunc(100*time.Millisecond, func() {
			if disableErr := p.MattermostPlugin.API.DisablePlugin(PluginID); disableErr != nil {
				p.MattermostPlugin.API.LogError(_OnlyofficeLoggerPrefix + "Could not disable the plugin via Mattermost API: " + disableErr.Message)
			}
		})
		return err
	}

	license := p.MattermostPlugin.API.GetLicense()
	serverConfig := p.MattermostPlugin.API.GetUnsanitizedConfig()
	serverConfig.FileSettings.SetDefaults(true)
	p.Filestore, err = filestore.NewFileBackend(filestore.NewFileBackendSettingsFromConfig(&serverConfig.FileSettings, (license != nil && *license.Features.Compliance), true))
	if err != nil {
		p.MattermostPlugin.API.LogError(_OnlyofficeLoggerPrefix + "Failed to initialize file backend: " + err.Error())
		time.AfterFunc(100*time.Millisecond, func() {
			if disableErr := p.MattermostPlugin.API.DisablePlugin(PluginID); disableErr != nil {
				p.MattermostPlugin.API.LogError(_OnlyofficeLoggerPrefix + "Could not disable the plugin via Mattermost API: " + disableErr.Message)
			}
		})
		return err
	}

	p.MattermostPlugin.API.LogInfo(_OnlyofficeLoggerPrefix + "Configuration updated successfully")
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
			DemoEnabled  bool
			DemoExpires  int64
			DemoAddress  string
			DemoHeader   string
			DemoPrefix   string
			DemoSecret   string
			MMAuthHeader string
		}{
			Address:      p.configuration.DESAddress,
			Secret:       p.configuration.DESJwt,
			Header:       p.configuration.DESJwtHeader,
			Prefix:       p.configuration.DESJwtPrefix,
			DemoEnabled:  p.configuration.DemoEnabled,
			DemoExpires:  p.configuration.DemoExpires,
			DemoAddress:  p.configuration.DemoAddress,
			DemoHeader:   p.configuration.DemoHeader,
			DemoPrefix:   p.configuration.DemoPrefix,
			DemoSecret:   p.configuration.DemoSecret,
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
		FormatManager:           p.FormatManager,
	}).ServeHTTP(w, r)
}

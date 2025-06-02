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
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"github.com/pkg/errors"
	"go.uber.org/fx"

	integration "github.com/ONLYOFFICE/onlyoffice-mattermost"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/callback"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/common"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/middleware"
)

var (
	PluginID      = integration.Manifest.Id
	PluginVersion = integration.Manifest.Version
)

type Plugin struct {
	plugin.MattermostPlugin
	app               *fx.App
	configuration     *configuration.Configuration
	configurationLock sync.RWMutex

	router        *mux.Router
	commandClient client.CommandClient
	jwtManager    crypto.JwtManager

	ready bool
}

func (p *Plugin) handleConfigError(config *configuration.Configuration, err error, message string) error {
	config.Error = err
	logMessage := common.OnlyofficeLoggerCmdPrefix + message
	if err != nil {
		logMessage += ": " + err.Error()
	}

	p.logError(logMessage)
	p.setConfiguration(config)

	time.AfterFunc(100*time.Millisecond, func() {
		if p.MattermostPlugin.API != nil {
			if disableErr := p.MattermostPlugin.API.DisablePlugin(PluginID); disableErr != nil {
				p.MattermostPlugin.API.LogError(
					common.OnlyofficeLoggerCmdPrefix + "Could not disable the plugin via Mattermost API: " + disableErr.Message,
				)
			}
		}
	})

	return err
}

func (p *Plugin) logError(message string) {
	if p.MattermostPlugin.API != nil {
		p.MattermostPlugin.API.LogError(message)
	} else {
		fmt.Println(message)
	}
}

func (p *Plugin) provideFormatManager() public.FormatManager {
	formatManager, err := public.NewMapFormatManager()
	if err != nil {
		panic(err)
	}

	return formatManager
}

func (p *Plugin) provideFileBackend() filestore.FileBackend {
	license := p.MattermostPlugin.API.GetLicense()
	serverConfig := p.MattermostPlugin.API.GetUnsanitizedConfig()
	serverConfig.FileSettings.SetDefaults(true)
	fs, err := filestore.NewFileBackend(
		filestore.NewFileBackendSettingsFromConfig(
			&serverConfig.FileSettings,
			license != nil && *license.Features.Compliance,
			true,
		),
	)

	if err != nil {
		p.logError(common.OnlyofficeLoggerCmdPrefix + "Failed to parse editor template: " + err.Error())
		return nil
	}

	return fs
}

func (p *Plugin) initializeContainer() *fx.App {
	return fx.New(
		fx.NopLogger,
		fx.Supply(p),
		fx.Provide(
			func() plugin.API { return p.MattermostPlugin.API },
			func() *configuration.Configuration { return p.configuration },
			func() middleware.AuthorizationMiddleware {
				return middleware.NewAuthorizationMiddleware(p.MattermostPlugin.API)
			},
			p.provideFormatManager,
			p.provideFileBackend,
			controller.NewCreateHandler,
			controller.NewConvertHandler,
			controller.NewEditorHandler,
			controller.NewCallbackHandler,
			controller.NewPermissionsHandler,
			controller.NewDownloadHandler,
			controller.NewCodeHandler,
			controller.NewNotFoundHandler,
			fx.Annotate(
				func() string {
					botID, err := p.EnsureBot()
					if err != nil {
						panic(err)
					}
					return botID
				},
				fx.ResultTags(`name:"bot_id"`),
			),
		),
		file.Module,
		callback.Module,
		crypto.Module,
		client.Module,
		converter.Module,
		bot.Module,
		web.Module,
		fx.Invoke(func(router *mux.Router) { p.router = router }),
		fx.Invoke(func(commandClient client.CommandClient) { p.commandClient = commandClient }),
		fx.Invoke(func(jwtManager crypto.JwtManager) { p.jwtManager = jwtManager }),
	)
}

func (p *Plugin) OnActivate() error {
	if p.configuration == nil {
		return errors.New("plugin configuration is not initialized")
	}

	if p.configuration.Error != nil {
		return p.configuration.Error
	}

	return nil
}

func (p *Plugin) OnDeactivate() error {
	if p.app != nil {
		return p.app.Stop(context.Background())
	}

	return nil
}

func (p *Plugin) OnConfigurationChange() error {
	defer func() {
		if r := recover(); r != nil {
			message := common.OnlyofficeLoggerCmdPrefix + fmt.Sprintf("Panic in OnConfigurationChange: %v", r)
			p.logError(message)
			time.AfterFunc(100*time.Millisecond, func() {
				if p.MattermostPlugin.API != nil {
					if disableErr := p.MattermostPlugin.API.DisablePlugin(PluginID); disableErr != nil {
						p.MattermostPlugin.API.LogError(
							common.OnlyofficeLoggerCmdPrefix +
								"Could not disable the plugin via Mattermost API: " +
								disableErr.Message,
						)
					}
				}
			})
		}
	}()

	configuration, err := p.prepareConfiguration()
	if err != nil {
		return err
	}

	p.setConfiguration(configuration)
	if err := p.reinitializeContainer(configuration); err != nil {
		return nil
	}

	return p.validateConfiguration()
}

func (p *Plugin) reinitializeContainer(config *configuration.Configuration) error {
	if p.app != nil {
		tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := p.app.Stop(tctx); err != nil {
			p.logError(common.OnlyofficeLoggerCmdPrefix + "Failed to stop existing fx container: " + err.Error())
			return err
		}
		p.app = nil
		p.commandClient = nil
		p.jwtManager = nil
		p.ready = false
	}

	p.app = p.initializeContainer()
	if p.app == nil {
		err := fmt.Errorf("failed to initialize fx container")
		p.handleConfigError(config, err, "Failed to initialize plugin dependencies")
		return err
	}

	tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := p.app.Start(tctx); err != nil {
		p.handleConfigError(config, err, "Failed to start plugin dependencies")
		return err
	}

	p.ready = true
	return nil
}

func (p *Plugin) prepareConfiguration() (*configuration.Configuration, error) {
	config := new(configuration.Configuration)
	if err := p.MattermostPlugin.API.LoadPluginConfiguration(config); err != nil {
		return nil, errors.Wrap(err, "failed to load plugin configuration")
	}

	config.SanitizeConfiguration()
	config.HandleDemoConfiguration(p.MattermostPlugin.API)
	return config, nil
}

func (p *Plugin) validateConfiguration() error {
	config, err := p.prepareConfiguration()
	if err != nil {
		return err
	}

	if err := config.IsValid(); err != nil {
		return p.handleConfigError(config, err, "Configuration validation failed")
	}

	if err := p.validateDependencies(config); err != nil {
		return err
	}

	if err := p.validateDocumentServer(config); err != nil {
		return err
	}

	p.setConfiguration(config)
	p.MattermostPlugin.API.LogInfo(common.OnlyofficeLoggerCmdPrefix + "Configuration updated successfully")
	return nil
}

func (p *Plugin) validateDependencies(config *configuration.Configuration) error {
	if p.commandClient == nil {
		err := fmt.Errorf("command client is nil after container initialization")
		return p.handleConfigError(config, err, "Plugin dependency initialization failed")
	}

	if p.jwtManager == nil {
		err := fmt.Errorf("jwt manager is nil after container initialization")
		return p.handleConfigError(config, err, "Plugin dependency initialization failed")
	}

	return nil
}

func (p *Plugin) validateDocumentServer(config *configuration.Configuration) error {
	token, err := p.createVersionToken(config)
	if err != nil {
		return p.handleConfigError(config, err, "Could not sign a JWT")
	}

	resp, err := p.sendVersionRequest(config, token)
	if err != nil {
		return p.handleConfigError(config, err, "Could not sign a JWT")
	}

	if resp.Error != 0 {
		err := &common.DocumentServerCommandResponseError{Code: resp.Error}
		return p.handleConfigError(
			config,
			err,
			"Could not retrieve document server version, please check your credentials and make sure that document server version is 8.2 or higher",
		)
	}

	return p.validateServerVersion(config, resp.Version)
}

func (p *Plugin) createVersionToken(config *configuration.Configuration) (string, error) {
	vreq := client.VersionRequest{Command: "version"}
	vreq.IssuedAt = jwt.NewNumericDate(time.Now())
	vreq.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Minute))
	return p.jwtManager.Sign([]byte(config.DESJwt), vreq)
}

func (p *Plugin) sendVersionRequest(config *configuration.Configuration, token string) (client.VersionResponse, error) {
	resp, err := p.commandClient.SendVersion(
		config.DESAddress+client.OnlyofficeCommandServicePath+"?shardkey="+uuid.New().String(),
		client.VersionRequest{
			Command: "version",
			Token:   token,
		},
		4*time.Second,
	)

	return resp, err
}

func (p *Plugin) validateServerVersion(config *configuration.Configuration, versionStr string) error {
	if versionStr == "" || len(versionStr) == 0 {
		p.logError(common.OnlyofficeLoggerCmdPrefix + "Received empty version from document server")
		return p.handleConfigError(
			config,
			common.ErrParseDocumentServerVersion,
			"Could not parse document server version, please check your credentials and make sure that document server version is 8.2 or higher",
		)
	}

	version, err := strconv.ParseInt(versionStr[0:1], 10, 64)
	if err != nil {
		return p.handleConfigError(
			config,
			common.ErrParseDocumentServerVersion,
			"Could not parse document server version, please check your credentials and make sure that document server version is 8.2 or higher",
		)
	}

	if version < 8 {
		return p.handleConfigError(
			config,
			common.ErrDeprecatedDocumentServerVersion,
			"Document server version is deprecated, please update your document server to version 8.2 or higher",
		)
	}

	return nil
}

func (p *Plugin) EnsureBot() (string, error) {
	bot := &model.Bot{
		Username:    "onlyoffice",
		DisplayName: "ONLYOFFICE",
		Description: "ONLYOFFICE Helper",
	}

	botID, err := p.MattermostPlugin.API.EnsureBotUser(bot)
	if err != nil {
		return "", common.ErrCreateBotProfile
	}

	bundlePath, err := p.MattermostPlugin.API.GetBundlePath()
	if err != nil {
		return "", err
	}

	profileImage, err := os.ReadFile(filepath.Join(bundlePath, "assets", "logo.png"))
	if err != nil {
		return "", common.ErrLoadBotProfileImage
	}

	if appErr := p.MattermostPlugin.API.SetProfileImage(botID, profileImage); appErr != nil {
		return "", common.ErrSetBotProfileImage
	}

	return botID, nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if p.router == nil {
		p.MattermostPlugin.API.LogError("Router not initialized")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p.router.ServeHTTP(w, r)
}

func (p *Plugin) getConfiguration() *configuration.Configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration.Configuration{}
	}

	return p.configuration
}

func (p *Plugin) setConfiguration(configuration *configuration.Configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		// Ignore assignment if the configuration struct is empty. Go will optimize the
		// allocation for same to point at the same memory address, breaking the check
		// above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing configuration")
	}

	configuration.SanitizeConfiguration()
	configuration.HandleDemoConfiguration(p.MattermostPlugin.API)
	p.configuration = configuration
}

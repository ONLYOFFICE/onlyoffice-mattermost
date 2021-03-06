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
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/security"
	"github.com/golang-jwt/jwt"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/models"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// configuration captures the plugin's external configuration as exposed in the Mattermost server
// configuration, as well as values computed from the configuration. Any public fields will be
// deserialized from the Mattermost server configuration in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and the plugin
// configuration can change at any time, access to the configuration must be synchronized. The
// strategy used in this plugin is to guard a pointer to the configuration, and clone the entire
// struct whenever it changes. You may replace this with whatever strategy you choose.
//
// If you add non-reference types to your configuration struct, be sure to rewrite Clone as a deep
// copy appropriate for your types.
type configuration struct {
	DESAddress   string
	DESJwt       string
	DESJwtHeader string
	DESJwtPrefix string
	TLS          bool
}

// Clone shallow coies the configuration. Your implementation may require a deep copy if
// your configuration has reference types.
func (c *configuration) Clone() *configuration {
	return &configuration{
		DESAddress:   c.DESAddress,
		DESJwt:       c.DESJwt,
		DESJwtHeader: c.DESJwtHeader,
		DESJwtPrefix: c.DESJwtPrefix,
		TLS:          c.TLS,
	}
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// reentrant. In particular, avoid using the plugin API entirely, as this may in turn trigger a
// hook back into the plugin. If that hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing configuration. This almost
// certainly means that the configuration was modified without being cloned and may result in
// an unsafe access.
func (p *Plugin) setConfiguration(configuration *configuration) {
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

	err := configuration.SanitizeConfiguration()

	if err != nil {
		p.API.LogError(err.Error())
		p.API.DisablePlugin(manifest.Id)
		return
	}

	p.configuration = configuration
}

// Sanitize config to get just the right format
func (c *configuration) SanitizeConfiguration() error {
	DESAddress := c.DESAddress
	if DESAddress[len(DESAddress)-1:] != "/" {
		c.DESAddress = fmt.Sprintf("%s%s", DESAddress, "/")
	}
	if c.DESJwtHeader != "" {
		if strings.TrimSpace(c.DESJwtHeader) == "Authorization" {
			return errors.New("[ONLYOFFICE]: Do not use 'Authorization' header")
		}
		c.DESJwtHeader = strings.TrimSpace(c.DESJwtHeader)
	}
	if c.DESJwtPrefix != "" {
		c.DESJwtPrefix = strings.TrimSpace(c.DESJwtPrefix)
	}
	return nil
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		p.API.DisablePlugin(manifest.Id)
		return errors.New("[ONLYOFFICE]: Failed to load ONLYOFFICE configuration")
	}

	p.setConfiguration(configuration)

	// Trying to connect to the Document Service
	var body = models.CommandBody{
		Command: models.ONLYOFFICE_COMMAND_VERSION,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(5 * time.Second).Unix(),
			Issuer:    "command",
		},
	}
	body.Token, _ = security.JwtSign(body, []byte(p.configuration.DESJwt))
	body.StandardClaims = jwt.StandardClaims{}
	var headers []Header = []Header{
		{
			Key:   configuration.DESJwtHeader,
			Value: configuration.DESJwtPrefix + " " + body.Token,
		},
	}

	var response = new(models.CommandResponse)

	p.GetHTTPClient().PostRequest(configuration.DESAddress+ONLYOFFICE_COMMAND_SERVICE, &body,
		headers, response)

	var err = response.ProcessResponse()

	if err != nil {
		p.API.DisablePlugin(manifest.Id)
		return err
	}

	bot_id, creationErr := p.Helpers.EnsureBot(&model.Bot{
		Username:    "onlyoffice",
		DisplayName: "ONLYOFFICE",
		Description: "ONLYOFFICE Helper",
	}, plugin.ProfileImagePath(filepath.Join("assets", "logo.png")))
	if creationErr != nil {
		p.API.DisablePlugin(manifest.Id)
		return errors.New("Failed to create an ONLYOFFICE bot")
	}

	p.onlyoffice_bot = ONLYOFFICE_BOT{
		Id:           bot_id,
		LoggerPrefix: "[ONLYOFFICE BOT]: ",
		P:            p,
	}

	p.API.LogInfo("[ONLYOFFICE]: The server responded without errors")

	return nil
}

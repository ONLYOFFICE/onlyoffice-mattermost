/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/client/model"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/validator"
	"github.com/golang-jwt/jwt"
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
	DESAddress     string
	DESJwt         string
	DESJwtHeader   string
	DESJwtPrefix   string
	InsecureClient bool
	Error          error
}

// Clone shallow copies the configuration. Your implementation may require a deep copy if
// your configuration has reference types.
func (c *configuration) Clone() *configuration {
	return &configuration{
		DESAddress:     c.DESAddress,
		DESJwt:         c.DESJwt,
		DESJwtHeader:   c.DESJwtHeader,
		DESJwtPrefix:   c.DESJwtPrefix,
		InsecureClient: c.InsecureClient,
	}
}

func (c *configuration) sanitizeConfiguration() {
	c.DESAddress = strings.TrimSpace(c.DESAddress)
	c.DESJwt = strings.TrimSpace(c.DESJwt)
	c.DESJwtHeader = strings.TrimSpace(c.DESJwtHeader)
	c.DESJwtPrefix = strings.TrimSpace(c.DESJwtPrefix)

	if strings.HasSuffix(c.DESAddress, "/") {
		for {
			c.DESAddress = strings.TrimSuffix(c.DESAddress, "/")
			if !strings.HasSuffix(c.DESAddress, "/") {
				break
			}
		}
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

	configuration.sanitizeConfiguration()

	p.configuration = configuration
}

func (c *configuration) IsValid() error {
	if c.DESAddress == "" {
		return &InvalidDocumentServerAddressError{
			Reason: "Document server address is empty",
		}
	}

	if !validator.IsValidUrl(c.DESAddress) {
		return &InvalidDocumentServerAddressError{
			Reason: "Document server address must match the following pattern: http(s)://<domain>.<domain_zone> or http(s)://<domain>.<domain_zone>/",
		}
	}

	if c.DESJwt == "" {
		return &BadConfigurationError{
			Property: "Document Server Secret",
			Reason:   "Please specify document server secret",
		}
	}

	if c.DESJwtPrefix == "" {
		return &BadConfigurationError{
			Property: "Document Serve Prefix",
			Reason:   "Please specify document server prefix",
		}
	}

	if c.DESJwtHeader == "" || strings.TrimSpace(c.DESJwtHeader) == "Authorization" {
		return &BadConfigurationError{
			Property: "Document Server Header",
			Reason:   "Please specify document server header (Note: do not use 'Authorization' header)",
		}
	}

	command := client.NewOnlyofficeCommandClient(c.InsecureClient, crypto.NewJwtManager([]byte(c.DESJwt)))
	resp, err := command.SendVersion(c.DESAddress+client.OnlyofficeCommandServicePath, model.CommandVersionRequest{
		Command: client.OnlyofficeCommandServiceVersion,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
		},
	}, 4*time.Second)

	if err != nil {
		return err
	}

	if resp.Error != 0 {
		return &DocumentServerCommandResponseError{
			Code: resp.Error,
		}
	}

	if resp.Version == "" {
		return ErrParseDocumentServerVersion
	}

	version, err := strconv.ParseInt(resp.Version[0:1], 10, 64)

	if err != nil {
		return ErrParseDocumentServerVersion
	}

	if version < 7 {
		return ErrDeprecatedDocumentServerVersion
	}

	return nil
}

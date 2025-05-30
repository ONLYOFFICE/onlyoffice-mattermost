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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/client/model"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/validator"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mattermost/mattermost/server/public/plugin"
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
	DemoEnabled  bool
	DemoExpires  int64
	DemoAddress  string
	DemoHeader   string
	DemoPrefix   string
	DemoSecret   string
	Error        error
}

// Clone shallow copies the configuration. Your implementation may require a deep copy if
// your configuration has reference types.
func (c *configuration) Clone() *configuration {
	return &configuration{
		DESAddress:   c.DESAddress,
		DESJwt:       c.DESJwt,
		DESJwtHeader: c.DESJwtHeader,
		DESJwtPrefix: c.DESJwtPrefix,
		DemoEnabled:  c.DemoEnabled,
		DemoExpires:  c.DemoExpires,
		DemoAddress:  c.DemoAddress,
		DemoHeader:   c.DemoHeader,
		DemoPrefix:   c.DemoPrefix,
		DemoSecret:   c.DemoSecret,
	}
}

func (c *configuration) sanitizeConfiguration() {
	c.DESAddress = strings.TrimSpace(c.DESAddress)
	c.DESJwt = strings.TrimSpace(c.DESJwt)
	c.DESJwtHeader = strings.TrimSpace(c.DESJwtHeader)
	c.DESJwtPrefix = strings.TrimSpace(c.DESJwtPrefix)

	c.DemoAddress = "https://onlinedocs.docs.onlyoffice.com"
	c.DemoHeader = "AuthorizationJWT"
	c.DemoPrefix = "Bearer "
	c.DemoSecret = "sn2puSUF7muF5Jas"

	if !c.DemoEnabled || c.DemoExpires <= time.Now().UnixMilli() {
		if c.DESAddress == c.DemoAddress {
			c.DESAddress = ""
		}
		if c.DESJwt == c.DemoSecret {
			c.DESJwt = ""
		}
		if c.DESJwtHeader == c.DemoHeader {
			c.DESJwtHeader = ""
		}
		if c.DESJwtPrefix == c.DemoPrefix {
			c.DESJwtPrefix = ""
		}
	}

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
//
//nolint:unused
func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

func (c *configuration) handleDemoConfiguration(api plugin.API) {
	if !c.DemoEnabled {
		return
	}

	hasUserCredentials := c.DESAddress != "" &&
		c.DESJwt != "" &&
		c.DESJwtHeader != "" &&
		c.DESJwtPrefix != "" &&
		(c.DESAddress != c.DemoAddress ||
			c.DESJwt != c.DemoSecret ||
			c.DESJwtHeader != c.DemoHeader ||
			c.DESJwtPrefix != c.DemoPrefix)

	if hasUserCredentials {
		return
	}

	now := time.Now().UnixMilli()

	start, kvErr := api.KVGet(DemoKey)
	if kvErr != nil || len(start) == 0 {
		if err := api.KVSet(DemoKey, []byte(strconv.FormatInt(now, 10))); err == nil {
			c.DemoExpires = now + _DemoPeriodMillis
			c.DESAddress = c.DemoAddress
			c.DESJwt = c.DemoSecret
			c.DESJwtHeader = c.DemoHeader
			c.DESJwtPrefix = c.DemoPrefix
		}
		return
	}

	startTime, parseErr := strconv.ParseInt(string(start), 10, 64)
	if parseErr != nil {
		if err := api.KVSet(DemoKey, []byte(strconv.FormatInt(now, 10))); err == nil {
			c.DemoExpires = now + _DemoPeriodMillis
			c.DESAddress = c.DemoAddress
			c.DESJwt = c.DemoSecret
			c.DESJwtHeader = c.DemoHeader
			c.DESJwtPrefix = c.DemoPrefix
		}
		return
	}

	expirationTime := startTime + _DemoPeriodMillis
	if now > expirationTime {
		c.DemoEnabled = false
		c.DemoExpires = expirationTime
		if c.DESAddress == c.DemoAddress {
			c.DESAddress = ""
		}
		if c.DESJwt == c.DemoSecret {
			c.DESJwt = ""
		}
		if c.DESJwtHeader == c.DemoHeader {
			c.DESJwtHeader = ""
		}
		if c.DESJwtPrefix == c.DemoPrefix {
			c.DESJwtPrefix = ""
		}
	} else {
		c.DemoExpires = expirationTime
		if !hasUserCredentials {
			c.DESAddress = c.DemoAddress
			c.DESJwt = c.DemoSecret
			c.DESJwtHeader = c.DemoHeader
			c.DESJwtPrefix = c.DemoPrefix
		}
	}
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
	configuration.handleDemoConfiguration(p.API)
	p.configuration = configuration
}

func (c *configuration) IsValid() error {
	// Check if demo is active
	demoActive := c.DemoEnabled && c.DemoExpires > time.Now().UnixMilli()

	// Check if we have valid credentials
	hasCredentials := c.DESAddress != "" && c.DESJwt != "" && c.DESJwtHeader != "" && c.DESJwtPrefix != ""

	// If no demo and no credentials, fail
	if !demoActive && !hasCredentials {
		return &BadConfigurationError{
			Property: "Document Server Configuration",
			Reason:   "No valid credentials provided and demo mode is not active",
		}
	}

	// If demo is active, allow it
	if demoActive {
		return nil
	}

	// Validate credentials
	if c.DESAddress == "" {
		return &InvalidDocumentServerAddressError{
			Reason: "Document server address is empty",
		}
	}

	if !validator.IsValidURL(c.DESAddress) {
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
			Property: "Document Server Prefix",
			Reason:   "Please specify document server prefix",
		}
	}

	if c.DESJwtHeader == "" || strings.TrimSpace(c.DESJwtHeader) == "Authorization" {
		return &BadConfigurationError{
			Property: "Document Server Header",
			Reason:   "Please specify document server header (Note: do not use 'Authorization' header)",
		}
	}

	command := client.NewOnlyofficeCommandClient(crypto.NewJwtManager([]byte(c.DESJwt)))
	resp, err := command.SendVersion(c.DESAddress+client.OnlyofficeCommandServicePath+"?shardkey="+uuid.New().String(), model.CommandVersionRequest{
		Command: client.OnlyofficeCommandServiceVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
		},
	}, 4*time.Second)

	if err != nil {
		return &BadConfigurationError{
			Property: "Document Server Connection",
			Reason:   "Could not retrieve document server version, please check your credentials and make sure that document server version is 8.2 or higher: " + err.Error(),
		}
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

	if version < 8 {
		return ErrDeprecatedDocumentServerVersion
	}

	return nil
}

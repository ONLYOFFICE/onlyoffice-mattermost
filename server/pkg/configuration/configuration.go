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
package configuration

import (
	"strconv"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/common"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// Configuration captures the plugin's external configuration as exposed in the Mattermost server
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
type Configuration struct {
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
func (c *Configuration) Clone() *Configuration {
	return &Configuration{
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

func (c *Configuration) SanitizeConfiguration() {
	c.DESAddress = strings.TrimSpace(c.DESAddress)
	c.DESJwt = strings.TrimSpace(c.DESJwt)
	c.DESJwtHeader = strings.TrimSpace(c.DESJwtHeader)
	c.DESJwtPrefix = strings.TrimSpace(c.DESJwtPrefix)

	c.DemoAddress = "https://onlinedocs.docs.onlyoffice.com"
	c.DemoHeader = "AuthorizationJWT"
	c.DemoPrefix = "Bearer "
	c.DemoSecret = "sn2puSUF7muF5Jas"

	if !c.DemoEnabled || c.DemoExpires <= time.Now().UnixMilli() {
		if strings.EqualFold(c.DESAddress, c.DemoAddress) {
			c.DESAddress = ""
			c.DESJwt = ""
			c.DESJwtHeader = ""
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

func (c *Configuration) HandleDemoConfiguration(api plugin.API) {
	if !c.DemoEnabled {
		return
	}

	now := time.Now().UnixMilli()
	start, kvErr := api.KVGet(DemoKey)
	if kvErr != nil || len(start) == 0 {
		if err := api.KVSet(DemoKey, []byte(strconv.FormatInt(now, 10))); err == nil {
			c.DemoExpires = now + _DemoPeriodMillis
		}

		start = []byte(strconv.FormatInt(now, 10))
	}

	startTime, parseErr := strconv.ParseInt(string(start), 10, 64)
	if parseErr != nil {
		if err := api.KVSet(DemoKey, []byte(strconv.FormatInt(now, 10))); err == nil {
			c.DemoExpires = now + _DemoPeriodMillis
		}

		startTime = now
	}

	expirationTime := startTime + _DemoPeriodMillis
	c.DemoExpires = expirationTime
	if now <= expirationTime {
		c.DESAddress = c.DemoAddress
		c.DESJwt = c.DemoSecret
		c.DESJwtHeader = c.DemoHeader
		c.DESJwtPrefix = c.DemoPrefix
	} else {
		if strings.EqualFold(c.DESAddress, c.DemoAddress) {
			c.DESAddress = ""
			c.DESJwt = ""
			c.DESJwtHeader = ""
			c.DESJwtPrefix = ""
		}
	}
}

func (c *Configuration) IsValid() error {
	demoActive := c.DemoEnabled && c.DemoExpires > time.Now().UnixMilli()
	hasCredentials := c.DESAddress != "" && c.DESJwt != "" && c.DESJwtHeader != "" && c.DESJwtPrefix != ""

	if demoActive {
		return nil
	}

	if c.DemoEnabled {
		if !hasCredentials {
			return &common.BadConfigurationError{
				Property: "Document Server Configuration",
				Reason:   "Demo mode has expired and no valid credentials provided",
			}
		}
	}

	if !hasCredentials {
		return &common.BadConfigurationError{
			Property: "Document Server Configuration",
			Reason:   "No valid credentials provided and demo mode is not active",
		}
	}

	if !tools.IsValidURL(c.DESAddress) {
		return &common.InvalidDocumentServerAddressError{
			Reason: "Document server address must match the following pattern: http(s)://<domain>.<domain_zone> or http(s)://<domain>.<domain_zone>/",
		}
	}

	if c.DESJwt == "" {
		return &common.BadConfigurationError{
			Property: "Document Server Secret",
			Reason:   "Please specify document server secret",
		}
	}

	if c.DESJwtPrefix == "" {
		return &common.BadConfigurationError{
			Property: "Document Server Prefix",
			Reason:   "Please specify document server prefix",
		}
	}

	if c.DESJwtHeader == "" || strings.TrimSpace(c.DESJwtHeader) == "Authorization" {
		return &common.BadConfigurationError{
			Property: "Document Server Header",
			Reason:   "Please specify document server header (Note: do not use 'Authorization' header)",
		}
	}

	return nil
}

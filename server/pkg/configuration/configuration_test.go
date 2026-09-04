/**
 *
 * (c) Copyright Ascensio System SIA 2026
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
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/common"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func validCredentials() *Configuration {
	return &Configuration{
		DESAddress:   "https://docs.example.com",
		DESJwt:       "secret",
		DESJwtHeader: "AuthorizationJWT",
		DESJwtPrefix: "Bearer ",
	}
}

func TestConfigurationClone(t *testing.T) {
	source := validCredentials()
	source.DemoEnabled = true
	source.Formats = "docx,xlsx"
	source.OwnerProtected = true

	clone := source.Clone()

	assert.Equal(t, source.DESAddress, clone.DESAddress)
	assert.Equal(t, source.Formats, clone.Formats)
	assert.True(t, clone.DemoEnabled)
	assert.True(t, clone.OwnerProtected)

	clone.DESAddress = "changed"

	assert.NotEqual(t, source.DESAddress, clone.DESAddress)
}

func TestSanitizeConfigurationTrimsAndStripsTrailingSlash(t *testing.T) {
	configuration := &Configuration{
		DESAddress:   "  https://docs.example.com///  ",
		DESJwt:       "  secret  ",
		DESJwtHeader: "  AuthorizationJWT  ",
		DESJwtPrefix: "  Bearer   ",
		Formats:      "  docx  ",
	}

	configuration.SanitizeConfiguration()

	assert.Equal(t, "https://docs.example.com", configuration.DESAddress)
	assert.Equal(t, "secret", configuration.DESJwt)
	assert.Equal(t, "AuthorizationJWT", configuration.DESJwtHeader)
	assert.Equal(t, "Bearer", configuration.DESJwtPrefix)
	assert.Equal(t, "docx", configuration.Formats)
	assert.Equal(t, "https://onlinedocs.docs.onlyoffice.com", configuration.DemoAddress)
}

func TestSanitizeConfigurationClearsExpiredDemoCredentials(t *testing.T) {
	configuration := &Configuration{
		DESAddress:  "https://onlinedocs.docs.onlyoffice.com",
		DESJwt:      "demo-secret",
		DemoEnabled: false,
	}

	configuration.SanitizeConfiguration()

	assert.Empty(t, configuration.DESAddress)
	assert.Empty(t, configuration.DESJwt)
}

func TestIsValidWithCredentials(t *testing.T) {
	configuration := validCredentials()
	require.NoError(t, configuration.IsValid())
}

func TestIsValidRejectsMissingCredentials(t *testing.T) {
	configuration := &Configuration{}
	err := configuration.IsValid()

	require.Error(t, err)

	var bad *common.BadConfigurationError
	assert.ErrorAs(t, err, &bad)
}

func TestIsValidRejectsAuthorizationHeader(t *testing.T) {
	configuration := validCredentials()
	configuration.DESJwtHeader = "Authorization"
	err := configuration.IsValid()
	require.Error(t, err)
}

func TestIsValidRejectsInvalidURL(t *testing.T) {
	configuration := validCredentials()
	configuration.DESAddress = "not-a-url"
	err := configuration.IsValid()
	require.Error(t, err)
	var inv *common.InvalidDocumentServerAddressError
	assert.ErrorAs(t, err, &inv)
}

func TestIsValidAllowsActiveDemo(t *testing.T) {
	configuration := &Configuration{
		DemoEnabled: true,
		DemoExpires: time.Now().UnixMilli() + 60_000,
	}

	require.NoError(t, configuration.IsValid())
}

func TestIsValidRejectsExpiredDemoWithoutCredentials(t *testing.T) {
	configuration := &Configuration{
		DemoEnabled: true,
		DemoExpires: time.Now().UnixMilli() - 1,
	}

	err := configuration.IsValid()
	require.Error(t, err)
}

func TestIsFormatAllowed(t *testing.T) {
	configuration := &Configuration{Formats: "docx, xlsx"}
	assert.True(t, configuration.IsFormatAllowedForViewing("docx"))
	assert.True(t, configuration.IsFormatAllowedForEditing("XLSX"))
	assert.False(t, configuration.IsFormatAllowedForViewing("pptx"))

	emptyConfiguration := &Configuration{Formats: ""}
	assert.True(t, emptyConfiguration.IsFormatAllowedForViewing("anything"))

	noneConfiguration := &Configuration{Formats: EmptyFormats}
	assert.False(t, noneConfiguration.IsFormatAllowedForViewing("docx"))
	assert.False(t, noneConfiguration.IsFormatAllowedForEditing("docx"))
}

func TestValidateFormatsInvalidName(t *testing.T) {
	configuration := validCredentials()
	configuration.Formats = "not-a-real-format"
	err := configuration.IsValid()

	require.Error(t, err)

	var bad *common.BadConfigurationError
	assert.ErrorAs(t, err, &bad)
}

func TestHandleDemoConfigurationSetsCredentials(t *testing.T) {
	api := &plugintest.API{}
	now := time.Now().UnixMilli()
	api.On("KVGet", DemoKey).Return([]byte(""), nil)
	api.On("KVSet", DemoKey, mock.AnythingOfType("[]uint8")).Return(nil)

	configuration := &Configuration{DemoEnabled: true}
	configuration.SanitizeConfiguration()
	configuration.HandleDemoConfiguration(api)

	assert.Equal(t, configuration.DemoAddress, configuration.DESAddress)
	assert.Equal(t, configuration.DemoSecret, configuration.DESJwt)
	assert.Equal(t, configuration.DemoHeader, configuration.DESJwtHeader)
	assert.Equal(t, configuration.DemoPrefix, configuration.DESJwtPrefix)
	assert.Greater(t, configuration.DemoExpires, now)
	api.AssertExpectations(t)
}

func TestHandleDemoConfigurationDisabledNoop(t *testing.T) {
	api := &plugintest.API{}
	configuration := &Configuration{DemoEnabled: false, DESAddress: "https://docs.example.com"}
	configuration.HandleDemoConfiguration(api)
	assert.Equal(t, "https://docs.example.com", configuration.DESAddress)
}

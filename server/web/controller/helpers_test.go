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
package controller

import (
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type stubBot struct {
	dms     []string
	replies []string
}

func (b *stubBot) BotCreateDM(message string, userID string) {
	b.dms = append(b.dms, userID+":"+message)
}

func (b *stubBot) BotCreatePost(message string, channelID string) {}

func (b *stubBot) BotCreateReply(message string, channelID string, parentID string) {
	b.replies = append(b.replies, message)
}

type stubCommandClient struct {
	convertResp client.ConvertResponse
	convertErr  error
	versionResp client.VersionResponse
	versionErr  error
}

func (c *stubCommandClient) SendVersion(commandURL string, request client.VersionRequest, timeout time.Duration) (client.VersionResponse, error) {
	return c.versionResp, c.versionErr
}

func (c *stubCommandClient) SendConvert(commandURL string, request client.ConvertRequest, timeout time.Duration) (client.ConvertResponse, error) {
	return c.convertResp, c.convertErr
}

func validDESConfig() *configuration.Configuration {
	return &configuration.Configuration{
		DESAddress:     "https://docs.example.com",
		DESJwt:         "secret",
		DESJwtHeader:   "AuthorizationJWT",
		DESJwtPrefix:   "Bearer ",
		DemoAddress:    "https://onlinedocs.docs.onlyoffice.com",
		PluginsEnabled: true,
		MacrosEnabled:  false,
	}
}

func siteURLConfig() *model.Config {
	config := &model.Config{}
	siteURL := "https://mm.example.com"
	config.ServiceSettings.SiteURL = &siteURL
	return config
}

func newFileHelper(t *testing.T) file.FileHelper {
	t.Helper()
	formatManager, err := public.NewMapFormatManager()

	require.NoError(t, err)

	return file.New(formatManager)
}

func newFormatManager(t *testing.T) public.FormatManager {
	t.Helper()
	formatManager, err := public.NewMapFormatManager()

	require.NoError(t, err)

	return formatManager
}

func loggingAPI() *plugintest.API {
	api := &plugintest.API{}
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogWarn", mock.Anything).Return().Maybe()
	api.On("LogInfo", mock.Anything).Return().Maybe()
	return api
}

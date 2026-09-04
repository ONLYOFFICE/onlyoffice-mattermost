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
package health

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type stubBot struct {
	dms []string
	mu  sync.Mutex
}

func (b *stubBot) BotCreateDM(message string, userID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dms = append(b.dms, userID+":"+message)
}

func (b *stubBot) BotCreatePost(message string, channelID string) {}

func (b *stubBot) BotCreateReply(message string, channelID string, parentID string) {}

type stubCommandClient struct {
	response client.VersionResponse
	err      error
}

func (c *stubCommandClient) SendVersion(commandURL string, request client.VersionRequest, timeout time.Duration) (client.VersionResponse, error) {
	return c.response, c.err
}

func (c *stubCommandClient) SendConvert(commandURL string, request client.ConvertRequest, timeout time.Duration) (client.ConvertResponse, error) {
	return client.ConvertResponse{}, nil
}

func validHealthConfig() *configuration.Configuration {
	return &configuration.Configuration{
		DESAddress:                 "https://docs.example.com",
		DESJwt:                     "secret",
		DESJwtHeader:               "AuthorizationJWT",
		DESJwtPrefix:               "Bearer ",
		HealthNotificationsEnabled: true,
	}
}

func TestHealthCheckerDefaultsHealthy(t *testing.T) {
	api := &plugintest.API{}
	handler := NewHealthChecker(validHealthConfig(), api, &stubBot{}, &stubCommandClient{}, crypto.NewJwtManager())
	assert.True(t, handler.IsHealthy())
}

func TestHealthCheckerCheckHealthSuccess(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogInfo", mock.Anything).Return().Maybe()
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("PublishWebSocketEvent", "health_status", mock.Anything, mock.Anything).Return().Maybe()

	cmd := &stubCommandClient{response: client.VersionResponse{Error: 0, Version: "8.2.0"}}
	handler := NewHealthChecker(validHealthConfig(), api, &stubBot{}, cmd, crypto.NewJwtManager()).(*healthChecker)
	handler.setHealthy(false)
	handler.checkHealth()

	assert.True(t, handler.IsHealthy())
}

func TestHealthCheckerCheckHealthInvalidConfig(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("PublishWebSocketEvent", "health_status", mock.Anything, mock.Anything).Return()
	api.On("GetUsers", mock.Anything).Return([]*model.User{}, nil)

	configuration := &configuration.Configuration{}
	bot := &stubBot{}
	handler := NewHealthChecker(configuration, api, bot, &stubCommandClient{}, crypto.NewJwtManager()).(*healthChecker)
	handler.checkHealth()

	assert.False(t, handler.IsHealthy())
}

func TestHealthCheckerCheckHealthCommandError(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("PublishWebSocketEvent", "health_status", mock.Anything, mock.Anything).Return()
	api.On("GetUsers", mock.Anything).Return([]*model.User{
		{Id: "admin-1", Roles: model.SystemAdminRoleId},
	}, nil)

	commandClient := &stubCommandClient{err: errors.New("unreachable")}
	bot := &stubBot{}
	handler := NewHealthChecker(validHealthConfig(), api, bot, commandClient, crypto.NewJwtManager()).(*healthChecker)
	handler.checkHealth()

	assert.False(t, handler.IsHealthy())
	require.NotEmpty(t, bot.dms)
}

func TestHealthCheckerStartStop(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogInfo", mock.Anything).Return().Maybe()
	api.On("LogDebug", mock.Anything).Return().Maybe()
	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("PublishWebSocketEvent", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	api.On("GetUsers", mock.Anything).Return([]*model.User{}, nil).Maybe()

	commandClient := &stubCommandClient{response: client.VersionResponse{Error: 0}}

	handler := NewHealthChecker(validHealthConfig(), api, &stubBot{}, commandClient, crypto.NewJwtManager())
	handler.Start()
	handler.Start()

	time.Sleep(50 * time.Millisecond)

	handler.Stop()
	handler.Stop()
}

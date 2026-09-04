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
	"testing"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type fxStubBot struct{}

func (b *fxStubBot) BotCreateDM(message string, userID string) {}

func (b *fxStubBot) BotCreatePost(message string, channelID string) {}

func (b *fxStubBot) BotCreateReply(message string, channelID string, parentID string) {}

func TestHealthModuleProvidesHealthChecker(t *testing.T) {
	api := &plugintest.API{}
	cfg := &configuration.Configuration{
		DESAddress:   "https://docs.example.com",
		DESJwt:       "secret",
		DESJwtHeader: "AuthorizationJWT",
		DESJwtPrefix: "Bearer ",
	}

	var checker HealthChecker
	app := fxtest.New(t,
		fx.NopLogger,
		fx.Provide(
			func() plugin.API { return api },
			func() *configuration.Configuration { return cfg },
			func() bot.Bot { return &fxStubBot{} },
		),
		crypto.Module,
		client.Module,
		Module,
		fx.Populate(&checker),
	)

	app.RequireStart().RequireStop()

	require.NotNil(t, checker)
	assert.True(t, checker.IsHealthy())
}

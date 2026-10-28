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
package bot

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestBotModuleProvidesBot(t *testing.T) {
	api := &plugintest.API{}
	var bot Bot

	app := fxtest.New(t,
		fx.NopLogger,
		fx.Provide(
			func() plugin.API { return api },
			fx.Annotate(
				func() string { return "bot-id" },
				fx.ResultTags(`name:"bot_id"`),
			),
		),
		Module,
		fx.Populate(&bot),
	)
	app.RequireStart().RequireStop()

	require.NotNil(t, bot)
}

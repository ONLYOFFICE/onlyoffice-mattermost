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
package callback

import (
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/converter"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/v8/platform/shared/filestore"
	"go.uber.org/fx"
)

type Config struct {
	PluginAPI plugin.API
	Converter converter.TimeConverter
	Filestore filestore.FileBackend
	Bot       bot.Bot
}

func NewConfig(
	pluginAPI plugin.API,
	converter converter.TimeConverter,
	filestore filestore.FileBackend,
	bot bot.Bot,
) *Config {
	return &Config{
		PluginAPI: pluginAPI,
		Converter: converter,
		Filestore: filestore,
		Bot:       bot,
	}
}

func New(config *Config) Handler {
	return newHandler(
		config.PluginAPI,
		config.Converter,
		config.Filestore,
		config.Bot,
	)
}

var Module = fx.Module("callback",
	fx.Provide(
		NewConfig,
		New,
	),
)

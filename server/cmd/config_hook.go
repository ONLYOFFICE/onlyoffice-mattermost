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
package cmd

import (
	"sort"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/common"
)

// ConfigurationWillBeSaved protects the plugin settings from partial
// configuration saves: some clients (scripts, older tools, other plugins)
// save a server config that does not contain this plugin's section, which
// would wipe the Document Server address and JWT secret and disable the plugin.
// Missing keys are restored from the current configuration; keys present in
// the new configuration are never overridden.
func (p *Plugin) ConfigurationWillBeSaved(newCfg *model.Config) (*model.Config, error) {
	if newCfg == nil || p.MattermostPlugin.API == nil {
		return nil, nil
	}

	current := p.MattermostPlugin.API.GetUnsanitizedConfig()
	if current == nil {
		return nil, nil
	}

	if newCfg.PluginSettings.Plugins == nil {
		newCfg.PluginSettings.Plugins = make(map[string]map[string]any)
	}

	merged, restored := restoreMissingPluginSettings(current.PluginSettings.Plugins[PluginID], newCfg.PluginSettings.Plugins[PluginID])
	if len(restored) == 0 {
		return nil, nil
	}

	newCfg.PluginSettings.Plugins[PluginID] = merged
	p.MattermostPlugin.API.LogWarn(
		common.OnlyofficeLoggerCmdPrefix+"Restored plugin settings missing from the saved configuration",
		"settings", strings.Join(restored, ","),
	)

	return newCfg, nil
}

// restoreMissingPluginSettings returns next with every key that exists only in
// current copied over, and the sorted list of restored keys.
func restoreMissingPluginSettings(current, next map[string]any) (map[string]any, []string) {
	if len(current) == 0 {
		return next, nil
	}

	if next == nil {
		next = make(map[string]any, len(current))
	}

	restored := make([]string, 0)
	for key, value := range current {
		if _, exists := next[key]; !exists {
			next[key] = value
			restored = append(restored, key)
		}
	}

	sort.Strings(restored)
	return next, restored
}

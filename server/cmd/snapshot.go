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
	"bytes"
	"encoding/json"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/mattermost/mattermost/server/public/plugin"

	integration "github.com/ONLYOFFICE/onlyoffice-mattermost"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/common"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
)

type snapshot struct {
	api       plugin.API
	restoring atomic.Bool
}

func newSnapshot(api plugin.API) *snapshot {
	return &snapshot{api: api}
}

func settingsKeys() map[string]struct{} {
	keys := make(map[string]struct{})
	if integration.Manifest.SettingsSchema == nil {
		return keys
	}

	for _, setting := range integration.Manifest.SettingsSchema.Settings {
		if setting != nil && setting.Key != "" {
			keys[strings.ToLower(setting.Key)] = struct{}{}
		}
	}

	return keys
}

func sortedKeys(settings map[string]any) []string {
	keys := make([]string, 0, len(settings))
	for key := range settings {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func mergeSettings(stored, missing map[string]any) map[string]any {
	merged := make(map[string]any, len(stored)+len(missing))
	for key, value := range stored {
		merged[key] = value
	}

	for key, value := range missing {
		merged[key] = value
	}

	return merged
}

func (s *snapshot) storedSettings() map[string]any {
	current := s.api.GetUnsanitizedConfig()
	if current == nil {
		return nil
	}

	return current.PluginSettings.Plugins[PluginID]
}

func (s *snapshot) loadSnapshot() map[string]any {
	data, appErr := s.api.KVGet(configuration.SettingsSnapshotKey)
	if appErr != nil || len(data) == 0 {
		return nil
	}

	snapshot := make(map[string]any)
	if err := json.Unmarshal(data, &snapshot); err != nil {
		s.api.LogError(common.OnlyofficeLoggerCmdPrefix + "Could not parse the plugin settings snapshot: " + err.Error())
		return nil
	}

	return snapshot
}

func (s *snapshot) storeSnapshot(settings map[string]any) {
	keys := settingsKeys()
	snapshot := make(map[string]any)
	for key, value := range settings {
		if _, known := keys[strings.ToLower(key)]; known {
			snapshot[key] = value
		}
	}

	if len(snapshot) == 0 {
		return
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		return
	}

	if existing, appErr := s.api.KVGet(configuration.SettingsSnapshotKey); appErr == nil && bytes.Equal(existing, data) {
		return
	}

	if appErr := s.api.KVSet(configuration.SettingsSnapshotKey, data); appErr != nil {
		s.api.LogError(common.OnlyofficeLoggerCmdPrefix + "Could not store the plugin settings snapshot: " + appErr.Error())
	}
}

func (s *snapshot) missingSettings() (map[string]any, map[string]any) {
	stored := s.storedSettings()
	keys := settingsKeys()
	missing := make(map[string]any)
	for key, value := range s.loadSnapshot() {
		if _, known := keys[strings.ToLower(key)]; !known {
			continue
		}

		if _, exists := stored[key]; !exists {
			missing[key] = value
		}
	}

	return stored, missing
}

func (s *snapshot) restoredSettings() map[string]any {
	return mergeSettings(s.missingSettings())
}

func (s *snapshot) applySnapshot(config *configuration.Configuration) {
	_, missing := s.missingSettings()
	if len(missing) == 0 {
		return
	}

	data, err := json.Marshal(missing)
	if err != nil {
		return
	}

	if err := json.Unmarshal(data, config); err != nil {
		s.api.LogError(common.OnlyofficeLoggerCmdPrefix + "Could not apply the plugin settings snapshot: " + err.Error())
	}
}

func (s *snapshot) reconcileSnapshot() {
	stored, missing := s.missingSettings()
	if len(missing) == 0 {
		s.storeSnapshot(stored)
		return
	}

	if !s.restoring.CompareAndSwap(false, true) {
		return
	}

	restored := mergeSettings(stored, missing)
	s.api.LogWarn(
		common.OnlyofficeLoggerCmdPrefix+"Restoring plugin settings missing from the stored configuration",
		"settings", strings.Join(sortedKeys(missing), ", "),
	)

	go func() {
		defer s.restoring.Store(false)
		if appErr := s.api.SavePluginConfig(restored); appErr != nil {
			s.api.LogError(common.OnlyofficeLoggerCmdPrefix + "Could not restore plugin settings: " + appErr.Error())
		}
	}()
}

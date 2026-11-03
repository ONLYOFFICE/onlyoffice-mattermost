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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
)

func snapshotBytes(t *testing.T, settings map[string]any) []byte {
	t.Helper()
	data, err := json.Marshal(settings)
	require.NoError(t, err)
	return data
}

func fullSettings(address string) map[string]any {
	return map[string]any{
		"desaddress":   address,
		"desjwt":       "secret",
		"desjwtheader": "AuthorizationJWT",
		"desjwtprefix": "Bearer ",
		"formats":      "docx",
	}
}

func snapshotAPI(t *testing.T, stored map[string]any, snapshot []byte) *plugintest.API {
	t.Helper()
	api := &plugintest.API{}

	api.On("LogError", mock.Anything).Return().Maybe()
	api.On("LogWarn", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	api.On("KVGet", configuration.SettingsSnapshotKey).Return(snapshot, nil)
	api.On("GetUnsanitizedConfig").Return(configWithPluginSettings(stored))

	return api
}

func TestReconcileSettingsSnapshotStoresCompleteSettings(t *testing.T) {
	stored := fullSettings("https://docs.example.com")
	stored["obsolete"] = "value"

	api := snapshotAPI(t, stored, nil)
	api.On("KVSet", configuration.SettingsSnapshotKey, mock.Anything).Return(nil).Once()

	plugin := &Plugin{}
	plugin.SetAPI(api)
	plugin.snapshot.reconcileSnapshot()

	api.AssertCalled(t, "KVSet", configuration.SettingsSnapshotKey, snapshotBytes(t, fullSettings("https://docs.example.com")))
	api.AssertNotCalled(t, "SavePluginConfig", mock.Anything)
}

func TestReconcileSettingsSnapshotSkipsUnchangedSnapshot(t *testing.T) {
	settings := fullSettings("https://docs.example.com")
	api := snapshotAPI(t, settings, snapshotBytes(t, settings))

	plugin := &Plugin{}
	plugin.SetAPI(api)
	plugin.snapshot.reconcileSnapshot()

	api.AssertNotCalled(t, "KVSet", mock.Anything, mock.Anything)
	api.AssertNotCalled(t, "SavePluginConfig", mock.Anything)
}

func TestReconcileSettingsSnapshotRestoresTruncatedSettings(t *testing.T) {
	snapshot := fullSettings("https://docs.example.com")
	api := snapshotAPI(t, map[string]any{"formats": "docx, xlsx", "desjwtprefix": ""}, snapshotBytes(t, snapshot))

	saved := make(chan map[string]any, 1)
	api.On("SavePluginConfig", mock.Anything).
		Run(func(args mock.Arguments) { saved <- args.Get(0).(map[string]any) }).
		Return((*model.AppError)(nil))

	plugin := &Plugin{}
	plugin.SetAPI(api)
	plugin.snapshot.reconcileSnapshot()

	select {
	case settings := <-saved:
		assert.Equal(t, "https://docs.example.com", settings["desaddress"])
		assert.Equal(t, "secret", settings["desjwt"])
		assert.Equal(t, "AuthorizationJWT", settings["desjwtheader"])
		assert.Equal(t, "", settings["desjwtprefix"])
		assert.Equal(t, "docx, xlsx", settings["formats"])
	case <-time.After(time.Second):
		t.Fatal("SavePluginConfig was not called")
	}

	api.AssertCalled(t, "LogWarn", mock.Anything, "settings", "desaddress, desjwt, desjwtheader")
	api.AssertNotCalled(t, "KVSet", mock.Anything, mock.Anything)
}

func TestPrepareConfigurationAppliesSettingsSnapshot(t *testing.T) {
	api := snapshotAPI(t, map[string]any{"formats": "docx, xlsx"}, snapshotBytes(t, fullSettings("https://docs.example.com")))
	api.On("LoadPluginConfiguration", mock.AnythingOfType("*configuration.Configuration")).
		Run(func(args mock.Arguments) {
			cfg := args.Get(0).(*configuration.Configuration)
			*cfg = configuration.Configuration{Formats: "docx, xlsx"}
		}).
		Return(nil)
	api.On("KVSet", mock.Anything, mock.Anything).Return(nil).Maybe()

	plugin := &Plugin{}
	plugin.SetAPI(api)

	cfg, err := plugin.prepareConfiguration()

	require.NoError(t, err)
	assert.Equal(t, "https://docs.example.com", cfg.DESAddress)
	assert.Equal(t, "secret", cfg.DESJwt)
	assert.Equal(t, "AuthorizationJWT", cfg.DESJwtHeader)
	assert.Equal(t, "Bearer", cfg.DESJwtPrefix)
	assert.Equal(t, "docx, xlsx", cfg.Formats)
	require.NoError(t, cfg.IsValid())
}

func TestOnConfigurationChangeRecoversTruncatedSettings(t *testing.T) {
	docs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(client.VersionResponse{Error: 0, Version: "8.2.0"})
	}))

	t.Cleanup(docs.Close)

	snapshot := fullSettings(docs.URL)
	snapshot["desallowprivate"] = true

	stored := localFileConfig(t)
	stored.PluginSettings.Plugins[PluginID] = map[string]any{"formats": "docx, xlsx"}

	api := &plugintest.API{}
	api.On("KVGet", configuration.SettingsSnapshotKey).Return(snapshotBytes(t, snapshot), nil)
	api.On("GetUnsanitizedConfig").Return(stored)
	api.On("LoadPluginConfiguration", mock.AnythingOfType("*configuration.Configuration")).
		Run(func(args mock.Arguments) {
			cfg := args.Get(0).(*configuration.Configuration)
			*cfg = configuration.Configuration{Formats: "docx, xlsx"}
		}).
		Return(nil)

	saved := make(chan struct{}, 1)

	api.On("SavePluginConfig", mock.Anything).
		Run(func(mock.Arguments) { saved <- struct{}{} }).
		Return((*model.AppError)(nil))

	for _, call := range []*mock.Call{
		api.On("LogError", mock.Anything).Return(),
		api.On("LogDebug", mock.Anything).Return(),
		api.On("LogInfo", mock.Anything).Return(),
		api.On("LogWarn", mock.Anything, mock.Anything, mock.Anything).Return(),
		api.On("PublishWebSocketEvent", mock.Anything, mock.Anything, mock.Anything).Return(),
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil),
		api.On("GetUsers", mock.Anything).Return([]*model.User{}, nil),
		api.On("GetLicense").Return((*model.License)(nil)),
		api.On("EnsureBotUser", mock.AnythingOfType("*model.Bot")).Return("bot-id", nil),
		api.On("GetBundlePath").Return(repoRoot(t), nil),
		api.On("SetProfileImage", "bot-id", mock.AnythingOfType("[]uint8")).Return(nil),
	} {
		call.Maybe()
	}

	plugin := &Plugin{}
	plugin.SetAPI(api)

	require.NoError(t, plugin.OnConfigurationChange())

	cfg := plugin.getConfiguration()

	assert.NoError(t, cfg.Error)
	assert.Equal(t, docs.URL, cfg.DESAddress)
	assert.Equal(t, "docx, xlsx", cfg.Formats)
	api.AssertNotCalled(t, "DisablePlugin", mock.Anything)

	select {
	case <-saved:
	case <-time.After(time.Second):
		t.Fatal("SavePluginConfig was not called")
	}

	require.NoError(t, plugin.OnDeactivate())
}

func TestConfigurationWillBeSavedRestoresFromSnapshot(t *testing.T) {
	api := snapshotAPI(t, map[string]any{"formats": "docx"}, snapshotBytes(t, fullSettings("https://docs.example.com")))

	plugin := &Plugin{}
	plugin.SetAPI(api)

	saved, err := plugin.ConfigurationWillBeSaved(configWithPluginSettings(map[string]any{
		"formats": "docx, xlsx",
	}))

	require.NoError(t, err)
	require.NotNil(t, saved)

	settings := saved.PluginSettings.Plugins[PluginID]

	assert.Equal(t, "https://docs.example.com", settings["desaddress"])
	assert.Equal(t, "secret", settings["desjwt"])
	assert.Equal(t, "docx, xlsx", settings["formats"])
}

func TestReconcileSnapshotSkipsRestoreInFlight(t *testing.T) {
	api := snapshotAPI(t, map[string]any{"formats": "docx"}, snapshotBytes(t, fullSettings("https://docs.example.com")))

	release := make(chan struct{})
	done := make(chan struct{})

	api.On("SavePluginConfig", mock.Anything).
		Run(func(mock.Arguments) {
			<-release
			close(done)
		}).
		Return((*model.AppError)(nil)).Once()

	plugin := &Plugin{}
	plugin.SetAPI(api)
	plugin.snapshot.reconcileSnapshot()
	plugin.snapshot.reconcileSnapshot()

	close(release)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("SavePluginConfig was not called")
	}

	assert.Eventually(t, func() bool { return !plugin.snapshot.restoring.Load() }, time.Second, 10*time.Millisecond)
	api.AssertNumberOfCalls(t, "SavePluginConfig", 1)
}

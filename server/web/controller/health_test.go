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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type stubHealthChecker struct {
	healthy bool
}

func (s *stubHealthChecker) Start() {}

func (s *stubHealthChecker) Stop() {}

func (s *stubHealthChecker) IsHealthy() bool { return s.healthy }

func TestHealthHandlerHealthy(t *testing.T) {
	api := &plugintest.API{}
	handler := NewHealthHandler(api, &stubHealthChecker{healthy: true})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	handler.GetHealthStatus(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response HealthStatusResponse

	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))

	assert.True(t, response.Healthy)
	assert.Greater(t, response.Time, int64(0))
}

func TestHealthHandlerUnhealthy(t *testing.T) {
	api := &plugintest.API{}
	handler := NewHealthHandler(api, &stubHealthChecker{healthy: false})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	handler.GetHealthStatus(recorder, req)

	var response HealthStatusResponse

	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))

	assert.False(t, response.Healthy)
}

func TestHealthHandlerNilCheckerDefaultsHealthy(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogError", mock.Anything).Return().Maybe()
	handler := NewHealthHandler(api, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	handler.GetHealthStatus(recorder, req)

	var response HealthStatusResponse

	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))

	assert.True(t, response.Healthy)
}

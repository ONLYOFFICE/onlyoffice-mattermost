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
package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPermissionsName(t *testing.T) {
	assert.Equal(t, `"edit"`, GetPermissionsName(model.Permissions{Edit: true}))
	assert.Equal(t, `"read only"`, GetPermissionsName(model.Permissions{Edit: false}))
}

func TestWriteJSONDefaultStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	WriteJSON(rr, map[string]int{"error": 0})

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var body map[string]int

	require.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
	assert.Equal(t, 0, body["error"])
}

func TestWriteJSONCustomStatus(t *testing.T) {
	recorder := httptest.NewRecorder()
	WriteJSON(recorder, map[string]int{"error": 1}, http.StatusForbidden)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

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
package api

import (
	"encoding/json"
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
)

func WriteJSON(w http.ResponseWriter, v interface{}, code ...int) {
	if len(code) == 1 {
		w.WriteHeader(code[0])
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func GetPermissionsName(permissions model.Permissions) string {
	if permissions.Edit {
		return "\"edit\""
	}
	return "\"read only\""
}

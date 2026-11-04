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
package callback

import (
	"encoding/json"

	"github.com/mattermost/mattermost/server/public/plugin"
)

// LastEdit describes the latest save of a document made through ONLYOFFICE.
type LastEdit struct {
	UserID string `json:"user_id"`
	At     int64  `json:"at"`
}

func lastEditKey(fileID string) string {
	return "onlyoffice_lastedit_" + fileID
}

// recordLastEdit stores who saved the file and when. Failures are logged only:
// they must never fail an otherwise successful save.
func recordLastEdit(api plugin.API, fileID string, edit LastEdit) {
	data, err := json.Marshal(edit)
	if err != nil {
		return
	}

	if appErr := api.KVSet(lastEditKey(fileID), data); appErr != nil {
		api.LogWarn(onlyofficeLoggerCallbackPrefix + "could not record last edit for file " + fileID + ": " + appErr.Error())
	}
}

// GetLastEdit returns the latest recorded save of the file, or nil if the file
// has not been saved through ONLYOFFICE since this information is recorded.
func GetLastEdit(api plugin.API, fileID string) *LastEdit {
	data, appErr := api.KVGet(lastEditKey(fileID))
	if appErr != nil || len(data) == 0 {
		return nil
	}

	var edit LastEdit
	if err := json.Unmarshal(data, &edit); err != nil || edit.At == 0 {
		return nil
	}

	return &edit
}

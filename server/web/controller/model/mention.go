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
package model

import "encoding/json"

type MentionUser struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MentionUsersResponse struct {
	C       string        `json:"c"`
	Users   []MentionUser `json:"users"`
	HasMore bool          `json:"hasMore,omitempty"`
}

func (m MentionUsersResponse) ToJSON() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}

	return data
}

type MentionNotifyRequest struct {
	FileID     string          `json:"fileId"`
	Emails     []string        `json:"emails"`
	Message    string          `json:"message"`
	ActionLink json.RawMessage `json:"actionLink,omitempty"`
}

type MentionNotifyResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func (m *MentionNotifyResponse) ToJSON() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}

	return data
}

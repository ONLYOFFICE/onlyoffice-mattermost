/**
 *
 * (c) Copyright Ascensio System SIA 2022
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

import "github.com/go-playground/validator/v10"

type Callback struct {
	Actions []struct {
		Type   int    `json:"type"`
		UserID string `json:"userid"`
	} `json:"actions"`
	Key    string   `json:"key" validate:"required"`
	Status int      `json:"status" validate:"required"`
	Users  []string `json:"users"`
	URL    string   `json:"url"`
	FileID string   `json:"-" validate:"required"`
	Token  string   `json:"token"`
}

type CallbackResponse struct {
	Error int8 `json:"error"`
}

func (c *Callback) Validate() error {
	return validator.New().Struct(c)
}

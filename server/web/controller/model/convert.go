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
package model

import (
	"encoding/json"

	validator "github.com/go-playground/validator/v10"
)

type ConvertFileRequest struct {
	FileID     string `json:"file_id" validate:"required,min=1"`
	OutputType string `json:"output_type,omitempty"`
	Password   string `json:"password,omitempty"`
}

func (c *ConvertFileRequest) Validate() error {
	return validator.New().Struct(c)
}

type ConvertFileResponse struct {
	Error int `json:"error"`
}

func (c *ConvertFileResponse) ToJSON() []byte {
	json, err := json.Marshal(c)
	if err != nil {
		return []byte{}
	}

	return json
}

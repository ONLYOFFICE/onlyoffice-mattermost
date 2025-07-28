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
package client

import jwt "github.com/golang-jwt/jwt/v5"

type VersionRequest struct {
	jwt.RegisteredClaims
	Command string `json:"c" mapstructure:"c"`
	Token   string `json:"token,omitempty" mapstructure:"token"`
	Header  string `json:"-"`
	Prefix  string `json:"-"`
}

type VersionResponse struct {
	Error   int    `json:"error" mapstructure:"error"`
	Version string `json:"version,omitempty" mapstructure:"version"`
}

type ConvertRequest struct {
	jwt.RegisteredClaims
	Async      bool   `json:"async"`
	Key        string `json:"key"`
	Filetype   string `json:"filetype"`
	Outputtype string `json:"outputtype"`
	Password   string `json:"password,omitempty"`
	URL        string `json:"url"`
	Token      string `json:"token,omitempty"`
	Region     string `json:"region,omitempty"`
}

type ConvertResponse struct {
	FileURL  string `json:"fileUrl"`
	FileType string `json:"fileType"`
	Error    int    `json:"error"`
}

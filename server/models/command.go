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

package models

import (
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type CommandBody struct {
	Command            string `json:"c"`
	Token              string `json:"token,omitempty"`
	jwt.StandardClaims `json:"-"`
}

type CommandResponse struct {
	Error        int    `json:"error"`
	Version      string `json:"version,omitempty"`
	IsSuccessful bool   `json:"-"`
}

func (dsr CommandResponse) ProcessResponse() error {
	if !dsr.IsSuccessful {
		var err error = errors.New("[ONLYOFFICE]: No connection to the Document Service")
		return err
	}
	if dsr.Error > 0 {
		var OnlyofficeError error = errors.New("[ONLYOFFICE]: The server responded with an error: " + strconv.Itoa(dsr.Error))
		return OnlyofficeError
	}
	return nil
}

func (dsr *CommandResponse) Succeeded() {
	dsr.IsSuccessful = true
}

func (dsr *CommandResponse) Failed() {
	dsr.IsSuccessful = false
}

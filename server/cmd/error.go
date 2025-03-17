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
package cmd

import (
	"errors"
	"fmt"
)

var ErrDeprecatedDocumentServerVersion = errors.New(_OnlyofficeLoggerPrefix + "old document server version")
var ErrParseDocumentServerVersion = errors.New(_OnlyofficeLoggerPrefix + "could not parse document server version")
var ErrCreateBotProfile = errors.New(_OnlyofficeLoggerPrefix + "could not create bot profile")
var ErrLoadBotProfileImage = errors.New(_OnlyofficeLoggerPrefix + "could not load bot profile image")
var ErrSetBotProfileImage = errors.New(_OnlyofficeLoggerPrefix + "could not set bot profile image")

type DocumentServerCommandResponseError struct {
	Code int
}

func (e *DocumentServerCommandResponseError) Error() string {
	return fmt.Sprintf(_OnlyofficeLoggerPrefix+"could not fetch document server's version (%d)", e.Code)
}

type InvalidDocumentServerAddressError struct {
	Reason string
}

func (e *InvalidDocumentServerAddressError) Error() string {
	return fmt.Sprintf(_OnlyofficeLoggerPrefix+"invalid document server address (%s)", e.Reason)
}

type BadConfigurationError struct {
	Property string
	Reason   string
}

func (e *BadConfigurationError) Error() string {
	return fmt.Sprintf(_OnlyofficeLoggerPrefix+"bad property '%s' configuration (%s)", e.Property, e.Reason)
}

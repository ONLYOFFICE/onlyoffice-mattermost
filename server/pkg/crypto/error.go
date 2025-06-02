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
package crypto

import (
	"errors"
)

var ErrJwtManagerSigning = errors.New(onlyofficeLoggerCryptoPrefix + "could not generate a new jwt")
var ErrJwtManagerEmptyToken = errors.New(onlyofficeLoggerCryptoPrefix + "could not verify an empty jwt")
var ErrJwtManagerEmptyDecodingBody = errors.New(onlyofficeLoggerCryptoPrefix + "could not decode a jwt. Got empty interface")
var ErrJwtManagerInvalidSigningMethod = errors.New(onlyofficeLoggerCryptoPrefix + "unexpected jwt signing method")
var ErrJwtManagerCastOrInvalidToken = errors.New(onlyofficeLoggerCryptoPrefix + "could not cast claims or invalid jwt")

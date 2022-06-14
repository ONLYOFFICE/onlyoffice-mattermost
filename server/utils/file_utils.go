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

package utils

import (
	"github.com/pkg/errors"
)

func IsExtensionSupported(fileExt string) bool {
	_, exists := ONLYOFFICE_EXTENSION_TYPE_MAP[fileExt]
	if exists {
		return true
	}
	return false
}

func IsExtensionEditable(fileExt string) bool {
	_, exists := ONLYOFFICE_EDITABLE_EXTENSION_MAP[fileExt]
	if exists {
		return true
	}
	return false
}

func GetFileType(fileExt string) (string, error) {
	fileType, exists := ONLYOFFICE_EXTENSION_TYPE_MAP[fileExt]
	if !exists {
		return "", errors.New(ONLYOFFICE_LOGGER_PREFIX + "This extension is not supported")
	}
	return fileType, nil
}

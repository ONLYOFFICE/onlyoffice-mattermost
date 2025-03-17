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
package onlyoffice

import "strings"

func (h helper) IsExtensionSupported(fileExt string) bool {
	_, exists := OnlyofficeFileExtensions[strings.ToLower(fileExt)]
	return exists
}

func (h helper) IsExtensionEditable(fileExt string) bool {
	_, exists := OnlyofficeEditableExtensions[strings.ToLower(fileExt)]
	return exists
}

func (h helper) GetFileType(fileExt string) (string, error) {
	fileType, exists := OnlyofficeFileExtensions[strings.ToLower(fileExt)]
	if !exists {
		return "", ErrOnlyofficeExtensionNotSupported
	}
	return fileType, nil
}

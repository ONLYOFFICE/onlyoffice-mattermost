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
package file

import (
	"strings"
)

func (h fileHelperImpl) IsExtensionSupported(fileExt string) bool {
	format, exists := h.formatManager.GetFormatByName(strings.ToLower(fileExt))
	return exists && format.IsViewable()
}

func (h fileHelperImpl) IsExtensionEditable(fileExt string) bool {
	format, exists := h.formatManager.GetFormatByName(strings.ToLower(fileExt))
	return exists && format.IsEditable()
}

func (h fileHelperImpl) GetFileType(fileExt string) (string, error) {
	format, exists := h.formatManager.GetFormatByName(strings.ToLower(fileExt))
	if !exists || !format.IsViewable() {
		return "", ErrExtensionNotSupported
	}
	return format.Type, nil
}

func (h fileHelperImpl) GetWordType() string {
	return onlyofficeWordType
}

func (h fileHelperImpl) GetCellType() string {
	return onlyofficeCellType
}

func (h fileHelperImpl) GetSlideType() string {
	return onlyofficeSlideType
}

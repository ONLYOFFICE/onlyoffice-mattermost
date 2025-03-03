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

const (
	_OnlyofficeLoggerPrefix            string = "[ONLYOFFICE Helper]: "
	OnlyofficeWordType                 string = "word"
	OnlyofficeCellType                 string = "cell"
	OnlyofficeSlideType                string = "slide"
	OnlyofficePermissionsPropSeparator string = "_"
	OnlyofficePermissionsProp          string = "ONLYOFFICE_PERMISSIONS"
	OnlyofficePermissionsWildcardKey   string = "*"
)

var OnlyofficeEditableExtensions = map[string]string{
	"xlsx": OnlyofficeCellType,
	"pptx": OnlyofficeSlideType,
	"docx": OnlyofficeWordType,
}

var OnlyofficeFileExtensions = map[string]string{
	"xls":  OnlyofficeCellType,
	"xlsx": OnlyofficeCellType,
	"xlsm": OnlyofficeCellType,
	"xlt":  OnlyofficeCellType,
	"xltx": OnlyofficeCellType,
	"xltm": OnlyofficeCellType,
	"ods":  OnlyofficeCellType,
	"fods": OnlyofficeCellType,
	"ots":  OnlyofficeCellType,
	"csv":  OnlyofficeCellType,
	"pps":  OnlyofficeSlideType,
	"ppsx": OnlyofficeSlideType,
	"ppsm": OnlyofficeSlideType,
	"ppt":  OnlyofficeSlideType,
	"pptx": OnlyofficeSlideType,
	"pptm": OnlyofficeSlideType,
	"pot":  OnlyofficeSlideType,
	"potx": OnlyofficeSlideType,
	"potm": OnlyofficeSlideType,
	"odp":  OnlyofficeSlideType,
	"fodp": OnlyofficeSlideType,
	"otp":  OnlyofficeSlideType,
	"doc":  OnlyofficeWordType,
	"docx": OnlyofficeWordType,
	"docm": OnlyofficeWordType,
	"dot":  OnlyofficeWordType,
	"dotx": OnlyofficeWordType,
	"dotm": OnlyofficeWordType,
	"odt":  OnlyofficeWordType,
	"fodt": OnlyofficeWordType,
	"ott":  OnlyofficeWordType,
	"rtf":  OnlyofficeWordType,
	"txt":  OnlyofficeWordType,
	"html": OnlyofficeWordType,
	"htm":  OnlyofficeWordType,
	"mht":  OnlyofficeWordType,
	"pdf":  OnlyofficeWordType,
	"djvu": OnlyofficeWordType,
	"fb2":  OnlyofficeWordType,
	"epub": OnlyofficeWordType,
	"xps":  OnlyofficeWordType,
}

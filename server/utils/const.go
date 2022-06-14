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

const ONLYOFFICE_LOGGER_PREFIX = "[ONLYOFFICE Utils]: "

const ONLYOFFICE_WORD_TYPE string = "word"
const ONLYOFFICE_CELL_TYPE string = "cell"
const ONLYOFFICE_SLIDE_TYPE string = "slide"

var ONLYOFFICE_EDITABLE_EXTENSION_MAP map[string]string = map[string]string{
	"xlsx": ONLYOFFICE_CELL_TYPE,
	"pptx": ONLYOFFICE_SLIDE_TYPE,
	"docx": ONLYOFFICE_WORD_TYPE,
}

var ONLYOFFICE_EXTENSION_TYPE_MAP map[string]string = map[string]string{
	"xls":  ONLYOFFICE_CELL_TYPE,
	"xlsx": ONLYOFFICE_CELL_TYPE,
	"xlsm": ONLYOFFICE_CELL_TYPE,
	"xlt":  ONLYOFFICE_CELL_TYPE,
	"xltx": ONLYOFFICE_CELL_TYPE,
	"xltm": ONLYOFFICE_CELL_TYPE,
	"ods":  ONLYOFFICE_CELL_TYPE,
	"fods": ONLYOFFICE_CELL_TYPE,
	"ots":  ONLYOFFICE_CELL_TYPE,
	"csv":  ONLYOFFICE_CELL_TYPE,
	"pps":  ONLYOFFICE_SLIDE_TYPE,
	"ppsx": ONLYOFFICE_SLIDE_TYPE,
	"ppsm": ONLYOFFICE_SLIDE_TYPE,
	"ppt":  ONLYOFFICE_SLIDE_TYPE,
	"pptx": ONLYOFFICE_SLIDE_TYPE,
	"pptm": ONLYOFFICE_SLIDE_TYPE,
	"pot":  ONLYOFFICE_SLIDE_TYPE,
	"potx": ONLYOFFICE_SLIDE_TYPE,
	"potm": ONLYOFFICE_SLIDE_TYPE,
	"odp":  ONLYOFFICE_SLIDE_TYPE,
	"fodp": ONLYOFFICE_SLIDE_TYPE,
	"otp":  ONLYOFFICE_SLIDE_TYPE,
	"doc":  ONLYOFFICE_WORD_TYPE,
	"docx": ONLYOFFICE_WORD_TYPE,
	"docm": ONLYOFFICE_WORD_TYPE,
	"dot":  ONLYOFFICE_WORD_TYPE,
	"dotx": ONLYOFFICE_WORD_TYPE,
	"dotm": ONLYOFFICE_WORD_TYPE,
	"odt":  ONLYOFFICE_WORD_TYPE,
	"fodt": ONLYOFFICE_WORD_TYPE,
	"ott":  ONLYOFFICE_WORD_TYPE,
	"rtf":  ONLYOFFICE_WORD_TYPE,
	"txt":  ONLYOFFICE_WORD_TYPE,
	"html": ONLYOFFICE_WORD_TYPE,
	"htm":  ONLYOFFICE_WORD_TYPE,
	"mht":  ONLYOFFICE_WORD_TYPE,
	"pdf":  ONLYOFFICE_WORD_TYPE,
	"djvu": ONLYOFFICE_WORD_TYPE,
	"fb2":  ONLYOFFICE_WORD_TYPE,
	"epub": ONLYOFFICE_WORD_TYPE,
	"xps":  ONLYOFFICE_WORD_TYPE,
}

const ONLYOFFICE_PERMISSIONS_PROP_SEPARATOR = "_"
const ONLYOFFICE_PERMISSIONS_PROP = "ONLYOFFICE_PERMISSIONS"
const ONLYOFFICE_PERMISSIONS_WILDCARD_KEY = "*"

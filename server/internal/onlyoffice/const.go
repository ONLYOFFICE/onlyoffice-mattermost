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

const (
	_OnlyofficeLoggerPrefix            string = "[ONLYOFFICE Helper]: "
	OnlyofficeWordType                 string = "word"
	OnlyofficeCellType                 string = "cell"
	OnlyofficeSlideType                string = "slide"
	OnlyofficePermissionsPropSeparator string = "_"
	OnlyofficePermissionsProp          string = "ONLYOFFICE_PERMISSIONS"
	OnlyofficePermissionsWildcardKey   string = "*"
)

var OnlyofficeLanguageMapping = map[string]string{
	"en": "en-US",
	"ru": "ru-RU",
	"de": "de-DE",
	"fr": "fr-FR",
	"es": "es-ES",
	"it": "it-IT",
	"pt": "pt-PT",
	"zh": "zh-CN",
	"ja": "ja-JP",
	"ko": "ko-KR",
	"ar": "ar-SA",
	"bg": "bg-BG",
	"ca": "ca-ES",
	"cs": "cs-CZ",
	"da": "da-DK",
	"el": "el-GR",
	"eu": "eu-ES",
	"fi": "fi-FI",
	"gl": "gl-ES",
	"he": "he-IL",
	"hu": "hu-HU",
	"hy": "hy-AM",
	"id": "id-ID",
	"lv": "lv-LV",
	"ms": "ms-MY",
	"nl": "nl-NL",
	"nb": "nb-NO",
	"pl": "pl-PL",
	"ro": "ro-RO",
	"si": "si-LK",
	"sk": "sk-SK",
	"sl": "sl-SI",
	"sq": "sq-AL",
	"sr": "sr-Latn-RS",
	"sv": "sv-SE",
	"tr": "tr-TR",
	"uk": "uk-UA",
	"ur": "ur-PK",
	"vi": "vi-VN",
}

func MapLanguageToTemplate(locale string) string {
	if strings.Contains(locale, "-") {
		return locale
	}

	if mapped, exists := OnlyofficeLanguageMapping[locale]; exists {
		return mapped
	}

	return "default"
}

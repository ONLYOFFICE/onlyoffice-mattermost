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

const (
	ONLYOFFICE_COMMAND_DROP      string = "drop"
	ONLYOFFICE_COMMAND_FORCESAVE string = "forcesave"
	ONLYOFFICE_COMMAND_INFO      string = "info"
	ONLYOFFICE_COMMAND_META      string = "meta"
	ONLYOFFICE_COMMAND_VERSION   string = "version"
)

var ONLYOFFICE_AUTHOR_PERMISSIONS Permissions = Permissions{
	Comment:  true,
	Copy:     true,
	Download: true,
	Edit:     true,
	Print:    true,
	Review:   true,
}

var ONLYOFFICE_DEFAULT_PERMISSIONS Permissions = Permissions{
	Edit: false,
}

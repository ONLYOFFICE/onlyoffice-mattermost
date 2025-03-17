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
package web

import (
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
)

const (
	_OnlyofficeRedirectULR   = "https://onlyoffice.com"
	_OnlyofficeLoggerPrefix  = "[ONLYOFFICE Handler]: "
	_OnlyofficeAPIRootSuffix = "plugins/com.onlyoffice.mattermost/api"
)

var (
	_CallbackErr = model.CallbackResponse{
		Error: 1,
	}
	_CallbackOK = model.CallbackResponse{
		Error: 0,
	}
)

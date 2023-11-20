/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
	"net/http"
	"net/url"
	"strings"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
)

func BuildDownloadHandler(plugin api.PluginAPI) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var jwt model.DownloadToken

		token := strings.ReplaceAll(r.Header.Get(plugin.Configuration.Header), "Bearer ", "")
		if token == "" {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not extract jwt with the header specified. Please validate your JWT Header settings")
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		err := plugin.Manager.Verify(token, &jwt)
		if err != nil {
			plugin.API.LogError(err.Error())
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		u, err := url.Parse(jwt.Payload.URL)
		if err != nil {
			plugin.API.LogError(err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		file, fileErr := plugin.API.GetFile(u.Query().Get("id"))
		if fileErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not download file. Reason: " + fileErr.Message)
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		plugin.API.LogDebug(_OnlyofficeLoggerPrefix + "downloading file " + u.Query().Get("id"))
		rw.Write(file)
	}
}

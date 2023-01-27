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

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
)

func BuildDownloadHandler(plugin api.PluginAPI) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		var jwt model.DownloadToken

		err := plugin.Manager.Verify(query.Get("token"), &jwt)
		if err != nil {
			plugin.API.LogError(err.Error())
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		file, fileErr := plugin.API.GetFile(jwt.FileID)
		if fileErr != nil {
			plugin.API.LogError(_OnlyofficeLoggerPrefix + "could not download file. Reason: " + fileErr.Message)
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		plugin.API.LogDebug(_OnlyofficeLoggerPrefix + "downloading file " + jwt.FileID)
		rw.Write(file)
	}
}

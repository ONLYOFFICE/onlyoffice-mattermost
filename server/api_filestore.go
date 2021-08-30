/**
 *
 * (c) Copyright Ascensio System SIA 2021
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

package main

import (
	"io"

	"github.com/mattermost/mattermost-server/v5/shared/filestore"
)

func getFilestore(p *Plugin) (filestore.FileBackend, error) {
	license := p.API.GetLicense()
	serverConfig := p.API.GetUnsanitizedConfig()
	filestore, err := filestore.NewFileBackend(serverConfig.FileSettings.ToFileBackendSettings(license != nil && *license.Features.Compliance))
	if err != nil {
		return nil, err
	}
	return filestore, nil
}

func (p *Plugin) WriteFile(fr io.Reader, path string) (int64, error) {
	filestore, err := getFilestore(p)
	if err != nil {
		return 0, err
	}

	result, err := filestore.WriteFile(fr, path)
	if err != nil {
		return result, err
	}
	return result, nil
}

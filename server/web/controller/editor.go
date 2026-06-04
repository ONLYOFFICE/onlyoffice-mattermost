/**
 *
 * (c) Copyright Ascensio System SIA 2026
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
package controller

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/file"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	common "github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/common"
)

type EditorHandler struct {
	api            plugin.API
	configuration  *configuration.Configuration
	fileHelper     file.FileHelper
	encoder        crypto.Encoder
	jwtManager     crypto.JwtManager
	editorTemplate *template.Template
}

func NewEditorHandler(
	api plugin.API,
	configuration *configuration.Configuration,
	fileHelper file.FileHelper,
	encoder crypto.Encoder,
	jwtManager crypto.JwtManager,
	editorTemplate *template.Template,
) EditorHandler {
	return EditorHandler{
		api:            api,
		configuration:  configuration,
		fileHelper:     fileHelper,
		encoder:        encoder,
		jwtManager:     jwtManager,
		editorTemplate: editorTemplate,
	}
}

func (h *EditorHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	h.api.LogDebug(common.LoggerPrefix + "got an editor request")

	config, docKey, statusCode, err := common.BuildEditorConfig(r, h.configuration, h.api, h.encoder, h.jwtManager, h.fileHelper)
	if err != nil {
		h.api.LogError(common.LoggerPrefix + err.Error())
		rw.WriteHeader(statusCode)
		return
	}

	userID := r.Header.Get(tools.MMAuthHeader)
	var mentionsCode string
	mentionsKey := "mentions:" + userID
	mentionsCodeBytes, err := h.api.KVGet(mentionsKey)
	if err != nil || len(mentionsCodeBytes) == 0 {
		mentionsCode = h.fileHelper.GenerateKey()
		if err := h.api.KVSetWithExpiry(mentionsKey, []byte(mentionsCode), 60*60*24); err != nil {
			h.api.LogError(common.LoggerPrefix + "could not set mentions code: " + err.Error())
		}
	} else {
		mentionsCode = string(mentionsCodeBytes)
	}

	if err := h.api.KVSetWithExpiry(mentionsCode, []byte(userID), 60*60*24); err != nil {
		h.api.LogError(common.LoggerPrefix + "could not set mentions code mapping: " + err.Error())
	}

	encodedConfig, cerr := json.Marshal(config)
	if cerr != nil {
		h.api.LogError(common.LoggerPrefix + "could not marshal config: " + cerr.Error())
		return
	}

	data := map[string]any{
		"apijs":        h.configuration.DESAddress + "/web-apps/apps/api/documents/api.js?shardkey=" + docKey,
		"config":       string(encodedConfig),
		"dark":         r.URL.Query().Get("dark"),
		"mentionscode": string(mentionsCode),
	}

	h.api.LogDebug(common.LoggerPrefix + "building an editor window")
	if err := h.editorTemplate.ExecuteTemplate(rw, "editor.html", data); err != nil {
		h.api.LogError(common.LoggerPrefix + "could not execute editor template: " + err.Error())
	}
}

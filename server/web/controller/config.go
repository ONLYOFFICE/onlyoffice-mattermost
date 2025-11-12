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
package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
)

type ConfigHandler struct {
	api           plugin.API
	configuration *configuration.Configuration
	formatManager public.FormatManager
}

func NewConfigHandler(
	api plugin.API,
	configuration *configuration.Configuration,
	formatManager public.FormatManager,
) ConfigHandler {
	return ConfigHandler{
		api:           api,
		configuration: configuration,
		formatManager: formatManager,
	}
}

func (h *ConfigHandler) getAllFormatsWhere(filterFunc func(public.Format) bool) []string {
	names := []string{}
	for name, format := range h.formatManager.GetAllFormats() {
		if filterFunc(format) {
			names = append(names, name)
		}
	}
	return names
}

func (h *ConfigHandler) parseFormats(rawFormats string, filterFunc func(public.Format) bool) []string {
	if rawFormats == "" {
		return h.getAllFormatsWhere(filterFunc)
	}

	if strings.ToLower(strings.TrimSpace(rawFormats)) == configuration.EmptyFormats {
		return []string{}
	}

	formats := []string{}
	for _, part := range strings.Split(rawFormats, ",") {
		if format := strings.TrimSpace(strings.ToLower(part)); format != "" {
			formats = append(formats, format)
		}
	}
	return formats
}

func (h *ConfigHandler) isSupported(f public.Format) bool {
	return f.IsViewable() || f.IsEditable() || f.IsLossyEditable() || f.IsAutoConvertable()
}

func (h *ConfigHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	response := model.FormatResponse{
		Formats: h.parseFormats(h.configuration.Formats, h.isSupported),
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not encode config response: " + err.Error())
	}
}

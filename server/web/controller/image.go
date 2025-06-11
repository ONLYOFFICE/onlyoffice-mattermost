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
	"net/http"

	"github.com/mattermost/mattermost/server/public/plugin"
)

type ImageHandler struct {
	api plugin.API
}

func NewImageHandler(api plugin.API) ImageHandler {
	return ImageHandler{api: api}
}

func (h *ImageHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.api.LogError(onlyofficeLoggerPrefix + "could not extract code")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	userID, err := h.api.KVGet(code)
	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get user id: " + err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	image, err := h.api.GetProfileImage(string(userID))
	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "could not get user image: " + err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	contentType := http.DetectContentType(image)
	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(http.StatusOK)
	rw.Write(image)
}

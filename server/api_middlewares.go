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

import "net/http"

func (p *Plugin) userAccessMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	authentication := &AuthenticationFilter{plugin: p}
	checkFile := &FileValidationFilter{plugin: p}
	postAccess := &PostAuthorizationFilter{plugin: p}

	authentication.SetNext(checkFile).SetNext(postAccess)

	return func(writer http.ResponseWriter, request *http.Request) {
		authentication.DoFilter(writer, request)

		if authentication.HasError() {
			authentication.Reset()
			writer.WriteHeader(403)
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) callbackMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	bodyJwtMiddleware := BodyJwtFilter{plugin: p}
	headerJwtMiddleware := HeaderJwtFilter{plugin: p}

	return func(writer http.ResponseWriter, request *http.Request) {
		bodyJwtMiddleware.DoFilter(writer, request)
		headerJwtMiddleware.DoFilter(writer, request)
		defer bodyJwtMiddleware.Reset()
		defer headerJwtMiddleware.Reset()

		if bodyJwtMiddleware.HasError() && headerJwtMiddleware.HasError() {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(403)
			writer.Write([]byte("{\"error\": 1}"))
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) permissionsMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	authenticationMiddleware := AuthenticationFilter{plugin: p}

	return func(writer http.ResponseWriter, request *http.Request) {

		authenticationMiddleware.DoFilter(writer, request)

		if authenticationMiddleware.HasError() {
			authenticationMiddleware.Reset()
			writer.WriteHeader(403)
			return
		}

		next(writer, request)
	}
}

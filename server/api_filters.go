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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/utils"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/security"
)

type Filter interface {
	DoFilter(writer http.ResponseWriter, request *http.Request)
	SetNext(Filter) Filter
	HasError() bool
	Reset()
}

//
type AuthenticationFilter struct {
	plugin   *Plugin
	next     Filter
	hasError bool
}

func (m *AuthenticationFilter) DoFilter(writer http.ResponseWriter, request *http.Request) {
	userId, cookieErr := request.Cookie(MATTERMOST_USER_COOKIE)
	user, userErr := m.plugin.API.GetUser(userId.Value)

	if userErr != nil || cookieErr != nil {
		m.hasError = true
		return
	}

	request.Header.Set(ONLYOFFICE_AUTHORIZATION_USERID_HEADER, user.Id)
	request.Header.Set(ONLYOFFICE_AUTHORIZATION_USERNAME_HEADER, user.Username)

	if m.next != nil {
		m.next.DoFilter(writer, request)
	}
}

func (m *AuthenticationFilter) SetNext(Next Filter) Filter {
	m.next = Next
	return m.next
}

func (m *AuthenticationFilter) HasError() bool {
	if m.next == nil {
		return m.hasError
	}
	return m.hasError || m.next.HasError()
}

func (m *AuthenticationFilter) Reset() {
	m.hasError = false

	if m.next != nil {
		m.next.Reset()
	}
}

//
type FileValidationFilter struct {
	plugin   *Plugin
	next     Filter
	hasError bool
}

func (m *FileValidationFilter) DoFilter(writer http.ResponseWriter, request *http.Request) {
	var fileId string = request.Header.Get(ONLYOFFICE_FILEVALIDATION_FILEID_HEADER)

	if request.Method == "POST" {
		formErr := request.ParseForm()
		if formErr != nil {
			m.hasError = true
			return
		}

		fileId = request.PostForm.Get("fileid")
	}

	fileInfo, fileInfoErr := m.plugin.API.GetFileInfo(fileId)

	if fileInfoErr != nil {
		m.hasError = true
		return
	}

	_, docTypeErr := utils.GetFileType(fileInfo.Extension)

	if docTypeErr != nil {
		m.hasError = true
		return
	}

	request.Header.Set(ONLYOFFICE_FILEVALIDATION_POSTID_HEADER, fileInfo.PostId)

	if m.next != nil {
		m.next.DoFilter(writer, request)
	}
}

func (m *FileValidationFilter) SetNext(Next Filter) Filter {
	m.next = Next
	return m.next
}

func (m *FileValidationFilter) HasError() bool {
	if m.next == nil {
		return m.hasError
	}
	return m.hasError || m.next.HasError()
}

func (m *FileValidationFilter) Reset() {
	m.hasError = false

	if m.next != nil {
		m.next.Reset()
	}
}

//
type PostAuthorizationFilter struct {
	plugin   *Plugin
	next     Filter
	hasError bool
}

func (m *PostAuthorizationFilter) DoFilter(writer http.ResponseWriter, request *http.Request) {
	var userId string = request.Header.Get(ONLYOFFICE_AUTHORIZATION_USERID_HEADER)
	var postId string = request.Header.Get(ONLYOFFICE_FILEVALIDATION_POSTID_HEADER)

	post, postErr := m.plugin.API.GetPost(postId)

	if postErr != nil {
		m.hasError = true
		return
	}

	_, membershipErr := m.plugin.API.GetChannelMember(post.ChannelId, userId)
	if membershipErr != nil {
		m.hasError = true
		return
	}

	request.Header.Set(ONYLOFFICE_CHANNELVALIDATION_CHANNELID_HEADER, post.ChannelId)

	if m.next != nil {
		m.next.DoFilter(writer, request)
	}
}

func (m *PostAuthorizationFilter) SetNext(Next Filter) Filter {
	m.next = Next
	return m.next
}

func (m *PostAuthorizationFilter) HasError() bool {
	if m.next == nil {
		return m.hasError
	}
	return m.hasError || m.next.HasError()
}

func (m *PostAuthorizationFilter) Reset() {
	m.hasError = false

	if m.next != nil {
		m.next.Reset()
	}
}

//
type BodyJwtFilter struct {
	plugin   *Plugin
	next     Filter
	hasError bool
}

func (m *BodyJwtFilter) DoFilter(writer http.ResponseWriter, request *http.Request) {
	if m.plugin.configuration.DESJwt != "" {
		type TokenBody struct {
			Token string `json:"token,omitempty"`
		}

		if request.Body == nil {
			m.hasError = true
			return
		}

		var tokenBody TokenBody
		var bodyBytes []byte

		bodyBytes, _ = ioutil.ReadAll(request.Body)
		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		decodingErr := json.Unmarshal(bodyBytes, &tokenBody)
		if decodingErr != nil {
			m.hasError = true
			return
		}

		if tokenBody.Token == "" {
			m.hasError = true
			return
		}

		_, jwtDecodingErr := security.JwtDecode(tokenBody.Token, []byte(m.plugin.configuration.DESJwt))
		if jwtDecodingErr != nil {
			m.hasError = true
			return
		}
	}

	if m.next != nil {
		m.next.DoFilter(writer, request)
	}
}

func (m *BodyJwtFilter) SetNext(Next Filter) Filter {
	m.next = Next
	return m.next
}

func (m *BodyJwtFilter) HasError() bool {
	if m.next != nil {
		return m.hasError || m.next.HasError()
	}
	return m.hasError
}

func (m *BodyJwtFilter) Reset() {
	m.hasError = false

	if m.next != nil {
		m.next.Reset()
	}
}

//
type HeaderJwtFilter struct {
	plugin   *Plugin
	next     Filter
	hasError bool
}

func (m *HeaderJwtFilter) DoFilter(writer http.ResponseWriter, request *http.Request) {
	if m.plugin.configuration.DESJwt != "" {

		jwtToken := request.Header.Get(m.plugin.configuration.DESJwtHeader)

		if jwtToken == "" {
			m.hasError = true
			return
		}

		jwtToken = strings.Split(jwtToken, m.plugin.configuration.DESJwtPrefix)[1]
		jwtToken = strings.TrimSpace(jwtToken)

		_, jwtErr := security.JwtDecode(jwtToken, []byte(m.plugin.configuration.DESJwt))

		if jwtErr != nil {
			m.hasError = true
			return
		}
	}

	if m.next != nil {
		m.next.DoFilter(writer, request)
	}
}

func (m *HeaderJwtFilter) SetNext(Next Filter) Filter {
	m.next = Next
	return m.next
}

func (m *HeaderJwtFilter) HasError() bool {
	if m.next != nil {
		return m.hasError || m.next.HasError()
	}
	return m.hasError
}

func (m *HeaderJwtFilter) Reset() {
	m.hasError = false

	if m.next != nil {
		m.next.Reset()
	}
}

//
type DecryptorFilter struct {
	plugin   *Plugin
	next     Filter
	hasError bool
}

func (m *DecryptorFilter) DoFilter(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	fileId := query.Get("fileId")

	decipheredFileid, decipherErr := security.EncryptorAES{}.Decrypt(fileId, m.plugin.internalKey)
	_, err := m.plugin.API.GetFileInfo(decipheredFileid)

	if err != nil || decipherErr != nil {
		m.hasError = true
		return
	}

	if m.next != nil {
		m.next.DoFilter(writer, request)
	}
}

func (m *DecryptorFilter) SetNext(Next Filter) Filter {
	m.next = Next
	return m.next
}

func (m *DecryptorFilter) HasError() bool {
	if m.next != nil {
		return m.hasError || m.next.HasError()
	}
	return m.hasError
}

func (m *DecryptorFilter) Reset() {
	m.hasError = false

	if m.next != nil {
		m.next.Reset()
	}
}

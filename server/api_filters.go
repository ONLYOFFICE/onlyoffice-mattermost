package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"security"
	"strings"
	"utils"
)

type Filter interface {
	DoFilter(writer http.ResponseWriter, request *http.Request)
	SetNext(Filter) Filter
	HasError() bool
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

	request.Header.Add(ONLYOFFICE_AUTHORIZATION_USERID_HEADER, user.Id)
	request.Header.Add(ONLYOFFICE_AUTHORIZATION_USERNAME_HEADER, user.Username)

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

//
type FileValidationFilter struct {
	plugin   *Plugin
	next     Filter
	hasError bool
}

func (m *FileValidationFilter) DoFilter(writer http.ResponseWriter, request *http.Request) {
	formErr := request.ParseForm()
	if formErr != nil {
		m.hasError = true
		return
	}

	var fileId string = request.PostForm.Get("fileid")
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

	request.Header.Add(ONLYOFFICE_FILEVALIDATION_POSTID_HEADER, fileInfo.PostId)

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

//
type ChannelAuthorizationFilter struct {
	plugin   *Plugin
	next     Filter
	hasError bool
}

func (m *ChannelAuthorizationFilter) DoFilter(writer http.ResponseWriter, request *http.Request) {
	var userId string = request.Header.Get(ONLYOFFICE_AUTHORIZATION_USERID_HEADER)
	var postId string = request.Header.Get(ONLYOFFICE_FILEVALIDATION_POSTID_HEADER)

	if postId == "" || userId == "" {
		m.hasError = true
		return
	}

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

	if m.next != nil {
		m.next.DoFilter(writer, request)
	}
}

func (m *ChannelAuthorizationFilter) SetNext(Next Filter) Filter {
	m.next = Next
	return m.next
}

func (m *ChannelAuthorizationFilter) HasError() bool {
	if m.next == nil {
		return m.hasError
	}
	return m.hasError || m.next.HasError()
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

//
type DecryptorFilter struct {
	plugin   *Plugin
	next     Filter
	hasError bool
}

func (m *DecryptorFilter) DoFilter(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	fileId := query.Get("fileId")

	var encryptor security.Encryptor = security.EncryptorAES{}
	decipheredFileid, decipherErr := encryptor.Decrypt(fileId, m.plugin.internalKey)
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

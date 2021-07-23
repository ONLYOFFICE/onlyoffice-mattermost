package main

import (
	"bytes"
	"encoding/json"
	"encryptors"
	"io/ioutil"
	"net/http"
	"strings"
	"utils"
)

type Middleware interface {
	Execute(writer http.ResponseWriter, request *http.Request)
	SetNext(Middleware) Middleware
	HasError() bool
}

type AuthenticationMiddleware struct {
	plugin   *Plugin
	next     Middleware
	hasError bool
}

func (m *AuthenticationMiddleware) Execute(writer http.ResponseWriter, request *http.Request) {
	userId, cookieErr := request.Cookie(utils.MMUserCookie)
	_, userErr := m.plugin.API.GetUser(userId.Value)
	if userErr != nil || cookieErr != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		m.hasError = true
		return
	}
	if m.next != nil {
		m.next.Execute(writer, request)
	}
}

func (m *AuthenticationMiddleware) SetNext(Next Middleware) Middleware {
	m.next = Next
	return m.next
}

func (m *AuthenticationMiddleware) HasError() bool {
	if m.next == nil {
		return m.hasError
	}
	return m.hasError || m.next.HasError()
}

type FileValidityMiddleware struct {
	plugin   *Plugin
	next     Middleware
	hasError bool
}

func (m *FileValidityMiddleware) Execute(writer http.ResponseWriter, request *http.Request) {
	formErr := request.ParseForm()
	if formErr != nil {
		http.Error(writer, "Forbidden", http.StatusForbidden)
		m.hasError = true
		return
	}

	var fileId string = request.PostForm.Get("fileid")
	fileInfo, fileInfoErr := m.plugin.API.GetFileInfo(fileId)
	if fileInfoErr != nil {
		http.Error(writer, "Forbidden", http.StatusForbidden)
		m.hasError = true
		return
	}

	_, docTypeErr := utils.GetFileType(fileInfo.Extension)

	if docTypeErr != nil {
		http.Error(writer, docTypeErr.Error(), http.StatusBadRequest)
		m.hasError = true
		return
	}

	request.Header.Add("ONLYOFFICE_POSTID", fileInfo.PostId)
	if m.next != nil {
		m.next.Execute(writer, request)
	}
}

func (m *FileValidityMiddleware) SetNext(Next Middleware) Middleware {
	m.next = Next
	return m.next
}

func (m *FileValidityMiddleware) HasError() bool {
	if m.next == nil {
		return m.hasError
	}
	return m.hasError || m.next.HasError()
}

type ChannelAuthorizationMiddleware struct {
	plugin   *Plugin
	next     Middleware
	hasError bool
}

func (m *ChannelAuthorizationMiddleware) Execute(writer http.ResponseWriter, request *http.Request) {
	userId, _ := request.Cookie(utils.MMUserCookie)
	postId := request.Header.Get("ONLYOFFICE_POSTID")
	if postId == "" || userId.Value == "" {
		http.Error(writer, "Forbidden", http.StatusForbidden)
		m.hasError = true
		return
	}

	post, postErr := m.plugin.API.GetPost(postId)

	if postErr != nil {
		http.Error(writer, "Forbidden", http.StatusForbidden)
		m.hasError = true
		return
	}
	_, membershipErr := m.plugin.API.GetChannelMember(post.ChannelId, userId.Value)
	if membershipErr != nil {
		http.Error(writer, "Forbidden", http.StatusForbidden)
		m.hasError = true
		return
	}
	if m.next != nil {
		m.next.Execute(writer, request)
	}
}

func (m *ChannelAuthorizationMiddleware) SetNext(Next Middleware) Middleware {
	m.next = Next
	return m.next
}

func (m *ChannelAuthorizationMiddleware) HasError() bool {
	if m.next == nil {
		return m.hasError
	}
	return m.hasError || m.next.HasError()
}

type BodyJwtMiddleware struct {
	plugin   *Plugin
	next     Middleware
	hasError bool
}

func (m *BodyJwtMiddleware) Execute(writer http.ResponseWriter, request *http.Request) {
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

		_, jwtDecodingErr := utils.JwtDecode(tokenBody.Token, []byte(m.plugin.configuration.DESJwt))
		if jwtDecodingErr != nil {
			m.hasError = true
			return
		}
	}

	if m.next != nil {
		m.next.Execute(writer, request)
	}
}

func (m *BodyJwtMiddleware) SetNext(Next Middleware) Middleware {
	m.next = Next
	return m.next
}

func (m *BodyJwtMiddleware) HasError() bool {
	if m.next != nil {
		return m.hasError || m.next.HasError()
	}
	return m.hasError
}

type HeaderJwtMiddleware struct {
	plugin   *Plugin
	next     Middleware
	hasError bool
}

func (m *HeaderJwtMiddleware) Execute(writer http.ResponseWriter, request *http.Request) {
	if m.plugin.configuration.DESJwt != "" {

		jwtToken := request.Header.Get(m.plugin.configuration.DESJwtHeader)

		if jwtToken == "" {
			m.hasError = true
			return
		}

		jwtToken = strings.Split(jwtToken, m.plugin.configuration.DESJwtPrefix)[1]
		jwtToken = strings.TrimSpace(jwtToken)

		_, jwtErr := utils.JwtDecode(jwtToken, []byte(m.plugin.configuration.DESJwt))

		if jwtErr != nil {
			m.hasError = true
			return
		}
	}

	if m.next != nil {
		m.next.Execute(writer, request)
	}
}

func (m *HeaderJwtMiddleware) SetNext(Next Middleware) Middleware {
	m.next = Next
	return m.next
}

func (m *HeaderJwtMiddleware) HasError() bool {
	if m.next != nil {
		return m.hasError || m.next.HasError()
	}
	return m.hasError
}

type DecryptorMiddleware struct {
	plugin   *Plugin
	next     Middleware
	hasError bool
}

func (m *DecryptorMiddleware) Execute(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	fileId := query.Get("fileId")

	m.plugin.encryptor = encryptors.EncryptorAES{}
	decipheredFileid, decipherErr := m.plugin.encryptor.Decrypt(fileId, m.plugin.internalKey)
	_, err := m.plugin.API.GetFileInfo(decipheredFileid)

	if err != nil || decipherErr != nil {
		m.hasError = true
		return
	}

	if m.next != nil {
		m.next.Execute(writer, request)
	}
}

func (m *DecryptorMiddleware) SetNext(Next Middleware) Middleware {
	m.next = Next
	return m.next
}

func (m *DecryptorMiddleware) HasError() bool {
	if m.next != nil {
		return m.hasError || m.next.HasError()
	}
	return m.hasError
}

func (p *Plugin) userAccessChain(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	authentication := &AuthenticationMiddleware{plugin: p}
	checkFile := &FileValidityMiddleware{plugin: p}
	channelAccess := &ChannelAuthorizationMiddleware{plugin: p}

	authentication.SetNext(checkFile).SetNext(channelAccess)

	return func(writer http.ResponseWriter, request *http.Request) {
		authentication.Execute(writer, request)

		if authentication.HasError() {
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) callbackChain(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	decryptorMiddleware := DecryptorMiddleware{plugin: p}
	bodyJwtMiddleware := BodyJwtMiddleware{plugin: p}

	decryptorMiddleware.SetNext(&bodyJwtMiddleware)

	return func(writer http.ResponseWriter, request *http.Request) {

		decryptorMiddleware.Execute(writer, request)

		if decryptorMiddleware.HasError() {
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) downloadChain(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	decryptorMiddleware := DecryptorMiddleware{plugin: p}
	headerJwtMiddleware := HeaderJwtMiddleware{plugin: p}

	decryptorMiddleware.SetNext(&headerJwtMiddleware)

	return func(writer http.ResponseWriter, request *http.Request) {

		decryptorMiddleware.Execute(writer, request)

		if decryptorMiddleware.HasError() {
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) permissionsChain(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	authenticationMiddleware := AuthenticationMiddleware{plugin: p}

	return func(writer http.ResponseWriter, request *http.Request) {

		authenticationMiddleware.Execute(writer, request)

		if authenticationMiddleware.HasError() {
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}

		next(writer, request)
	}
}

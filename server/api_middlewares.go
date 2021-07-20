package main

import (
	"encryptors"
	"net/http"
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

func (p *Plugin) userAccessMiddlewareChain(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	authentication := &AuthenticationMiddleware{plugin: p}
	checkFile := &FileValidityMiddleware{plugin: p}
	channelAccess := &ChannelAuthorizationMiddleware{plugin: p}

	authentication.SetNext(checkFile).SetNext(channelAccess)

	return func(writer http.ResponseWriter, request *http.Request) {
		authentication.Execute(writer, request)

		var hasError bool = authentication.HasError()

		if hasError {
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) docServerOnlyMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		fileId := query.Get("fileId")

		p.encryptor = encryptors.EncryptorAES{}
		decipheredFileid, decipherErr := p.encryptor.Decrypt(fileId, p.internalKey)

		_, err := p.API.GetFileInfo(decipheredFileid)
		if err != nil || decipherErr != nil {
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}

		next(writer, request)
	}
}

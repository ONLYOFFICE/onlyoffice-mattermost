package main

import (
	"net/http"
)

func (p *Plugin) authenticationMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		//TODO: Check whether auth token is valid or not
		//TODO: Add channel checks
		userId, cookieErr := request.Cookie("MMUSERID")
		_, userErr := p.API.GetUser(userId.Value)

		if userErr != nil || cookieErr != nil {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) docServerOnlyMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		fileId := query.Get("fileId")
		decipheredFileid, decipherErr := p.decryptAES(fileId, p.internalKey)
		_, err := p.API.GetFileInfo(decipheredFileid)
		if err != nil || decipherErr != nil {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(writer, request)
	}
}

func (p *Plugin) fileAuthorizationMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}

		var fileId string = request.PostForm.Get("fileid")
		userId, _ := request.Cookie("MMUSERID")

		if !p.checkFilePermissions(userId.Value, fileId, &writer) {
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) checkFilePermissions(userId string, fileId string, writer *http.ResponseWriter) bool {
	_, fileErr := p.API.GetFileInfo(fileId)

	if fileErr != nil {
		http.Error(*writer, "Forbidden", http.StatusForbidden)
		return false
	}

	return true
}

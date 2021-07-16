package main

import (
	"encoders"
	"net/http"
)

func (p *Plugin) authenticationMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
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

		p.encoder = encoders.EncoderAES{}
		decipheredFileid, decipherErr := p.encoder.Decode(fileId, p.internalKey)

		_, err := p.API.GetFileInfo(decipheredFileid)
		if err != nil || decipherErr != nil {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(writer, request)
	}
}

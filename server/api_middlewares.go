package main

import "net/http"

func (p *Plugin) userAccessMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	authentication := &AuthenticationFilter{plugin: p}
	checkFile := &FileValidationFilter{plugin: p}
	channelAccess := &ChannelAuthorizationFilter{plugin: p}

	authentication.SetNext(checkFile).SetNext(channelAccess)

	return func(writer http.ResponseWriter, request *http.Request) {
		authentication.DoFilter(writer, request)

		if authentication.HasError() {
			writer.WriteHeader(403)
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) callbackMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	decryptorMiddleware := DecryptorFilter{plugin: p}
	bodyJwtMiddleware := BodyJwtFilter{plugin: p}

	decryptorMiddleware.SetNext(&bodyJwtMiddleware)

	return func(writer http.ResponseWriter, request *http.Request) {

		decryptorMiddleware.DoFilter(writer, request)

		if decryptorMiddleware.HasError() {
			writer.WriteHeader(403)
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) downloadMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	decryptorMiddleware := DecryptorFilter{plugin: p}
	headerJwtMiddleware := HeaderJwtFilter{plugin: p}

	decryptorMiddleware.SetNext(&headerJwtMiddleware)

	return func(writer http.ResponseWriter, request *http.Request) {

		decryptorMiddleware.DoFilter(writer, request)

		if decryptorMiddleware.HasError() {
			writer.WriteHeader(403)
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
			writer.WriteHeader(403)
			return
		}

		next(writer, request)
	}
}

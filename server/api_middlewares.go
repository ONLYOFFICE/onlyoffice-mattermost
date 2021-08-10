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
			writer.WriteHeader(403)
			return
		}

		next(writer, request)
	}
}

func (p *Plugin) callbackMiddleware(next func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	decryptorMiddlewareBody := DecryptorFilter{plugin: p}
	bodyJwtMiddleware := BodyJwtFilter{plugin: p}

	decryptorMiddlewareBody.SetNext(&bodyJwtMiddleware)

	decryptorMiddlewareHeader := DecryptorFilter{plugin: p}
	headerJwtMiddleware := HeaderJwtFilter{plugin: p}

	decryptorMiddlewareHeader.SetNext(&headerJwtMiddleware)

	return func(writer http.ResponseWriter, request *http.Request) {

		decryptorMiddlewareBody.DoFilter(writer, request)
		decryptorMiddlewareHeader.DoFilter(writer, request)

		if decryptorMiddlewareBody.HasError() && decryptorMiddlewareHeader.HasError() {
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

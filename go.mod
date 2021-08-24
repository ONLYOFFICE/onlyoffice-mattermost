module github.com/ONLYOFFICE/onlyoffice-mattermost

go 1.12

require (
	github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/mattermost/mattermost-server/v5 v5.38.1
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	models v0.0.1
	security v0.0.1
	utils v0.0.1
)

replace models => ./server/models

replace security => ./server/security

replace utils => ./server/utils

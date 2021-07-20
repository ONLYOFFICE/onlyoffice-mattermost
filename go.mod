module github.com/ONLYOFFICE/onlyoffice-mattermost

go 1.12

require (
	encryptors v0.0.1
	github.com/gorilla/mux v1.8.0
	github.com/mattermost/mattermost-server/v5 v5.37.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	models v0.0.1
	utils v0.0.1
)

replace models => ./server/models

replace encryptors => ./server/encryptors

replace utils => ./server/utils

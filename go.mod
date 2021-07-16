module github.com/ONLYOFFICE/onlyoffice-mattermost

go 1.12

require (
	dto v0.0.1
	encoders v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/mattermost/mattermost-server/v5 v5.36.1
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	utils v0.0.1
)

replace dto => ./server/dto

replace encoders => ./server/encoders

replace utils => ./server/utils

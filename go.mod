module github.com/ONLYOFFICE/onlyoffice-mattermost

go 1.12

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-hclog v0.16.2 // indirect
	github.com/hashicorp/go-plugin v1.4.2 // indirect
	github.com/hashicorp/yamux v0.0.0-20210707203944-259a57b3608c // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/mattermost/mattermost-server/v5 v5.37.1
	github.com/minio/minio-go/v7 v7.0.12 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/tinylib/msgp v1.1.6 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/zap v1.19.0 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	google.golang.org/genproto v0.0.0-20210809142519-0135a39c2737 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	models v0.0.1
	security v0.0.1
	utils v0.0.1
)

replace models => ./server/models

replace security => ./server/security

replace utils => ./server/utils

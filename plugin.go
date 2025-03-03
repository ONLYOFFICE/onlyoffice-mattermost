package root

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

//go:embed plugin.json
var manifestString string

var Manifest model.Manifest

func init() {
	_ = json.NewDecoder(strings.NewReader(manifestString)).Decode(&Manifest)
}

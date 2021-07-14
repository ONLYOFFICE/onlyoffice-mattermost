package dto

type CallbackBody struct {
	Actions []struct {
		Type   int    `json:"type"`
		UserID string `json:"userid"`
	} `json:"actions"`
	Key    string   `json:"key"`
	Status int      `json:"status"`
	Users  []string `json:"users"`
	Url    string   `json:"url"`
	FileId string   `json:"-"`
}

// CommandRequest and CommandResponse DTOs
const (
	DROP      string = "drop"
	FORCESAVE string = "forcesave"
	INFO      string = "info"
	META      string = "meta"
	VERSION   string = "version"
)

type CommandBody struct {
	Command string `json:"c"`
}

type CommandResponse struct {
	Error      int    `json:"error"`
	Version    string `json:"version,omitempty"`
	Connection bool   `json:"-"`
}

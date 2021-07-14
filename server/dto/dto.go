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

type Document struct {
	FileType string `json:"fileType"`
	Key      string `json:"key"`
	Title    string `json:"title"`
	Url      string `json:"url"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type EditorConfig struct {
	User        User   `json:"user"`
	CallbackUrl string `json:"callbackUrl"`
}

type Config struct {
	Document     Document     `json:"document"`
	DocumentType string       `json:"documentType"`
	EditorConfig EditorConfig `json:"editorConfig"`
}

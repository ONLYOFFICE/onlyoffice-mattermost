package models

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

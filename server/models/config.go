package models

import "github.com/golang-jwt/jwt"

type Config struct {
	Document           Document     `json:"document"`
	DocumentType       string       `json:"documentType"`
	EditorConfig       EditorConfig `json:"editorConfig"`
	Token              string       `json:"token,omitempty"`
	jwt.StandardClaims `json:"-"`
}

type Document struct {
	FileType string      `json:"fileType"`
	Key      string      `json:"key"`
	Title    string      `json:"title"`
	Url      string      `json:"url"`
	P        Permissions `json:"permissions"`
}

type Permissions struct {
	Comment                 bool `json:"comment"`
	Copy                    bool `json:"copy"`
	DeleteCommentAuthorOnly bool `json:"deleteCommentAuthorOnly"`
	Download                bool `json:"download"`
	Edit                    bool `json:"edit"`
	EditCommentAuthorOnly   bool `json:"editCommentAuthorOnly"`
	FillForms               bool `json:"fillForms"`
	ModifyContentControl    bool `json:"modifyContentControl"`
	ModifyFilter            bool `json:"modifyFilter"`
	Print                   bool `json:"print"`
	Review                  bool `json:"review"`
}

type EditorConfig struct {
	User        User   `json:"user"`
	CallbackUrl string `json:"callbackUrl"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

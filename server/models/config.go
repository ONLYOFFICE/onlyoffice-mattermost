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
	Comment                 bool `json:"comment,omitempty"`
	Copy                    bool `json:"copy,omitempty"`
	DeleteCommentAuthorOnly bool `json:"deleteCommentAuthorOnly,omitempty"`
	Download                bool `json:"download,omitempty"`
	Edit                    bool `json:"edit"`
	EditCommentAuthorOnly   bool `json:"editCommentAuthorOnly,omitempty"`
	FillForms               bool `json:"fillForms,omitempty"`
	ModifyContentControl    bool `json:"modifyContentControl,omitempty"`
	ModifyFilter            bool `json:"modifyFilter,omitempty"`
	Print                   bool `json:"print,omitempty"`
	Review                  bool `json:"review,omitempty"`
}

type EditorConfig struct {
	User          User          `json:"user"`
	CallbackUrl   string        `json:"callbackUrl"`
	Customization Customization `json:"customization,omitempty"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Customization struct {
	Goback Goback `json:"goback"`
}

type Goback struct {
	RequestClose bool `json:"requestClose"`
}

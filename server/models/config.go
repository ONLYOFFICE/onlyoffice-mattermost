package models

import "github.com/golang-jwt/jwt"

//TODO: Method to generate template strings
type Config struct {
	Document           Document     `json:"document"`
	DocumentType       string       `json:"documentType"`
	EditorConfig       EditorConfig `json:"editorConfig"`
	Token              string       `json:"token,omitempty"`
	jwt.StandardClaims `json:"-"`
}

type EditorConfig struct {
	User        User   `json:"user"`
	CallbackUrl string `json:"callbackUrl"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Document struct {
	FileType string `json:"fileType"`
	Key      string `json:"key"`
	Title    string `json:"title"`
	Url      string `json:"url"`
}

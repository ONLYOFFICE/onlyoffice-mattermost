package model

import (
	"encoding/json"

	validator "github.com/go-playground/validator/v10"
)

type ConvertFile struct {
	FileID     string `json:"file_id" validate:"required,min=1"`
	OutputType string `json:"output_type,omitempty"`
	Password   string `json:"password,omitempty"`
}

func (c *ConvertFile) Validate() error {
	return validator.New().Struct(c)
}

type ConvertFileResponse struct {
	Error int `json:"error"`
}

func (c *ConvertFileResponse) ToJSON() []byte {
	json, err := json.Marshal(c)
	if err != nil {
		return []byte{}
	}

	return json
}

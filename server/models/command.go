package models

import (
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type CommandBody struct {
	Command            string `json:"c"`
	Token              string `json:"token,omitempty"`
	jwt.StandardClaims `json:"-"`
}

type CommandResponse struct {
	Error        int    `json:"error"`
	Version      string `json:"version,omitempty"`
	IsSuccessful bool   `json:"-"`
}

func (dsr *CommandResponse) ProcessResponse() error {
	if !dsr.IsSuccessful {
		var err error = errors.New("[ONLYOFFICE]: No connection to the Document Service")
		return errors.Wrap(err, err.Error())
	}
	if dsr.Error > 0 {
		var OnlyofficeError error = errors.New("[ONLYOFFICE]: The server responded with an error: " + strconv.Itoa(dsr.Error))
		return errors.Wrap(OnlyofficeError, OnlyofficeError.Error())
	}
	return nil
}

func (dsr *CommandResponse) Succeeded() {
	dsr.IsSuccessful = true
}

func (dsr *CommandResponse) Failed() {
	dsr.IsSuccessful = false
}

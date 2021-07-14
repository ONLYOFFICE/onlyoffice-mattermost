package dto

import (
	"strconv"

	"github.com/pkg/errors"
)

func (dsr *CommandResponse) CheckResponse() error {
	if !dsr.Connection {
		var err error = errors.New("[ONLYOFFICE]: No connection to the Document Service")
		return errors.Wrap(err, err.Error())
	}
	if dsr.Error > 0 {
		var OnlyofficeError error = errors.New("[ONLYOFFICE]: The server responded with an error: " + strconv.Itoa(dsr.Error))
		return errors.Wrap(OnlyofficeError, OnlyofficeError.Error())
	}
	return nil
}

func (dsr *CommandResponse) Connected() {
	dsr.Connection = true
}

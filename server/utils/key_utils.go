package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateKey() string {
	uuidWithHyphen := uuid.New()
	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}

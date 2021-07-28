package utils

import (
	"github.com/pkg/errors"
)

func IsExtensionSupported(fileExt string) bool {
	_, exists := ONLYOFFICE_EXTENSION_TYPE_MAP[fileExt]
	if exists {
		return true
	}
	return false
}

func GetFileType(fileExt string) (string, error) {
	fileType, exists := ONLYOFFICE_EXTENSION_TYPE_MAP[fileExt]
	if !exists {
		return "", errors.New(ONLYOFFICE_LOGGER_PREFIX + "This extension is not supported")
	}
	return fileType, nil
}

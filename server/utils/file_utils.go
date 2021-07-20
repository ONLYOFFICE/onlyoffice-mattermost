package utils

import (
	"github.com/pkg/errors"
)

type extensionsConfig struct {
	Viewed    []string
	Edited    []string
	Converted []string
}

type extensionTypesConfig struct {
	Spreadsheet  []string
	Presentation []string
	Document     []string
}

const DOCUMENT_TYPE string = "word"
const SPREADSHEET_TYPE string = "cell"
const PRESENTATION_TYPE string = "slide"

var ExtensionToType map[string]string = map[string]string{
	"xls":  SPREADSHEET_TYPE,
	"xlsx": SPREADSHEET_TYPE,
	"xlsm": SPREADSHEET_TYPE,
	"xlt":  SPREADSHEET_TYPE,
	"xltx": SPREADSHEET_TYPE,
	"xltm": SPREADSHEET_TYPE,
	"ods":  SPREADSHEET_TYPE,
	"fods": SPREADSHEET_TYPE,
	"ots":  SPREADSHEET_TYPE,
	"csv":  SPREADSHEET_TYPE,
	"pps":  PRESENTATION_TYPE,
	"ppsx": PRESENTATION_TYPE,
	"ppsm": PRESENTATION_TYPE,
	"ppt":  PRESENTATION_TYPE,
	"pptx": PRESENTATION_TYPE,
	"pptm": PRESENTATION_TYPE,
	"pot":  PRESENTATION_TYPE,
	"potx": PRESENTATION_TYPE,
	"potm": PRESENTATION_TYPE,
	"odp":  PRESENTATION_TYPE,
	"fodp": PRESENTATION_TYPE,
	"otp":  PRESENTATION_TYPE,
	"doc":  DOCUMENT_TYPE,
	"docx": DOCUMENT_TYPE,
	"docm": DOCUMENT_TYPE,
	"dot":  DOCUMENT_TYPE,
	"dotx": DOCUMENT_TYPE,
	"dotm": DOCUMENT_TYPE,
	"odt":  DOCUMENT_TYPE,
	"fodt": DOCUMENT_TYPE,
	"ott":  DOCUMENT_TYPE,
	"rtf":  DOCUMENT_TYPE,
	"txt":  DOCUMENT_TYPE,
	"html": DOCUMENT_TYPE,
	"htm":  DOCUMENT_TYPE,
	"mht":  DOCUMENT_TYPE,
	"pdf":  DOCUMENT_TYPE,
	"djvu": DOCUMENT_TYPE,
	"fb2":  DOCUMENT_TYPE,
	"epub": DOCUMENT_TYPE,
	"xps":  DOCUMENT_TYPE,
}

var extensions = extensionsConfig{
	Viewed: []string{".pdf", ".djvu", ".xps"},
	Edited: []string{".docx", ".xlsx", ".csv", ".pptx", ".txt"},
}

func IsExtensionSupported(fileExt string) bool {
	_, exists := ExtensionToType[fileExt]
	if !exists {
		return false
	}
	return true
}

func GetFileType(fileExt string) (string, error) {
	fileType, exists := ExtensionToType[fileExt]
	if !exists {
		return "", errors.Errorf("This extension is not supported")
	}
	return fileType, nil
}

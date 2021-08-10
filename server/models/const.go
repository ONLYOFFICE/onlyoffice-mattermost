package models

const (
	ONLYOFFICE_COMMAND_DROP      string = "drop"
	ONLYOFFICE_COMMAND_FORCESAVE string = "forcesave"
	ONLYOFFICE_COMMAND_INFO      string = "info"
	ONLYOFFICE_COMMAND_META      string = "meta"
	ONLYOFFICE_COMMAND_VERSION   string = "version"
)

var ONLYOFFICE_AUTHOR_PERMISSIONS Permissions = Permissions{
	Comment:  true,
	Copy:     true,
	Download: true,
	Edit:     true,
	Print:    true,
	Review:   true,
}

var ONLYOFFICE_DEFAULT_PERMISSIONS Permissions = Permissions{
	Edit: false,
}

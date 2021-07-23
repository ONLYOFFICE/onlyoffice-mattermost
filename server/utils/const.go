package utils

import "models"

const RC4Key string = "ONLYOFFICE"
const MMUserCookie = "MMUSERID"

const MMPluginApi string = "plugins/com.onlyoffice.mattermost-plugin/onlyofficeapi"

const DESCommandService string = "coauthoring/CommandService.ashx"
const DESApijs string = "web-apps/apps/api/documents/api.js"
const DESConverter string = ""

const ONLYOFFICE_PERMISSIONS_PROP = "ONLYOFFICE_PERMISSIONS"
const ONLYOFFICE_PERMISSIONS_WILDCARD_KEY = "*"

var ONLYOFFICE_AUTHOR_PERMISSIONS models.Permissions = models.Permissions{
	Comment:  true,
	Copy:     true,
	Download: true,
	Edit:     true,
	Print:    true,
	Review:   true,
}

var ONLYOFFICE_ALL_USERS_PERMISSIONS models.Permissions = models.Permissions{}

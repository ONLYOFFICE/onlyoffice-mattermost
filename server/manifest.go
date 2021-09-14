// This file is automatically generated. Do not modify it manually.

package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "com.onlyoffice.mattermost",
  "name": "ONLYOFFICE",
  "description": "This plugin allows users to edit office documents from Mattermost using ONLYOFFICE Docs.",
  "homepage_url": "https://github.com/ONLYOFFICE/onlyoffice-mattermost",
  "support_url": "https://github.com/ONLYOFFICE/onlyoffice-mattermost/issues",
  "version": "1.0.0",
  "release_notes_url": "https://github.com/ONLYOFFICE/onlyoffice-mattermost/releases/latest",
  "icon_path": "assets/logo.svg",
  "min_server_version": "5.37.2",
  "server": {
    "executables": {
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "webapp": {
    "bundle_path": "webapp/dist/main.js"
  },
  "settings_schema": {
    "header": "ONLYOFFICE Docs is an open-source office suite which comprises powerful collaborative editors for text documents, spreadsheets, and presentations highly compatible with OOXML formats.",
    "footer": "Check https://www.onlyoffice.com/office-suite.aspx for more information.",
    "settings": [
      {
        "key": "DESAddress",
        "display_name": "Document Editing Service address: ",
        "type": "text",
        "help_text": "ONLYOFFICE Document Service Location specifies the address of the server with the document services installed. Please change the '\u003cdocumentserver\u003e' for the server address in the below line.",
        "placeholder": "https://\u003cdocumentserver\u003e:\u003cport\u003e/",
        "default": "https://\u003cdocumentserver\u003e:\u003cport\u003e/"
      },
      {
        "key": "DESJwt",
        "display_name": "Secret key (leave blank to disable): ",
        "type": "text",
        "help_text": "Document server JWT secret.",
        "placeholder": "Enter your secret key",
        "default": ""
      },
      {
        "key": "DESJwtHeader",
        "display_name": "JWT Header: ",
        "type": "text",
        "help_text": "",
        "placeholder": "Enter your secret key",
        "default": "Authorization"
      },
      {
        "key": "DESJwtPrefix",
        "display_name": "JWT Prefix: ",
        "type": "text",
        "help_text": "",
        "placeholder": "",
        "default": "Bearer "
      }
    ]
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}

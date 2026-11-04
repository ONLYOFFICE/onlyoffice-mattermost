# Change Log

## 2.4.4
### Changed
- file preview: large labelled "Open in ONLYOFFICE" (primary) and "Share" buttons instead of small icons in the corner

## 2.4.3
### Added
- file details in the preview panel: who uploaded the file and when, who saved it last in ONLYOFFICE and when, the actual file size and the current user's access (can edit / view only, owner)
- `GET /api/fileinfo` endpoint (channel access is required)

### Fixed
- the preview shows the real file size: Mattermost does not update the stored size after ONLYOFFICE saves a new version

## 2.4.2
### Changed
- documents always open in a new browser tab instead of a modal window inside Mattermost; one tab per file (clicking the same file again focuses it), the tab title is the file name, and Close closes the tab

### Fixed
- untranslated "document server unavailable" messages in de/es/fr/it/ru

## 2.4.1
### Added
- files are downloaded only from the configured ONLYOFFICE Docs address (same scheme, host and port), including redirects; this prevents server-side request forgery without blocking Document Servers on private networks
- plugin settings are restored when a configuration save omits them
- unit tests for saving, access checks, URL checks and configuration

### Fixed
- unknown or removed formats in the "Formats" setting are ignored and logged instead of disabling the plugin

## 2.4.0
### Added
- support for ONLYOFFICE Docs 10.x (the version check read "10.x" as "1" and disabled the plugin)
- compatibility with Mattermost 12 (React 19 web app, v12 server SDK)

### Changed
- the minimum ONLYOFFICE Docs version 8.2 is now actually enforced (8.0/8.1 used to pass)
- web app no longer depends on the unmaintained `mattermost-redux` package
- updated dependencies

### Fixed
- edits were silently lost when saving failed: errors are now reported to ONLYOFFICE Docs so that it retries
- a failed download from ONLYOFFICE Docs no longer overwrites the stored file with an error page
- saving to S3 storage without SSL (e.g. MinIO on an internal network)
- TLS certificate verification for S3/Azure storage was always disabled
- users who are not members of a channel could open its files, list its members and send bot notifications
- users could create file posts in channels they cannot post to
- converted files were posted even if downloading the result failed
- crash on save callbacks without users


## 2.3.0
### Added
- configuration flag to toggle health-check notifications.
- live view support

### Changed
- disabled administrator notifications by default.
- modern light/dark themes
- refresh editor on outdated file version
- updated dependencies.

## 2.2.0
### Added
- user mentions via Mattermost bot
- periodic health-checks for document server
- plugins, macros and protection flags in settings

### Changed
- updated copyrights
- updated dependencies

### Fixed
- webapp websocket notifications
- disabled chat by default
- fixed real-time extensions issue

## 2.1.0
### Added
- selectable file extensions for opening documents in the editor

## Changed
- replace `open file location` with `close`

### Fixed
- fixed `shardkey` query parameter in API.js

## 2.0.0
## Added
- support for new document formats
- conversion to OOXML format
- manual XML conversion to document/spreadsheet
- conversion of protected documents
- document creation via new upload option
- demo server support
- onyx and Indigo theme support (including editor support)

## Changed
- refactored server to use fx modules
- replaced formats service with formats submodule
- added empty document templates submodule
- updated document icons
- improved shard key generation for document editors
- enhanced configuration validation
- updated permissions window styles

## Fixed
- general optimization and bug fixes

## 1.2.0
## Added
- support mattermost v9 and v10

## Changed
- fixed styles
- dependency bump for enhanced security

## 1.1.3
## Added
- link to docs cloud

## Changed
- update of vulnerable dependencies

## 1.1.2
## Changed
- update JWT lib

## 1.1.1
## Changed
- fix work with s3 bucket
- editor loader fix
- dependency update

## 1.1.0
## Added
- support mattermost v6 and v7

## 1.0.1
## Added
- disable certificate verification
- warning in the settings about downloading and executing the code of the editors

## Changed
- improved JWT validation

## 1.0.0
## Added
- config page for urls and JWT settings
- file sharing for editing
- coediting docx, xlsx, pptx
- viewing xls, xlsx, csv, xlsm, xlt, xltm, ods, fods, ots, pps, ppsx, ppsm, ppt, pptx, pptm, pot, potx, potm, odp, fodp, otp, doc, docx, docm, dot, dotx, dotm, odt, fodt, ott, rtf

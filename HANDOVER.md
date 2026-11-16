# ONLYOFFICE for Mattermost - Handover

## State at handover

| | |
| :--- | :--- |
| `develop` | **3.0.0**, unreleased, 25 commits ahead of `main` |
| `main` | 2.4.2 |
| Latest release | v2.4.2 (2026-11-05). Before that: v2.4.1 (2026-10-29), v2.4.0 (2026-10-28), v2.3.0 (2026-06-04) |
| Tests | 17 Go packages and 107 Jest tests, all passing on `develop` |
| Open GitHub issues | 12, mostly old install and config questions (#37 to #85). None blocks a release |

---

## 1. Integration type

**An in-host plugin.** It is a Mattermost server plugin (Go, run by Mattermost as an RPC subprocess via
hashicorp/go-plugin) plus a webapp bundle (React 19, registered through `window.registerPlugin`). There is
no separate service and no database of its own.

- **Plugin ID:** `com.onlyoffice.mattermost`. Never change it, or existing installs and settings break.
- **Files:** stay in Mattermost's file store (local, S3 or Azure). The plugin builds its own
  `filestore.FileBackend` from `GetUnsanitizedConfig()`.
- **Per-file permissions:** stored in post props `ONLYOFFICE_PERMISSIONS_<fileId>_<userId>`, where `*` is
  the default for all users.
- **Plugin KV store:**
  - short-lived auth codes (`/code` 120 s, editor 10 s, mentions 24 h)
  - `onlyoffice_demo_start`
  - `onlyoffice_settings_snapshot`
- **Bot account:** `onlyoffice`, created by `EnsureBot`. Its avatar is `assets/logo.png`.

### ONLYOFFICE products used

Only ONLYOFFICE Docs (Document Server), either self-hosted or Docs Cloud.

| Docs API | Used for | Code |
| :--- | :--- | :--- |
| `api.js` editor embed | The editor iframe in `public/editor.html`, with a JWT-signed config | `server/web/common/config.go` |
| Callback | `POST /api/callback`. Status 2 (closed with changes) and 6 (force-save) save the file. Status 1 and 4 are informational | `server/pkg/callback/`, `server/web/controller/callback.go` |
| Command service | `/command` `version` call, run on activation and every 30 min | `server/cmd/plugin.go`, `server/pkg/health/` |
| Conversion service | `/converter`, to OOXML (legacy formats, XML to docx/xlsx) | `server/web/controller/convert.go` |
| JWT | HS256 with a shared secret. Configurable header, which must not be `Authorization` | `server/pkg/crypto/` |
| Demo server | `onlinedocs.docs.onlyoffice.com`, a 7-day trial. The public demo credentials are in `configuration.go`; the same public demo is used by all ONLYOFFICE connectors | `server/pkg/configuration/` |

### HTTP API

All routes live under `/plugins/com.onlyoffice.mattermost/api` (`server/web/router.go`).

| Route | Auth | Purpose |
| :--- | :--- | :--- |
| `POST /callback` | Docs JWT | Save callback |
| `GET /download` | Docs JWT | Docs fetches the file |
| `GET /image` | KV code | Images for the editor |
| `GET /config`, `GET /health` | none | Format catalogue/settings for the webapp, health state |
| `GET /editor`, `GET /editor/config` | Mattermost user or KV code, 5 s timeout | Editor page, refreshed config for an expired key |
| `GET/POST /permissions` | Mattermost user | Read or set per-file rights (setting is author only) |
| `POST /create`, `POST /convert` | Mattermost user | New file from a template, conversion |
| `GET /code` | Mattermost user | Issues a short-lived code for iframes and windows |
| `GET /mentions/users`, `POST /mentions/notify` | Mattermost user or KV code | Comment mentions |

### Supported versions

| Component | Version | Source |
| :--- | :--- | :--- |
| Mattermost | >= 9.11.0 | `plugin.json` `min_server_version` |
| Mattermost in e2e | 11.10.1 by default (a workflow input) | `.github/workflows/e2e.yml` |
| Mattermost 12 | Targeted by 3.0.0: React 19, `mattermost-redux` 11, `@mattermost/types` 11.9. Verified on a local MM 12 docker stand | `CHANGELOG.md` |
| ONLYOFFICE Docs | >= 8.2. Older versions make the plugin disable itself on activation | `validateServerVersion` in `server/cmd/plugin.go` |
| Docs in e2e | `onlyoffice/documentserver:latest` | `e2e.yml` |
| Server platforms | linux amd64/arm64, darwin amd64/arm64, windows amd64 | `plugin.json` |

**Unconfirmed:** there is no written compatibility matrix for Mattermost ESR, Mattermost Cloud, or Docs
Enterprise/Developer. In practice, the e2e run (MM 11.10.1 + Docs latest) is the only guaranteed
combination.

---

## 2. Main features (user scenarios)

1. **Open an attachment.** From the file preview "Open" button or the file menu "Open file in ONLYOFFICE".
   - The editor opens inside Mattermost, with real-time co-editing.
   - Whether a file opens for view or edit depends on three things: the format catalogue (`public/formats`),
     the admin `Formats` allowlist, and the user's permissions.
2. **Save back to Mattermost.** On callback status 2 or 6:
   - The plugin downloads to a temporary file and moves it over the attachment.
   - It bumps `post.UpdateAt`, which rotates the document key.
   - The bot replies in the thread: "File X was updated by @user".
3. **Sharing settings** (post author only, via "Change access rights"):
   - Personal chats: a drop-down.
   - Group chats: default rights plus per-user rights.
   - Changes are announced by the bot in the channel or by DM.
4. **Create a file.** Paperclip, then ONLYOFFICE, then a name and type (document, spreadsheet,
   presentation, and others from `public/template`).
   - The new file is posted to the channel.
   - File names are validated on the client and the server.
5. **Convert.** Legacy or lossy formats (doc, xls, ppt, odt and so on) and XML convert to OOXML. The result
   is posted as a new file.
6. **Mentions in comments.** An @mention in a document comment sends a bot notification with a quote.
   - The notification has an "Open the comment" link.
   - The link uses the `/comment` team route, which opens the editor at that comment.
7. **Admin settings in the System Console:**
   - Docs address, JWT secret, header and prefix
   - Allow private address, demo mode
   - Format allowlist (a custom React setting)
   - Owner-only protection, plugins and macros toggles, health DMs
8. **Health monitoring.**
   - A check runs every 30 minutes and once on start.
   - When Docs is unhealthy, a `health_status` websocket event disables the editor buttons for everyone.
     Sysadmins optionally get a DM.
   - A `config_changed` event refreshes the webapp after settings change.
9. **Archived channels.** Files open read-only, and save callbacks into archived channels are refused
   (new in 3.0.0).
10. **Live viewer and themes.** The editor theme follows the Mattermost light or dark theme.

---

## 3. Build

### Toolchain

| Tool | Version | Pinned in |
| :--- | :--- | :--- |
| Go | 1.26.7 (`go` directive). 1.26.9 also works | `go.mod`, `licenses.yml` |
| Node | 22 in CI. 23 also works locally | workflows |
| golangci-lint | v2.13.2 (config format v2) | workflows, `.golangci.yml` |
| make, git, jq, tar | any | - |

### Commands

```sh
git submodule update --init --recursive   # required: Go embeds public/formats and public/template
make             # check-style + test + dist, produces dist/com.onlyoffice.mattermost-<version>.tar.gz
make dist        # build only: 5 server binaries + webpack + tar
make test        # go test -race ./server/... and jest
make check-style # eslint + tsc + golangci-lint
make deploy      # dist + upload via build/bin/pluginctl
```

`make deploy` needs `MM_SERVICESETTINGS_SITEURL`, plus either `MM_ADMIN_TOKEN` or
`MM_ADMIN_USERNAME`/`MM_ADMIN_PASSWORD`. Mattermost must have plugin uploads enabled.

### Non-obvious points

- **`make dist` moves the submodules to upstream HEAD.** It depends on `submodules`, which runs
  `git submodule update --remote --merge`. Every build can therefore change `public/formats`,
  `webapp/src/public/formats` and `public/template`, which is why the working tree often shows
  `M public/formats`. Treat it like this:
  - Do not commit a submodule bump by accident.
  - Bump deliberately before a release, so the new formats ship.
  - To return to the pinned state: `git submodule update --init --recursive`.
- **The formats catalogue is a submodule twice.**
  - `public/formats` is used by the server and embedded with `go:embed`.
  - `webapp/src/public/formats` is imported by the webapp.
  - Both must point at the same commit, or server and UI disagree about formats.
- **`public/template` tracks branch `main/new`** of `ONLYOFFICE/document-templates`.
- **Packing on macOS.** The Makefile tars with `COPYFILE_DISABLE=1`. When packing by hand, also pass
  `--no-xattrs`, or Mattermost rejects the upload with "Unable to find manifest".
- **Faster dev builds.**
  - With `MM_SERVICESETTINGS_ENABLEDEVELOPER` set, only the host OS/arch binary is built.
  - `MM_DEBUG=1` adds `-gcflags all=-N -l`, for `make attach` (dlv).
- **`build/`** comes from `mattermost-plugin-starter-template` (`pluginctl`, `manifest`, `setup.mk`).
  Do not edit it by hand.
- **The webapp cannot run standalone.** Test UI changes by deploying into a Mattermost instance.
- **The release version is defined twice:** in `plugin.json` `version` and in the first heading of
  `CHANGELOG.md`. The tag workflow reads the changelog and the release workflow reads `plugin.json`, so
  they must match.

### Repositories and access

| What | Where | Notes |
| :--- | :--- | :--- |
| Repository, CI, releases | `github.com/ONLYOFFICE/onlyoffice-mattermost` | GitHub Actions run here, releases are published here, and users file issues here |
| E2E harness | `github.com/ONLYOFFICE/integration-tests` (public) | Reusable workflow `e2e.yml@develop` with `system: mattermost` |
| License checker | `ONLYOFFICE/check-licenses@v1` | Approved exceptions live in `.check-licenses.yml` |
| Formats and templates | `ONLYOFFICE/document-formats`, `ONLYOFFICE/document-templates` (public) | Submodules |

There are no private Go or npm dependencies, and no CI secrets beyond the default `GITHUB_TOKEN`.

### CI (`.github/workflows/`)

| Workflow | Trigger | What it does |
| :--- | :--- | :--- |
| `artifact.yml` | PR to `main`, manual | Runs `make` and uploads the bundle |
| `e2e.yml` | PR to `main`/`develop`, manual | Runs `make dist`, then the `integration-tests` e2e against MM 11.10.1 + Docs latest |
| `licenses.yml` | push/PR to `main`/`develop` | Dependency license check (root and webapp) |
| `create-tag.yml` | push to `main` | Creates tag `v<first version in CHANGELOG.md>` |
| `release.yml` | tag `v*` | Runs `make`, then creates a GitHub release. Notes come from `make docs` (the top CHANGELOG section); the tar.gz is attached |

### Git flow

1. Branch `fix/*` or `feature/*` from `develop`.
2. Open a PR into `develop`.
3. For a release, open a PR from `develop` into `main` (recent ones were merged by Sergey Linnik).
4. The tag and GitHub release are then created automatically.

Commit messages follow Conventional Commits with a scope, for example `fix(server): ...`,
`feat(mentions): ...`, `chore: bumped plugin version`, `docs: added changelog entries`.

---

## 4. Publishing and distribution

- **GitHub Releases is the distribution channel.** The release is automatic from the tag and carries
  `com.onlyoffice.mattermost-<version>.tar.gz`. Admins install it via System Console, Plugin Management,
  Upload Plugin.
- **The plugin is not in the official Mattermost Marketplace.** `mattermost/mattermost-marketplace`
  `plugins.json` was checked on 2026-11-16: 176 plugins, no ONLYOFFICE entry. The Marketplace button in the
  System Console therefore does not offer it. The README text "App Marketplace, find ONLYOFFICE app" is
  misleading, unless the plugin was already uploaded.
  - **Unconfirmed:** whether a Mattermost Marketplace listing was ever requested or promised.
  - Mattermost only lists plugins it reviews and signs. Getting listed means contacting the Mattermost
    partner/integrations team and passing their review.
- **Other places to update on a release** (all manual):
  - the ONLYOFFICE Help Center article linked from `plugin.json`:
    `helpcenter.onlyoffice.com/integration/mattermost.aspx`
  - the onlyoffice.com integrations / app directory page
  - the docs team usually handles both
  - **Unconfirmed:** who the contact on the docs team is.
- **Mattermost plugin signing:** not used. Self-uploaded plugins need no signature unless the server
  enforces `PluginSettings.RequirePluginSignature`. If it does, the customer must add their own key.

---

## 5. Implicit knowledge

- **Never use the `Authorization` header.** Mattermost strips it on plugin routes. The JWT header is
  configurable (default `AuthorizationJWT`), and `IsValid()` rejects `Authorization`. Docs must be
  configured with the same header.
- **The callback trusts only signed claims.** `controller/callback.go` builds the payload from
  `claims.Signed()`; unsigned body fields are ignored, and `?file=` is the only unsigned input.
  - Before 3.0.0, the editor-config JWT that every user's browser receives could be replayed as a callback
    to overwrite any file from any URL (SSRF).
  - Do not simplify this back to decoding the body.
- **Safe save** (`callback/save.go`):
  - The download is written to `<path>.onlyoffice-<id>.tmp`.
  - It is checked against `Content-Length` and `MaxFileSize`, and must not be empty.
  - Only then does `MoveFile` replace the original. A failed download must never touch the original.
- **Document key** = `encode(fileId + post.UpdateAt)`. Every save must bump `UpdateAt`, otherwise editors
  reconnect to a stale key and see "opened from a server backup copy".
- **Plugin RPC turns false bools into nil.** The plugin API config arrives over gob, which drops zero
  values, so a `*bool` set to false arrives as nil.
  - In `provideFileBackend`, `AmazonS3SSL` and `AzureSSL` are pinned before `SetDefaults`, because there nil
    means true. That bug broke S3-without-SSL saves (reported on MM 11.7.2).
  - In plugin code, always read server bools with `model.SafeDereference` (nil means false).
- **API calls can return nil.**
  - `GetUnsanitizedConfig()` returns nil on RPC failure, and `license.Features` can be nil. Both caused a
    panic on MM 11.7, so guard them.
  - Do not replace `GetUnsanitizedConfig()` with `GetConfig()`: S3 secrets are masked there.
- **Settings snapshot** (`cmd/snapshot.go`, KV `onlyoffice_settings_snapshot`). In HA clusters, and on
  partial System Console saves, Mattermost could drop plugin keys. The plugin snapshots known keys and
  restores missing ones in `OnConfigurationChange` and `ConfigurationWillBeSaved`.
- **The fx container is rebuilt on every config change** (`reinitializeContainer`).
  - Handlers capture the `*Configuration` pointer when they are built.
  - Add new dependencies through fx providers, not globals.
- **An invalid config disables the plugin on purpose.** `handleConfigError` calls `DisablePlugin` after
  100 ms when the config is invalid, Docs is unreachable, Docs is older than 8.2, or the JWT is bad. Admins
  see the reason in the System Console and the server log (`[ONLYOFFICE]` prefix).
- **SSRF guard** (`client/http.go`, built on `doyensec/safeurl`):
  - Private and loopback targets, including redirect hops, are blocked unless `DESAllowPrivate` is on.
  - `HTTP_PROXY`, `HTTPS_PROXY` and `NO_PROXY` are honoured. Proxied targets are validated by resolving
    them locally first.
  - All outgoing requests must use this client, never `http.DefaultClient`.
- **Auth for iframes and windows.**
  - Routes normally rely on the `Mattermost-User-Id` header injected by Mattermost.
  - Where the browser cannot send it (editor iframe, images, mentions), a short-lived KV code is passed as
    `?code=`.
- **The `Formats` setting.**
  - Empty means all formats; `none` blocks everything.
  - Unknown formats are dropped and logged, not fatal, because formats get removed upstream.
- **`/editor` and `/editor/config` are wrapped in a 5 s `http.TimeoutHandler`.** Panics in any route are
  recovered and return 500.
- **`EnsureBot` reads `assets/logo.png` from the bundle.** Removing that asset breaks activation.
- **Do not edit `build/`** or the generated `server/manifest` code.

---

## 6. Known issues, tech debt, unmerged branches

### Known issues and debt

- **README is out of date:**
  - "Known issues: S3 saves fail on MM 11.7.2" is fixed in 3.0.0. Remove it when 3.0.0 is released.
  - The demo period is 7 days in code but "one month" in the README.
  - The README says to install from the Marketplace (see section 4).
- **`make dist` pulls submodules from remote** (section 3). It is surprising and makes builds
  non-reproducible. Consider dropping `--remote` from the `dist` dependency chain.
- **Raw KV keys as codes.** The auth middleware and the image handler look up `?code=` as a raw KV key, so
  any existing KV key is accepted as a code. Codes are random and short-lived, but a key prefix (for example
  `code_`) would be safer.
- **View and edit allowlists are not separate.** `IsFormatAllowedForViewing` and
  `IsFormatAllowedForEditing` are identical.
- **`Configuration` is messy.** It has a shallow `Clone()`, an `Error` field mixed into the settings
  struct, and demo constants rewritten on every sanitize.
- **The webapp lint config uses a prerelease:** `@mattermost/eslint-plugin` `1.1.0-0`.
- **Old open GitHub issues** (#37 to #85, mostly install and JWT-header questions) could be answered and
  closed. #37 (Boards integration) and #62 (version history menu) are feature requests.

---

## 7. Plans and commitments

- **Next step: release 3.0.0** from `develop`. Its contents are in `CHANGELOG.md`: mention links, the
  React 19 / MM 12 migration, safe save, security fixes, S3 without SSL, archived channels, file name
  validation, the permalink sharing fix and icon fixes.
  - Before releasing, decide whether to include the pending formats submodule bump.
  - Update the README (section 6).

---

## Testing

### Automated tests

- **Go unit tests** sit next to the code (`*_test.go`). They use `plugintest.API` mocks.
- **`server/web/http_test.go`** is the smoke suite. It builds the full router with the fx container and
  drives every route over HTTP:
  - health/config
  - callback (editing, unauthorized, save persists the file, header token)
  - download round trip
  - auth middleware and `/code`
  - permissions (author only, set)
  - convert, create, image, mentions
- **Jest** (`webapp/src/**/*.test.ts(x)`, Testing Library) covers the components, redux and API client.
- **Browser e2e:** `.github/workflows/e2e.yml` runs on every PR into `develop`/`main`, against real
  Mattermost and Docs containers.

Run them with `make test`, or `go test ./server/web/ -run 'Test(Callback|Download|Create|Convert|SetPermissions)'`
for just the smoke subset.

### Manual smoke checklist

Run this after `make deploy`, against real Mattermost + Docs (for example
`onlyoffice/documentserver` in docker, with `DESAllowPrivate` on for a local address).

1. **Activation**
   - [ ] Valid settings: the plugin enables.
   - [ ] A wrong JWT secret or an unreachable address: the plugin disables itself, and the reason shows in
     the System Console.
   - [ ] Demo mode with empty settings: the plugin works.
2. **Open**
   - [ ] Upload docx, xlsx, pptx and pdf. Preview, then Open, shows the editor.
   - [ ] The author can edit.
   - [ ] Another channel member gets view only by default.
3. **Save**
   - [ ] Edit, then close the editor. About 10 s later, the downloaded file has the change.
   - [ ] The bot posts "File ... was updated by @user" in the thread.
   - [ ] Reopening the file shows no "backup copy" warning.
4. **Co-editing**
   - [ ] Two users with edit rights see each other's changes live.
5. **Sharing**
   - [ ] The author gives a user Edit in a group channel. That user can edit, and the bot DMs them.
   - [ ] A non-author does not see "Change access rights".
6. **Create**
   - [ ] Each file type posts a new file.
   - [ ] Names like `a/b` or `<b>` show a validation error.
7. **Convert**
   - [ ] A `.doc` or `.odt` file converts to a new `.docx` post.
8. **Mentions**
   - [ ] An @user in a document comment sends that user a bot notification.
   - [ ] "Open the comment" opens the editor at that comment.
9. **Formats allowlist**
   - [ ] With only `docx` allowed, an xlsx file has no Open action.
   - [ ] Opening the xlsx editor URL directly returns 403.
10. **Archived channel**
    - [ ] After archiving, the channel's files open read-only.
11. **Health**
    - [ ] Stop Docs and restart the plugin: the editor buttons are disabled.
    - [ ] With health DMs on, sysadmins get a DM.
12. **S3 storage** (if relevant)
    - [ ] Saves work with `AmazonS3SSL=false`.
13. **UI**
    - [ ] The dark theme follows Mattermost.
    - [ ] Sharing settings opened from a permalink show the search field and the Add button.
14. **Upgrade**
    - [ ] Install over the previous release: settings and existing file permissions are preserved.

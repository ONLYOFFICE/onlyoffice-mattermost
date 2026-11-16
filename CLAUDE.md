# ONLYOFFICE Mattermost plugin

A Mattermost plugin (Go server + React webapp) that opens chat attachments in ONLYOFFICE Docs. Plugin ID
`com.onlyoffice.mattermost`. Supports Mattermost >= 9.11 and Docs >= 8.2.

For background (features, decisions, branches, release flow), read `HANDOVER.md`.

## Commands

```sh
git submodule update --init --recursive  # required before any Go build (go:embed of formats/templates)
make                 # check-style + test + dist
make test            # go test -race ./server/... and jest
make check-style     # eslint + tsc + golangci-lint v2
make dist            # bundle: dist/com.onlyoffice.mattermost-<version>.tar.gz
make deploy          # dist + upload (MM_SERVICESETTINGS_SITEURL, MM_ADMIN_TOKEN)
cd webapp && npm run lint && npm run check-types && npx jest <path>
go test ./server/web/ -run TestCallback   # single Go test
```

Toolchain: Go 1.26.7, Node 22, golangci-lint v2.13.2.

The webapp cannot run standalone. Check UI changes by deploying into Mattermost.

## Layout

- `server/cmd/plugin.go`: plugin lifecycle (`OnConfigurationChange` rebuilds the fx container, validates
  Docs, disables the plugin on bad config). `snapshot.go`: HA settings restore.
- `server/web/router.go`: all routes under `/api`. `controller/`: handlers. `common/config.go`: the editor
  config builder. `middleware/auth.go`: `Mattermost-User-Id` header or a KV `?code=`.
- `server/pkg/`:
  - `callback/`: save callbacks, with safe save in `save.go`
  - `client/`: Docs command and convert client, plus the SSRF-safe HTTP client
  - `configuration/`: settings, sanitizing, demo mode, format allowlist
  - `crypto/`: JWT and document key
  - `file/`: post-prop permissions
  - `bot/`: the `onlyoffice` bot
  - `health/`: the 30 min cron
- `webapp/src/`: `index.tsx` registers everything; `components/` holds editor, permissions, manager
  (create), converter, preview, settings and comment.
- `public/`: `editor.html` template; `formats` and `template` are submodules embedded into Go.
- `build/`: from the Mattermost starter template; do not edit.

## Conventions

- Commits: `type(scope): lowercase summary`, for example `fix(server): ...`, `feat(webapp): ...`.
- Branch from `develop`; `main` is release only.
- Every source file starts with the Ascensio System SIA Apache-2.0 header. Copy it from a neighbouring
  file.
- Go imports: stdlib, third party, then `github.com/ONLYOFFICE/onlyoffice-mattermost` (goimports
  local-prefixes).
- Tests sit next to the code. Server tests use `plugintest.API` mocks; `server/web/http_test.go` is the
  route-level smoke suite. Add a test with every fix.
- User-visible webapp strings go through i18n (`webapp/src/i18n`).
- A release bumps `plugin.json` `version` and adds a top `## x.y.z` section to `CHANGELOG.md`. Both must
  match.

## Do not

- Do not use the `Authorization` header for Docs JWT; Mattermost strips it. Use the configured
  `DESJwtHeader`.
- Do not read callback fields from the unsigned body. Use only verified JWT claims (`claims.Signed()`).
- Do not write directly over an attachment. Save via a temp file plus `MoveFile` (`callback/save.go`), and
  bump `post.UpdateAt` after a save.
- Do not use `http.DefaultClient` or a new client for outgoing requests. Use the injected safe client
  (`pkg/client`).
- Do not trust server config bools or pointers from the plugin API:
  - A false `*bool` arrives as nil, so read it with `model.SafeDereference` (nil means false).
  - `GetUnsanitizedConfig()` and `license.Features` can be nil.
- Do not write files to local disk. Use the Mattermost `filestore.FileBackend`.
- Do not commit accidental submodule bumps. `make dist` pulls submodules from remote; reset them with
  `git submodule update --init --recursive`.
- Do not change the plugin ID, KV key names, or the post-prop permission key format; existing installs
  depend on them.
- Never write secrets (JWT secrets, tokens, passwords) into files or commits.

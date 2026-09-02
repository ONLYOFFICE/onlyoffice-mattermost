# ONLYOFFICE Mattermost Integration

## Project Architecture
A Mattermost plugin integrating ONLYOFFICE Docs. Files live in Mattermost's file store (`filestore.FileBackend`), not a local folder. Saves overwrite attachments in place and bump the post's `UpdateAt` timestamp, which rotates the document's unique tracking `key`.

- **Go Server (`server/`)**: HTTP API at `/plugins/com.onlyoffice.mattermost/api/*`, callback handling, JWT management, file I/O, bot alerts, and background health crons.
- **TypeScript Webapp (`webapp/`)**: Mattermost UI (previews, editor launcher, permissions dialog). Built with Webpack and registered via `window.registerPlugin`.

### Environment Specs
- **Mattermost Server**: >= 9.11.0 (`plugin.json`)
- **ONLYOFFICE Document Server**: >= 8.2 (Validated via Command Service `version` call)
- **Deep Reference Docs**: Core specifications reside in `.claude/skills/**` and `.claude/skills/SKILL.md`.

---

## Development Commands

### Compilation & Quality Control
- **Full Release**: `make` (Runs check-style, test, and dist steps)
- **Bundle Production**: `make dist` (Builds submodules, server, webapp, and output zip)
- **Go Binary**: `make server` (Cross-compiles binaries into `server/dist/`)
- **Webapp Compilation**: `make webapp` (Webpack production build to `webapp/dist/main.js`)
- **Linter Checks**: `make check-style` (Runs eslint/tsc on webapp, golangci-lint on server)
- **Test Suites**: `make test` (Runs `go test ./server/...` and Jest tests)
- **Clean Artifacts**: `make clean` (Removes dist folders and node_modules)

### Active Iteration & Live Debugging
- **Hot-Watch**: `MM_DEBUG=1 make watch` (Webpack watch + Go debug flags `-gcflags "all=-N -l"`)
- **Deploy Artifacts**: `make deploy` or `make deploy-from-watch` (Pushes bundle to local server via `pluginctl`)
- **Instance Cycle**: `make reset` (Disables and re-enables the plugin on Mattermost)
- **Attach Debugger**: `make attach` (Attaches `dlv` to the running plugin process on port 2346)
- **Update Formats**: `make submodules` (Pulls default templates and definitions from `ONLYOFFICE/document-formats`)

*Note: Webapp cannot run standalone. For frontend iteration, use `cd webapp && npm run lint` / `npm run check-types`.*

---

## Configuration Layout (`plugin.json`)
Settings deserialize into `server/pkg/configuration.Configuration` during `OnConfigurationChange`. Configuration failures disable the plugin automatically via `API.DisablePlugin`.

| Key | Purpose |
| :--- | :--- |
| `DESAddress` | Document Server base URL (No trailing slash; sanitized on load) |
| `DESJwt` | Shared JWT secret key used for signing payloads and headers |
| `DESJwtHeader` | Custom header name (**Must not** be `Authorization`; blocked by Mattermost) |
| `DESJwtPrefix` | Authorization header prefix string (Defaults to `Bearer `) |
| `DemoEnabled` | Forces 7-day cloud trial override tracked via KV store key `onlyoffice_demo_start` |
| `Formats` | Comma-separated extension allowlist; empty = all; `none` = block all |
| `OwnerProtected` | Restricts editor document protection features exclusively to the post author |

---

## Codebase Routing & Package Map

### HTTP Router (`server/web/router.go` & `controller/`)
Endpoints use the sub-route prefix `/plugins/com.onlyoffice.mattermost/api/*`.

| Handler | Route | Auth Context | Function |
| :--- | :--- | :--- | :--- |
| `CallbackHandler` | `POST /callback` | JWT (Header/Body) | Handles saving states (`1` editing, `2`/`6` save, `4` closed) |
| `DownloadHandler` | `GET /download` | JWT Header | Streams workspace binary files directly to ONLYOFFICE |
| `EditorHandler` | `GET /editor` | Session / Code | Parses `public/editor.html` with signed template configuration |
| `RefreshHandler` | `GET /editor/config` | Session / Code | Recovers and regenerates configurations for expired document keys |
| `PermissionsHandler`| `GET/POST /permissions` | Mattermost User | Evaluates per-user permissions stored in post metadata |
| `CreateHandler` | `POST /create` | Mattermost User | Spawns a new file using templates located in `public/template` |
| `ConvertHandler` | `POST /convert` | Mattermost User | Handles format transformations via Document Server APIs |
| `CodeHandler` | `GET /code` | Mattermost User | Emits short-lived fallback authentication tokens for isolated windows |

*Authentication Context is resolved via `middleware/auth.go`, validating the `Mattermost-User-Id` header or falling back to custom short-lived KV codes.*

### Package Registry (`server/pkg/`)
- `configuration/`: Sanitizes settings and implements format allowlist filters.
- `callback/`: Downloads binary assets from callback URLs and commits saves to `filestore.FileBackend`.
- `crypto/`: Manages JWT claims signing and document token keys via `JwtManager` and `Encoder`.
- `client/`: Interfaces with Document Server Command Services (`/command`) and Conversion APIs (`/converter`).
- `file/`: Performs permission mappings (`ONLYOFFICE_PERMISSIONS_<fileId>_<userId>`) using post props.
- `bot/`: Controls the `onlyoffice` bot account to stream channel notification logs.
- `health/`: Runs background loops (`*/30 * * * *`) tracking core Document Server states.

---

## Model Operational Directives
* **Reference Requirements**: Always parse `.claude/SKILL.md` and context parameters inside nested skills before updating callback handling or editor configuration signing algorithms.
* **Storage Constraint**: Never generate localized disk storage steps. Pass all active payloads directly to Mattermost storage abstraction boundaries via `filestore.FileBackend`.
* **Security Context**: Ensure token authorization utilizes unique client parameters. Never allow the use of standard HTTP `Authorization` headers to prevent data blockage within Mattermost proxies.

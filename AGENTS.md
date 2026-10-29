# Agent guide

## Development and verification

- Go **1.27.2+** is required by `go.mod`. Build with `go build`; run the TUI with `go run .` (both root Go files are needed).
- Before committing Go changes: `go fmt ./... && go vet ./... && go test ./...`. Use `go test -race ./...` for concurrency changes.
- Tests in the root, `ldapapi`, `internal/tabs` and `internal/ui` need no LDAP server. Focused example: `go test -v -run '^TestFormatDNFilter$' ./ldapapi`; form/dialog validation: `go test ./internal/tabs`.
- `scripts/tag.sh 1.2.3` requires a clean worktree and authenticated `gh`; it verifies/builds before pushing `v1.2.3-<10-character hash>` and publishing the latest release (Linux amd64/arm64, Windows amd64, checksums). It replaces `dist/`. Set `RELEASE_NOTES_FILE` to use custom notes instead of generated notes.

## Local LDAP setup

- `scripts/local-test-server.sh start` requires a running Podman/Docker engine **and host `ldapadd`** for seeding. It prefers Podman; use `CONTAINER_APP=docker scripts/local-test-server.sh start` to override. Use `status` to inspect or `stop` to stop/remove the container; no arguments prints usage and exits 1.
- The script starts `goldap-test` and seeds `scripts/0-ous.ldif` then `scripts/1-uids.ldif` only when starting the container. An already-running container is not reseeded.
- Test connection: `ldap://127.0.0.1:8389`, base DN `dc=goldap,dc=sh`, admin DN `cn=admin,dc=goldap,dc=sh`, password `admin123`.
- Create `~/.goldapconfig.json` manually; there is no config-path flag. JSON keys are `LDAP_URL`, `LDAP_BASE_DN`, `LDAP_ADMIN_DN`, with optional `Users`, `Groups`, `OrgUnits` attribute lists (see `README.md`). Missing/invalid config or an empty URL prints an example and exits; it does not create the file.
- `MIN_UID_NUM`/`MIN_GID_NUM` optionally set positive lower bounds for ID suggestions. Without them, use the first gap above the lowest existing ID; empty lists fall back to 1.

## Execution model

- The entrypoint is root `main.go`, with table loading in `load_tabs.go`. `internal/ui.Session` owns one persistent `tea.Program`; screen runners use its `Run` method. A child's `tea.Quit` returns the screen to the outer loop without leaving the alternate screen. Close the session before stderr fallback output.
- All tables load concurrently up front, preserving `TableNames` order. `reload_model` stages fresh tables and an LDAP cache after attempted writes or Ctrl+R; wait for every loader and replace them only on success. Failed reloads retain the existing model/cache; ordinary tab switches reuse loaded data.
- Use Charm **v2** imports (`charm.land/bubbletea/v2`, `bubbles/v2`, `lipgloss/v2`). `View()` returns `tea.View`; screens set `AltScreen` and `MouseModeCellMotion` there and handle `tea.KeyPressMsg`.
- Mouse hit-testing uses rendered text and terminal-cell widths (`internal/mouse`); layout/label changes must keep click targets aligned, including ANSI styling and wide characters.
- Add/confirmation/message runners return errors separately from cancellation. Route failures through `showMessage`; if the dialog itself fails, print the original result to stderr and leave the TUI so it stays visible.

## Coupled data contracts

- Table changes span `ldapapi/constants.go` (`TableNames`, `TableFilters`, column names/widths and attribute-to-column maps) and `GetTableInfo` in `ldapapi/ldap_search.go`. Adding a table also requires reviewing table-name switches in `config.GetLdapAttributes`, `GetNextIdNumber`, `ConstructDnFromUpdates` and `UpdateHasRequiredObjectClass`, plus default/required attributes.
- `ldapapi/tables.go` prepends a 1-based row-number column and keeps `TableInfo.Rows` aligned with `DN`. The Groups “Members” column merges `member` DNs (displayed as uids) and `memberUid` values.
- `ListUsers`/`ListGroups` populate `internal/cache` with uid→DN, gidNumber→DN and next-ID entries. Group writes resolve/validate member uids against this cache, so loading users is a prerequisite for those writes.
- Forms return attribute names plus `map[int]string` updates keyed by indices into that name slice. Keep indices aligned when reordering or adding fields. Configured attribute lists supply new-entry templates **and** edit-field order; `OrderAttributes` inserts missing editable attributes with empty values.
- LDAP operations open/bind/close a connection per call (`defer l.Close()` immediately after a successful dial). Preserve the existing contextual error style, e.g. `"Bind Error; %v"`.
- `AddEntry`/`ModifyAttr` SSHA-hash `userPassword`; existing `{SSHA}` values pass through. Other ordinary attributes use `SplitAttributeValues` for comma-separated input; membership attributes use cache-backed conversion. Modify uses LDAP `Replace`, not append.

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

GoLDAP is a Bubbletea TUI for managing POSIX accounts, groups and OUs on an LDAP server. `AGENTS.md` holds the full style guide (import grouping, error format `"OperationName Error; %v"`, table-driven tests); follow it.

## Commands

```bash
go build                                   # builds ./goldap (gitignored)
go run .                                   # run the TUI
go fmt ./... && go vet ./...               # required before committing
go test ./...                              # only ldapapi has tests
go test -v -run TestFormatDNFilter ./ldapapi   # single test
scripts/local-test-server.sh               # OpenLDAP in podman/docker, seeded from scripts/*.ldif
scripts/tag.sh 1.2.3                       # creates and pushes tag v1.2.3-<short hash>
```

The dev server listens on `ldap://127.0.0.1:8389`, with base DN `dc=goldap,dc=sh`, admin `cn=admin,dc=goldap,dc=sh` and password `admin123`. The app reads its config only from `~/.goldapconfig.json`. If `LDAP_URL` is empty, it prints an example config and exits.

## Architecture

**Each screen is its own `tea.Program`, run one after another from the outer loop in `main.go`.** There is no single root model. The flow:

1. `login.Run()` prompts for the admin password, which becomes `LdapApi.Secret`.
2. `createNewTabsModel` (`load_tabs.go`) calls `api.GetTableInfo(name)` for every entry in `ldapapi.TableNames`, so all tables load up front.
3. `tabs.RunTabs(m)` runs until the user quits or triggers an action. It then returns a `*tabs.State` with `Cmd` (`QuitCmd`/`AddCmd`/`UpdateCmd`/`DeleteCmd`) and `FormInfo` (table name, DN).
4. `main.go` runs `RunAddForm` / `RunUpdateForm` / `RunConfirmBox` and performs the LDAP write. It shows `RunMessageBox` with the result, saves the tab index and per-tab cursor rows, and sets `reload_model`. That clears the cache and refetches every table before the tabs program runs again.

**`ldapapi` (`LdapApi{Config, Secret, Cache}`):** every operation dials, binds as admin, and closes with `defer l.Close()` (see `Search` in `ldap_search.go`). No connection is kept open.
- `constants.go` defines the tables: `TableNames`, `TableFilters`, and per-table `*Cols` / `*Attr` (LDAP attr → column name) / `*ColsWidth`, plus `DefaultAttributes` (new-entry templates) and `RequiredAttributes`. Adding a table means adding entries here **and** a case in `GetTableInfo` (`ldap_search.go`). Cases keyed on the table-name strings also appear in `config.GetLdapAttributes`, `GetNextIdNumber`, `ConstructDnFromUpdates` and `UpdateHasRequiredObjectClass`.
- `tables.go` converts search results into `TableInfo{Cols, Rows, DN}`. It prepends a row-number column, joins multi-valued attributes with `ValueDelimeter` (`", "`), and merges `member` (DNs → uids) and `memberUid` into a single Groups "Members" column.
- `internal/cache` is a string→string map filled as a side effect of `ListUsers`/`ListGroups`. It holds `uid=<uid>` → DN, `gidNumber=<n>` → DN, and `nextUidNumber`/`nextGidNumber`, which are the first gap in existing IDs. Add/Modify depend on it to turn uid lists in `member`/`memberUid` into DNs or to validate them. Those writes fail if a uid is missing from the cache.
- In `AddEntry`/`ModifyAttr`, `userPassword` gets SSHA-hashed (`passwords.go`; values already prefixed `{SSHA}` pass through). Every other value is split on commas (`SplitAttributeValues`). Modify uses `Replace` for each changed attribute.

**Forms (`internal/tabs/form.go`)** return `(attrNames []string, updates map[int]string)`, where the map keys are indices into `attrNames`. In the add flow, initial values come from the config's `Users`/`Groups`/`OrgUnits` lists, or from `DefaultAttributes` if those are absent. `{{attr}}` placeholders are filled live from other fields. The next free uid/gidNumber is injected automatically. After confirm, the form validates required attributes and objectClass. It builds the DN from `ou` (comma list, nested) + `uid`/`cn` + base DN, and for Users it appends `cn = givenName + " " + sn`.

## Roadmap

- `PLAN.md` is the to-do list: `-c` config flag, default password hashing, SSHA256, DN move/rename, sorting.

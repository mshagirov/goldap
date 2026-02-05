# GoLDAP : a Friendly Face for Your LDAP Server

A TUI app for managing LDAP POSIX accounts and groups

<p>
    <img src="./demo.gif" width="100%" alt="goldap demo">
</p>

## Motivation

- Are you tired of using bash scripts or outdated desktop apps for managing your LDAP users?
- You don't want to write a [`ldif`](./scripts/0-ous.ldif) or a long query in the terminal.
- You need a quick way to update your LDAP entries and get on with your life?

Me too. So I wrote `goldap` to do the above using a user-friendly app in a terminal.

## Features

- Browsing users and groups (POSIX accounts and groups in LDAP).
- Modifying entry attributes, e.g., passwords, emails, etc.
    - Updating (moving) DN records not yet supported.
- Adding new entries to POSIX accounts, groups, and OrgUnits.
- Deleting entries.

## Installation

1. Install Go from [webinstall.dev](https://webinstall.dev/golang/)
or [go.dev](https://go.dev/doc/install). Go is only
used during the installation process and not required after building `goldap`.
    ```bash
    curl -sS https://webi.sh/golang | sh; \
    source ~/.config/envman/PATH.env
    ```
1. Use Go to install `goldap`:
    ```bash
    go install github.com/mshagirov/goldap@latest
    ```
1. Start GoLDAP by simply entering:
    ```bash
    goldap
    ```

## Navigation commands and keybindings

|               Keys            |          Command         |
|:-----------------------------:|:------------------------:|
|  `tab`, `n`, `shift-tab`, `p` |     Next/previous tab    |
|  `down` and `up` arrows, or `j` and `k` | Row navigation |
|             `enter`           | Open an entry (row) for viewing and editing   |
|            `ctrl-a`           | Add a new entry to LDAP (current tab)         |
|            `ctrl-d`           | Delete an entry from LDAP (current tab)       |
|         `/` or `?`   | Search (press `enter` to change the focus to the table)|
|    `esc` and `ctrl-c`         | Exit program or cancel search |

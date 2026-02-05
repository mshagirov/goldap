# GoLDAP

A TUI app for managing LDAP POSIX accounts and groups

<p>
    <img src="./demo.gif" width="100%" alt="goldap demo">
</p>

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

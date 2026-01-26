# PLAN.md

- [x] UI for Viewing and Editing Entries: Done
- [x] Saving confirmation UI: "Save changes?"
- [x] Implement viewport for forms TUI
- [ ] Save updates to ldap (ldap API)
  - Modify password (hashed), email, and other non-DN entries
  - Password modification may need a separate step
  - Modify DN: uid, cn, ou, dc etc.
  - (may need to break done update steps for complex DN+non-DN updates)
- [ ] Create new LDAP entry using forms UI backend (UI)
- [ ] Save newly created entry (ldap API)

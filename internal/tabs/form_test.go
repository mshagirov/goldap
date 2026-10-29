package tabs

import (
	"slices"
	"strings"
	"testing"

	"github.com/mshagirov/goldap/internal/config"
	"github.com/mshagirov/goldap/ldapapi"
)

func TestPrepareAddEntry(t *testing.T) {
	cases := []struct {
		name    string
		omit    string
		remove  string
		replace string
		value   string
		wantErr string
	}{
		{name: "valid user"},
		{name: "missing required value", omit: "uid", wantErr: "missing required attribute 'uid'"},
		{name: "empty required value", replace: "gidNumber", value: "", wantErr: "missing required attribute 'gidNumber'"},
		{name: "whitespace required value", replace: "sn", value: "  ", wantErr: "missing required attribute 'sn'"},
		{name: "missing password", omit: "userPassword", wantErr: "missing required attribute 'userPassword'"},
		{name: "missing object class", omit: "objectClass", wantErr: "'objectClass' is missing"},
		{name: "wrong object class", replace: "objectClass", value: "inetOrgPerson", wantErr: "doesn't contain required value 'posixAccount'"},
		{name: "missing DN attribute", remove: "uid", wantErr: "dn entry must include"},
		{name: "missing common name sources", remove: "names", wantErr: "givenName"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			attrs := []string{"ou", "objectClass", "uid", "givenName", "sn", "uidNumber", "gidNumber", "userPassword"}
			values := []string{"People", "posixAccount, inetOrgPerson, top", "alice", "Alice", "Example", "1000", "1000", "test-password"}
			var names []string
			updates := make(map[int]string)
			for i, attr := range attrs {
				if attr == c.remove || (c.remove == "names" && (attr == "givenName" || attr == "sn")) {
					continue
				}
				names = append(names, attr)
				if attr == c.omit {
					continue
				}
				value := values[i]
				if attr == c.replace {
					value = c.value
				}
				updates[len(names)-1] = value
			}

			fi := FormInfo{
				DN:        "dc=goldap,dc=sh",
				TableName: "Users",
				Api: &ldapapi.LdapApi{
					Config: config.Config{LdapBaseDn: "dc=goldap,dc=sh"},
				},
			}
			got, err := prepareAddEntry(&fi, names, updates)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("error = %v, want containing %q", err, c.wantErr)
				}
				if fi.DN != "dc=goldap,dc=sh" {
					t.Errorf("failed validation changed DN to %q", fi.DN)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if fi.DN != "uid=alice,ou=People,dc=goldap,dc=sh" {
				t.Errorf("DN = %q", fi.DN)
			}
			if cn := updates[slices.Index(got, "cn")]; cn != "Alice Example" {
				t.Errorf("cn = %q, want Alice Example", cn)
			}
		})
	}
}

func TestPrepareGroupAndOUEntries(t *testing.T) {
	cases := []struct {
		name    string
		table   string
		names   []string
		updates map[int]string
		wantDN  string
		wantErr string
	}{
		{
			name: "valid group without optional members", table: "Groups",
			names:   []string{"cn", "objectClass", "gidNumber", "member", "memberUid"},
			updates: map[int]string{0: "staff", 1: "top, posixGroup", 2: "1000"},
			wantDN:  "cn=staff,dc=goldap,dc=sh",
		},
		{
			name: "missing group name", table: "Groups",
			names: []string{"cn", "objectClass"}, updates: map[int]string{1: "posixGroup"},
			wantErr: "missing required attribute 'cn'",
		},
		{
			name: "blank group name", table: "Groups",
			names: []string{"cn", "objectClass"}, updates: map[int]string{0: " ", 1: "posixGroup"},
			wantErr: "missing required attribute 'cn'",
		},
		{
			name: "wrong group class", table: "Groups",
			names: []string{"cn", "objectClass"}, updates: map[int]string{0: "staff", 1: "top"},
			wantErr: "required value 'posixGroup'",
		},
		{
			name: "valid OU", table: "OrgUnits",
			names: []string{"ou", "objectClass"}, updates: map[int]string{0: "People", 1: "top, organizationalUnit"},
			wantDN: "ou=People,dc=goldap,dc=sh",
		},
		{
			name: "missing OU name", table: "OrgUnits",
			names: []string{"ou", "objectClass"}, updates: map[int]string{1: "organizationalUnit"},
			wantErr: "missing required attribute 'ou'",
		},
		{
			name: "blank OU name", table: "OrgUnits",
			names: []string{"ou", "objectClass"}, updates: map[int]string{0: "", 1: "organizationalUnit"},
			wantErr: "missing required attribute 'ou'",
		},
		{
			name: "missing OU class", table: "OrgUnits",
			names: []string{"ou", "objectClass"}, updates: map[int]string{0: "People"},
			wantErr: "'objectClass' is missing",
		},
		{
			name: "wrong OU class", table: "OrgUnits",
			names: []string{"ou", "objectClass"}, updates: map[int]string{0: "People", 1: "top"},
			wantErr: "required value 'organizationalUnit'",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fi := FormInfo{
				DN: "dc=goldap,dc=sh", TableName: c.table,
				Api: &ldapapi.LdapApi{Config: config.Config{LdapBaseDn: "dc=goldap,dc=sh"}},
			}
			_, err := prepareAddEntry(&fi, c.names, c.updates)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("error = %v, want containing %q", err, c.wantErr)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if fi.DN != c.wantDN {
					t.Errorf("DN = %q, want %q", fi.DN, c.wantDN)
				}
			}
		})
	}
}

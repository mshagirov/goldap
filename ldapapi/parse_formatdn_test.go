package ldapapi

import (
	"testing"
)

func TestFormatDNFilter(t *testing.T) {
	cases := []struct {
		name      string
		dn        string
		tableName string
		want      string
	}{
		{
			name:      "Valid DN for Users table",
			dn:        "uid=testuser,ou=users,dc=goldap,dc=sh",
			tableName: "Users",
			want:      "(&(objectClass=PosixAccount)(uid=testuser))",
		},
		{
			name:      "Valid DN for Groups table",
			dn:        "cn=testgroup,ou=groups,dc=goldap,dc=sh",
			tableName: "Groups",
			want:      "(&(objectClass=PosixGroup)(cn=testgroup))",
		},
		{
			name:      "Valid DN for OrgUnits table",
			dn:        "ou=testou,dc=goldap,dc=sh",
			tableName: "OrgUnits",
			want:      "(&(objectClass=OrganizationalUnit)(ou=testou))",
		},
		{
			name:      "Valid DN for unknown table",
			dn:        "cn=unknown,dc=goldap,dc=sh",
			tableName: "UnknownTable",
			want:      "(cn=unknown)",
		},
		{
			name:      "DN without comma",
			dn:        "uid=testuser",
			tableName: "Users",
			want:      "",
		},
		{
			name:      "Empty DN",
			dn:        "",
			tableName: "Users",
			want:      "",
		},
		{
			name:      "Complex DN for Users",
			dn:        "uid=salvor.hardin,ou=people,ou=users,dc=goldap,dc=sh",
			tableName: "Users",
			want:      "(&(objectClass=PosixAccount)(uid=salvor.hardin))",
		},
		{
			name:      "DN with special characters",
			dn:        "cn=test+group,ou=groups,dc=goldap,dc=sh",
			tableName: "Groups",
			want:      "(&(objectClass=PosixGroup)(cn=test+group))",
		},
		{
			name:      "Test with trailing comma",
			dn:        "uid=testuser,",
			tableName: "Users",
			want:      "(&(objectClass=PosixAccount)(uid=testuser))",
		},
		{
			name:      "Comma at the beginning",
			dn:        ",uid=testuser",
			tableName: "Users",
			want:      "(&(objectClass=PosixAccount)())",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := FormatDNFilter(c.dn, c.tableName)
			if got != c.want {
				t.Errorf("FormatDNFilter(%q, %q) = %q, want %q", c.dn, c.tableName, got, c.want)
			}
		})
	}
}

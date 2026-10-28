package ldapapi

import (
	"slices"
	"testing"
)

func TestOrderAttributes(t *testing.T) {
	cases := []struct {
		name       string
		names      []string
		values     []string
		order      []string
		wantNames  []string
		wantValues []string
	}{
		{
			name:       "no config keeps entry order",
			names:      []string{"cn", "gidNumber"},
			values:     []string{"admins", "100"},
			wantNames:  []string{"cn", "gidNumber"},
			wantValues: []string{"admins", "100"},
		},
		{
			name:       "config order first, rest as they come",
			names:      []string{"description", "gidNumber", "cn", "extra"},
			values:     []string{"desc", "100", "admins", "x"},
			order:      []string{"cn", "gidNumber"},
			wantNames:  []string{"cn", "gidNumber", "description", "extra"},
			wantValues: []string{"admins", "100", "desc", "x"},
		},
		{
			name:       "missing config attributes are empty fields",
			names:      []string{"cn", "gidNumber", "memberUid"},
			values:     []string{"admins", "100", "asimov"},
			order:      []string{"cn", "gidNumber", "member", "memberUid", "description"},
			wantNames:  []string{"cn", "gidNumber", "member", "memberUid", "description"},
			wantValues: []string{"admins", "100", "", "asimov", ""},
		},
		{
			name:       "case-insensitive match keeps entry's name",
			names:      []string{"memberUid", "cn"},
			values:     []string{"asimov", "admins"},
			order:      []string{"CN", "memberuid"},
			wantNames:  []string{"cn", "memberUid"},
			wantValues: []string{"admins", "asimov"},
		},
		{
			name:       "objectClass and missing DN components skipped",
			names:      []string{"uid", "cn"},
			values:     []string{"asimov", "Isaac Asimov"},
			order:      []string{"ou", "objectClass", "dn", "dc", "uid", "mail"},
			wantNames:  []string{"uid", "mail", "cn"},
			wantValues: []string{"asimov", "", "Isaac Asimov"},
		},
		{
			name:       "present ou is ordered",
			names:      []string{"description", "ou"},
			values:     []string{"people", "People"},
			order:      []string{"ou", "objectClass", "description"},
			wantNames:  []string{"ou", "description"},
			wantValues: []string{"People", "people"},
		},
		{
			name:       "duplicate config names added once",
			names:      []string{"cn"},
			values:     []string{"admins"},
			order:      []string{"cn", "description", "Description", "cn"},
			wantNames:  []string{"cn", "description"},
			wantValues: []string{"admins", ""},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotNames, gotValues := OrderAttributes(c.names, c.values, c.order)
			if !slices.Equal(gotNames, c.wantNames) || !slices.Equal(gotValues, c.wantValues) {
				t.Errorf("OrderAttributes() = %q, %q; want %q, %q", gotNames, gotValues, c.wantNames, c.wantValues)
			}
		})
	}
}

package ldapapi

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

var (
	TableNames = []string{
		"Users",
		"Groups",
		"OrgUnits",
	} // must match cases in *ldapapi.GetTableInfo(s string)

	TableFilters = map[string]string{
		"Users":    "(objectClass=PosixAccount)",
		"Groups":   "(objectClass=PosixGroup)",
		"OrgUnits": "(objectClass=OrganizationalUnit)",
	}

	UsrCols = []string{"Username", "uid", "Name", "Group"}
	UsrAttr = map[string]string{
		"uid":       "Username",
		"uidNumber": "uid",
		"cn":        "Name",
		"gidNumber": "Group",
	}
	UsrColsWidth = []int{15, 5, 20, 25}

	GrpCols = []string{"Name", "gid", "Members", "Description"}
	GrpAttr = map[string]string{
		"cn":          "Name",
		"gidNumber":   "gid",
		"memberUid":   "Members",
		"member":      "Members",
		"description": "Description",
	}
	GrpColsWidth = []int{15, 5, 30, 30}

	OUCols = []string{"Name", "dn", "Description"}
	OUAttr = map[string]string{
		"ou":          "Name",
		"dn":          "dn",
		"description": "Description",
	}
	OUColsWidth = []int{15, 25, 25}
)

func FormatRDNFilter(tableFilter, rdn string) string {
	return fmt.Sprintf("(&%s(%s))", tableFilter, ldap.EscapeFilter(rdn))
}

func FormatDNFilter(dn, tableName string) string {
	rdn, _, found := strings.Cut(dn, ",")
	if !found {
		return ""
	}
	tableFilter, ok := TableFilters[tableName]
	if !ok {
		return fmt.Sprintf("(%s)", ldap.EscapeFilter(rdn))
	}
	return FormatRDNFilter(tableFilter, rdn)
}

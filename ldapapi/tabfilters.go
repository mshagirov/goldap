package ldapapi

import (
	"fmt"
	"strings"
)

const (
	UserFilter  = "(objectClass=PosixAccount)"
	GroupFilter = "(objectClass=PosixGroup)"
	OUsFilter   = "(objectClass=OrganizationalUnit)"
)

var (
	TableNames = []string{
		"Users",
		"Groups",
		"OrgUnits",
	} // must match cases in *ldapapi.GetTableInfo(s string)

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

func FormatDNFilter(dn, tableName string) string {
	// TableNames must match cases in *ldapapi.GetTableInfo(s string)
	entryName, _, found := strings.Cut(dn, ",")
	if !found {
		return ""
	}
	switch tableName {
	case "Users":
		return fmt.Sprintf("(&%s(%s))", UserFilter, entryName)
	case "Groups":
		return fmt.Sprintf("(&%s(%s))", GroupFilter, entryName)
	case "OrgUnits":
		return fmt.Sprintf("(&%s(%s))", OUsFilter, entryName)
	default:
		return fmt.Sprintf("(%s)", entryName)
	}
}

package ldapapi

const (
	ValueDelimeter = ", "
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

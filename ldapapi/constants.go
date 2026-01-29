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

	DefaultFields = map[string][]struct {
		Name string
		Val  []string
	}{
		"Users": {
			{Name: "ou", Val: []string{"People"}},
			{Name: "objectClass", Val: []string{"posixAccount", "inetOrgPerson", "top"}},
			{Name: "employeeType", Val: []string{"Staff"}},
			{Name: "uid", Val: []string{"enter username ..."}},
			{Name: "givenName", Val: []string{"name ..."}},
			{Name: "sn", Val: []string{"surname ..."}},
			{Name: "mail", Val: []string{"user's email ..."}},
			{Name: "uidNumber", Val: []string{"1234"}},
			{Name: "gidNumber", Val: []string{"123"}},
			{Name: "homeDirectory", Val: []string{"path to home folder ..."}},
			{Name: "userPassword", Val: []string{"password"}},
			{Name: "description"},
		},
		"Groups": {
			{Name: "ou"},
			{Name: "objectClass", Val: []string{"top", "posixGroup"}},
			{Name: "cn", Val: []string{"Group's name"}},
			{Name: "gidNumber"},
			{Name: "member"},
			{Name: "memberUid"},
		},
		"OrgUnits": {},
	}

	NonDefaultTabFields = []struct {
		Name string
		Val  []string
	}{
		{Name: "dn", Val: []string{""}},
		{Name: "ou", Val: []string{""}},
		{Name: "cn", Val: []string{""}},
		{Name: "objectClass", Val: []string{"top"}},
		{Name: "description", Val: []string{""}},
	}

	RequiredFields = map[string]map[string]struct{}{
		"Users": {
			"uid":           {},
			"givenName":     {},
			"sn":            {},
			"mail":          {},
			"uidNumber":     {},
			"gidNumber":     {},
			"homeDirectory": {},
			"userPassword":  {},
		},
		"Groups":   {},
		"OrgUnits": {},
	}
)

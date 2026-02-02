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

	DefaultAttributes = map[string][]LdapAttribute{
		"Users": {
			{Name: "ou", Value: []string{"People"}},
			{Name: "objectClass", Value: []string{"posixAccount", "inetOrgPerson", "top"}},
			{Name: "employeeType", Value: []string{"Staff"}},
			{Name: "uid", Value: []string{"USERNAME"}},
			{Name: "givenName", Value: []string{"NAME"}},
			{Name: "sn", Value: []string{"SURNAME"}},
			{Name: "mail", Value: []string{"{{uid}}@goldap.sh"}},
			{Name: "uidNumber", Value: []string{"1234"}},
			{Name: "gidNumber", Value: []string{"Enter group gidNumber, e.g. 100"}},
			{Name: "homeDirectory", Value: []string{"/home/{{uid}}"}},
			{Name: "userPassword", Value: []string{"password"}},
			{Name: "description"},
		},
		"Groups": {
			{Name: "ou"},
			{Name: "objectClass", Value: []string{"top", "posixGroup"}},
			{Name: "cn", Value: []string{"Group's name"}},
			{Name: "gidNumber"},
			{Name: "member", Value: []string{"asimov", "hseldon", "rdolivaw", "..."}},
			{Name: "memberUid", Value: []string{"asimov", "hseldon", "rdolivaw", "..."}},
		},
		"OrgUnits": {
			{Name: "ou"},
			{Name: "objectClass", Value: []string{"top", "organizationalUnit"}},
		},
	}

	UnknownTableAttributes = []LdapAttribute{
		{Name: "dn", Value: []string{""}},
		{Name: "ou", Value: []string{""}},
		{Name: "cn", Value: []string{""}},
		{Name: "objectClass", Value: []string{"top"}},
		{Name: "description", Value: []string{""}},
	}

	// RequiredAttributes members suggestions are erased when updated in Add forms
	RequiredAttributes = map[string]map[string]struct{}{
		"Users": {
			"uid":          {},
			"givenName":    {},
			"sn":           {},
			"gidNumber":    {},
			"userPassword": {},
		},
		"Groups": {
			"cn":        {},
			"member":    {},
			"memberUid": {},
		},
		"OrgUnits": {
			"ou": {},
			"dn": {},
		},
	}
)

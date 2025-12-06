package ldapapi

const (
	UserFilter  = "(objectClass=PosixAccount)"
	GroupFilter = "(objectClass=PosixGroup)"
	OUsFilter   = "(objectClass=OrganizationalUnit)"
)

var (
	UsrCols = []string{"Username", "uid", "Name", "Group"}
	UsrAttr = map[string]string{
		"uid":       "Username",
		"uidNumber": "uid",
		"cn":        "Name",
		"gidNumber": "Group",
	}
	UsrColsWidth = []int{15, 5, 15, 10, 25}

	GrpCols = []string{"Name", "gid", "Description"}
	GrpAttr = map[string]string{
		"cn":          "Name",
		"gidNumber":   "gid",
		"description": "Description",
	}
	GrpColsWidth = []int{15, 5, 25}

	OUCols = []string{"Name", "dn", "Description"}
	OUAttr = map[string]string{
		"ou":          "Name",
		"dn":          "dn",
		"description": "Description",
	}
	OUColsWidth = []int{15, 25, 25}
)

package ldapapi

import (
	"fmt"
	"slices"
	"strings"
)

func (api *LdapApi) ConstructDnFromUpdates(attrNames []string, updates map[int]string, tableName string) (string, error) {
	dn_str := strings.TrimSpace(api.Config.LdapBaseDn)
	if attrs, ok := updates[slices.Index(attrNames, "ou")]; ok {
		attrs = strings.Trim(attrs, ValueDelimeter)
		attrs_slice := strings.Split(attrs, ValueDelimeter)
		slices.Reverse(attrs_slice)
		for _, val := range attrs_slice {
			dn_str = fmt.Sprintf("ou=%s,%s", val, dn_str)
		}
	}

	switch tableName {
	case "Users":
		attrs, ok := updates[slices.Index(attrNames, "uid")]
		if !ok {
			return "", fmt.Errorf("User's dn entry must include \"uid\"!")
		}

		dn_str = fmt.Sprintf("uid=%s,%s", strings.TrimSpace(attrs), dn_str)
	case "Groups":
		attrs, ok := updates[slices.Index(attrNames, "cn")]
		if !ok {
			return "", fmt.Errorf("Group's dn entry must include \"cn\"!")
		}

		dn_str = fmt.Sprintf("cn=%s,%s", strings.TrimSpace(attrs), dn_str)
	}

	return dn_str, nil
}

func (api *LdapApi) AppendCnIfUserForm(attrNames *[]string, updates *map[int]string, tableName string) error {
	if tableName == "Users" {

		givenName, nameOk := (*updates)[slices.Index(*attrNames, "givenName")]
		if !nameOk {
			givenName = ""
		}

		sn, snOk := (*updates)[slices.Index(*attrNames, "sn")]
		if !snOk {
			sn = ""
		}

		if !nameOk && !snOk {
			return fmt.Errorf("Either one or both \"givenName\" and \"sn\" are missing!")
		}

		cn := fmt.Sprintf("%s %s", strings.TrimSpace(givenName), strings.TrimSpace(sn))
		*attrNames = append(*attrNames, "cn")
		(*updates)[len(*attrNames)-1] = cn
	}

	return nil
}

func (api *LdapApi) UpdateHasRequiredObjectClass(attrNames []string, updates map[int]string, tableName string) error {
	attrs, ok := updates[slices.Index(attrNames, "objectClass")]
	if !ok {
		return fmt.Errorf("*LdapApi.UpdateHasRequiredObjectClass: 'objectClass' is missing")
	}
	attrs = strings.Trim(attrs, ValueDelimeter)
	values := strings.Split(attrs, ValueDelimeter)
	required := "?(unknown entry type)"
	switch tableName {
	case "Users":
		required = "posixAccount"
	case "Groups":
		required = "posixGroup"
	case "OrgUnits":
		required = "organizationalUnit"
	}
	for _, val := range values {
		if strings.Contains(val, required) {
			return nil
		}
	}

	return fmt.Errorf("*LdapApi.UpdateHasRequiredObjectClass: %s entry 'objectClass' doesn't contain required value '%s'", tableName, required)
}

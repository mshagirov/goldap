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

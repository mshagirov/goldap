package ldapapi

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func (api *LdapApi) AddEntry(dn string, attr []string, updates map[int]string) error {
	l, err := ldap.DialURL(api.Config.LdapUrl)
	if err != nil {
		return fmt.Errorf("DialURL Error; %v", err)
	}
	defer l.Close()

	if err := l.Bind(api.Config.LdapAdminDn, api.Secret); err != nil {
		return fmt.Errorf("Bind Error; %v", err)
	}

	addReq := ldap.NewAddRequest(dn, []ldap.Control{})

	var values []string

	for id, val := range updates {
		attr_name := attr[id]
		switch strings.ToLower(attr_name) {
		case "userpassword":
			val, err = HashPasswordSSHA(val, 4)
			if err != nil {
				return fmt.Errorf("Error hashing password; %v", err)
			}
			values = []string{val}
		case "member":
			val = strings.Trim(val, ValueDelimeter)
			values, err = api.uidsStringToDnSlice(val)
			if err != nil {
				return err
			}
		case "memberuid":
			val = strings.Trim(val, ValueDelimeter)
			values, err = api.uidsVerifyWithDnCache(val)
			if err != nil {
				return err
			}
		case "dn", "ou", "dc":
			continue
		default:
			values = SplitAttributeValues(val)
		}
		addReq.Attribute(attr_name, values)
	}

	if err = l.Add(addReq); err != nil {
		return fmt.Errorf("Add request error; %v", err)
	}
	return nil
}

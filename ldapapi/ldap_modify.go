package ldapapi

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func (api *LdapApi) ModifyAttr(dn string, attr []string, updates map[int]string) error {
	l, err := ldap.DialURL(api.Config.LdapUrl)
	if err != nil {
		return fmt.Errorf("DialURL Error; %v", err)
	}
	defer l.Close()

	if err := l.Bind(api.Config.LdapAdminDn, api.Secret); err != nil {
		return fmt.Errorf("Bind Error; %v", err)
	}

	modReq := ldap.NewModifyRequest(dn, []ldap.Control{})

	var values []string

	for id, val := range updates {
		attr_name := attr[id]
		switch strings.ToLower(attr_name) {
		case "userpassword":
			val, err = HashPasswordSSHA(val, 4)
			if err != nil {
				return fmt.Errorf("Error hashing password; %v", err)
			}
		case "member":
			val = strings.TrimRight(val, ValueDelimeter)
			values, err = uidsStringToDnSlice(val, api)
			if err != nil {
				return err
			}
		case "memberuid":
			val = strings.TrimRight(val, ValueDelimeter)
			values, err = uidsVerifyWithDnCache(val, api)
			if err != nil {
				return err
			}
		default:
			val = strings.TrimRight(val, ValueDelimeter)
			values = strings.Split(val, ValueDelimeter)
		}
		modReq.Replace(attr_name, values)
	}

	if err = l.Modify(modReq); err != nil {
		return fmt.Errorf("Modify request error; %v", err)
	}
	return nil
}

func uidsStringToDnSlice(cleanValueString string, api *LdapApi) ([]string, error) {
	values := strings.Split(cleanValueString, ValueDelimeter)

	for k := range values {
		dn, found := api.Cache.Get(fmt.Sprintf("uid=%s", values[k]))
		if found {
			values[k] = dn
		} else {
			return []string{}, fmt.Errorf("Group member update: uid=%s not found", values[k])
		}
	}
	return values, nil
}

func uidsVerifyWithDnCache(cleanValueString string, api *LdapApi) ([]string, error) {
	values := strings.Split(cleanValueString, ValueDelimeter)
	for _, uid := range values {
		_, found := api.Cache.Get(fmt.Sprintf("uid=%s", uid))
		if !found {
			return []string{}, fmt.Errorf("Group memberUid update: uid=%s not found", uid)
		}
	}
	return values, nil
}

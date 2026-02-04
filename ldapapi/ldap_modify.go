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
				return fmt.Errorf("SSHA error: %v", err)
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
		default:
			val = strings.Trim(val, ValueDelimeter)
			values = strings.Split(val, ValueDelimeter)
		}
		modReq.Replace(attr_name, values)
	}

	if err = l.Modify(modReq); err != nil {
		return err
	}
	return nil
}

func (api *LdapApi) uidsStringToDnSlice(cleanValueString string) ([]string, error) {
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

func (api *LdapApi) uidsVerifyWithDnCache(cleanValueString string) ([]string, error) {
	values := strings.Split(cleanValueString, ValueDelimeter)
	for _, uid := range values {
		_, found := api.Cache.Get(fmt.Sprintf("uid=%s", uid))
		if !found {
			return []string{}, fmt.Errorf("Group memberUid update: uid=%s not found", uid)
		}
	}
	return values, nil
}

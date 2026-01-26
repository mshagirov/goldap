package ldapapi

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func (api *LdapApi) ModifyPassword(dn, newPassword string) error {
	l, err := ldap.DialURL(api.Config.LdapUrl)
	if err != nil {
		return fmt.Errorf("DialURL Error; %v", err)
	}
	defer l.Close()

	if err := l.Bind(api.Config.LdapAdminDn, api.Secret); err != nil {
		return fmt.Errorf("Bind Error; %v", err)
	}

	passwordModifyRequest := ldap.NewPasswordModifyRequest(dn, "", newPassword)
	if _, err = l.PasswordModify(passwordModifyRequest); err != nil {
		return fmt.Errorf("Password could not be changed: %s", err.Error())
	}
	return nil
}

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
	for id, val := range updates {
		attr_name := attr[id]
		if strings.ToLower(attr_name) == "userpassword" {
			val, err = HashPasswordSSHA(val, 4)
			if err != nil {
				return fmt.Errorf("Error hashing password; %v", err)
			}
		}
		values := strings.Split(val, ValueDelimeter)
		modReq.Replace(attr_name, values)
	}

	if err = l.Modify(modReq); err != nil {
		return fmt.Errorf("Modify request error; %v", err)
	}
	return nil
}

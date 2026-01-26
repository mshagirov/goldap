package ldapapi

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

func (api *LdapApi) ModifyPassword(dn, newPassword string) (*ldap.PasswordModifyResult, error) {
	l, err := ldap.DialURL(api.Config.LdapUrl)
	if err != nil {
		return nil, fmt.Errorf("DialURL Error; %v", err)
	}
	defer l.Close()

	if err := l.Bind(api.Config.LdapAdminDn, api.Secret); err != nil {
		return nil, fmt.Errorf("Bind Error; %v", err)
	}

	passwordModifyRequest := ldap.NewPasswordModifyRequest(dn, "", newPassword)
	res, err := l.PasswordModify(passwordModifyRequest)
	if err != nil {
		return nil, fmt.Errorf("Password could not be changed: %s", err.Error())
	}
	return res, nil
}

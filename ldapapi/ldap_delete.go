package ldapapi

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

func (api *LdapApi) DeleteEntry(dn string) error {
	l, err := ldap.DialURL(api.Config.LdapUrl)
	if err != nil {
		return fmt.Errorf("DialURL Error; %v", err)
	}
	defer l.Close()

	if err := l.Bind(api.Config.LdapAdminDn, api.Secret); err != nil {
		return fmt.Errorf("Bind Error; %v", err)
	}

	req := ldap.NewDelRequest(dn, []ldap.Control{})
	if err := l.Del(req); err != nil {
		return err
	}
	return nil
}

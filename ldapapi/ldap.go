package ldapapi

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/mshagirov/goldap/internal/config"
)

type LdapApi struct {
	Config config.Config
	Secret string
}

func (api *LdapApi) TryConnecting() error {
	l, err := ldap.DialURL(api.Config.LdapUrl)
	if err != nil {
		return fmt.Errorf("DialURL Error; %v", err)
	}
	defer l.Close()
	return nil
}

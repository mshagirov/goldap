package config

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type LdapApi struct {
	Config *Config
	Secret string
}

func (api *LdapApi) Search(filter string) (*ldap.SearchResult, error) {
	l, err := ldap.DialURL(api.Config.LdapUrl)
	if err != nil {
		return nil, fmt.Errorf("DialURL Error; %v", err)
	}
	defer l.Close()

	if err := l.Bind(api.Config.LdapAdminDn, api.Secret); err != nil {
		return nil, fmt.Errorf("Bind Error; %v", err)
	}

	searchRequest := ldap.NewSearchRequest(
		api.Config.LdapBaseDn,
		ldap.ScopeWholeSubtree,
		0, 0, 0, false,
		filter,
		[]string{},
		nil,
	)

	res, err := l.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("Search Error; %v", err)
	}
	return res, err
}

func (api *LdapApi) ListUsers() (*ldap.SearchResult, error) {
	return api.Search(UserFilter)
}

func (api *LdapApi) ListGroups() (*ldap.SearchResult, error) {
	return api.Search(GroupFilter)
}

func (api *LdapApi) Users() [2]TableInfo {
	var t [2]TableInfo
	usrRes, err := api.ListUsers()
	if err == nil {
		LoadTableInfoFromSearchResults(&t[0], UsrCols, UsrAttr, UsrColsWidth, usrRes)
	}

	grpRes, err := api.ListGroups()
	if err == nil {
		LoadTableInfoFromSearchResults(&t[1], GrpCols, GrpAttr, GrpColsWidth, grpRes)
	}
	return t
}

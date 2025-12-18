package ldapapi

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/mshagirov/goldap/internal/config"
)

type LdapApi struct {
	Config *config.Config
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

func (api *LdapApi) FindUser(dn string) (*ldap.SearchResult, error) {
	return api.Search(fmt.Sprintf(UserFilterTemplate, dn))
}

func (api *LdapApi) ListGroups() (*ldap.SearchResult, error) {
	return api.Search(GroupFilter)
}

func (api *LdapApi) ListOUs() (*ldap.SearchResult, error) {
	return api.Search(OUsFilter)
}

func (api *LdapApi) GetTableInfo(s string) (TableInfo, error) {
	var t TableInfo
	switch s {
	case "Users":
		usrRes, err := api.ListUsers()
		if err != nil {
			return t, err
		}

		LoadTableInfoFromSearchResults(&t, UsrCols, UsrAttr, UsrColsWidth, usrRes)
		return t, nil
	case "Groups":
		grpRes, err := api.ListGroups()
		if err != nil {
			return t, err
		}

		LoadTableInfoFromSearchResults(&t, GrpCols, GrpAttr, GrpColsWidth, grpRes)
		return t, nil
	case "OrgUnits":
		ouRes, err := api.ListOUs()
		if err != nil {
			return t, err
		}

		LoadTableInfoFromSearchResults(&t, OUCols, OUAttr, OUColsWidth, ouRes)
		return t, nil
	default:
		return t, fmt.Errorf("LdapApi.GetTableInfo: the input '%v' value not recognised", s)
	}
}

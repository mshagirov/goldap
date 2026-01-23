package ldapapi

import (
	"fmt"
	"strings"

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
	return api.Search(TableFilters["Users"])
}

func (api *LdapApi) ListGroups() (*ldap.SearchResult, error) {
	return api.Search(TableFilters["Groups"])
}

func (api *LdapApi) ListOUs() (*ldap.SearchResult, error) {
	return api.Search(TableFilters["OrgUnits"])
}

func (api *LdapApi) GetTableInfo(tableName string) (TableInfo, error) {
	var t TableInfo
	switch tableName {
	case "Users":
		usrRes, err := api.ListUsers()
		if err != nil {
			return t, err
		}
		LoadTableInfoFromSearchResults(&t, tableName, UsrCols, UsrAttr, UsrColsWidth, usrRes)
		return t, nil
	case "Groups":
		grpRes, err := api.ListGroups()
		if err != nil {
			return t, err
		}
		LoadTableInfoFromSearchResults(&t, tableName, GrpCols, GrpAttr, GrpColsWidth, grpRes)
		return t, nil
	case "OrgUnits":
		ouRes, err := api.ListOUs()
		if err != nil {
			return t, err
		}
		LoadTableInfoFromSearchResults(&t, tableName, OUCols, OUAttr, OUColsWidth, ouRes)
		return t, nil
	default:
		return t, fmt.Errorf("LdapApi.GetTableInfo: the input '%v' value not recognised", tableName)
	}
}

func (api *LdapApi) SearchDN(dn, tableName string) (*ldap.SearchResult, error) {
	filter := FormatDNFilter(dn, tableName)
	if len(filter) == 0 {
		return nil, fmt.Errorf("Error formatting: %s", dn)
	}
	return api.Search(filter)
}

func (api *LdapApi) GetAttrWithDN(dn, tableName string) ([]string, []string) {
	var (
		attrNames []string
		attrVals  []string
	)
	sr, err := api.SearchDN(dn, tableName)
	if err != nil || len(sr.Entries) == 0 {
		return attrNames, attrVals
	}
	e := sr.Entries[0]
	for _, a := range e.Attributes {
		attrNames = append(attrNames, a.Name)
		vals := a.Values
		if strings.ToLower(a.Name) == "member" {
			if _, uidVals, ok := GetFirstDnAttrs(vals); ok {
				vals = uidVals
			}
		}
		if len(vals) > 1 {
			attrVals = append(attrVals, strings.Join(vals, ", "))
		} else {
			attrVals = append(attrVals, vals[0])
		}
	}

	return attrNames, attrVals
}

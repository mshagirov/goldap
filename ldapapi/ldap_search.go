package ldapapi

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

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

func getFirstMissingInt(nums []int) int {
	slices.Sort(nums)
	missing := nums[0] + 1 // candidate
	for _, n := range nums {
		if n < missing {
			continue
		}
		if n == missing {
			missing++
		}
		if n > missing {
			return missing
		}
	}
	return missing
}

func (api *LdapApi) ListUsers() (*ldap.SearchResult, error) {
	r, err := api.Search(TableFilters["Users"])
	if err != nil {
		return r, err
	}

	var uid string
	uidNumbers := make([]int, len(r.Entries))
	for i, entry := range r.Entries {
		uid = entry.GetAttributeValue("uid")
		uidNumber := strings.TrimSpace(entry.GetAttributeValue("uidNumber"))
		uidNumbers[i], _ = strconv.Atoi(uidNumber)
		api.Cache.Add(fmt.Sprintf("uid=%v", uid), entry.DN)
	}

	api.Cache.Add("nextUidNumber", strconv.Itoa(getFirstMissingInt(uidNumbers)))

	return r, err
}

func (api *LdapApi) ListGroups() (*ldap.SearchResult, error) {
	r, err := api.Search(TableFilters["Groups"])
	if err != nil {
		return r, err
	}

	var gidNumber string
	gidNumbers := make([]int, len(r.Entries))
	for i, entry := range r.Entries {
		gidNumber = strings.TrimSpace(entry.GetAttributeValue("gidNumber"))
		gidNumbers[i], _ = strconv.Atoi(gidNumber)
		api.Cache.Add(fmt.Sprintf("gidNumber=%v", gidNumber), entry.DN)
	}
	nextAvailableId := strconv.Itoa(getFirstMissingInt(gidNumbers))
	api.Cache.Add("nextGidNumber", nextAvailableId)
	return r, err
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
		if strings.ToLower(a.Name) == "objectclass" {
			continue
		}
		attrNames = append(attrNames, a.Name)
		vals := a.Values
		if strings.ToLower(a.Name) == "member" {
			if _, uidVals, ok := GetFirstDnAttrs(vals); ok {
				vals = uidVals
			}
		}
		if len(vals) > 1 {
			attrVals = append(attrVals, strings.Join(vals, ValueDelimeter))
		} else {
			attrVals = append(attrVals, vals[0])
		}
	}

	return attrNames, attrVals
}

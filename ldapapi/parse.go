package ldapapi

import (
	"fmt"
	"slices"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func GetFirstDnAttr(dn string) (attr, value string, ok bool) {
	firstDN, _, found := strings.Cut(dn, ",")
	if !found {
		return "", "", false
	}

	attr, value, found = strings.Cut(firstDN, "=")
	if !found || attr == "" || value == "" {
		return "", "", false
	}

	return attr, value, true
}

func GetFirstDnAttrs(dns []string) (attrs, values []string, ok bool) {
	attrs = make([]string, len(dns))
	values = make([]string, len(dns))
	for i, dn := range dns {
		attrs[i], values[i], ok = GetFirstDnAttr(dn)
		if !ok {
			return []string{}, []string{}, false
		}
	}
	return attrs, values, true
}

func GetUsersColId(attr string) int {
	colName, ok := UsrAttr[attr]
	gidColId := -1
	if ok {
		gidColId = slices.Index(UsrCols, colName)
	}
	return gidColId
}

func FormatRDNFilter(tableFilter, rdn string) string {
	return fmt.Sprintf("(&%s(%s))", tableFilter, ldap.EscapeFilter(rdn))
}

func FormatDNFilter(dn, tableName string) string {
	rdn, _, found := strings.Cut(dn, ",")
	if !found {
		return ""
	}
	tableFilter, ok := TableFilters[tableName]
	if !ok {
		return fmt.Sprintf("(%s)", ldap.EscapeFilter(rdn))
	}
	return FormatRDNFilter(tableFilter, rdn)
}

package ldapapi

import (
	"strings"
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

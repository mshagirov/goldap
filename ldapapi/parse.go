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

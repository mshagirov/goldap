package ldapapi

import (
	"slices"
	"strings"
)

type LdapAttribute struct {
	Name  string
	Value []string
}

func GetDefaultAttributes(t string) ([]string, []string, error) {
	defaultAttr, ok := DefaultAttributes[t]
	if !ok {
		defaultAttr = UnknownTableAttributes
	}

	names := make([]string, len(defaultAttr))
	values := make([]string, len(defaultAttr))
	for i := range defaultAttr {
		names[i] = defaultAttr[i].Name
		values[i] = strings.Join(defaultAttr[i].Value, ValueDelimeter)
	}

	return names, values, nil
}

func GetRequiredAttributesSet(names []string, t string) map[int]struct{} {
	requiredAttrSet := map[int]struct{}{}
	if requiredAttr, ok := RequiredAttributes[t]; ok {
		for a := range requiredAttr {
			requiredAttrSet[slices.Index(names, a)] = struct{}{}
		}
	}
	return requiredAttrSet
}

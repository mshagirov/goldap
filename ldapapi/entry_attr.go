package ldapapi

import (
	"slices"
	"strings"
)

type LdapAttribute struct {
	Name  string
	Value []string
}

func (api *LdapApi) GetNewEntryAttributes(tableName string) ([]string, []string, error) {
	attrs, ok := api.Config.GetLdapAttributes(tableName)
	if !ok || len(attrs) < 1 {
		return GetDefaultAttributes(tableName)
	}

	names := make([]string, len(attrs))
	values := make([]string, len(attrs))
	for i := range attrs {
		names[i] = attrs[i].Name
		values[i] = strings.Join(attrs[i].Value, ValueDelimeter)
	}

	return names, values, nil
}

func GetDefaultAttributes(tableName string) ([]string, []string, error) {
	defaultAttr, ok := DefaultAttributes[tableName]
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

func GetRequiredAttributesSet(names []string, tableName string) map[int]struct{} {
	requiredAttrSet := map[int]struct{}{}
	if requiredAttr, ok := RequiredAttributes[tableName]; ok {
		for a := range requiredAttr {
			requiredAttrSet[slices.Index(names, a)] = struct{}{}
		}
	}
	return requiredAttrSet
}

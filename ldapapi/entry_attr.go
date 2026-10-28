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

// OrderAttributes orders an entry's attributes by the configured attribute
// names (matched case-insensitively), then appends the remaining entry
// attributes in their original order. Configured attributes missing from the
// entry are added with empty values so they can be set later; objectClass and
// DN components (dn, ou, dc) are not added since they're not editable here.
func OrderAttributes(names, values, order []string) ([]string, []string) {
	orderedNames := make([]string, 0, len(names)+len(order))
	orderedValues := make([]string, 0, len(names)+len(order))
	used := make([]bool, len(names))

	for _, attr := range order {
		if strings.EqualFold(attr, "objectClass") ||
			slices.ContainsFunc(orderedNames, func(n string) bool { return strings.EqualFold(n, attr) }) {
			continue
		}
		i := slices.IndexFunc(names, func(n string) bool { return strings.EqualFold(n, attr) })
		switch {
		case i > -1:
			orderedNames = append(orderedNames, names[i])
			orderedValues = append(orderedValues, values[i])
			used[i] = true
		case !slices.Contains([]string{"dn", "ou", "dc"}, strings.ToLower(attr)):
			orderedNames = append(orderedNames, attr)
			orderedValues = append(orderedValues, "")
		}
	}

	for i := range names {
		if !used[i] {
			orderedNames = append(orderedNames, names[i])
			orderedValues = append(orderedValues, values[i])
		}
	}
	return orderedNames, orderedValues
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

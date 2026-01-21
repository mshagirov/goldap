package ldapapi

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/go-ldap/ldap/v3"
)

type TableInfo struct {
	Cols []table.Column
	Rows []table.Row
	DN   []string
}

func ToIdMap(colNames []string) map[string]int {
	idMap := make(map[string]int)
	for id, a := range colNames {
		idMap[a] = id
	}
	return idMap
}

func MakeColumns(names []string, widths []int) []table.Column {
	cols := []table.Column{}
	for id, n := range names {
		cols = append(cols, table.Column{Title: n, Width: widths[id]})
	}
	return cols
}

func processMemberValues(attrs []*ldap.EntryAttribute, colAttr map[string]string) string {
	var allValues []string

	for _, attr := range attrs {
		columnName, exists := colAttr[attr.Name]
		if !exists || columnName != "Members" {
			continue
		}

		switch attr.Name {
		case "member":
			// Parse DNs and extract values
			_, extractedValues, ok := GetFirstDnAttrs(attr.Values)
			if ok {
				allValues = append(allValues, extractedValues...)
			}
		case "memberUid":
			allValues = append(allValues, attr.Values...)
		}
	}

	// Sort alphanumerically
	slices.Sort(allValues)

	// Join with comma and space
	return strings.Join(allValues, ", ")
}

func needsMultiAttrMerging(tableName string, colAttr map[string]string) bool {
	return tableName == "Groups" && hasMultipleAttrsForColumn(colAttr, "Members")
}

func hasMultipleAttrsForColumn(colAttr map[string]string, columnName string) bool {
	count := 0
	for _, col := range colAttr {
		if col == columnName {
			count++
			if count > 1 {
				return true
			}
		}
	}
	return false
}

func LoadTableInfoFromSearchResults(
	ti *TableInfo,
	tableName string,
	colNames []string,
	colAtrr map[string]string,
	widths []int,
	sr *ldap.SearchResult,
) {
	colIds := ToIdMap(colNames)
	ti.Cols = MakeColumns(append([]string{""}, colNames...), append([]int{4}, widths...))
	ti.Rows = []table.Row{}
	ti.DN = []string{}

	for i, entry := range sr.Entries {
		row_i := make([]string, len(colNames)+1)
		row_i[0] = fmt.Sprintf("%v", i+1)

		if val, ok := colAtrr["dn"]; ok {
			row_i[colIds[val]+1] = entry.DN
		}

		if tableName == "Groups" && hasMultipleAttrsForColumn(colAtrr, "Members") {
			// Groups table might use member or memberUid attribute
			memberValues := processMemberValues(entry.Attributes, colAtrr)
			if memberValues != "" {
				row_i[colIds["Members"]+1] = memberValues
			}

			// Process other attributes
			for _, attr := range entry.Attributes {
				if attr.Name == "memberUid" || attr.Name == "member" {
					continue
				}

				if columnName, ok := colAtrr[attr.Name]; ok {
					id := colIds[columnName]
					if len(attr.Values) > 1 {
						row_i[id+1] = strings.Join(attr.Values, ", ")
					} else if len(attr.Values) > 0 {
						row_i[id+1] = attr.Values[0]
					}
				}
			}
		} else {
			// Use original logic for non-Groups tables or when no merging needed
			for _, attr := range entry.Attributes {
				if columnName, ok := colAtrr[attr.Name]; ok {
					id := colIds[columnName]
					if len(attr.Values) > 1 {
						row_i[id+1] = strings.Join(attr.Values, ", ")
					} else if len(attr.Values) > 0 {
						row_i[id+1] = attr.Values[0]
					}
				}
			}
		}

		ti.Rows = append(ti.Rows, row_i)
		ti.DN = append(ti.DN, entry.DN)
	}
}

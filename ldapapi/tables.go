package ldapapi

import (
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

func LoadTableInfoFromSearchResults(
	ti *TableInfo,
	colNames []string,
	colAtrr map[string]string,
	widths []int,
	sr *ldap.SearchResult,
) {
	colIds := ToIdMap(colNames)
	ti.Cols = MakeColumns(colNames, widths)
	ti.Rows = []table.Row{}
	ti.DN = []string{}

	for _, entry := range sr.Entries {
		row_i := make([]string, len(colNames))
		if val, ok := colAtrr["dn"]; ok {
			row_i[colIds[val]] = entry.DN
		}
		for _, attr := range entry.Attributes {
			_, ok := colAtrr[attr.Name]
			if !ok {
				continue
			}
			id := colIds[colAtrr[attr.Name]]
			if len(attr.Values) > 1 {
				row_i[id] = strings.Join(attr.Values, "; ")
			} else {
				row_i[id] = attr.Values[0]
			}
		}
		ti.Rows = append(ti.Rows, row_i)
		ti.DN = append(ti.DN, entry.DN)
	}
}

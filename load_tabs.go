package main

import (
	"fmt"
	"os"

	"github.com/mshagirov/goldap/internal/tabs"
	"github.com/mshagirov/goldap/ldapapi"
)

func createNewTabsModel(tableIdx int, rowIndices []int, api *ldapapi.LdapApi) tabs.Model {
	contents := make([]ldapapi.TableInfo, 0)
	dn := make([][]string, 0)

	for _, tabName := range ldapapi.TableNames {
		t, err := api.GetTableInfo(tabName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		contents = append(contents, t)
		dn = append(dn, t.DN)
	}

	m := tabs.NewTabsModel(ldapapi.TableNames, contents, dn, api)

	m.State.TabId = tableIdx
	m.State.Table = tabs.NewTable(contents[tableIdx])
	for i, rowId := range rowIndices {
		m.State.TabSates[i].Cursor = min(len(dn[i]), rowId)
	}
	m.SetTable()
	return m
}

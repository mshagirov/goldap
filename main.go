package main

import (
	"fmt"
	"os"

	"github.com/mshagirov/goldap/internal/config"
	"github.com/mshagirov/goldap/internal/login"
	"github.com/mshagirov/goldap/internal/tabs"
	"github.com/mshagirov/goldap/ldapapi"
)

func main() {
	ldapConfig := config.Read()
	if ldapConfig.LdapUrl == "" {
		fmt.Printf("%v", config.ExampleJson())
		os.Exit(1)
	}

	secret, err := login.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	LdapApi := &ldapapi.LdapApi{
		Config: ldapConfig,
		Secret: secret,
	}

	var (
		contents     []ldapapi.TableInfo
		dn           [][]string
		reload_model = false
		tableIndex   = 0
		rowIndices   = make([]int, len(ldapapi.TableNames))
	)

	for true {
		for _, tabName := range ldapapi.TableNames {
			t, err := LdapApi.GetTableInfo(tabName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			contents = append(contents, t)
			dn = append(dn, t.DN)
		}

		m := tabs.NewTabsModel(ldapapi.TableNames, contents, dn, LdapApi)

		if reload_model {
			m.State.TabId = tableIndex
			m.State.Table = tabs.NewTable(contents[tableIndex])
			for i, rowId := range rowIndices {
				m.State.TabSates[i].Cursor = min(len(dn[i]), rowId)
			}
			m.SetTable()

			reload_model = false
		}

		state, quit := tabs.RunTabs(m)

		if !quit {
			state.FormInfo.Api = LdapApi

			for i := range rowIndices {
				rowIndices[i] = m.State.TabSates[i].Cursor
			}
			tableIndex = state.TabId

			attrNames, updates := tabs.RunForm(state.FormInfo)

			// Updates can be accessed this way
			for i, val := range updates {
				fmt.Println("Updated:", state.FormInfo.DN, attrNames[i], "->", val)
			}

			reload_model = true
		} else {
			break
		}
	}
}

package main

import (
	"fmt"
	"log"
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

		state := tabs.RunTabs(m)
		switch state.Cmd {
		case tabs.QuitCmd:
			return
		case tabs.UpdateCmd:
			state.FormInfo.Api = LdapApi

			for i := range rowIndices {
				rowIndices[i] = m.State.TabSates[i].Cursor
			}
			tableIndex = state.TabId

			attrNames, updates := tabs.RunForm(state.FormInfo)

			for k, val := range updates {
				if attrNames[k] == "userPassword" {
					res, err := LdapApi.ModifyPassword(state.FormInfo.DN, val)
					if err != nil {
						log.Println("Error updating password for", state.FormInfo.DN)
					} else {
						log.Println(*res)
					}
				}
				log.Println("Updated:", state.FormInfo.DN, attrNames[k])
			}
			reload_model = true
		}
	}
}

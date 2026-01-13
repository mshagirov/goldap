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
		rowIndices   []int
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
			m.ActiveTab = tableIndex
			m.ActiveTable = tabs.NewTable(contents[tableIndex])
			m.ActiveRows = rowIndices
			for i, rowId := range rowIndices {
				m.ActiveRows[i] = min(len(dn[i]), rowId)
			}
			m.SetCursor()

			reload_model = false
		}

		fi, quit := tabs.RunTabs(m)

		if !quit {
			fi.Api = LdapApi

			rowIndices = fi.RowIndices
			tableIndex = fi.TableIndex

			tabs.RunForm(fi)
			reload_model = true
		} else {
			break
		}
	}
}

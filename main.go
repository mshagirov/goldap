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

	ldap := ldapapi.LdapApi{
		Config: &ldapConfig,
		Secret: secret,
	}

	var (
		contents   []ldapapi.TableInfo
		dn         [][]string
		reloaded   = false
		tableIndex = 0
		rowIndices []int
	)

	for true {
		for _, tabName := range ldapapi.TableNames {
			t, err := ldap.GetTableInfo(tabName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			contents = append(contents, t)
			dn = append(dn, t.DN)
		}

		m := tabs.NewTabsModel(ldapapi.TableNames, contents, dn, &ldap)

		if reloaded {
			m.ActiveTab = tableIndex
			m.ActiveTable = tabs.NewTable(contents[tableIndex])
			m.ActiveRows = rowIndices
			for i, rowId := range rowIndices {
				m.ActiveRows[i] = min(len(dn[i]), rowId)
			}
			m.SetCursor()

			reloaded = false
		}

		fi, quit := tabs.RunTabs(m)

		if !quit {
			fi.Api = &ldap

			rowIndices = fi.RowIndices
			tableIndex = fi.TableIndex

			tabs.RunForm(fi)
			reloaded = true
		} else {
			break
		}
	}
}

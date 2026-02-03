package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mshagirov/goldap/internal/cache"
	"github.com/mshagirov/goldap/internal/config"
	"github.com/mshagirov/goldap/internal/login"
	"github.com/mshagirov/goldap/internal/tabs"
	"github.com/mshagirov/goldap/ldapapi"
)

func main() {
	cfg := config.Read()
	if cfg.LdapUrl == "" {
		fmt.Printf("%v", config.ExampleJson())
		os.Exit(1)
	}

	// fmt.Println("Users: ", len(cfg.Users), " attributes {")
	// for _, attr := range cfg.Users {
	// 	fmt.Printf("%15v : %20v\n", attr.Name, attr.Val)
	// }
	// fmt.Println("}")
	//
	// fmt.Println("Groups: ", len(cfg.Groups), " attributes {")
	// for _, attr := range cfg.Groups {
	// 	fmt.Printf("%15v : %20v\n", attr.Name, attr.Val)
	// }
	// fmt.Println("}")
	//
	// fmt.Println("OrgUnits: ", len(cfg.OrgUnits), " attributes {")
	// for _, attr := range cfg.OrgUnits {
	// 	fmt.Printf("%15v : %20v\n", attr.Name, attr.Val)
	// }
	// fmt.Println("}")

	secret, err := login.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	LdapApi := &ldapapi.LdapApi{
		Config: cfg,
		Secret: secret,
		Cache:  cache.NewCache(),
	}

	var (
		reload_model = false
		tableIndex   = 0
		rowIndices   = make([]int, len(ldapapi.TableNames))
		m            = createNewTabsModel(tableIndex, rowIndices, LdapApi)
	)

	for true {
		if reload_model {
			LdapApi.Cache.Clear()
			m = createNewTabsModel(tableIndex, rowIndices, LdapApi)
			reload_model = false
		}

		state := tabs.RunTabs(m)

		switch state.Cmd {
		case tabs.QuitCmd:
			return
		case tabs.AddCmd:
			state.FormInfo.Api = LdapApi
			state.FormInfo.DN = cfg.LdapBaseDn
			for i := range rowIndices {
				rowIndices[i] = m.State.TabSates[i].Cursor
			}
			tableIndex = state.TabId
			attrNames, updates := tabs.RunAddForm(state)
			if updates != nil {
				log.Println("Added entry to", state.FormInfo.TableName, fmt.Sprintf("\"%s\"", state.FormInfo.DN))
				for k := range updates {
					log.Println(attrNames[k], updates[k])
				}
				reload_model = true
			}
		case tabs.UpdateCmd:
			state.FormInfo.Api = LdapApi
			for i := range rowIndices {
				rowIndices[i] = m.State.TabSates[i].Cursor
			}
			tableIndex = state.TabId
			attrNames, updates := tabs.RunUpdateForm(state)

			if err := LdapApi.ModifyAttr(state.FormInfo.DN, attrNames, updates); err != nil {
				log.Println(err)
			} else {
				for k := range updates {
					log.Println("Updated:", state.FormInfo.DN, attrNames[k])
				}
			}
			reload_model = true
		}
		state.Cmd = tabs.QuitCmd
	}
}

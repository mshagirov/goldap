package main

import (
	"fmt"
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
		SUCCESS_MSG  = "Success【✓】"
		FAILURE_MSG  = "Failure【✗】"
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
			if len(updates) > 0 {
				if err := LdapApi.AddEntry(state.FormInfo.DN, attrNames, updates); err != nil {
					tabs.RunMessageBox(FAILURE_MSG,
						fmt.Sprintf("Failed to add entry to '%v': %v", state.FormInfo.TableName, err))
				} else {
					tabs.RunMessageBox(SUCCESS_MSG,
						fmt.Sprintf("Successfully added %s to %s", state.FormInfo.DN, state.FormInfo.TableName))
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

			if len(updates) > 0 {
				if err := LdapApi.ModifyAttr(state.FormInfo.DN, attrNames, updates); err != nil {
					tabs.RunMessageBox(FAILURE_MSG,
						fmt.Sprintf("Failed to update entry: %v", err))
				} else {
					tabs.RunMessageBox(SUCCESS_MSG,
						fmt.Sprintf("Successfully updated %s", state.FormInfo.DN))
				}
				reload_model = true
			}
		case tabs.DeleteCmd:
			title := "Deleting Entry"
			message := fmt.Sprintf("Delete LDAP entry: '%s'?", state.FormInfo.DN)
			if tabs.RunConfirmBox(title, message) == tabs.ResultConfirm {
				if err := LdapApi.DeleteEntry(state.FormInfo.DN); err != nil {
					tabs.RunMessageBox(FAILURE_MSG,
						fmt.Sprintf("Failed to delete entry from '%v': %v", state.FormInfo.TableName, err))
				} else {
					tabs.RunMessageBox(SUCCESS_MSG,
						fmt.Sprintf("Successfully deleted %s", state.FormInfo.DN))
				}
				reload_model = true
			}
		}
	}
}

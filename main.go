package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mshagirov/goldap/internal/cache"
	"github.com/mshagirov/goldap/internal/config"
	"github.com/mshagirov/goldap/internal/login"
	"github.com/mshagirov/goldap/internal/tabs"
	"github.com/mshagirov/goldap/internal/ui"
	"github.com/mshagirov/goldap/ldapapi"
)

func main() {
	cfg := config.Read()
	if cfg.LdapUrl == "" {
		fmt.Printf("%v", config.ExampleJson())
		os.Exit(1)
	}

	session := ui.NewSession()
	defer session.Close()

	secret, err := login.Run(session)
	if err != nil {
		session.Close()
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
	)
	m, err := createNewTabsModel(tableIndex, rowIndices, LdapApi)
	if err != nil {
		session.Close()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for true {
		if reload_model {
			reload_model = false
			next, err := createNewTabsModel(tableIndex, rowIndices, LdapApi)
			if err != nil {
				if !showMessage(session, "Failure", fmt.Sprintf("Failed to reload tables: '%v'", err)) {
					return
				}
			} else {
				m = next
			}
		}

		state, err := tabs.RunTabs(session, m)
		if err != nil {
			session.Close()
			fmt.Fprintln(os.Stderr, err)
			return
		}

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
			attrNames, updates, err := tabs.RunAddForm(session, state)
			if err != nil {
				if !showMessage(session, FAILURE_MSG,
					fmt.Sprintf("Failed to add entry to '%v': %v", state.FormInfo.TableName, err)) {
					return
				}
				continue
			}
			if len(updates) > 0 {
				if err := LdapApi.AddEntry(state.FormInfo.DN, attrNames, updates); err != nil {
					if !showMessage(session, FAILURE_MSG,
						fmt.Sprintf("Failed to add entry to '%v': %v", state.FormInfo.TableName, err)) {
						return
					}
				} else {
					if !showMessage(session, SUCCESS_MSG,
						fmt.Sprintf("Added '%s' to '%s'", state.FormInfo.DN, state.FormInfo.TableName)) {
						return
					}
				}
				reload_model = true
			}
		case tabs.UpdateCmd:
			state.FormInfo.Api = LdapApi
			for i := range rowIndices {
				rowIndices[i] = m.State.TabSates[i].Cursor
			}
			tableIndex = state.TabId
			attrNames, updates, err := tabs.RunUpdateForm(session, state)
			if err != nil {
				if !showMessage(session, FAILURE_MSG, fmt.Sprintf("Failed to update: %v", err)) {
					return
				}
				continue
			}

			if len(updates) > 0 {
				if err := LdapApi.ModifyAttr(state.FormInfo.DN, attrNames, updates); err != nil {
					if !showMessage(session, FAILURE_MSG,
						fmt.Sprintf("Failed to update: %v", err)) {
						return
					}
				} else {
					if !showMessage(session, SUCCESS_MSG,
						fmt.Sprintf("Updated %s", state.FormInfo.DN)) {
						return
					}
				}
				reload_model = true
			}
		case tabs.ReloadCmd:
			tableIndex = state.TabId
			for i := range rowIndices {
				rowIndices[i] = state.TabSates[i].Cursor
			}
			rowIndices[tableIndex] = state.Table.Cursor()
			reload_model = true
		case tabs.DeleteCmd:
			title := "Deleting Entry"
			message := fmt.Sprintf("Delete '%s'?", state.FormInfo.DN)
			result, err := tabs.RunConfirmBox(session, title, message)
			if err != nil {
				if !showMessage(session, FAILURE_MSG,
					fmt.Sprintf("Failed to confirm deletion of '%s': %v", state.FormInfo.DN, err)) {
					return
				}
				continue
			}
			if result == tabs.ResultConfirm {
				if err := LdapApi.DeleteEntry(state.FormInfo.DN); err != nil {
					if !showMessage(session, FAILURE_MSG,
						fmt.Sprintf("Failed deleting '%s' from '%v': %v", state.FormInfo.DN, state.FormInfo.TableName, err)) {
						return
					}
				} else {
					if !showMessage(session, SUCCESS_MSG,
						fmt.Sprintf("Deleted %s", state.FormInfo.DN)) {
						return
					}
				}
				reload_model = true
			}
		}
	}
}

// Return to the shell if the dialog fails, so the next TUI cannot hide stderr.
func showMessage(session *ui.Session, title, message string) bool {
	run := func(title, message string) (tabs.MessageBoxResult, error) {
		result, err := tabs.RunMessageBox(session, title, message)
		if err != nil {
			session.Close()
		}
		return result, err
	}
	return reportMessage(title, message, run, os.Stderr)
}

func reportMessage(title, message string, run func(string, string) (tabs.MessageBoxResult, error), output io.Writer) bool {
	if _, err := run(title, message); err != nil {
		fmt.Fprintf(output, "%s: %s\nUnable to display result: %v\n", title, message, err)
		return false
	}
	return true
}

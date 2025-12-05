package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/mshagirov/goldap/internal/config"
	"github.com/mshagirov/goldap/internal/login"
	"github.com/mshagirov/goldap/internal/tabs"
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

	ldap := config.LdapApi{
		Config: &ldapConfig,
		Secret: secret,
	}

	// "(objectClass=*)" // all classes

	var (
		tabnames = []string{"Users", "Groups", "OrgUnits"}
		contents []table.Model
		dn       [][]string
	)

	w, h := tabs.GetTableDimensions()

	for _, tabName := range tabnames {
		t := ldap.GetTableInfo(tabName)
		contents = append(contents,
			table.New(table.WithColumns(t.Cols),
				table.WithRows(t.Rows),
				table.WithFocused(true),
				table.WithHeight(h),
				table.WithWidth(w),
				table.WithStyles(tabs.GetTableStyle()),
			),
		)
		dn = append(dn, t.DN)
	}

	tabs.Run(tabnames, contents, dn)
}

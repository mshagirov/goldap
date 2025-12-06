package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
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

	if err := ldap.TryConnecting(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// "(objectClass=*)" // all classes

	var (
		tabnames = []string{"Users", "Groups", "OrgUnits"}
		contents []table.Model
		dn       [][]string
	)

	w, h := tabs.GetTableDimensions()

	for _, tabName := range tabnames {
		t, _ := ldap.GetTableInfo(tabName)
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

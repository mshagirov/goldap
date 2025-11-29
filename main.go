package main

import (
	"fmt"
	"os"

	"github.com/mshagirov/goldap/internal/config"
	"github.com/mshagirov/goldap/internal/login"
)

func main() {
	// need to have config file
	ldapConfig := config.Read()
	if ldapConfig.LdapUrl == "" {
		fmt.Printf("%v", config.ExampleJson())
		os.Exit(1)
	}

	// enter ldap admin password
	secret, err := login.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ldap := config.LdapApi{
		Config: &ldapConfig,
		Secret: secret,
	}

	filters := []string{
		"PosixAccount",       // all PosixGroups
		"PosixGroup",         // all PosixGroups
		"OrganizationalUnit", // ou's
	}
	// "(objectClass=*)" // all classes
	// "(uid=*)" // all ldap users
	// "(cn=*)" // all ldap users
	// fmt.Sprintf("(uid=%s)", "jbourne") // find user

	var contents []string

	for i, filter := range filters {
		sr, err := ldap.Search(fmt.Sprintf("(objectClass=%v)", filter))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		contents = append(contents, "")
		for _, entry := range sr.Entries {
			contents[i] += fmt.Sprintln("dn:", entry.DN)
			for _, attr := range entry.Attributes {
				for _, val := range attr.Values {
					contents[i] += fmt.Sprintf(" | %v: '%v'\n", attr.Name, val)
				}
			}
		}
	}

	runTabs(
		// []string{"Users", "Groups", "Orgs"},
		//[]string{"User1, user2,...", "group1, group2,...", "Managers, Devs"},
		filters,
		contents,
	)

}

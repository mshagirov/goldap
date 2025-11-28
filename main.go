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

	// filter := "(objectClass=*)" // all classes
	filter := "(objectClass=PosixAccount)" // all PosixGroups
	// filter := "(objectClass=PosixGroup)" // all PosixGroups
	// filter := "(uid=*)" // all ldap users
	// filter := "(cn=*)" // all ldap users
	// filter := fmt.Sprintf("(uid=%s)", "jbourne") // find user

	sr, err := ldap.Search(filter)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, entry := range sr.Entries {
		fmt.Println("dn:", entry.DN)
		if len(entry.GetAttributeValue("uid")) > 1 {
			fmt.Printf("%v\n", entry.GetAttributeValue("uid"))
		}
		// for _, attr := range entry.Attributes {
		// 	fmt.Printf("    |-- %v: %v\n", attr.Name, attr.Values)
		// }
	}

	runTabs(
		[]string{"Users", "Groups", "Orgs"},
		[]string{"User1, user2,...", "group1, group2,...", "Managers, Devs"},
	)

}

package main

import (
	"fmt"
	"os"

	"github.com/mshagirov/goldap/internal/config"
	"github.com/mshagirov/goldap/login"
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

	filter := "(objectClass=*)" // all classes
	// filter := fmt.Sprintf("(uid=%s)", "jbourne") // find user

	sr, err := ldap.Search(filter)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	choices := []string{}
	for _, entry := range sr.Entries {
		choices = append(choices, entry.DN)
	}

	p, m := NewInitialModel(choices)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	for selected := range m.selected {
		fmt.Println("dn:", m.choices[selected])
		for _, attr := range sr.Entries[selected].Attributes {
			fmt.Printf("    |-- %v: %v\n", attr.Name, attr.Values)
		}

	}
}

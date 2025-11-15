package main

import (
	"fmt"
	"os"

	"github.com/mshagirov/goldap/internal/config"
)

func main() {
	ldapConfig := config.Read()
	if ldapConfig.LdapUrl == "" {
		fmt.Printf("%v", config.ExampleJson())
		os.Exit(1)
	}

	fmt.Printf("%#v\n", ldapConfig)
	// p := NewInitialModel()
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Alas, there's been an error: %v", err)
	// 	os.Exit(1)
	// }
}

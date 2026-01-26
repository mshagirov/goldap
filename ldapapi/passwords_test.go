package ldapapi

import (
	"fmt"
	"testing"
)

func TestHashPasswordSSHA_salt4(t *testing.T) {
	cases := []string{"", " ", "   ", "p@ssw0rd", "secret_123", "~!@#$%^&*()-_=+ "}

	for id, c := range cases {
		t.Run(
			fmt.Sprintf("Test %v", id),
			func(t *testing.T) {
				out, err := HashPasswordSSHA(c, 4)
				if err != nil {
					t.Errorf("error using HashPasswordSSHA %v", err)
					return
				}
				ok, err := VerifyHashSSHA(c, out)
				if err != nil {
					t.Errorf("error using VerifyHashSSHA %v", err)
					return
				}
				if !ok {
					t.Errorf("%v did not match the hash %v", c, out)
					return
				}
			})
	}
}

func TestHashPasswordSSHA_salt8(t *testing.T) {
	cases := []string{"", " ", "   ", "p@ssw0rd", "secret_123", "~!@#$%^&*()-_=+ "}

	for id, c := range cases {
		t.Run(
			fmt.Sprintf("Test %v", id),
			func(t *testing.T) {
				out, err := HashPasswordSSHA(c, 8)
				if err != nil {
					t.Errorf("error using HashPasswordSSHA %v", err)
					return
				}
				ok, err := VerifyHashSSHA(c, out)
				if err != nil {
					t.Errorf("error using VerifyHashSSHA %v", err)
					return
				}
				if !ok {
					t.Errorf("%v did not match the hash %v", c, out)
					return
				}
			})
	}
}

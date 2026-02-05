package ldapapi

import (
	"fmt"
	"testing"
)

func TestHashPasswordSSHA_salt4(t *testing.T) {
	cases := []struct {
		correct  string
		passwd   string
		expected bool
	}{
		{correct: "", passwd: "", expected: true},
		{correct: " ", passwd: "", expected: false},
		{correct: " ", passwd: " ", expected: true},
		{correct: "       ", passwd: "password", expected: false},
		{correct: "       ", passwd: "       ", expected: true},
		{correct: "p@ssw0rd", passwd: "p@ssword", expected: false},
		{correct: "p@ssw0rd", passwd: "p@ssw0rd", expected: true},
		{correct: "secret_123", passwd: "secret_123", expected: true},
		{correct: "~!@#$%^&*()-_=+ ", passwd: " +=_-)(*&^%$#@!~", expected: false},
		{correct: "~!@#$%^&*()-_=+ ", passwd: "~!@#$%^&*()-_=+ ", expected: true},
	}

	for id, c := range cases {
		t.Run(
			fmt.Sprintf("Test %v", id),
			func(t *testing.T) {
				out, err := HashPasswordSSHA(c.correct, 4)
				if err != nil {
					t.Errorf("error using HashPasswordSSHA %v", err)
					return
				}
				got, err := VerifyHashSSHA(c.passwd, out)
				if err != nil {
					t.Errorf("error using VerifyHashSSHA %v", err)
					return
				}
				if got != c.expected {
					t.Errorf("Hashing '%v' verify with %v expected %v got %v",
						c.correct, c.passwd, c.expected, got)
				}
			})
	}
}

func TestHashPasswordSSHA_salt8(t *testing.T) {
	cases := []struct {
		correct  string
		passwd   string
		expected bool
	}{
		{correct: "", passwd: "", expected: true},
		{correct: " ", passwd: "", expected: false},
		{correct: " ", passwd: " ", expected: true},
		{correct: "       ", passwd: "password", expected: false},
		{correct: "       ", passwd: "       ", expected: true},
		{correct: "p@ssw0rd", passwd: "p@ssword", expected: false},
		{correct: "p@ssw0rd", passwd: "p@ssw0rd", expected: true},
		{correct: "secret_123", passwd: "secret_123", expected: true},
		{correct: "~!@#$%^&*()-_=+ ", passwd: " +=_-)(*&^%$#@!~", expected: false},
		{correct: "~!@#$%^&*()-_=+ ", passwd: "~!@#$%^&*()-_=+ ", expected: true},
	}

	for id, c := range cases {
		t.Run(
			fmt.Sprintf("Test %v", id),
			func(t *testing.T) {
				out, err := HashPasswordSSHA(c.correct, 8)
				if err != nil {
					t.Errorf("error using HashPasswordSSHA %v", err)
				}
				got, err := VerifyHashSSHA(c.passwd, out)
				if err != nil {
					t.Errorf("error using VerifyHashSSHA %v", err)
				}
				if got != c.expected {
					t.Errorf("Hashing '%v' verify with '%v' expected %v got %v",
						c.correct, c.passwd, c.expected, got)

					return
				}
			})
	}
}

func TestHashPasswordSSHA_skip(t *testing.T) {
	cases := []struct {
		ssha_hash string
	}{
		{ssha_hash: "{SSHA}zEnHeLNDZS71VG7n3ONEIu/vzOuA1Dk7ioAoEw=="},
		{ssha_hash: "{SSHA}j0Y57BC4HIFt+FciKN4uPSaxRUocZz81"},
	}

	for id, c := range cases {
		t.Run(
			fmt.Sprintf("Test %v", id),
			func(t *testing.T) {
				out, err := HashPasswordSSHA(c.ssha_hash, 4)
				if err != nil {
					t.Errorf("error using HashPasswordSSHA %v", err)
					return
				}
				if out != c.ssha_hash {
					t.Errorf("in:'%s' != out:'%s'", c.ssha_hash, out)
				}
			})
	}
}

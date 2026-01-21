package ldapapi

import (
	"fmt"
	"testing"
)

func TestGetFirstDnAttr(t *testing.T) {
	type expected struct {
		attr  string
		value string
		ok    bool
	}

	cases := []struct {
		dn  string
		out expected
	}{
		{dn: "", out: expected{attr: "", value: "", ok: false}},
		{dn: " ", out: expected{attr: "", value: "", ok: false}},
		{dn: ",uid=asimov", out: expected{attr: "", value: "", ok: false}},
		{dn: "uid=,asimov", out: expected{attr: "", value: "", ok: false}},
		{dn: "uid=asimov", out: expected{attr: "", value: "", ok: false}},
		{dn: "uid=asimov,", out: expected{attr: "uid", value: "asimov", ok: true}},
		{dn: "uid=asimov,ou=People,dc=goldap,dc=sh", out: expected{attr: "uid", value: "asimov", ok: true}},
		{dn: "ou=People,dc=goldap,dc=sh", out: expected{attr: "ou", value: "People", ok: true}},
	}

	for id, c := range cases {
		t.Run(
			fmt.Sprintf("Test %v", id),
			func(t *testing.T) {
				attr, value, ok := GetFirstDnAttr(c.dn)
				if attr != c.out.attr || value != c.out.value || ok != c.out.ok {
					t.Errorf("expected %+v doesn't match output (attr:%v value:%v ok:%v)", c.out, attr, value, ok)
					return
				}
			})
	}
}

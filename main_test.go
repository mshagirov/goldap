package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/mshagirov/goldap/internal/tabs"
)

func TestReportMessage(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "dialog succeeds", true: "stderr fallback"}[fail], func(t *testing.T) {
			var output bytes.Buffer
			run := func(title, message string) (tabs.MessageBoxResult, error) {
				if title != "Failure" || message != "LDAP delete failed: insufficient access" {
					t.Fatal("original operation result was not passed to the dialog")
				}
				if fail {
					return tabs.ResultCancel, errors.New("terminal unavailable")
				}
				return tabs.ResultConfirm, nil
			}
			if got := reportMessage("Failure", "LDAP delete failed: insufficient access", run, &output); got == fail {
				t.Errorf("continue = %v, want %v", got, !fail)
			}
			if fail {
				for _, text := range []string{"Failure", "LDAP delete failed: insufficient access", "terminal unavailable"} {
					if !strings.Contains(output.String(), text) {
						t.Errorf("fallback output missing %q", text)
					}
				}
			} else if output.Len() != 0 {
				t.Errorf("unexpected stderr: %s", output.String())
			}
		})
	}
}

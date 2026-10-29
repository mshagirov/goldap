package tabs

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestDialogResults(t *testing.T) {
	failure := errors.New("terminal unavailable")
	dialogs := []struct {
		name    string
		run     func(func() (tea.Model, error)) (MessageBoxResult, error)
		confirm tea.Model
		cancel  tea.Model
	}{
		{"confirmation", runConfirmation, ConfirmBoxModel{Result: ResultConfirm}, ConfirmBoxModel{Result: ResultCancel}},
		{"message", runMessage, MessageBoxModel{Result: ResultConfirm}, MessageBoxModel{Result: ResultCancel}},
	}
	for _, dialog := range dialogs {
		t.Run(dialog.name, func(t *testing.T) {
			cases := []struct {
				name   string
				model  tea.Model
				err    error
				want   MessageBoxResult
				errMsg string
			}{
				{name: "confirmed", model: dialog.confirm, want: ResultConfirm},
				{name: "cancelled", model: dialog.cancel, want: ResultCancel},
				{name: "program failure", err: failure, want: ResultCancel, errMsg: "terminal unavailable"},
				{name: "nil result", want: ResultCancel, errMsg: "unexpected result"},
				{name: "unexpected model", model: formModel{}, want: ResultCancel, errMsg: "unexpected result"},
			}
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					got, err := dialog.run(func() (tea.Model, error) { return c.model, c.err })
					if got != c.want {
						t.Errorf("result = %v, want %v", got, c.want)
					}
					if c.errMsg == "" {
						if err != nil {
							t.Fatal(err)
						}
					} else if err == nil || !strings.Contains(err.Error(), c.errMsg) {
						t.Errorf("error = %v, want containing %q", err, c.errMsg)
					}
					if c.err != nil && !errors.Is(err, c.err) {
						t.Errorf("error = %v, does not wrap original failure", err)
					}
				})
			}
		})
	}
}

func TestConfirmBoxCancellation(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{
		{Code: tea.KeyEscape},
		{Code: 'c', Mod: tea.ModCtrl},
		{Code: tea.KeyEnter}, // Cancel is selected by default.
	} {
		t.Run(key.String(), func(t *testing.T) {
			model := NewConfirmBoxModel("Delete", "test entry")
			if model.Result != ResultCancel {
				t.Fatal("new dialog must default to cancellation")
			}
			result, cmd := model.Update(key)
			if result.(ConfirmBoxModel).Result != ResultCancel || cmd == nil {
				t.Fatal("cancel key must quit without confirming deletion")
			}
		})
	}
}

func TestConfirmBoxSave(t *testing.T) {
	for _, c := range []struct {
		name string
		keys []tea.KeyPressMsg
	}{
		{name: "save shortcut", keys: []tea.KeyPressMsg{{Code: 's', Text: "s"}}},
		{name: "select save then enter", keys: []tea.KeyPressMsg{{Code: tea.KeyTab}, {Code: tea.KeyEnter}}},
	} {
		t.Run(c.name, func(t *testing.T) {
			var model tea.Model = NewConfirmBoxModel("Delete", "test entry")
			var cmd tea.Cmd
			for _, key := range c.keys {
				model, cmd = model.Update(key)
			}
			if model.(ConfirmBoxModel).Result != ResultConfirm || cmd == nil {
				t.Fatal("explicit confirmation must quit with ResultConfirm")
			}
		})
	}
}

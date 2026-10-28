package tabs

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func typeText(model tea.Model, s string) tea.Model {
	for _, r := range s {
		model, _ = model.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	return model
}

func TestUpdateFormEmptyField(t *testing.T) {
	cases := []struct {
		name        string
		typed       string
		wantUpdated map[int]string
	}{
		{name: "fill empty field", typed: "asimov, hseldon", wantUpdated: map[int]string{1: "asimov, hseldon"}},
		{name: "leave empty field empty", typed: "", wantUpdated: map[int]string{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			updated := make(map[int]string)
			var model tea.Model = initialFormModel("cn=admins", "Save", "cn=admins",
				[]string{"admins", ""}, []string{"cn", "memberUid"}, &updated)
			model, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 20})

			model, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyDown})
			model, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			if v := model.(formModel).inputs[1].Value(); v != "" {
				t.Fatalf("editing empty field starts with %q, want empty", v)
			}
			model = typeText(model, c.typed)
			model, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

			if len(updated) != len(c.wantUpdated) || updated[1] != c.wantUpdated[1] {
				t.Errorf("updated = %v, want %v", updated, c.wantUpdated)
			}
		})
	}
}

func TestAddFormEmptyDefaultUsesName(t *testing.T) {
	updated := make(map[int]string)
	m := InitialAddFormModel("new", "Add", "", []string{"", "x"}, []string{"description", "cn"}, map[int]struct{}{}, &updated)
	if p, v := m.inputs[0].Placeholder, m.inputs[0].Value(); p != "description" || v != "description" {
		t.Errorf("placeholder, value = %q, %q; want attribute name for both", p, v)
	}
}

package login

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/mshagirov/goldap/internal/mouse"
)

func TestLoginMouse(t *testing.T) {
	cases := []struct {
		name      string
		startOn   int // focus index before clicking: 0 field, 1 submit button
		target    string
		xOffset   int // cells right of the target text's start
		wantFocus int
		wantQuit  bool
	}{
		{name: "click field focuses it", startOn: 1, target: "> Password", wantFocus: 0},
		{name: "click empty part of field", startOn: 1, target: "> Password", xOffset: 14, wantFocus: 0},
		{name: "click focused field keeps focus", startOn: 0, target: "> Password", wantFocus: 0},
		{name: "click submit", startOn: 0, target: submitText, wantFocus: 0, wantQuit: true},
		{name: "click title does nothing", startOn: 1, target: "LDAP Administrator", wantFocus: 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := initialModel()
			m.focus(c.startOn)
			r, ok := mouse.Find(m.render(), c.target)
			if !ok {
				t.Fatalf("%q not visible", c.target)
			}
			res, cmd := m.Update(tea.MouseClickMsg{X: r.X + c.xOffset, Y: r.Y, Button: tea.MouseLeft})
			got := res.(model)
			if got.focusIndex != c.wantFocus {
				t.Errorf("focusIndex = %d, want %d", got.focusIndex, c.wantFocus)
			}
			if focused := got.inputs[0].Focused(); focused != (c.wantFocus == 0) {
				t.Errorf("input focused = %v, want %v", focused, c.wantFocus == 0)
			}
			quit := cmd != nil && func() bool { _, ok := cmd().(tea.QuitMsg); return ok }()
			if quit != c.wantQuit {
				t.Errorf("quit = %v, want %v", quit, c.wantQuit)
			}
		})
	}
}

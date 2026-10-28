package tabs

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mshagirov/goldap/internal/mouse"
	"github.com/mshagirov/goldap/ldapapi"
)

func newTestTabsModel(t *testing.T, nRows int) Model {
	t.Helper()
	ti := ldapapi.TableInfo{
		Cols: ldapapi.MakeColumns([]string{"", "Name"}, []int{4, 10}),
	}
	for i := range nRows {
		ti.Rows = append(ti.Rows, []string{fmt.Sprint(i + 1), fmt.Sprintf("group%02d", i+1)})
		ti.DN = append(ti.DN, fmt.Sprintf("cn=group%02d,dc=goldap,dc=sh", i+1))
	}
	m := NewTabsModel([]string{"Groups"}, []ldapapi.TableInfo{ti}, [][]string{ti.DN}, nil)
	model, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 20})
	return model.(Model)
}

// screenLine returns the screen line where pattern matches, or -1.
func screenLine(view, pattern string) int {
	re := regexp.MustCompile(pattern)
	for y, line := range strings.Split(view, "\n") {
		if re.MatchString(ansi.Strip(line)) {
			return y
		}
	}
	return -1
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

func TestTabsMouse(t *testing.T) {
	cases := []struct {
		name       string
		scroll     int    // wheel steps before clicking (negative: up)
		clickRow   string // row number to click
		clicks     int
		wantCursor int
		wantQuit   bool
	}{
		{name: "first click selects", clickRow: "3", clicks: 1, wantCursor: 2},
		{name: "second click opens entry", clickRow: "3", clicks: 2, wantCursor: 2, wantQuit: true},
		{name: "wheel down moves cursor", scroll: 4, wantCursor: 4},
		{name: "wheel up stops at top", scroll: -3, wantCursor: 0},
		{name: "click after scrolling", scroll: 25, clickRow: "20", clicks: 1, wantCursor: 19},
		{name: "click above after scrolling", scroll: 25, clickRow: "16", clicks: 2, wantCursor: 15, wantQuit: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestTabsModel(t, 40)
			var model tea.Model = m
			for range max(c.scroll, -c.scroll) {
				button := wheelMoveDown
				if c.scroll < 0 {
					button = wheelMoveUp
				}
				model, _ = model.Update(tea.MouseWheelMsg{Button: button})
			}

			var cmd tea.Cmd
			for range c.clicks {
				y := screenLine(model.(Model).viewport.View(), `^│\s*`+c.clickRow+`\s`)
				if y < 0 {
					t.Fatalf("row %s not visible:\n%s", c.clickRow, ansi.Strip(model.(Model).viewport.View()))
				}
				model, cmd = model.Update(tea.MouseClickMsg{Y: y, Button: tea.MouseLeft})
			}

			got := model.(Model)
			if cursor := got.State.Table.Cursor(); cursor != c.wantCursor {
				t.Errorf("cursor = %d, want %d", cursor, c.wantCursor)
			}
			if quit := isQuit(cmd); quit != c.wantQuit {
				t.Errorf("quit = %v, want %v", quit, c.wantQuit)
			}
			if c.wantQuit && got.State.Cmd != UpdateCmd {
				t.Errorf("State.Cmd = %v, want UpdateCmd", got.State.Cmd)
			}
		})
	}
}

func TestTabsMouseClickOutsideRows(t *testing.T) {
	m := newTestTabsModel(t, 5)
	for _, y := range []int{0, 1, 2, 3, 100} { // tab row, table header, out of screen
		model, cmd := m.Update(tea.MouseClickMsg{Y: y, Button: tea.MouseLeft})
		if cursor := model.(Model).State.Table.Cursor(); cursor != 0 || isQuit(cmd) {
			t.Errorf("click at y=%d: cursor = %d, quit = %v", y, cursor, isQuit(cmd))
		}
	}
}

func newTestFormModel(nFields int) formModel {
	names := make([]string, nFields)
	values := make([]string, nFields)
	for i := range nFields {
		names[i] = fmt.Sprintf("attr%02d", i)
		values[i] = fmt.Sprintf("value%02d", i)
	}
	updated := make(map[int]string)
	m := initialFormModel("title", "msg title", "msg", values, names, &updated)
	model, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 20})
	return model.(formModel)
}

func TestFormMouse(t *testing.T) {
	cases := []struct {
		name       string
		scroll     int
		clickField string // field name to click; "" for none
		valueLine  bool   // click the value line instead of the name line
		clicks     int
		wantIndex  int
		wantActive bool
	}{
		{name: "first click selects", clickField: "attr02", clicks: 1, wantIndex: 2},
		{name: "click on value line selects", clickField: "attr03", valueLine: true, clicks: 1, wantIndex: 3},
		{name: "second click edits", clickField: "attr02", clicks: 2, wantIndex: 2, wantActive: true},
		{name: "third click keeps editing", clickField: "attr02", clicks: 3, wantIndex: 2, wantActive: true},
		{name: "wheel down moves selection", scroll: 3, wantIndex: 3},
		{name: "wheel up stops at top", scroll: -2, wantIndex: 0},
		{name: "click after scrolling", scroll: 20, clickField: "attr17", clicks: 1, wantIndex: 17},
		{name: "click above after scrolling", scroll: 20, clickField: "attr19", clicks: 2, wantIndex: 19, wantActive: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var model tea.Model = newTestFormModel(30)
			for range max(c.scroll, -c.scroll) {
				button := wheelMoveDown
				if c.scroll < 0 {
					button = wheelMoveUp
				}
				model, _ = model.Update(tea.MouseWheelMsg{Button: button})
			}

			for range c.clicks {
				y := screenLine(model.(formModel).render(), `^\s*`+c.clickField+`\s*$`)
				if y < 0 {
					t.Fatalf("field %s not visible:\n%s", c.clickField, ansi.Strip(model.(formModel).render()))
				}
				if c.valueLine {
					y++
				}
				model, _ = model.Update(tea.MouseClickMsg{Y: y, Button: tea.MouseLeft})
			}

			got := model.(formModel)
			if got.index != c.wantIndex {
				t.Errorf("index = %d, want %d", got.index, c.wantIndex)
			}
			if _, active := got.active[got.index]; active != c.wantActive {
				t.Errorf("active = %v, want %v", active, c.wantActive)
			}
		})
	}
}

func TestTabsMouseSelectTab(t *testing.T) {
	ti := ldapapi.TableInfo{Cols: ldapapi.MakeColumns([]string{"", "Name"}, []int{4, 10})}
	// short names: without a terminal the screen falls back to 20 columns
	names := []string{"Grp", "OU"}
	m := NewTabsModel(names, []ldapapi.TableInfo{ti, ti}, [][]string{nil, nil}, nil)
	var model tea.Model
	model, _ = m.Update(tea.WindowSizeMsg{Width: 40, Height: 20})

	cases := []struct {
		name    string
		tab     string // tab label to click; "" clicks the filler right of the tabs
		yOffset int    // click the tab's top/bottom border instead of its label
		wantTab int
	}{
		{name: "click second tab", tab: "OU", wantTab: 1},
		{name: "click same tab", tab: "OU", wantTab: 1},
		{name: "click first tab border", tab: "Grp", yOffset: -1, wantTab: 0},
		{name: "click second tab bottom", tab: "OU", yOffset: 1, wantTab: 1},
		{name: "click filler", wantTab: 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			x, y := 18, 1
			if c.tab != "" {
				r, ok := mouse.Find(model.(Model).viewport.View(), c.tab)
				if !ok {
					t.Fatalf("tab %s not visible", c.tab)
				}
				x, y = r.X, r.Y+c.yOffset
			}
			var cmd tea.Cmd
			model, cmd = model.Update(tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseLeft})
			if got := model.(Model).State.TabId; got != c.wantTab {
				t.Errorf("TabId = %d, want %d", got, c.wantTab)
			}
			if isQuit(cmd) {
				t.Errorf("tab click should not quit")
			}
		})
	}
}

func TestConfirmBoxMouse(t *testing.T) {
	cases := []struct {
		name       string
		button     string // "" clicks outside the buttons
		wantResult MessageBoxResult
		wantQuit   bool
	}{
		{name: "cancel", button: "[C]ancel", wantResult: ResultCancel, wantQuit: true},
		{name: "save", button: "[S]ave", wantResult: ResultConfirm, wantQuit: true},
		{name: "outside", wantResult: ResultConfirm},
	}
	for _, c := range cases {
		for _, inForm := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/inForm=%v", c.name, inForm), func(t *testing.T) {
				var model tea.Model = NewConfirmBoxModel("title", "message")
				if inForm {
					f := newTestFormModel(3)
					f.focused = false
					model = f
				}
				x, y := 0, 0
				if c.button != "" {
					r, ok := mouse.Find(NewConfirmBoxModel("title", "message").render(), c.button)
					if !ok {
						t.Fatalf("button %s not visible", c.button)
					}
					x, y = r.X, r.Y
				}
				model, cmd := model.Update(tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseLeft})
				if quit := isQuit(cmd); quit != c.wantQuit {
					t.Fatalf("quit = %v, want %v", quit, c.wantQuit)
				}
				if !c.wantQuit {
					return
				}
				box, ok := model.(ConfirmBoxModel)
				if !ok {
					t.Fatalf("final model is %T, want ConfirmBoxModel", model)
				}
				if box.Result != c.wantResult {
					t.Errorf("Result = %v, want %v", box.Result, c.wantResult)
				}
			})
		}
	}
}

func TestMessageBoxMouse(t *testing.T) {
	model, cmd := NewMessageBoxModel("title", "message").Update(tea.MouseClickMsg{Button: tea.MouseLeft})
	if !isQuit(cmd) || model.(MessageBoxModel).Result != ResultConfirm {
		t.Errorf("left click should confirm and close the message box")
	}
	if _, cmd := NewMessageBoxModel("title", "message").Update(tea.MouseClickMsg{Button: tea.MouseRight}); isQuit(cmd) {
		t.Errorf("right click should not close the message box")
	}
}

func TestWheelDirection(t *testing.T) {
	cases := []struct {
		name   string
		button tea.MouseButton
		want   int // selected row/field after one wheel event from row/field 1
	}{
		{name: "wheel up moves down", button: tea.MouseWheelUp, want: 2},
		{name: "wheel down moves up", button: tea.MouseWheelDown, want: 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tabs := newTestTabsModel(t, 5)
			tabs.State.Table.MoveDown(1)
			model, _ := tabs.Update(tea.MouseWheelMsg{Button: c.button})
			if got := model.(Model).State.Table.Cursor(); got != c.want {
				t.Errorf("table cursor = %d, want %d", got, c.want)
			}

			var form tea.Model = newTestFormModel(5)
			form, _ = form.Update(tea.KeyPressMsg{Code: tea.KeyDown})
			form, _ = form.Update(tea.MouseWheelMsg{Button: c.button})
			if got := form.(formModel).index; got != c.want {
				t.Errorf("form index = %d, want %d", got, c.want)
			}
		})
	}
}

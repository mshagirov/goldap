package tabs

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

const longDN = "uid=rgreventlov,ou=People,ou=Engineering,ou=Departments,dc=goldap,dc=sh"

// joinedText strips styling and line breaks so wrapped text can be compared.
func joinedText(s string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(ansi.Strip(s), "│", " ")), "")
}

func TestWrapText(t *testing.T) {
	cases := []struct {
		name  string
		input string
		width int
		want  string
	}{
		{name: "fits", input: "short text", width: 20, want: "short text"},
		{name: "breaks at spaces", input: "Failed to update entry", width: 10, want: "Failed to\nupdate\nentry"},
		{name: "breaks DN at commas", input: "uid=asimov,ou=People,dc=goldap", width: 12, want: "uid=asimov,\nou=People,\ndc=goldap"},
		{name: "breaks long word", input: "abcdefghijkl", width: 5, want: "abcde\nfghij\nkl"},
		{name: "no width", input: "unchanged text", width: 0, want: "unchanged text"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := wrapText(c.input, c.width); got != c.want {
				t.Errorf("wrapText(%q, %d) = %q, want %q", c.input, c.width, got, c.want)
			}
		})
	}
}

func TestBoxesWrapText(t *testing.T) {
	// without a terminal the screen falls back to 20 columns, so text wraps
	// to the minimum width or the width of the box's footer/buttons
	frame := msgBoxStyle.GetHorizontalFrameSize()
	footerWidth := lipgloss.Width(msgBlurredBtnStyle.Render(messageFooter))
	buttonsWidth := lipgloss.Width(msgFocusedBtnStyle.Render("[C]ancel")) + lipgloss.Width(msgBlurredBtnStyle.Render("[S]ave")) + 5
	message := "Failed to update entry: LDAP Result Code 65 \"Object Class Violation\": attribute 'member' not allowed on " + longDN
	cases := []struct {
		name     string
		render   string
		maxWidth int
	}{
		{
			name:     "message box",
			render:   NewMessageBoxModel("Failure", message).render(),
			maxWidth: max(minBoxTextWidth, footerWidth) + frame,
		},
		{
			name:     "confirm box",
			render:   NewConfirmBoxModel("Deleting Entry", "Delete LDAP entry: '"+longDN+"'?").render(),
			maxWidth: max(minBoxTextWidth, buttonsWidth) + frame,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if w := lipgloss.Width(c.render); w > c.maxWidth {
				t.Errorf("box is %d wide, want at most %d:\n%s", w, c.maxWidth, ansi.Strip(c.render))
			}
		})
	}
	if got := joinedText(cases[0].render); !strings.Contains(got, joinedText(message)) {
		t.Errorf("wrapped message lost text:\n%s", ansi.Strip(cases[0].render))
	}
}

func TestFormHeaderWraps(t *testing.T) {
	updated := make(map[int]string)
	var model tea.Model = initialFormModel(longDN, "Save", longDN,
		[]string{"a", "b"}, []string{"attr00", "attr01"}, &updated)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 40, Height: 30})
	m := model.(formModel)

	header := m.headerView()
	if w := lipgloss.Width(header); w > 40 {
		t.Errorf("header is %d wide, want at most 40:\n%s", w, ansi.Strip(header))
	}
	if !strings.Contains(joinedText(header), joinedText(longDN)) {
		t.Errorf("wrapped header lost text:\n%s", ansi.Strip(header))
	}
	if h := lipgloss.Height(header); m.viewport.YPosition != h || h < 4 {
		t.Errorf("viewport YPosition = %d, header height = %d; want equal and wrapped", m.viewport.YPosition, h)
	}

	// clicks still land on the right field below the taller header
	y := screenLine(m.render(), `^\s*attr01\s*$`)
	model, _ = m.Update(tea.MouseClickMsg{Y: y, Button: tea.MouseLeft})
	if got := model.(formModel).index; got != 1 {
		t.Errorf("click on attr01 selected field %d, want 1", got)
	}
}

package tabs

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

func initSearch() textinput.Model {
	ti := textinput.New()
	ti.Prompt = "/"
	ti.Placeholder = "search"
	ti.CharLimit = 80
	ti.SetWidth(50)
	// v2 defaults color the prompt and blurred text; keep them unstyled as in v1
	s := ti.Styles()
	s.Focused.Prompt = lipgloss.NewStyle()
	s.Blurred.Prompt = lipgloss.NewStyle()
	s.Blurred.Text = lipgloss.NewStyle()
	ti.SetStyles(s)
	ti.Focus()
	return ti
}

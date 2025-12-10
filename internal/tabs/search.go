package tabs

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func initSearch() textinput.Model {
	ti := textinput.New()
	ti.Prompt = "/"
	ti.Placeholder = "search"
	ti.CharLimit = 80
	ti.Width = 50
	ti.Focus()
	return ti
}

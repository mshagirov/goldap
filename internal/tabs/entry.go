package tabs

import (
	// "fmt"

	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/go-ldap/ldap/v3"
	"github.com/mshagirov/goldap/ldapapi"
)

var (
	focusedColor = lipgloss.AdaptiveColor{Light: "#DAA520", Dark: "#FFD700"} //lipgloss.Color("215")
	focusedStyle = lipgloss.NewStyle().Foreground(focusedColor).PaddingTop(2)
	blurredStyle = lipgloss.NewStyle().Foreground(blurredColor)

	cursorStyle = focusedStyle
	noStyle     = lipgloss.NewStyle().PaddingTop(2)

	titleStyle = focusedStyle.Bold(true)
	footer     = blurredStyle.Render("( press esc or q to exit )")

	// focusedButton = lipgloss.NewStyle().Foreground(focusedColor).Render("[ Save ]")
	// blurredButton = fmt.Sprintf("%s", lipgloss.NewStyle().Foreground(blurredColor).Render("[ Save ]"))

	contentStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(blurredColor)
)

type entryModel struct {
	entries  []string
	names    []string
	selected int
}

func (m entryModel) Init() tea.Cmd {
	return nil
}

func initialEntryModel(entries, names []string) entryModel {
	m := entryModel{
		entries: entries,
		names:   names,
	}
	return m
}

func runEntryModel(filter string, api *ldapapi.LdapApi) {
	sr, err := api.Search(filter)

	m := initialEntryModel(entries, names)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s", err)
		return
	}
}

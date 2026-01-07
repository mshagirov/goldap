package tabs

import (
	// "fmt"

	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// "github.com/go-ldap/ldap/v3"
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
	title     string
	attrNames []string
	attrVals  []string
	selected  int
}

func (m entryModel) Init() tea.Cmd {
	return nil
}

func initialEntryModel(title string, attrNames, attrVals []string) entryModel {
	m := entryModel{
		title:     title,
		attrNames: attrNames,
		attrVals:  attrVals,
	}
	return m
}

func runEntryModel(dn, tableName string, api *ldapapi.LdapApi) {
	attrNames, attrVals := api.GetAttrWithDN(dn, tableName)
	m := initialEntryModel(dn, attrNames, attrVals)
	fmt.Println(m)
	// if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
	// 	fmt.Printf("could not start program: %s\n", err)
	// 	return
	// }
}

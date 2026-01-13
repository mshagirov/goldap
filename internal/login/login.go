package login

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	focusedColor = lipgloss.AdaptiveColor{Light: "#DAA520", Dark: "#FFD700"} //lipgloss.Color("215")
	blurredColor = lipgloss.Color("241")

	focusedStyle = lipgloss.NewStyle().Foreground(focusedColor).PaddingTop(2)
	blurredStyle = lipgloss.NewStyle().Foreground(blurredColor)

	cursorStyle = focusedStyle
	noStyle     = lipgloss.NewStyle().PaddingTop(2)

	titleText string = `LDAP Administrator Password`
	titleSyle        = blurredStyle.Bold(true)
	helpStyle        = blurredStyle
	helpText  string = `(press esc or ctrl+c to exit)`

	focusedButton = lipgloss.NewStyle().Foreground(focusedColor).Render("[ Submit ]")
	blurredButton = fmt.Sprintf("%s", lipgloss.NewStyle().Foreground(blurredColor).Render("[ Submit ]"))

	contentStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(blurredColor)
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 1),
	}

	var t textinput.Model
	t = textinput.New()
	t.Cursor.Style = cursorStyle
	t.Width = 15
	t.CharLimit = 30

	t.Placeholder = "Password"
	t.EchoMode = textinput.EchoPassword
	t.EchoCharacter = '•'
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	t.Focus()

	m.inputs[0] = t

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("GoLDAP | login"), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	physicalWidth, physicalHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		physicalWidth, physicalHeight = 20, 20
	}

	var b strings.Builder

	b.WriteString(titleSyle.Render(titleText))

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render(helpText))

	alignedContents := contentStyle.Render(b.String())

	renderedContent := lipgloss.Place(
		physicalWidth,
		physicalHeight,
		lipgloss.Center, // Horizontal alignment
		lipgloss.Center, // Vertical alignment
		alignedContents,
	)
	return renderedContent
}

func Run() (string, error) {
	m := initialModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return "", fmt.Errorf("could not start program: %s", err)
	}
	secret := m.inputs[0].Value()
	if len(secret) < 1 {
		return "", fmt.Errorf("LDAP password not provided")
	}
	return secret, nil
}

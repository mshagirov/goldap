package login

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"golang.org/x/term"
)

var (
	focusedColor = compat.AdaptiveColor{Light: lipgloss.Color("#DAA520"), Dark: lipgloss.Color("#FFD700")} //lipgloss.Color("215")
	blurredColor = lipgloss.Color("241")

	focusedStyle = lipgloss.NewStyle().Foreground(focusedColor).PaddingTop(2)
	blurredStyle = lipgloss.NewStyle().Foreground(blurredColor)

	noStyle = lipgloss.NewStyle().PaddingTop(2)

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
	t.SetWidth(15)
	t.CharLimit = 30

	t.Placeholder = "Password"
	t.EchoMode = textinput.EchoPassword
	t.EchoCharacter = '•'
	s := t.Styles()
	s.Cursor.Color = focusedColor
	s.Focused.Prompt = focusedStyle
	s.Focused.Text = focusedStyle
	s.Blurred.Prompt = noStyle
	s.Blurred.Text = noStyle
	t.SetStyles(s)
	t.Focus()

	m.inputs[0] = t

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
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
					continue
				}
				m.inputs[i].Blur()
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

func (m model) View() tea.View {
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
	v := tea.NewView(renderedContent)
	v.AltScreen = true
	v.WindowTitle = "GoLDAP | login"
	return v
}

func Run() (string, error) {
	m := initialModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return "", fmt.Errorf("could not start program: %s", err)
	}
	secret := m.inputs[0].Value()
	if len(secret) < 1 {
		return "", fmt.Errorf("LDAP password not provided")
	}
	return secret, nil
}

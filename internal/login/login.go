package login

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"

	"github.com/mshagirov/goldap/internal/mouse"
	"github.com/mshagirov/goldap/internal/ui"
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

	submitText    = "[ Submit ]"
	focusedButton = lipgloss.NewStyle().Foreground(focusedColor).Render(submitText)
	blurredButton = fmt.Sprintf("%s", lipgloss.NewStyle().Foreground(blurredColor).Render(submitText))

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

			return m, m.focus(m.focusIndex)
		}
	case tea.MouseClickMsg:
		if msg.Button != tea.MouseLeft {
			break
		}
		view := m.render()
		for i := range m.inputs {
			if m.inputRegion(view, i).Contains(msg.X, msg.Y) {
				return m, m.focus(i)
			}
		}
		if mouse.Hit(view, submitText, msg.X, msg.Y) {
			return m, tea.Quit
		}
		return m, nil
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

// focus focuses input i (or the submit button when i == len(m.inputs)).
func (m *model) focus(i int) tea.Cmd {
	m.focusIndex = i
	cmds := make([]tea.Cmd, len(m.inputs))
	for j := range m.inputs {
		if j == i {
			cmds[j] = m.inputs[j].Focus()
			continue
		}
		m.inputs[j].Blur()
	}
	return tea.Batch(cmds...)
}

// inputRegion returns the clickable area of input i in the rendered view:
// its prompt followed by the full field width.
func (m model) inputRegion(view string, i int) mouse.Region {
	lines := strings.Split(strings.TrimSpace(ansi.Strip(m.inputs[i].View())), "\n")
	r, ok := mouse.Find(view, strings.TrimSpace(lines[len(lines)-1]))
	if !ok {
		return mouse.Region{}
	}
	r.Width = max(r.Width, lipgloss.Width(m.inputs[i].Prompt)+m.inputs[i].Width()+1)
	return r
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.WindowTitle = "GoLDAP | login"
	return v
}

func (m model) render() string {
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

func Run(runner ui.Runner) (string, error) {
	m := initialModel()
	result, err := runner.Run(m)
	if err != nil {
		return "", fmt.Errorf("could not start program: %s", err)
	}
	secret := result.(model).inputs[0].Value()
	if len(secret) < 1 {
		return "", fmt.Errorf("LDAP password not provided")
	}
	return secret, nil
}

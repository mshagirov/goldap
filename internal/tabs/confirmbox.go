package tabs

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type ConfirmBoxModel struct {
	title      string
	message    string
	confirm    bool
	cancelBtn  string
	confirmBtn string
	Result     MessageBoxResult
}

type MessageBoxResult int

const (
	ResultCancel MessageBoxResult = iota
	ResultConfirm
)

func (m ConfirmBoxModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmBoxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "s", "S":
			m.Result = ResultConfirm
			return m, tea.Quit

		case "c", "C":
			m.Result = ResultCancel
			return m, tea.Quit

		case "enter":
			if m.confirm {
				m.Result = ResultConfirm
			} else {
				m.Result = ResultCancel
			}
			return m, tea.Quit

		case "tab":
			m.confirm = !m.confirm
			return m, nil

		case "shift+tab":
			m.confirm = !m.confirm
			return m, nil

		case "esc":
			m.Result = ResultCancel
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ConfirmBoxModel) View() string {
	physicalWidth, physicalHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		physicalWidth, physicalHeight = 20, 20
	}

	var b strings.Builder
	title := msgTitleStyle.Height(1).Render(m.title)
	b.WriteString(title)
	b.WriteString("\n\n")

	msgtext := msgTextStyle.Height(1).Render(m.message)
	b.WriteString(msgtext)
	b.WriteString("\n\n")

	w := max(lipgloss.Width(title), lipgloss.Width(msgtext))
	cancelStyle := &msgFocusedBtnStyle
	confirmStyle := &msgBlurredBtnStyle
	if m.confirm {
		cancelStyle = &msgBlurredBtnStyle
		confirmStyle = &msgFocusedBtnStyle
	}

	cancelBtn := cancelStyle.Render(m.cancelBtn)
	confirmBtn := confirmStyle.Render(m.confirmBtn)
	buttonSpacing := 5
	minWidth := lipgloss.Width(cancelBtn) + lipgloss.Width(confirmBtn) + buttonSpacing
	padding := (w - minWidth) / 2
	if padding < 0 {
		w = minWidth
		padding = 0
	}

	buttons := lipgloss.Place(w, 1, lipgloss.Center, lipgloss.Center,
		cancelBtn+strings.Repeat(" ", buttonSpacing)+confirmBtn,
	)
	// strings.Repeat(" ", padding)+
	// cancelBtn+strings.Repeat(" ", buttonSpacing)+confirmBtn,
	// +strings.Repeat(" ", padding),
	// )
	w = lipgloss.Width(buttons) + 4

	b.WriteString(buttons)
	b.WriteString("\n")

	renderedContent := lipgloss.Place(
		physicalWidth,
		physicalHeight,
		lipgloss.Center,
		lipgloss.Center,
		msgBoxStyle.Width(w).Render(b.String()),
	)
	return renderedContent
}

func NewMessageBox(title, message string) ConfirmBoxModel {
	return ConfirmBoxModel{
		title:      title,
		message:    message,
		confirm:    false,
		cancelBtn:  "[C]ancel",
		confirmBtn: "[S]ave",
		Result:     ResultConfirm,
	}
}

func RunMessageBox(title, message string) MessageBoxResult {
	m := NewMessageBox(title, message)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if result, err := p.Run(); err != nil {
		return ResultCancel // Default to cancel on error
	} else {
		return result.(ConfirmBoxModel).Result
	}
}

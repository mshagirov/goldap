package tabs

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type MessageBoxModel struct {
	title      string
	message    string
	confirm    bool
	cancelBtn  string
	confirmBtn string
	Width      int
	Result     MessageBoxResult
}

type MessageBoxResult int

const (
	ResultCancel MessageBoxResult = iota
	ResultConfirm
)

// func NewConfirmSaveBox(message string) MessageBoxModel {
// 	return NewMessageBox("Save changes?", message, "[C]ancel", "[S]ave")
// }

func (m MessageBoxModel) Init() tea.Cmd {
	return nil
}

func (m MessageBoxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		// S/s: Trigger save result
		case "s", "S":
			m.Result = ResultConfirm
			return m, tea.Quit

		// C/c: Trigger cancel result
		case "c", "C":
			m.Result = ResultCancel
			return m, tea.Quit

		// Enter: Activate focused button
		case "enter":
			if m.confirm {
				m.Result = ResultConfirm
			} else {
				m.Result = ResultCancel
			}
			return m, tea.Quit

		// Tab: Navigate to next button
		case "tab":
			m.confirm = !m.confirm
			return m, nil

		// Shift+Tab: Navigate to previous button
		case "shift+tab":
			m.confirm = !m.confirm
			return m, nil

		// Esc: Cancel
		case "esc":
			m.Result = ResultCancel
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MessageBoxModel) View() string {
	physicalWidth, physicalHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		physicalWidth, physicalHeight = 20, 20
	}

	var b strings.Builder

	b.WriteString(msgTitleStyle.Render(m.title))
	b.WriteString("\n\n")
	b.WriteString(msgTextStyle.Render(m.message))
	b.WriteString("\n\n")

	cancelStyle := &msgFocusedBtnStyle
	confirmStyle := &msgBlurredBtnStyle
	if m.confirm {
		cancelStyle = &msgBlurredBtnStyle
		confirmStyle = &msgFocusedBtnStyle
	}

	buttons := lipgloss.PlaceHorizontal(
		m.Width,
		lipgloss.Center,
		fmt.Sprintf("%s     %s", cancelStyle.Render(m.cancelBtn), confirmStyle.Render(m.confirmBtn)),
	)
	b.WriteString(buttons)
	b.WriteString("\n")

	renderedContent := lipgloss.Place(
		physicalWidth,
		physicalHeight,
		lipgloss.Center,
		lipgloss.Center,
		msgBoxStyle.Render(b.String()),
	)
	return renderedContent
}

func NewMessageBox(title, message string) MessageBoxModel {
	return MessageBoxModel{
		title:      title,
		message:    message,
		confirm:    true,
		cancelBtn:  "[C]ancel",
		confirmBtn: "[S]ave",
		Width:      max(lipgloss.Width(title), lipgloss.Width(message), 20),
		Result:     ResultConfirm,
	}
}

func RunMessageBox(title, message string) MessageBoxResult {
	m := NewMessageBox(title, message)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if result, err := p.Run(); err != nil {
		return ResultCancel // Default to cancel on error
	} else {
		return result.(MessageBoxModel).Result
	}
}

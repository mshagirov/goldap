package tabs

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	messageFooter = `Press ANY key/button to continue`
)

type MessageBoxModel struct {
	title      string
	message    string
	confirmBtn string
	Result     MessageBoxResult
}

func (m MessageBoxModel) Init() tea.Cmd {
	return nil
}

func (m MessageBoxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		default:
			m.Result = ResultConfirm
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
	title := msgTitleStyle.Height(1).Render(m.title)
	b.WriteString(title)
	b.WriteString("\n\n")

	msgtext := msgTextStyle.Height(1).Render(m.message)
	b.WriteString(msgtext)
	b.WriteString("\n\n")

	w := max(lipgloss.Width(title), lipgloss.Width(msgtext))

	footerTxt := msgBlurredBtnStyle.Render(messageFooter)
	minWidth := lipgloss.Width(footerTxt)
	if ((w - minWidth) / 2) < 0 {
		w = minWidth
	}

	rendered_footer := lipgloss.Place(w, 1, lipgloss.Center, lipgloss.Center, footerTxt)
	w = lipgloss.Width(rendered_footer) + 4

	b.WriteString(rendered_footer)
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

func NewMessageBoxModel(title, message string) MessageBoxModel {
	return MessageBoxModel{
		title:      title,
		message:    message,
		confirmBtn: "[ OK ]",
		Result:     ResultConfirm,
	}
}

func RunMessageBox(title, message string) MessageBoxResult {
	m := NewMessageBoxModel(title, message)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if result, err := p.Run(); err != nil {
		return ResultCancel // cancel on error
	} else {
		return result.(MessageBoxModel).Result
	}
}

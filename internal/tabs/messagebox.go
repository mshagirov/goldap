package tabs

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"golang.org/x/term"

	"github.com/mshagirov/goldap/internal/ui"
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
	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		default:
			m.Result = ResultConfirm
			return m, tea.Quit
		}
	case tea.MouseClickMsg:
		if msg.Button == tea.MouseLeft {
			m.Result = ResultConfirm
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MessageBoxModel) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m MessageBoxModel) render() string {
	physicalWidth, physicalHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		physicalWidth, physicalHeight = 20, 20
	}

	var b strings.Builder
	footerTxt := msgBlurredBtnStyle.Render(messageFooter)
	minWidth := lipgloss.Width(footerTxt)

	// the box is at least as wide as the footer, so let the text use it
	textWidth := max(boxTextWidth(msgBoxStyle, physicalWidth), minWidth)
	title := msgTitleStyle.Render(wrapText(m.title, textWidth))
	b.WriteString(title)
	b.WriteString("\n\n")

	msgtext := msgTextStyle.Render(wrapText(m.message, textWidth))
	b.WriteString(msgtext)
	b.WriteString("\n\n")

	w := max(lipgloss.Width(title), lipgloss.Width(msgtext))
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
		msgBoxStyle.Width(w+msgBoxStyle.GetHorizontalBorderSize()).Render(b.String()),
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

func RunMessageBox(runner ui.Runner, title, message string) (MessageBoxResult, error) {
	m := NewMessageBoxModel(title, message)
	return runMessage(func() (tea.Model, error) { return runner.Run(m) })
}

func runMessage(run func() (tea.Model, error)) (MessageBoxResult, error) {
	result, err := run()
	if err != nil {
		return ResultCancel, fmt.Errorf("Message dialog Error; %w", err)
	}
	box, ok := result.(MessageBoxModel)
	if !ok {
		return ResultCancel, fmt.Errorf("Message dialog Error; unexpected result %T", result)
	}
	return box.Result, nil
}

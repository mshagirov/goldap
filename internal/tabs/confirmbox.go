package tabs

import (
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"golang.org/x/term"

	"github.com/mshagirov/goldap/internal/mouse"
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
	case tea.KeyPressMsg:
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
	case tea.MouseClickMsg:
		if msg.Button != tea.MouseLeft {
			break
		}
		view := m.render()
		if mouse.Hit(view, m.cancelBtn, msg.X, msg.Y) {
			m.Result = ResultCancel
			return m, tea.Quit
		}
		if mouse.Hit(view, m.confirmBtn, msg.X, msg.Y) {
			m.Result = ResultConfirm
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ConfirmBoxModel) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m ConfirmBoxModel) render() string {
	physicalWidth, physicalHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		physicalWidth, physicalHeight = 20, 20
	}

	var b strings.Builder
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

	// the box is at least as wide as the buttons, so let the text use it
	textWidth := max(boxTextWidth(msgBoxStyle, physicalWidth), minWidth)
	title := msgTitleStyle.Render(wrapText(m.title, textWidth))
	b.WriteString(title)
	b.WriteString("\n\n")

	msgtext := msgTextStyle.Render(wrapText(m.message, textWidth))
	b.WriteString(msgtext)
	b.WriteString("\n\n")

	w := max(lipgloss.Width(title), lipgloss.Width(msgtext))
	padding := (w - minWidth) / 2
	if padding < 0 {
		w = minWidth
		padding = 0
	}

	buttons := lipgloss.Place(w, 1, lipgloss.Center, lipgloss.Center,
		cancelBtn+strings.Repeat(" ", buttonSpacing)+confirmBtn,
	)
	w = lipgloss.Width(buttons) + 4

	b.WriteString(buttons)
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

func NewConfirmBoxModel(title, message string) ConfirmBoxModel {
	return ConfirmBoxModel{
		title:      title,
		message:    message,
		confirm:    false,
		cancelBtn:  "[C]ancel",
		confirmBtn: "[S]ave",
		Result:     ResultConfirm,
	}
}

func RunConfirmBox(title, message string) MessageBoxResult {
	m := NewConfirmBoxModel(title, message)
	p := tea.NewProgram(m)
	if result, err := p.Run(); err != nil {
		return ResultCancel // Default to cancel on error
	} else {
		return result.(ConfirmBoxModel).Result
	}
}

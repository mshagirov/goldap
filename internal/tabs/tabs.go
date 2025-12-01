package tabs

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type Model struct {
	TabNames  []string
	Tables    []table.Model
	ActiveTab int
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "n", "tab":
			m.ActiveTab = min(m.ActiveTab+1, len(m.TabNames)-1)
			return m, nil
		case "p", "shift+tab":
			m.ActiveTab = max(m.ActiveTab-1, 0)
			return m, nil
			// case "enter", "l" , "left":
			//   selected info --> m.Tables[m.ActiveTab].SelectedRow() : 1xN slice/array
			//   return m, tea.Batch(...)
		}
	}
	m.Tables[m.ActiveTab], cmd = m.Tables[m.ActiveTab].Update(msg)
	return m, cmd
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(0, 0, 2, 0)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#DAA520", Dark: "#FFD700"}
	blurredColor      = lipgloss.Color("241")
	inactiveTabStyle  = lipgloss.NewStyle().Foreground(blurredColor).Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Foreground(highlightColor).BorderForeground(highlightColor).Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().
				BorderForeground(highlightColor).
				Align(lipgloss.Left).
				BorderStyle(lipgloss.NormalBorder()).
				Border(lipgloss.NormalBorder()).
				UnsetBorderTop()
	fillerBorderStyle = lipgloss.NewStyle().Border(
		lipgloss.Border{Bottom: "─", BottomRight: "┐"}, false, true, true, false).
		BorderForeground(highlightColor)
)

func GetTableStyle() table.Styles {
	s := table.DefaultStyles()
	hlColor := lipgloss.AdaptiveColor{Light: "#0014a8", Dark: "#265ef7"}
	s.Header = s.Header.Foreground(hlColor)
	s.Selected = s.Selected.Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}).Background(hlColor)
	return s
}

func GetTabledDimensions() (int, int) {
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth, termHeight = 20, 20
	}

	w, h := windowStyle.GetHorizontalFrameSize(), windowStyle.GetVerticalFrameSize()
	return (termWidth - w), (termHeight - 7*h)
}

func (m Model) View() string {
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth, termHeight = 20, 20
	}

	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.TabNames {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.TabNames)-1, i == m.ActiveTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	if remainingWidth := termWidth - lipgloss.Width(row); remainingWidth > 0 {
		fillStyle := fillerBorderStyle.Width(remainingWidth - 1)
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, fillStyle.Render(""))
	}
	tab_h, tab_w := windowStyle.GetVerticalFrameSize(), windowStyle.GetHorizontalFrameSize()
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(
		windowStyle.Width(termWidth - tab_w).
			Height(termHeight - tab_h*7).
			Render(m.Tables[m.ActiveTab].View()),
	)
	return docStyle.Width(termWidth).Height(termHeight - tab_h*2).Render(doc.String())
}

func Run(names []string, tables []table.Model) {

	m := Model{TabNames: names, Tables: tables}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

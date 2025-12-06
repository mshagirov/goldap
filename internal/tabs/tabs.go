package tabs

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type Model struct {
	TabNames  []string
	Tables    []table.Model
	DN        [][]string
	ActiveTab int
	Searches  map[int]textinput.Model
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	_, insearch := m.Searches[m.ActiveTab]

	var searchFocus bool
	if insearch {
		searchFocus = m.Searches[m.ActiveTab].Focused()
	} else {
		searchFocus = false
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			if insearch {
				delete(m.Searches, m.ActiveTab)
				return m, nil
			} else {
				return m, tea.Quit
			}
		case "esc":
			if insearch {
				delete(m.Searches, m.ActiveTab)
				return m, nil
			}
		case "n", "tab":
			if !insearch || !searchFocus || msg.String() != "n" {
				m.ActiveTab = (m.ActiveTab + 1) % len(m.TabNames)
				return m, nil
			}
		case "p", "shift+tab":
			if !insearch || !searchFocus || msg.String() != "p" {
				m.ActiveTab = (m.ActiveTab - 1 + len(m.TabNames)) % len(m.TabNames)
				return m, nil
			}
		case "/":
			if !insearch {
				m.Searches[m.ActiveTab] = initialSearch()
				return m, nil
			} else if !searchFocus && insearch {
				ti := m.Searches[m.ActiveTab]
				cmd = ti.Focus()
				m.Searches[m.ActiveTab] = ti
				return m, cmd
			}
		case "enter":
			if insearch && searchFocus {
				ti := m.Searches[m.ActiveTab]
				ti.Blur()
				m.Searches[m.ActiveTab] = ti
				return m, nil
			} else {
				// expand entry
				// search entry disabled
			}
		}
	}
	if insearch && searchFocus {
		m.Searches[m.ActiveTab], cmd = m.Searches[m.ActiveTab].Update(msg)
	} else {
		m.Tables[m.ActiveTab], cmd = m.Tables[m.ActiveTab].Update(msg)
	}
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
	docStyle          = lipgloss.NewStyle().Padding(0, 0, 0, 0)
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
	infoBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#3B3B3B", Dark: "#ADADAD"}).
			Align(lipgloss.Right)
	searchBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}).
			Align(lipgloss.Left)
)

func GetTableStyle() table.Styles {
	s := table.DefaultStyles()
	hlColor := lipgloss.AdaptiveColor{Light: "#0014a8", Dark: "#265ef7"}
	s.Header = s.Header.Foreground(hlColor)
	s.Selected = s.Selected.Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}).Background(hlColor)
	return s
}

func GetTableDimensions() (int, int) {
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth, termHeight = 20, 20
	}
	w, h := windowStyle.GetHorizontalFrameSize(), windowStyle.GetVerticalFrameSize()
	return (termWidth - w), (termHeight - 6*h)
}

func (m Model) View() string {
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth = 20
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

	w, h := GetTableDimensions()
	m.Tables[m.ActiveTab].SetWidth(w)
	m.Tables[m.ActiveTab].SetHeight(h)

	dn := m.DN[m.ActiveTab][m.Tables[m.ActiveTab].Cursor()]

	var searchField string
	if s, ok := m.Searches[m.ActiveTab]; ok {
		searchField = searchBarStyle.Render(fmt.Sprintf("%v", s.View()))
	}

	dnInfo := infoBarStyle.Width(w - lipgloss.Width(searchField)).Render(fmt.Sprintf("%v", dn))
	infoBar := lipgloss.JoinHorizontal(lipgloss.Top, searchField, dnInfo)

	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width(w).Height(h).
		Render(m.Tables[m.ActiveTab].View() + "\n" + infoBar),
	)
	return docStyle.Width(termWidth).Height(h).Render(doc.String())
}

func Run(names []string, tables []table.Model, dn [][]string) {
	m := Model{TabNames: names, Tables: tables, DN: dn}
	m.Searches = make(map[int]textinput.Model)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

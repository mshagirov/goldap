package tabs

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"github.com/mshagirov/goldap/ldapapi"
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

type Model struct {
	TabNames    []string
	Contents    []ldapapi.TableInfo
	DN          [][]string
	ActiveTable table.Model
	ActiveRows  []int
	ActiveTab   int
	Searches    map[int]textinput.Model
	LdapApi     *ldapapi.LdapApi
	FormInfo    *FormInfo
	quit        *bool
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) CurrentRowId() int {
	rowId, err := strconv.Atoi(m.ActiveTable.SelectedRow()[0])
	if err != nil {
		return 1
	}
	return rowId
}

func (m Model) CurrentDN() string {
	rowId := m.CurrentRowId()
	if (rowId) > len(m.DN[m.ActiveTab]) {
		return fmt.Sprintf("row %v is out of range", rowId+1)
	}
	return m.DN[m.ActiveTab][rowId-1]
}

func (m *Model) nextTab() (tea.Model, tea.Cmd) {
	m.ActiveRows[m.ActiveTab] = m.ActiveTable.Cursor()

	// next tab
	m.ActiveTab = (m.ActiveTab + 1) % len(m.TabNames)
	if _, ok := m.Searches[m.ActiveTab]; ok {
		m.ActiveTable = newTableWithFilter(m.Contents[m.ActiveTab],
			m.Searches[m.ActiveTab].Value())
	} else {
		m.ActiveTable = NewTable(m.Contents[m.ActiveTab])
	}
	m.SetCursor()

	return m, nil
}

func (m *Model) prevTab() (tea.Model, tea.Cmd) {
	m.ActiveRows[m.ActiveTab] = m.ActiveTable.Cursor()

	// previous tab
	m.ActiveTab = (m.ActiveTab - 1 + len(m.TabNames)) % len(m.TabNames)
	if _, ok := m.Searches[m.ActiveTab]; ok {
		m.ActiveTable = newTableWithFilter(m.Contents[m.ActiveTab],
			m.Searches[m.ActiveTab].Value())
	} else {
		m.ActiveTable = NewTable(m.Contents[m.ActiveTab])
	}
	m.SetCursor()
	return m, nil
}

func (m *Model) SetCursor() {
	m.ActiveTable.SetCursor(m.ActiveRows[m.ActiveTab])
}

func (m *Model) setFormInfo() {
	*m.FormInfo = FormInfo{
		DN:         m.CurrentDN(),
		TableName:  m.TabNames[m.ActiveTab],
		TableIndex: m.ActiveTab,
		RowIndices: m.ActiveRows}
	*m.quit = false
}

func (m Model) getSearchState() (bool, bool) {
	_, insearch := m.Searches[m.ActiveTab]

	var searchFocus bool
	if insearch {
		searchFocus = m.Searches[m.ActiveTab].Focused()
	} else {
		searchFocus = false
	}
	return insearch, searchFocus
}

func (m *Model) startSearch(insearch bool) (tea.Model, tea.Cmd) {
	if !insearch {
		m.Searches[m.ActiveTab] = initSearch()
		return m, nil
	}
	ti := m.Searches[m.ActiveTab]
	cmd := ti.Focus()
	m.Searches[m.ActiveTab] = ti
	return m, cmd
}

func (m *Model) blurSearch() (tea.Model, tea.Cmd) {
	ti := m.Searches[m.ActiveTab]
	ti.Blur()
	m.Searches[m.ActiveTab] = ti
	return m, nil
}

func (m *Model) stopSearch() (tea.Model, tea.Cmd) {
	rowId := m.CurrentRowId()

	delete(m.Searches, m.ActiveTab)

	m.ActiveTable = NewTable(m.Contents[m.ActiveTab])
	m.ActiveRows[m.ActiveTab] = rowId - 1
	m.SetCursor()
	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	insearch, searchFocus := m.getSearchState()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			if insearch && msg.String() != "q" {
				return m.stopSearch()
			} else if !searchFocus {
				return m, tea.Quit
			}
		case "esc":
			if insearch {
				return m.stopSearch()
			}
		case "n", "tab":
			if !insearch || !searchFocus || msg.String() != "n" {
				return m.nextTab()
			}
		case "p", "shift+tab":
			if !insearch || !searchFocus || msg.String() != "p" {
				return m.prevTab()
			}
		case "/", "?":
			if !insearch || !searchFocus {
				return m.startSearch(insearch)
			}
		case "enter":
			if insearch && searchFocus {
				return m.blurSearch()
			} else {
				m.setFormInfo()
				return m, tea.Quit
			}
		}
	}
	if insearch && searchFocus {
		m.Searches[m.ActiveTab], cmd = m.Searches[m.ActiveTab].Update(msg)
		m.ActiveTable = newTableWithFilter(m.Contents[m.ActiveTab], m.Searches[m.ActiveTab].Value())
	} else {
		m.ActiveTable, cmd = m.ActiveTable.Update(msg)
	}
	return m, cmd
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
	m.ActiveTable.SetWidth(w)
	m.ActiveTable.SetHeight(h)

	var searchField string
	if s, ok := m.Searches[m.ActiveTab]; ok {
		searchField = searchBarStyle.Render(fmt.Sprintf("%v", s.View()))
	}

	dnInfo := infoBarStyle.
		Width(w - lipgloss.Width(searchField)).
		Render(fmt.Sprintf("%v", m.CurrentDN()))
	infoBar := lipgloss.JoinHorizontal(lipgloss.Top, searchField, dnInfo)

	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width(w).Height(h).
		Render(m.ActiveTable.View() + "\n" + infoBar),
	)
	return docStyle.Width(termWidth).Height(h).Render(doc.String())
}

func NewTabsModel(names []string, contents []ldapapi.TableInfo, dn [][]string, api *ldapapi.LdapApi) Model {
	m := Model{TabNames: names, Contents: contents, DN: dn}
	m.Searches = make(map[int]textinput.Model, len(names))
	m.ActiveTable = NewTable(contents[0])
	m.ActiveRows = make([]int, len(names))
	m.LdapApi = api
	quit := true
	m.quit = &quit
	m.FormInfo = &FormInfo{}

	return m
}

func RunTabs(m Model) (FormInfo, bool) {
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return *m.FormInfo, *m.quit
}

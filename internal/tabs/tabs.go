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

type TabState struct {
	Cursor       int
	Search       textinput.Model
	IsFiltered   bool
	CachedFilter string
}

type CmdType int

const (
	QuitCmd CmdType = iota
	UpdateCmd
	AddCmd
	DeleteCmd
)

type State struct {
	Table    table.Model
	TabId    int
	TabSates []TabState
	FormInfo FormInfo
	Cmd      CmdType
}

type Model struct {
	TabNames []string
	Contents []ldapapi.TableInfo
	DN       [][]string
	State    *State
	LdapApi  *ldapapi.LdapApi
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("GoLDAP")
}

func (m Model) CurrentRowId() int {
	rowId, err := strconv.Atoi(m.State.Table.SelectedRow()[0])
	if err != nil {
		return 1
	}
	return rowId
}

func (m Model) CurrentDN() string {
	rowId := m.CurrentRowId()
	if (rowId) > len(m.DN[m.State.TabId]) {
		return fmt.Sprintf("row %v is out of range", rowId+1)
	}
	return m.DN[m.State.TabId][rowId-1]
}

func (m *Model) SetCursor() {
	m.State.Table.SetCursor(m.State.TabSates[m.State.TabId].Cursor)
}

func (m *Model) SetTable() {
	if m.State.TabSates[m.State.TabId].IsFiltered {
		m.State.Table = newTableWithFilter(m.Contents[m.State.TabId],
			m.State.TabSates[m.State.TabId].Search.Value())
	} else {
		m.State.Table = NewTable(m.Contents[m.State.TabId])
	}
	m.SetCursor()
}

func (m *Model) nextTab() (tea.Model, tea.Cmd) {
	m.State.TabSates[m.State.TabId].Cursor = m.State.Table.Cursor()
	// next tab
	m.State.TabId = (m.State.TabId + 1) % len(m.TabNames)
	m.SetTable()
	return m, nil
}

func (m *Model) prevTab() (tea.Model, tea.Cmd) {
	m.State.TabSates[m.State.TabId].Cursor = m.State.Table.Cursor()
	// previous tab
	m.State.TabId = (m.State.TabId - 1 + len(m.TabNames)) % len(m.TabNames)
	m.SetTable()
	return m, nil
}

func (m *Model) setFormInfo() { // simplify and consolidate state snapshot and reloading
	m.State.FormInfo = FormInfo{
		DN:         m.CurrentDN(),
		TableName:  m.TabNames[m.State.TabId],
		TableIndex: m.State.TabId,
	}
}

func (m Model) getSearchState() (bool, bool) {
	insearch := m.State.TabSates[m.State.TabId].IsFiltered
	searchFocus := false
	if insearch {
		searchFocus = m.State.TabSates[m.State.TabId].Search.Focused()
	}
	return insearch, searchFocus
}

func (m *Model) startSearch() (tea.Model, tea.Cmd) {
	insearch, _ := m.getSearchState()
	if !insearch {
		m.State.TabSates[m.State.TabId].Search = initSearch()
		m.State.TabSates[m.State.TabId].IsFiltered = true
		return m, nil
	}
	ti := m.State.TabSates[m.State.TabId].Search
	cmd := ti.Focus()
	m.State.TabSates[m.State.TabId].Search = ti

	return m, cmd
}

func (m *Model) blurSearch() (tea.Model, tea.Cmd) {
	ti := m.State.TabSates[m.State.TabId].Search
	ti.Blur()
	m.State.TabSates[m.State.TabId].Search = ti
	return m, nil
}

func (m *Model) stopSearch() (tea.Model, tea.Cmd) {
	rowId := m.CurrentRowId()
	m.State.TabSates[m.State.TabId].IsFiltered = false
	m.State.TabSates[m.State.TabId].Search = textinput.Model{}
	m.State.TabSates[m.State.TabId].Cursor = rowId
	m.SetTable()

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
				return m.startSearch()
			}
		case "ctrl+n":
			m.State.Cmd = AddCmd
			return m, tea.Quit
		case "ctrl+d":
			m.State.Cmd = DeleteCmd
			m.setFormInfo()
			return m, tea.Quit
		case "enter":
			if insearch && searchFocus {
				return m.blurSearch()
			} else {
				m.State.Cmd = UpdateCmd
				m.setFormInfo()
				return m, tea.Quit
			}
		}
	}
	if insearch && searchFocus {
		m.State.TabSates[m.State.TabId].Search, cmd = m.State.TabSates[m.State.TabId].Search.Update(msg)
		s := m.State.TabSates[m.State.TabId].Search.Value()
		if s != m.State.TabSates[m.State.TabId].CachedFilter {
			m.State.Table = newTableWithFilter(m.Contents[m.State.TabId], s)
		}
	} else {
		m.State.Table, cmd = m.State.Table.Update(msg)
		m.State.TabSates[m.State.TabId].Cursor = m.State.Table.Cursor()
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
		isFirst, isLast, isActive := i == 0, i == len(m.TabNames)-1, i == m.State.TabId
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
	m.State.Table.SetWidth(w)
	m.State.Table.SetHeight(h)

	var searchField string
	if m.State.TabSates[m.State.TabId].IsFiltered {
		s := m.State.TabSates[m.State.TabId].Search
		searchField = searchBarStyle.Render(s.View())
	}

	dnInfo := infoBarStyle.
		Width(w - lipgloss.Width(searchField)).
		Render(fmt.Sprintf("%v", m.CurrentDN()))
	infoBar := lipgloss.JoinHorizontal(lipgloss.Top, searchField, dnInfo)

	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width(w).Height(h).
		Render(m.State.Table.View() + "\n" + infoBar),
	)
	return docStyle.Width(termWidth).Height(h).Render(doc.String())
}

func NewTabsModel(names []string, contents []ldapapi.TableInfo, dn [][]string, api *ldapapi.LdapApi) Model {
	tabStates := make([]TabState, len(names))
	for i := range tabStates {
		tabStates[i] = TabState{
			Cursor:     0,
			Search:     textinput.Model{},
			IsFiltered: false,
		}
	}

	if colId := ldapapi.GetUsersColId("gidNumber"); colId > -1 {
		for i, tableName := range names {
			if tableName == "Users" {
				for _, row := range contents[i].Rows {
					grpDN, gotCache := api.Cache.Get(fmt.Sprintf("gidNumber=%v", row[colId+1]))
					if _, grpName, ok := ldapapi.GetFirstDnAttr(grpDN); gotCache && ok {
						row[colId+1] = fmt.Sprintf("%v(%v)", row[colId+1], grpName)
					}
				}
			}
		}
	}

	state := &State{
		Table:    NewTable(contents[0]),
		TabId:    0,
		TabSates: tabStates,
		FormInfo: FormInfo{},
		Cmd:      QuitCmd, // quit after exit by default
	}
	return Model{
		TabNames: names,
		Contents: contents,
		DN:       dn,
		State:    state,
		LdapApi:  api,
	}
}

func RunTabs(m Model) *State {
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return m.State
}

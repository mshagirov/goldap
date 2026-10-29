package tabs

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"

	"github.com/mshagirov/goldap/internal/ui"
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
	ReloadCmd
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

	viewport viewport.Model
	ready    bool
}

func (m Model) Init() tea.Cmd {
	return nil
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

func (m *Model) selectTab(tabId int) (tea.Model, tea.Cmd) {
	m.State.TabSates[m.State.TabId].Cursor = m.State.Table.Cursor()
	m.State.TabId = tabId
	m.SetTable()
	return m, nil
}

func (m *Model) nextTab() (tea.Model, tea.Cmd) {
	return m.selectTab((m.State.TabId + 1) % len(m.TabNames))
}

func (m *Model) prevTab() (tea.Model, tea.Cmd) {
	return m.selectTab((m.State.TabId - 1 + len(m.TabNames)) % len(m.TabNames))
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
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.updateTabsView(msg)
	case tea.MouseMsg:
		return m.updateMouse(msg)
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height))
			m.viewport.SetContent(m.viewContent())
			m.ready = true
		} else {
			m.viewport.SetContent(m.viewContent())
			m.State.Table.UpdateViewport()
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height)
		}
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) updateTabsView(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := m.updateTabs(msg)
	m.State.Table.UpdateViewport()
	m.viewport.SetContent(m.viewContent())
	return m, cmd
}

// Mouse wheel direction for tables and forms: a wheel-up event moves the
// selection down the list and wheel-down moves it up.
const (
	wheelMoveDown = tea.MouseWheelUp
	wheelMoveUp   = tea.MouseWheelDown
)

// updateMouse scrolls the table with the wheel; a left click selects a tab
// or a row, and a click on the already selected row acts like "enter".
func (m Model) updateMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	mouse := msg.Mouse()
	switch msg.(type) {
	case tea.MouseWheelMsg:
		switch mouse.Button {
		case wheelMoveUp:
			m.State.Table.MoveUp(1)
		case wheelMoveDown:
			m.State.Table.MoveDown(1)
		}
	case tea.MouseClickMsg:
		if mouse.Button != tea.MouseLeft {
			break
		}
		if tabId, ok := m.tabAt(mouse.X, mouse.Y); ok {
			if tabId != m.State.TabId {
				m.selectTab(tabId)
			}
			break
		}
		idx, ok := m.rowAt(mouse.Y)
		if !ok {
			break
		}
		if _, searchFocus := m.getSearchState(); searchFocus {
			m.blurSearch()
		}
		// move instead of SetCursor to keep the table's scroll offset
		if delta := idx - m.State.Table.Cursor(); delta > 0 {
			m.State.Table.MoveDown(delta)
		} else if delta < 0 {
			m.State.Table.MoveUp(-delta)
		} else {
			m.State.Cmd = UpdateCmd
			m.setFormInfo()
			return m, tea.Quit
		}
	}
	m.State.TabSates[m.State.TabId].Cursor = m.State.Table.Cursor()
	m.State.Table.UpdateViewport()
	m.viewport.SetContent(m.viewContent())
	return m, nil
}

// rowAt returns the table row index rendered at screen line y. Rows are
// identified by their number in the first column.
func (m Model) rowAt(y int) (int, bool) {
	lines := strings.Split(m.viewport.View(), "\n")
	if y < 0 || y >= len(lines) {
		return 0, false
	}
	line := strings.TrimLeft(ansi.Strip(lines[y]), "│ ")
	rowNum, _, _ := strings.Cut(line, " ")
	for i, row := range m.State.Table.Rows() {
		if len(row) > 0 && row[0] == rowNum {
			return i, true
		}
	}
	return 0, false
}

func (m *Model) updateTabs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	insearch, searchFocus := m.getSearchState()

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			if insearch && msg.String() != "q" {
				return m.stopSearch()
			} else if !searchFocus {
				m.State.Cmd = QuitCmd
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
		case "ctrl+n", "ctrl+a":
			m.State.Cmd = AddCmd
			m.setFormInfo()
			return m, tea.Quit
		case "ctrl+d":
			m.State.Cmd = DeleteCmd
			m.setFormInfo()
			return m, tea.Quit
		case "ctrl+r":
			m.State.Cmd = ReloadCmd
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

func (m *Model) viewTabs() string {
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth = 20
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, m.renderTabs()...)
	if remainingWidth := termWidth - lipgloss.Width(row); remainingWidth > 0 {
		fillStyle := fillerBorderStyle.Width(remainingWidth)
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, fillStyle.Render(""))
	}

	return row
}

// renderTabs renders each tab box; they are drawn left to right from x=0.
func (m *Model) renderTabs() []string {
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
	return renderedTabs
}

// tabAt returns the index of the tab drawn at screen cell x, y.
func (m *Model) tabAt(x, y int) (int, bool) {
	y += m.viewport.YOffset()
	if x < 0 || y < 0 {
		return 0, false
	}
	for i, t := range m.renderTabs() {
		if y >= lipgloss.Height(t) {
			return 0, false
		}
		w := lipgloss.Width(t)
		if x < w {
			return i, true
		}
		x -= w
	}
	return 0, false
}

func (m Model) View() tea.View {
	s := "\n  Initializing..."
	if m.ready {
		s = m.viewport.View()
	}
	v := tea.NewView(s)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.WindowTitle = "GoLDAP"
	return v
}

func (m *Model) viewContent() string {
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth = 20
	}

	doc := strings.Builder{}

	w, h := GetTableDimensions()
	m.State.Table.SetWidth(w)
	m.State.Table.SetHeight(h)

	var searchField string
	if m.State.TabSates[m.State.TabId].IsFiltered {
		s := m.State.TabSates[m.State.TabId].Search
		searchField = searchBarStyle.Render(s.View())
	}

	// truncate: wrapping would grow the fixed-height table frame
	dnWidth := w - lipgloss.Width(searchField)
	dnInfo := infoBarStyle.
		Width(dnWidth).
		Render(ansi.Truncate(m.CurrentDN(), dnWidth, "…"))
	infoBar := lipgloss.JoinHorizontal(lipgloss.Top, searchField, dnInfo)
	tabsRow := m.viewTabs()
	doc.WriteString(tabsRow)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width(w + windowStyle.GetHorizontalBorderSize()).Height(h + windowStyle.GetVerticalBorderSize()).
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

	if i, colId := slices.Index(names, "Users"), ldapapi.GetUsersColId("gidNumber"); (colId > -1) && (i > -1) {
		for _, row := range contents[i].Rows {
			grpDN, gotCache := api.Cache.Get(fmt.Sprintf("gidNumber=%v", row[colId+1]))
			if _, grpName, ok := ldapapi.GetFirstDnAttr(grpDN); gotCache && ok {
				row[colId+1] = fmt.Sprintf("%v(%v)", row[colId+1], grpName)
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
		viewport: viewport.New(),
	}
}

func RunTabs(runner ui.Runner, m Model) (*State, error) {
	if _, err := runner.Run(m); err != nil {
		return nil, fmt.Errorf("Tabs Error; %w", err)
	}
	return m.State, nil
}

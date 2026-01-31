package tabs

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mshagirov/goldap/ldapapi"
)

const (
	passwordPlaceholder = "••••••••"
	fieldWidth          = 50
	inputLimit          = 100
	formFooter          = `
tab/shift-tab/up/down: navigation             enter: edit/update entry
esc : cancel edit (esc twice)/exit and save   ctrl-c: exit without saving`
)

type (
	errMsg        error
	editingStatus int
)

const (
	editCANCELLED  editingStatus = iota // 0
	editCANCELLING                      // 1
	editACTIVE                          // 2
)

type formModel struct {
	title      string
	inputs     []textinput.Model
	inputNames []string
	index      int
	updated    *map[int]string
	active     map[int]struct{}
	editing    editingStatus
	err        error

	recordOnMove     bool             // record before moving to next entry
	eraseOnEdit      map[int]struct{} // erase default suggestion on edit
	alwaysRecordEdit bool             // always record entries incl. empty

	msgBox    ConfirmBoxModel
	focused   bool // true when form fields are active else activate msgBox
	updateMsg bool // use uid/cn/ou for msg if true

	viewport viewport.Model
	ready    bool // for syncing viewport dimensions
}

type FormInfo struct {
	DN         string
	TableName  string
	TableIndex int
	Api        *ldapapi.LdapApi
}

func initialFormModel(title string, attrValues, attrNames []string) formModel {
	var inputs []textinput.Model = make([]textinput.Model, len(attrNames))
	var inputNames []string = make([]string, len(attrNames))
	for i := range attrNames {
		inputs[i] = textinput.New()
		if len(attrValues[i]) > 0 {
			inputs[i].Placeholder = attrValues[i]
		} else {
			inputs[i].Placeholder = attrNames[i]
		}
		inputs[i].CharLimit = inputLimit + lipgloss.Width(attrValues[i])
		inputs[i].Width = fieldWidth
		inputs[i].Prompt = ""
		inputNames[i] = attrNames[i]
		if strings.Contains(strings.ToLower(attrNames[i]), "password") {
			inputs[i].EchoMode = textinput.EchoPassword
			inputs[i].EchoCharacter = '•'
			inputs[i].Placeholder = passwordPlaceholder
			attrValues[i] = ""
		}
	}

	inputs[0].Focus()

	return formModel{
		title:      title,
		inputs:     inputs,
		inputNames: inputNames,
		index:      0,
		active:     make(map[int]struct{}),
		err:        nil,
		focused:    true,
		msgBox:     NewMessageBox("Save changes for ...", title),
		viewport:   viewport.New(0, 0),
	}
}

func (m formModel) Init() tea.Cmd {
	return nil
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.focused {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "ctrl+c":
				return m, tea.Quit
			case "esc", "enter", "tab", "down", "up", "shift+tab", "-":
				return m.updateViewport(msg)
			}
		case tea.WindowSizeMsg:
			headerHeight := lipgloss.Height(m.headerView())
			footerHeight := lipgloss.Height(m.footerView())
			verticalMarginHeight := headerHeight + footerHeight
			if !m.ready {
				m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
				m.viewport.YPosition = headerHeight
				m.viewport.SetContent(m.viewForm())
				m.ready = true
			} else {
				m.viewport.Width = msg.Width
				m.viewport.Height = msg.Height - verticalMarginHeight
			}
		}

		if m.editing != editCANCELLED {
			return m.updateViewport(msg)
		}

		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	return m.msgBox.Update(msg)
}

func (m formModel) updateViewport(msg tea.Msg) (tea.Model, tea.Cmd) {
	mnew, cmd := m.updateForm(msg)
	m = mnew.(formModel)

	m.viewport.SetContent(m.viewForm())

	if (2*m.index + 2) > m.viewport.Height {
		m.viewport.SetYOffset(m.viewport.YOffset + 2)
	} else if 2*m.index < m.viewport.YOffset {
		m.viewport.SetYOffset(m.viewport.YOffset - 2)
	}

	return m, cmd
}

func (m formModel) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	_, isActiveField := m.active[m.index]
	hasUpdates := len(*m.updated) > 0

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if _, ok := m.active[m.index]; !ok {
				m.startEditing()
				return m, textinput.Blink
			}
			m.recordInput()
			return m, nil
		case "ctrl+c", "esc":
			if isActiveField {
				m.cancelEditing()
				return m, nil
			}
			if hasUpdates {
				m.focused = false
				m.updateConfirmMsg()
				return m, nil
			}
			return m, tea.Quit
		case "-":
			if !isActiveField && !hasUpdates {
				return m, tea.Quit
			} else if !isActiveField && hasUpdates {
				m.focused = false
				m.updateConfirmMsg()
				return m, nil
			}
		case "up", "shift+tab":
			if isActiveField && m.recordOnMove {
				m.recordInput()
			}
			m.prevInput()
		case "down", "tab":
			if isActiveField && m.recordOnMove {
				m.recordInput()
			}
			m.nextInput()
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.index].Focus()

	// handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	if _, ok := m.active[m.index]; ok || m.recordOnMove {
		m.inputs[m.index], cmds[m.index] = m.inputs[m.index].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *formModel) updateConfirmMsg() {
	if m.updateMsg {
		// dialog box
		uid_id := slices.Index(m.inputNames, "uid")
		cn_id := slices.Index(m.inputNames, "cn")
		ou_id := slices.Index(m.inputNames, "ou")

		var msgboxMsg string
		if uid_id > -1 {
			msgboxMsg = "uid=" + m.inputs[uid_id].Value()
		} else if cn_id > -1 {
			msgboxMsg = "cn=" + m.inputs[cn_id].Value()
		} else if ou_id > -1 {
			msgboxMsg = "ou" + m.inputs[ou_id].Value()
		}
		m.msgBox.message = msgboxMsg
	}
}

func (m formModel) View() string {
	if !m.focused {
		return m.msgBox.View()
	}

	if !m.ready {
		return "\n  Initializing..."
	}
	doc := strings.Builder{}
	doc.WriteString(m.headerView())
	doc.WriteString("\n")
	// doc.WriteString(m.viewForm())
	doc.WriteString(m.viewport.View())
	doc.WriteString(m.footerView())

	return doc.String()
}

func (m formModel) viewForm() string {
	doc := strings.Builder{}

	for i, val := range m.inputs {
		doc.WriteString(formFieldNameStyle.Width(30).Render(m.inputNames[i]))
		doc.WriteString("\n")
		if _, ok := m.active[i]; ok {
			val.TextStyle = formActiveStyle
			doc.WriteString(formInputPadding.Render(val.View()))
		} else if _, ok := (*m.updated)[i]; ok {
			val.TextStyle = formModifiedStyle
			doc.WriteString(formInputPadding.Render(val.View()))
		} else {
			val.TextStyle = formBlurredStyle
			doc.WriteString(formInputPadding.Render(val.View()))
		}
		doc.WriteString("\n")
	}
	return doc.String()
}

func (m formModel) headerView() string {
	return formHeaderStyle.Render(m.title)
}

func (m formModel) footerView() string {
	return formFooterStyle.Render(formFooter)
}

// nextInput focuses the next input field
func (m *formModel) nextInput() {
	m.index = min((m.index + 1), len(m.inputs)-1)
}

// prevInput focuses the previous input field
func (m *formModel) prevInput() {
	m.index = max(m.index-1, 0)
}

func (m *formModel) startEditing() {
	m.editing = editACTIVE
	m.active[m.index] = struct{}{}
	if strings.Contains(strings.ToLower(m.inputNames[m.index]), "password") {
		m.inputs[m.index].SetValue("")
		m.inputs[m.index].Placeholder = "Enter password"
		return
	}

	if _, ok := m.eraseOnEdit[m.index]; !ok {
		m.inputs[m.index].SetValue(m.inputs[m.index].Placeholder)
	}
}

func (m *formModel) cancelEditing() {
	switch m.editing {
	case editACTIVE:
		m.editing = editCANCELLING // wait for second call
	default:
		m.editing = editCANCELLED
		delete(m.active, m.index)

		if strings.Contains(strings.ToLower(m.inputNames[m.index]), "password") {
			m.inputs[m.index].Placeholder = passwordPlaceholder
		}
		m.inputs[m.index].SetValue("")
	}
}

func (m *formModel) recordInput() {
	// NEED TO ADD ENTRY VALIDATION
	delete(m.active, m.index)
	m.editing = editCANCELLED

	old_entry := m.inputs[m.index].Placeholder
	new_entry := m.inputs[m.index].Value()

	changed := old_entry != new_entry

	if changed || m.alwaysRecordEdit {
		(*m.updated)[m.index] = new_entry
		m.inputs[m.index].Placeholder = new_entry
	} else if _, ok := (*m.updated)[m.index]; !ok {
		m.inputs[m.index].SetValue("")
	}
}

func RunUpdateForm(s *State) ([]string, map[int]string) {
	fi := s.FormInfo

	var updateResult MessageBoxResult

	attrNames, attrVals := fi.Api.GetAttrWithDN(fi.DN, fi.TableName)

	updates := make(map[int]string)

	m := initialFormModel(fi.DN, attrVals, attrNames)
	m.updated = &updates
	p := tea.NewProgram(m, tea.WithAltScreen())

	result, err := p.Run()
	if err != nil {
		log.Fatal(err)
		return []string{}, nil
	}

	if msgBox, ok := result.(ConfirmBoxModel); ok {
		updateResult = msgBox.Result
	} else {
		updateResult = ResultCancel
	}
	if updateResult == ResultConfirm {
		// return updates
		return attrNames, updates
	}

	return []string{}, nil
}

func RunAddForm(s *State) ([]string, map[int]string) {
	updates := make(map[int]string)
	// fi.DN == "" (empty; AUTO GENERATED); if table not user or group allow DN editing
	// fi.TableName=from	ldapapi.TableNames AUTO GEN use defaults for each table if possible
	// TableIndex= tab ID:int
	// Api=ptr to ldapapi.LdapApi
	// auto: dn, cn
	// auto: objectClass: top, posixAccount, inetOrgPerson
	// suggest: homeDirectory
	// 	DefaultFields

	defaultAttr, ok := ldapapi.DefaultAttributes[s.FormInfo.TableName]
	if !ok {
		defaultAttr = ldapapi.UnknownTableAttributes
	}

	attrNames := make([]string, len(defaultAttr))
	attrVals := make([]string, len(defaultAttr))
	for i := range defaultAttr {
		attrNames[i] = defaultAttr[i].Name
		attrVals[i] = strings.Join(defaultAttr[i].Val, ldapapi.ValueDelimeter)
	}

	eraseOnEdit := map[int]struct{}{}
	if requiredAttr, ok := ldapapi.RequiredAttributes[s.FormInfo.TableName]; ok {
		for attr := range requiredAttr {
			eraseOnEdit[slices.Index(attrNames, attr)] = struct{}{}
		}
	}

	m := initialFormModel(fmt.Sprintf("%s: new entry", s.FormInfo.TableName), attrVals, attrNames)
	m.recordOnMove = true
	m.alwaysRecordEdit = true
	m.eraseOnEdit = eraseOnEdit
	m.updated = &updates
	m.msgBox.title = fmt.Sprintf("Adding new entry to %s ...", s.FormInfo.TableName)
	m.updateMsg = true
	p := tea.NewProgram(m, tea.WithAltScreen())

	result, err := p.Run()
	if err != nil {
		log.Fatal(err)
		return []string{}, nil
	}

	// report and confirm entries
	var updateResult MessageBoxResult
	if msgBox, ok := result.(ConfirmBoxModel); ok {
		updateResult = msgBox.Result
	} else {
		updateResult = ResultCancel
	}

	if updateResult == ResultConfirm {
		// all updated;
		for id := range attrNames {
			_, ok := updates[id]
			_, req := eraseOnEdit[id]
			if req && !ok {
				log.Printf("Error when ADDING new entry to \"%v\": missing required attribute \"%v\"", s.FormInfo.TableName, attrNames[id])
				return []string{}, nil
			}
			if !ok && !req {
				// copy default if not updated (shared attributes)
				updates[id] = attrVals[id]
			}
		}

		dn_str, err := ldapapi.ConstructDnFromUpdates(attrNames, updates, s.FormInfo.Api.Config.LdapBaseDn, s.FormInfo.TableName)
		if err != nil {
			log.Println(err)
			return []string{}, nil
		}
		s.FormInfo.DN = dn_str

		return attrNames, updates
	}

	return []string{}, nil
}

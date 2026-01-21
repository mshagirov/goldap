package tabs

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mshagirov/goldap/ldapapi"
)

const (
	passwordPlaceholder = "••••••••"
	fieldWidth          = 50
	inputLimit          = 100
	formHelpText        = `
enter : start/stop editing selection     esc/ctrl-c  : cancel editing
↑/tab/↓/shift-tab : navigation           esc/-/ctrl-c: exit
`
)

type (
	errMsg error
)

type formModel struct {
	title      string
	inputs     []textinput.Model
	inputNames []string
	index      int
	updated    *map[int]string
	active     map[int]struct{}
	focused    bool // true when form fields are active else activate msgBox
	msgBox     MessageBoxModel
	err        error
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
		inputs[i].Placeholder = attrValues[i]
		inputs[i].CharLimit = inputLimit + lipgloss.Width(attrValues[i])
		inputs[i].Width = fieldWidth
		inputs[i].Prompt = ""
		inputNames[i] = attrNames[i]
		if strings.Contains(strings.ToLower(attrNames[i]), "password") {
			inputs[i].EchoMode = textinput.EchoPassword
			inputs[i].EchoCharacter = '•'
			inputs[i].Placeholder = passwordPlaceholder
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
	}
}

func (m formModel) Init() tea.Cmd {
	return nil
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.focused {
		return m.updateForm(msg)
	}
	return m.msgBox.Update(msg)
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
			} else {
				// TO-DO: FOR PASSWORDS RE-ENTER TO CONFIRM with CORRECT ECHO
				m.recordInput()
			}
			return m, nil
		case "ctrl+c", "esc":
			if isActiveField {
				m.cancelEditing()
				return m, nil
			}
			if hasUpdates {
				m.focused = false
				return m, nil
			}
			return m, tea.Quit
		case "-":
			if !isActiveField && !hasUpdates {
				return m, tea.Quit
			} else if !isActiveField && hasUpdates {
				m.focused = false
				return m, nil
			}
		case "up", "shift+tab":
			m.prevInput()
		case "down", "tab":
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

	if _, ok := m.active[m.index]; ok {
		m.inputs[m.index], cmds[m.index] = m.inputs[m.index].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m formModel) View() string {
	if m.focused {
		return m.viewForm()
	}
	return m.msgBox.View()
}

func (m formModel) viewForm() string {
	doc := strings.Builder{}
	doc.WriteString(formTitleStyle.Render(m.title))
	doc.WriteString("\n")

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

	doc.WriteString(formHelpStyle.Render(formHelpText))

	return doc.String()
}

// nextInput focuses the next input field
func (m *formModel) nextInput() {
	m.index = (m.index + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *formModel) prevInput() {
	m.index--
	// Wrap around
	if m.index < 0 {
		m.index = len(m.inputs) - 1
	}
}

func (m *formModel) startEditing() {
	m.active[m.index] = struct{}{}
	if strings.Contains(strings.ToLower(m.inputNames[m.index]), "password") {
		m.inputs[m.index].SetValue("")
		m.inputs[m.index].Placeholder = "Enter password"
		return
	}
	m.inputs[m.index].SetValue(m.inputs[m.index].Placeholder)
}

func (m *formModel) cancelEditing() {
	delete(m.active, m.index)
	if strings.Contains(strings.ToLower(m.inputNames[m.index]), "password") {
		m.inputs[m.index].Placeholder = passwordPlaceholder
	}
	m.inputs[m.index].SetValue("")
}

func (m *formModel) recordInput() {
	// NEED TO ADD ENTRY VALIDATION
	delete(m.active, m.index)

	old_entry := m.inputs[m.index].Placeholder
	new_entry := m.inputs[m.index].Value()

	if old_entry != new_entry {
		(*m.updated)[m.index] = new_entry
		m.inputs[m.index].Placeholder = new_entry
	} else if _, ok := (*m.updated)[m.index]; !ok {
		m.inputs[m.index].SetValue("")
	}
}

func RunForm(fi FormInfo) ([]string, map[int]string) {
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

	// confirmation, if needed
	if msgBox, ok := result.(MessageBoxModel); ok {
		updateResult = msgBox.Result
	} else {
		updateResult = ResultCancel // Default fallback
	}
	if updateResult == ResultConfirm {
		return attrNames, updates
	}

	return []string{}, nil
}

package tabs

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mshagirov/goldap/ldapapi"
)

type (
	errMsg error
)

type formModel struct {
	title      string
	inputs     []textinput.Model
	inputNames []string
	index      int
	updated    map[int]struct{}
	active     map[int]struct{}
	err        error
}

type FormInfo struct {
	DN         string
	TableName  string
	TableIndex int
	RowIndices []int
	Api        *ldapapi.LdapApi
}

func RunForm(fi FormInfo) {
	//func runForm(dn, tableName string, api *ldapapi.LdapApi) {
	attrNames, attrVals := fi.Api.GetAttrWithDN(fi.DN, fi.TableName)
	p := tea.NewProgram(initialFormModel(fi.DN, attrVals, attrNames), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initialFormModel(title string, attrValues, attrNames []string) formModel {
	var inputs []textinput.Model = make([]textinput.Model, len(attrNames))
	var inputNames []string = make([]string, len(attrNames))
	for i := range attrNames {
		inputs[i] = textinput.New()
		inputs[i].Placeholder = attrValues[i]
		inputs[i].CharLimit = 45
		inputs[i].Width = 40
		inputs[i].Prompt = ""
		// inputs[i].SetValue(attrValues[i])
		inputNames[i] = attrNames[i]
	}

	inputs[0].Focus()

	return formModel{
		title:      title,
		inputs:     inputs,
		inputNames: inputNames,
		index:      0,
		active:     make(map[int]struct{}),
		updated:    make(map[int]struct{}),
		err:        nil,
	}
}

func (m formModel) Init() tea.Cmd {
	return nil
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.active[m.index] = struct{}{}
			return m, textinput.Blink
		case "ctrl+c", "esc":
			if _, ok := m.active[m.index]; ok {
				delete(m.active, m.index)
				return m, nil
			}
			return m, tea.Quit
		case "up", "shift+tab":
			m.prevInput()
		case "down", "tab":
			m.nextInput()
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.index].Focus()

	// We handle errors just like any other message
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
	doc := strings.Builder{}
	doc.WriteString(formTitleStyle.Render(m.title))
	doc.WriteString("\n")

	for i, val := range m.inputs {
		doc.WriteString(formFieldNameStyle.Width(30).Render(m.inputNames[i]))
		doc.WriteString("\n")
		if _, ok := m.active[i]; ok {
			doc.WriteString(formActiveStyle.Render(val.View()))
		} else {
			doc.WriteString(formBlurredStyle.Render(val.View()))
		}
		doc.WriteString("\n")
	}

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

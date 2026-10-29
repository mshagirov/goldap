// Package ui keeps one terminal session alive while switching between screens.
package ui

import (
	"fmt"
	"sync"

	tea "charm.land/bubbletea/v2"
)

type Runner interface {
	Run(tea.Model) (tea.Model, error)
}

type Session struct {
	program *tea.Program
	done    chan struct{}
	err     error // Read only after done closes.
	close   sync.Once
}

// NewSession starts the renderer once. Run displays one screen at a time;
// screens still use tea.Quit to return to their caller, not to exit the terminal.
func NewSession(options ...tea.ProgramOption) *Session {
	s := &Session{done: make(chan struct{})}
	options = append(options, tea.WithFilter(func(_ tea.Model, msg tea.Msg) tea.Msg {
		switch msg.(type) {
		case closeSession:
			return tea.QuitMsg{}
		}
		return msg
	}))
	s.program = tea.NewProgram(host{}, options...)
	go func() {
		_, s.err = s.program.Run()
		close(s.done)
	}()
	return s
}

func (s *Session) Run(model tea.Model) (tea.Model, error) {
	reply := make(chan screenResult, 1)
	s.program.Send(screenRequest{model: model, reply: reply})
	select {
	case result := <-reply:
		return result.model, result.err
	case <-s.done:
		if s.err != nil {
			return nil, s.err
		}
		return nil, fmt.Errorf("UI session is closed")
	}
}

func (s *Session) Close() {
	s.close.Do(func() { s.program.Send(closeSession{}) })
	<-s.done
}

type screenRequest struct {
	model tea.Model
	reply chan screenResult
}

type screenResult struct {
	model tea.Model
	err   error
}

type screenDone struct{ generation uint64 }
type closeSession struct{}

type host struct {
	child      tea.Model
	reply      chan screenResult
	size       tea.WindowSizeMsg
	view       tea.View
	generation uint64
}

func (m host) Init() tea.Cmd { return nil }

func (m host) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenRequest:
		if m.child != nil {
			msg.reply <- screenResult{err: fmt.Errorf("UI screen is already active")}
			return m, nil
		}
		m.child, m.reply = msg.model, msg.reply
		m.generation++
		init := m.child.Init()
		var resize tea.Cmd
		if m.size.Width > 0 && m.size.Height > 0 {
			m.child, resize = m.child.Update(m.size)
		}
		return m, screenCmd(m.generation, tea.Batch(init, resize))
	case screenDone:
		if m.child != nil && msg.generation == m.generation {
			// Freeze the last frame before returning ownership to the caller.
			// During LDAP work we neither expose the shell nor touch the old model.
			m.view = m.child.View()
			result := m.child
			m.child = nil
			m.reply <- screenResult{model: result}
			m.reply = nil
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.size = msg
	}
	if m.child == nil {
		return m, nil
	}
	var cmd tea.Cmd
	m.child, cmd = m.child.Update(msg)
	return m, screenCmd(m.generation, cmd)
}

// Scope quit commands (including nested batches) to the screen that issued
// them. A delayed quit from the previous screen must not close the next screen.
func screenCmd(generation uint64, cmd tea.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	return func() tea.Msg {
		switch msg := cmd().(type) {
		case tea.QuitMsg:
			return screenDone{generation: generation}
		case tea.BatchMsg:
			cmds := make(tea.BatchMsg, len(msg))
			for i, cmd := range msg {
				cmds[i] = screenCmd(generation, cmd)
			}
			return cmds
		default:
			return msg
		}
	}
}

func (m host) View() tea.View {
	if m.child != nil {
		return m.child.View()
	}
	v := m.view
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

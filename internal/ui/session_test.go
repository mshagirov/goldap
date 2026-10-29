package ui

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

type testScreen struct {
	text string
	size tea.WindowSizeMsg
}

func (m testScreen) Init() tea.Cmd { return tea.Quit }

func (m testScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.size = size
	}
	return m, nil
}

func (m testScreen) View() tea.View {
	v := tea.NewView(m.text)
	v.AltScreen = true
	return v
}

func TestSessionKeepsProgramAliveBetweenScreens(t *testing.T) {
	session := NewSession(
		tea.WithInput(nil), tea.WithOutput(io.Discard), tea.WithoutRenderer(),
		tea.WithoutSignals(), tea.WithWindowSize(80, 24),
	)
	t.Cleanup(session.Close)
	for _, name := range []string{"tabs", "add form", "confirmation", "result", "tabs"} {
		t.Run(name, func(t *testing.T) {
			done := make(chan screenResult, 1)
			go func() {
				model, err := session.Run(testScreen{text: name})
				done <- screenResult{model: model, err: err}
			}()
			select {
			case result := <-done:
				if result.err != nil {
					t.Fatal(result.err)
				}
				if result.model.(testScreen).text != name {
					t.Fatal("returned the wrong screen")
				}
			case <-time.After(5 * time.Second):
				t.Fatal("screen did not return")
			}
			select {
			case <-session.done:
				t.Fatal("screen exit stopped the terminal session")
			default:
			}
		})
	}
	session.Close()
	session.Close() // Shutdown is safe to repeat on error/fallback paths.
	if _, err := session.Run(testScreen{}); err == nil {
		t.Fatal("running a screen after shutdown must fail")
	}
}

func TestHostFreezesScreenAndResizesNextScreen(t *testing.T) {
	var model tea.Model = host{}
	model, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	reply := make(chan screenResult, 1)
	model, _ = model.Update(screenRequest{model: testScreen{text: "last frame"}, reply: reply})
	if model.(host).child.(testScreen).size.Width != 80 {
		t.Fatal("screen did not receive the current terminal size")
	}
	generation := model.(host).generation
	model, _ = model.Update(screenDone{generation: generation})
	<-reply
	if model.(host).child != nil {
		t.Fatal("host retained ownership of completed screen")
	}
	if !model.View().AltScreen {
		t.Fatal("alternate screen must stay active between screens")
	}
	if model.View().Content != "last frame" {
		t.Fatal("last screen must remain visible between screens")
	}
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model, _ = model.Update(screenRequest{model: testScreen{text: "next frame"}, reply: reply})
	model, _ = model.Update(screenDone{generation: generation})
	if model.(host).child == nil {
		t.Fatal("a delayed quit from the last screen closed the next screen")
	}
	if size := model.(host).child.(testScreen).size; size.Width != 100 || size.Height != 30 {
		t.Errorf("next screen size = %v, want 100x30", size)
	}
}

type terminalOutput struct {
	sync.Mutex
	bytes.Buffer
}

func (b *terminalOutput) Write(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	return b.Buffer.Write(p)
}

func (b *terminalOutput) snapshot() string {
	b.Lock()
	defer b.Unlock()
	return b.Buffer.String()
}

func TestSessionDoesNotLeaveAlternateScreenBetweenScreens(t *testing.T) {
	var output terminalOutput
	session := NewSession(
		tea.WithInput(nil), tea.WithOutput(&output), tea.WithoutSignals(),
		tea.WithEnvironment([]string{"TERM=xterm-256color"}), tea.WithWindowSize(80, 24),
	)
	t.Cleanup(session.Close)
	for _, name := range []string{"tabs", "form", "delete confirmation", "result", "tabs"} {
		if _, err := session.Run(testScreen{text: name}); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(output.snapshot(), "\x1b[?1049l") {
			t.Fatal("screen transition restored the shell")
		}
	}
	session.Close()
	text := output.snapshot()
	if enters, exits := strings.Count(text, "\x1b[?1049h"), strings.Count(text, "\x1b[?1049l"); enters != 1 || exits != 1 {
		t.Errorf("alternate screen entered %d times and exited %d times, want once each", enters, exits)
	}
}

func TestScreenCmdScopesNestedQuits(t *testing.T) {
	cmd := screenCmd(7, tea.Batch(tea.Quit, tea.Batch(tea.Quit, tea.Quit)))
	var check func(tea.Msg)
	quits := 0
	check = func(msg tea.Msg) {
		switch msg := msg.(type) {
		case tea.BatchMsg:
			for _, cmd := range msg {
				check(cmd())
			}
		case screenDone:
			if msg.generation != 7 {
				t.Errorf("generation = %d, want 7", msg.generation)
			}
			quits++
		default:
			t.Fatalf("unexpected command result %T", msg)
		}
	}
	check(cmd())
	if quits != 3 {
		t.Errorf("scoped quit count = %d, want 3", quits)
	}
}

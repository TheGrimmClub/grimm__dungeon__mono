package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/TheGrimmClub/grimm__dungeon__mono/content"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/engine"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/session"
)

func testModel(t *testing.T) model {
	t.Helper()
	w, err := world.Load(content.FS, content.WorldGlob)
	if err != nil {
		t.Fatalf("load world: %v", err)
	}
	return newModel(session.New(engine.New(w), ""), "INTRO")
}

// step applies a message and returns the concrete model back.
func step(t *testing.T, m model, msg tea.Msg) (model, tea.Cmd) {
	t.Helper()
	next, cmd := m.Update(msg)
	return next.(model), cmd
}

func TestModelSurvivesInputBeforeAndZeroSize(t *testing.T) {
	// A key (and a submit) arriving before any WindowSizeMsg, plus a 0x0 size
	// from a dumb pty, must not panic (regression: viewport GotoBottom).
	m := testModel(t)
	m.input.SetValue("look")
	m, _ = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	m, _ = step(t, m, tea.WindowSizeMsg{Width: 0, Height: 0})
	if !m.ready {
		t.Fatal("model should be ready after a (zero) WindowSizeMsg")
	}
	_ = m.View() // must not panic
}

func TestModelSubmitHistoryAndQuit(t *testing.T) {
	m := testModel(t)

	m, _ = step(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	if !m.ready {
		t.Fatal("model not ready after WindowSizeMsg")
	}

	// Submit a movement command.
	m.input.SetValue("go north")
	m, _ = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if joined := strings.Join(m.transcript, "\n"); !strings.Contains(joined, "Halle der schlafenden Maschinen") {
		t.Errorf("transcript missing engine output:\n%s", joined)
	}

	// Up recalls the last submission into a cleared input.
	m, _ = step(t, m, tea.KeyMsg{Type: tea.KeyUp})
	if got := m.input.Value(); got != "go north" {
		t.Errorf("history up = %q, want 'go north'", got)
	}

	// /quit returns a command that yields tea.QuitMsg.
	m.input.SetValue("/quit")
	_, cmd := step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("/quit produced no command")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Error("/quit command did not yield tea.QuitMsg")
	}
}

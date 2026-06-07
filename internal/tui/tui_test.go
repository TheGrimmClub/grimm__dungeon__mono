package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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

func testModelWithWorkDir(t *testing.T, dir string) model {
	t.Helper()
	w, err := world.Load(content.FS, content.WorldGlob)
	if err != nil {
		t.Fatalf("load world: %v", err)
	}
	s := session.New(engine.New(w), "")
	s.SetWorkDir(dir)
	return newModel(s, "INTRO")
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

func TestHUDAppearsWithHelmet(t *testing.T) {
	m := testModel(t)
	m, _ = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	if strings.Contains(m.View(), "INVENTAR") {
		t.Error("HUD should be hidden before the helmet is worn")
	}

	m.input.SetValue("wear 1") // take + wear the helmet (hud:true)
	m, _ = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	v := m.View()
	if !strings.Contains(v, "INVENTAR") || !strings.Contains(v, "KARTE") {
		t.Errorf("HUD boxes should appear once the helmet is worn:\n%s", v)
	}
	if !strings.Contains(v, "Helm mit Stirnlampe") {
		t.Errorf("inventory box should list the helmet:\n%s", v)
	}
}

func TestHUDHiddenOnNarrowTerminal(t *testing.T) {
	m := testModel(t)
	m, _ = step(t, m, tea.WindowSizeMsg{Width: 40, Height: 30})
	m.input.SetValue("wear 1")
	m, _ = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if strings.Contains(m.View(), "INVENTAR") {
		t.Error("HUD should stay hidden on a narrow terminal even with the helmet")
	}
}

func TestTerminalSuspendsTUI(t *testing.T) {
	m := testModelWithWorkDir(t, t.TempDir())
	m, _ = step(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})

	m.input.SetValue("/terminal")
	_, cmd := step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("/terminal should produce a tea.ExecProcess command")
	}
}

func TestExecDoneAppendsNote(t *testing.T) {
	m := testModel(t)
	m, _ = step(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	m, _ = step(t, m, execDoneMsg{after: "zurück im Verlies"})
	if joined := strings.Join(m.transcript, "\n"); !strings.Contains(joined, "zurück im Verlies") {
		t.Errorf("resume note not appended:\n%s", joined)
	}
}

// TestNoLineReachesLastColumn guards the line-deletion bug: NO rendered line of
// the full View may reach the terminal's last column, or it auto-wraps and
// Bubble Tea's renderer desyncs and erases the line above. We exercise the long
// verb hint / room text (look) with and without the HUD, plus a long input.
func TestNoLineReachesLastColumn(t *testing.T) {
	for _, width := range []int{40, 64, 80, 100} {
		// Without and with the helmet (HUD narrows the viewport).
		for _, helmet := range []bool{false, true} {
			m := testModelWithWorkDir(t, t.TempDir())
			m, _ = step(t, m, tea.WindowSizeMsg{Width: width, Height: 24})
			if helmet {
				m.input.SetValue("wear 1")
				m, _ = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
			}
			// Produce long output and a long pending input.
			m.input.SetValue("look")
			m, _ = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})
			m.input.SetValue(strings.Repeat("go north and then ", 12))
			m.relayout()

			for i, l := range strings.Split(m.View(), "\n") {
				if w := lipgloss.Width(l); w >= width {
					t.Errorf("width=%d helmet=%v: line %d reaches the last column (%d>=%d): %q",
						width, helmet, i, w, width, l)
				}
			}
		}
	}
}

// tallModel builds a model whose transcript overflows a short terminal, so the
// viewport is genuinely scrollable.
func tallModel(t *testing.T) model {
	t.Helper()
	w, err := world.Load(content.FS, content.WorldGlob)
	if err != nil {
		t.Fatalf("load world: %v", err)
	}
	intro := strings.TrimRight(strings.Repeat("Eine Zeile im Verlies.\n", 40), "\n")
	return newModel(session.New(engine.New(w), ""), intro)
}

// TestTypingDoesNotScrollViewport is the real line-deletion repro: typing a
// keyword that contains a viewport scroll-key letter (k/j/u/d/f/b) must not move
// the transcript. Previously the keystroke was forwarded to the viewport too.
func TestTypingDoesNotScrollViewport(t *testing.T) {
	m := tallModel(t)
	m, _ = step(t, m, tea.WindowSizeMsg{Width: 80, Height: 8})

	before := m.vp.YOffset
	if before == 0 {
		t.Fatal("test setup: viewport is not scrolled (nothing to scroll)")
	}
	for _, r := range "look" { // the trailing 'k' used to scroll up
		m, _ = step(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	if m.vp.YOffset != before {
		t.Errorf("typing scrolled the viewport: YOffset %d -> %d", before, m.vp.YOffset)
	}
	if m.input.Value() != "look" {
		t.Errorf("typed text should land in the input, got %q", m.input.Value())
	}
}

func TestPgUpStillScrolls(t *testing.T) {
	m := tallModel(t)
	m, _ = step(t, m, tea.WindowSizeMsg{Width: 80, Height: 8})
	before := m.vp.YOffset
	m, _ = step(t, m, tea.KeyMsg{Type: tea.KeyPgUp})
	if m.vp.YOffset >= before {
		t.Errorf("PgUp should scroll up: YOffset %d -> %d", before, m.vp.YOffset)
	}
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

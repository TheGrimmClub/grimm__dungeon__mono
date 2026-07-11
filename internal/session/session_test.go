package session_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/TheGrimmClub/grimm__dungeon__mono/content"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/engine"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/state"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/session"
)

func newSession(t *testing.T, savePath string) *session.Session {
	t.Helper()
	w, err := world.Load(content.FS, content.WorldGlob)
	if err != nil {
		t.Fatalf("load world: %v", err)
	}
	return session.New(engine.New(w), savePath)
}

func TestHelpListsCommandsAndVerbs(t *testing.T) {
	s := newSession(t, "")
	out := s.Submit("/help").Output
	for _, want := range []string{"/help", "/quit", "/save", "look", "go <dir>", "wear <n|name>"} {
		if !strings.Contains(out, want) {
			t.Errorf("/help missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "/antigravity") {
		t.Error("hidden command leaked into /help")
	}
}

func TestEmptyInputShowsInfoBlock(t *testing.T) {
	s := newSession(t, "")
	if out := s.Submit("   ").Output; !strings.Contains(out, "/help") {
		t.Errorf("empty input did not show the info block:\n%s", out)
	}
}

func TestUnknownCommand(t *testing.T) {
	s := newSession(t, "")
	out := s.Submit("/fly").Output
	if !strings.Contains(out, "/fly") || !strings.Contains(out, "/help") {
		t.Errorf("unknown command message missing name or hint:\n%s", out)
	}
}

func TestQuitSignals(t *testing.T) {
	s := newSession(t, "")
	res := s.Submit("/quit")
	if !res.Quit {
		t.Error("/quit did not signal Quit")
	}
	if !strings.Contains(res.Output, "Fackeln") {
		t.Errorf("/quit missing goodbye:\n%s", res.Output)
	}
}

func TestEasterEggViaText(t *testing.T) {
	s := newSession(t, "")
	if out := s.Submit("import antigravity").Output; !strings.Contains(out, "schwerelos") {
		t.Errorf("easter egg did not fire:\n%s", out)
	}
}

func TestFreeTextRoutesToEngine(t *testing.T) {
	s := newSession(t, "")
	if out := s.Submit("go north").Output; !strings.Contains(out, "Halle der schlafenden Maschinen") {
		t.Errorf("free text did not reach the engine:\n%s", out)
	}
}

func TestSaveWritesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "save.syon")
	s := newSession(t, path)
	s.Submit("go north")
	if out := s.Submit("/save").Output; !strings.Contains(out, "versiegelt") {
		t.Fatalf("save did not confirm:\n%s", out)
	}
	if !state.Exists(path) {
		t.Errorf("save file not written at %s", path)
	}
}

func TestSaveDisabledWhenNoPath(t *testing.T) {
	s := newSession(t, "")
	if out := s.Submit("/save").Output; !strings.Contains(out, "deaktiviert") {
		t.Errorf("expected save-disabled message:\n%s", out)
	}
}

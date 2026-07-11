package session_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/TheGrimmClub/grimm__toolbox__mono/tools/alchemist"
)

func TestAlchemistNoWorkDir(t *testing.T) {
	s := newSession(t, "")
	if out := s.Submit("/alchemist look").Output; !strings.Contains(out, "keinen Kessel") {
		t.Errorf("expected no-cauldron message without a work dir:\n%s", out)
	}
}

func TestAlchemistInitInGame(t *testing.T) {
	if !alchemist.Available() {
		t.Skip("git not installed")
	}
	dir := filepath.Join(t.TempDir(), "work")
	s := newSession(t, "")
	s.SetWorkDir(dir)

	if out := s.Submit("/alchemist init").Output; !strings.Contains(out, "Repository ist bereit") {
		t.Errorf("/alchemist init did not create the repo:\n%s", out)
	}
	if out := s.Submit("/alchemist look").Output; !strings.Contains(out, "git status") && !strings.Contains(out, "ruht") {
		t.Errorf("/alchemist look failed:\n%s", out)
	}
	// The grimoire is shown with no verb.
	if out := s.Submit("/alchemist").Output; !strings.Contains(out, "Grimoire") {
		t.Errorf("/alchemist (no verb) should show the grimoire:\n%s", out)
	}
}

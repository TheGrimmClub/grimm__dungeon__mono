package session_test

import (
	"strings"
	"testing"
)

func TestClassListingAndChoice(t *testing.T) {
	s := newSession(t, "")

	// No arg lists the paths.
	if out := s.Submit("/class").Output; !strings.Contains(out, "alchemist") || !strings.Contains(out, "Human") {
		t.Errorf("/class did not list the paths:\n%s", out)
	}
	// An unknown path is rejected.
	if out := s.Submit("/class wizard").Output; !strings.Contains(out, "kein Pfad") {
		t.Errorf("unknown class should be rejected:\n%s", out)
	}
	// Choosing changes the title (which the prompt reads).
	if out := s.Submit("/class alchemist").Output; !strings.Contains(out, "Alchemist") {
		t.Errorf("choosing alchemist failed:\n%s", out)
	}
	if got := s.Game().Title(); got != "Alchemist" {
		t.Errorf("title after choosing = %q, want Alchemist", got)
	}
}

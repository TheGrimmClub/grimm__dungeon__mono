package app

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestNewSessionIntro(t *testing.T) {
	_, intro, err := NewSession("")
	if err != nil {
		t.Fatalf("NewSession: %v", err)
	}
	for _, want := range []string{
		`/ ___|`,                  // ASCII banner fragment
		"Willkommen, Human",       // welcome addresses the class-less player
		"Das verwunschene Tor",    // starting room
		"[1] Helm mit Stirnlampe", // numbered loot
		"Sprich mit dem Verlies",  // verb hint
	} {
		if !strings.Contains(intro, want) {
			t.Errorf("intro missing %q:\n%s", want, intro)
		}
	}
}

func TestSaveThenContinue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "save.yaml")

	// First session: move north and save.
	s1, _, err := NewSession(path)
	if err != nil {
		t.Fatal(err)
	}
	s1.Submit("go north")
	if out := s1.Submit("/save").Output; !strings.Contains(out, "versiegelt") {
		t.Fatalf("save did not confirm:\n%s", out)
	}

	// Second session over the same path resumes in the hall.
	_, intro2, err := NewSession(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(intro2, "Du nimmst deinen Weg wieder auf") {
		t.Errorf("second intro did not report continuing:\n%s", intro2)
	}
	if !strings.Contains(intro2, "Halle der schlafenden Maschinen") {
		t.Errorf("second session did not resume in the saved room:\n%s", intro2)
	}
}

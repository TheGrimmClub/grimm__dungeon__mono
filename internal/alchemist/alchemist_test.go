package alchemist

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newRepo(t *testing.T) (*Alchemist, string) {
	t.Helper()
	if !Available() {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	a := New(dir)
	if _, err := a.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	// Make commits possible without relying on global git identity.
	a.git(0, "config", "user.email", "tester@grimm.local")
	a.git(0, "config", "user.name", "Grimm Tester")
	return a, dir
}

func TestInitCreatesRepo(t *testing.T) {
	_, dir := newRepo(t)
	if _, err := os.Stat(filepath.Join(dir, ".git", "HEAD")); err != nil {
		t.Errorf(".git/HEAD missing after init: %v", err)
	}
}

func TestBrewNeedsIngredients(t *testing.T) {
	a, _ := newRepo(t)
	msg, err := a.Brew("leer")
	if err != nil {
		t.Fatalf("Brew: %v", err)
	}
	if !strings.Contains(msg, "nichts zu brauen") {
		t.Errorf("empty brew should report nothing to brew, got: %s", msg)
	}
}

func TestBrewCommitsAndLook(t *testing.T) {
	a, dir := newRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "trank.txt"), []byte("ein Tropfen"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Untracked file shows up in look.
	if out, _ := a.Look(); !strings.Contains(out, "trank.txt") {
		t.Errorf("look should show the new file:\n%s", out)
	}
	msg, err := a.Brew("Erster Trank")
	if err != nil {
		t.Fatalf("Brew: %v", err)
	}
	if !strings.Contains(msg, "Erster Trank") {
		t.Errorf("brew message should name the potion: %s", msg)
	}
	// After brewing, the tree is clean.
	if out, _ := a.Look(); !strings.Contains(out, "ruht") {
		t.Errorf("look after brew should be clean:\n%s", out)
	}
}

func TestDispatchHelp(t *testing.T) {
	if !Available() {
		t.Skip("git not installed")
	}
	out, err := Dispatch(New(t.TempDir()), nil)
	if err != nil {
		t.Fatalf("Dispatch help: %v", err)
	}
	if !strings.Contains(out, "Grimoire") || !strings.Contains(out, "brew") {
		t.Errorf("help should show the grimoire:\n%s", out)
	}
}

func TestDispatchUnknownVerb(t *testing.T) {
	if !Available() {
		t.Skip("git not installed")
	}
	if _, err := Dispatch(New(t.TempDir()), []string{"teleport"}); err == nil {
		t.Error("unknown verb should error")
	}
}

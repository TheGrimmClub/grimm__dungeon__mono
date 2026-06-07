package session_test

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestTerminalReturnsExec(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "work")
	s := newSession(t, "")
	s.SetWorkDir(dir)

	res := s.Submit("/terminal")
	if res.Exec == nil {
		t.Fatalf("/terminal should return an exec request; output=%q", res.Output)
	}
	if res.Exec.Dir != dir {
		t.Errorf("exec dir = %q, want %q", res.Exec.Dir, dir)
	}
	if res.Exec.Name == "" {
		t.Error("exec should name a shell binary")
	}
}

func TestTerminalNeedsWorkDir(t *testing.T) {
	s := newSession(t, "")
	res := s.Submit("/terminal")
	if res.Exec != nil {
		t.Error("/terminal without a work dir must not exec")
	}
	if !strings.Contains(res.Output, "keinen Kessel") {
		t.Errorf("expected no-cauldron message:\n%s", res.Output)
	}
}

func TestBookUsageAndExec(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "work")
	s := newSession(t, "")
	s.SetWorkDir(dir)
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "true") // a real, harmless binary on unix

	if res := s.Submit("/book"); res.Exec != nil || !strings.Contains(res.Output, "Welches Buch") {
		t.Errorf("/book with no file should show usage, got exec=%v out=%q", res.Exec, res.Output)
	}

	res := s.Submit("/book zauber.py")
	if res.Exec == nil {
		t.Fatalf("/book <file> should exec an editor; out=%q", res.Output)
	}
	if len(res.Exec.Args) != 1 || res.Exec.Args[0] != "zauber.py" {
		t.Errorf("editor args = %v, want [zauber.py]", res.Exec.Args)
	}
	if res.Exec.Dir != dir {
		t.Errorf("editor dir = %q, want %q", res.Exec.Dir, dir)
	}
}

func TestBookKeepsFileInsideWorkDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "work")
	s := newSession(t, "")
	s.SetWorkDir(dir)
	t.Setenv("EDITOR", "true")

	res := s.Submit("/book ../../etc/passwd")
	if res.Exec == nil {
		t.Fatal("expected exec")
	}
	if res.Exec.Args[0] != "passwd" {
		t.Errorf("path traversal not stripped: %q", res.Exec.Args[0])
	}
}

func TestLaunchersListedInHelp(t *testing.T) {
	s := newSession(t, "")
	out := s.Submit("/help").Output
	if !strings.Contains(out, "/terminal") || !strings.Contains(out, "/book") {
		t.Errorf("/help should list the new surfaces:\n%s", out)
	}
}

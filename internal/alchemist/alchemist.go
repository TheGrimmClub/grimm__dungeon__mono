// Package alchemist is the "git as potion-brewing" tool: a thin, German-flavoured
// wrapper over git that teaches version control as crafting (req R006). It is the
// shared library behind both the standalone cmd/alchemist binary and grimm's
// in-game /alchemist command (decision D006). Each potion verb maps to real git.
package alchemist

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/runner"
)

// Verb is one potion action and the git reality it teaches.
type Verb struct {
	Name    string
	Git     string
	Summary string
}

// Grimoire is the canonical potion -> git mapping (shown by help).
var Grimoire = []Verb{
	{"init", "git init", "erschaffe deinen Kessel (ein neues Repository)"},
	{"add", "git add", "lege eine Zutat in den Kessel (stage)"},
	{"brew", "git add -A + git commit", "braue den Trank — halte den Stand fest"},
	{"bottle", "git push", "fülle den Trank ab und schicke ihn fort"},
	{"discard", "git restore", "schütte die offenen Änderungen weg"},
	{"clean", "git clean", "fege unversionierte Reste vom Tisch"},
	{"look", "git status", "betrachte den Kessel — was brodelt gerade?"},
}

// Alchemist runs git in a fixed working directory (the student's repo).
type Alchemist struct {
	dir string
}

// New binds an Alchemist to a working directory.
func New(dir string) *Alchemist { return &Alchemist{dir: dir} }

// Available reports whether git is installed.
func Available() bool { return runner.Available("git") }

func (a *Alchemist) git(timeout time.Duration, args ...string) runner.Result {
	return runner.Run(context.Background(), a.dir, "", timeout, "git", args...)
}

// Init creates the repository (the cauldron).
func (a *Alchemist) Init() (string, error) {
	res := a.git(0, "init")
	if err := fail(res); err != nil {
		return "", err
	}
	return "Ein neuer Kessel entsteht — dein Repository ist bereit (git init).", nil
}

// Add stages ingredients. With no paths, stages everything.
func (a *Alchemist) Add(paths []string) (string, error) {
	args := append([]string{"add"}, defaultPaths(paths)...)
	res := a.git(0, args...)
	if err := fail(res); err != nil {
		return "", err
	}
	return "Zutaten in den Kessel gelegt (git add).", nil
}

// Brew stages everything and commits (records the state).
func (a *Alchemist) Brew(message string) (string, error) {
	if strings.TrimSpace(message) == "" {
		message = "Ein namenloser Trank"
	}
	if res := a.git(0, "add", "-A"); fail(res) != nil {
		return "", fail(res)
	}
	res := a.git(0, "commit", "-m", message)
	if res.ExitCode != 0 {
		// The most common cause: nothing staged.
		return "Es gibt nichts zu brauen — füge erst Zutaten hinzu (git add).", nil
	}
	if err := fail(res); err != nil {
		return "", err
	}
	return fmt.Sprintf("Trank gebraut: »%s« (git commit).", message), nil
}

// Bottle pushes to the remote (a longer timeout for the network).
func (a *Alchemist) Bottle() (string, error) {
	res := a.git(30*time.Second, "push")
	if res.TimedOut {
		return "", fmt.Errorf("das Abfüllen dauerte zu lange (Netzwerk?)")
	}
	if res.ExitCode != 0 {
		return "Der Trank ließ sich nicht abfüllen — gibt es schon ein Fass (remote)?\n" +
			strings.TrimSpace(res.Stderr), nil
	}
	if err := fail(res); err != nil {
		return "", err
	}
	return "Trank abgefüllt und fortgeschickt (git push).", nil
}

// Discard throws away unstaged changes.
func (a *Alchemist) Discard() (string, error) {
	res := a.git(0, "restore", ".")
	if err := fail(res); err != nil {
		return "", err
	}
	return "Die offenen Änderungen sind weggeschüttet (git restore).", nil
}

// Clean removes untracked files.
func (a *Alchemist) Clean() (string, error) {
	res := a.git(0, "clean", "-fd")
	if err := fail(res); err != nil {
		return "", err
	}
	out := strings.TrimSpace(res.Stdout)
	if out == "" {
		return "Der Tisch war schon sauber (git clean).", nil
	}
	return "Unversionierte Reste vom Tisch gefegt (git clean):\n" + out, nil
}

// Look reports the working-tree status.
func (a *Alchemist) Look() (string, error) {
	res := a.git(0, "status", "--short", "--branch")
	if err := fail(res); err != nil {
		return "", err
	}
	out := strings.TrimSpace(res.Stdout)
	// `--branch` always prints a "## branch" header; the tree is clean when that
	// is the only line.
	if changeLines(out) == 0 {
		return "Der Kessel ruht — nichts brodelt (git status).", nil
	}
	return "Im Kessel (git status):\n" + out, nil
}

// changeLines counts status lines that describe actual changes (not the "##"
// branch header).
func changeLines(status string) int {
	n := 0
	for _, line := range strings.Split(status, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "##") {
			continue
		}
		n++
	}
	return n
}

func defaultPaths(paths []string) []string {
	if len(paths) == 0 {
		return []string{"."}
	}
	return paths
}

// fail turns a runner result into an error for git invocations that should
// simply succeed (start failure, timeout, or non-zero exit).
func fail(res runner.Result) error {
	switch {
	case res.Err != nil:
		return fmt.Errorf("git ließ sich nicht ausführen: %w", res.Err)
	case res.TimedOut:
		return fmt.Errorf("git brauchte zu lange")
	case res.ExitCode != 0:
		msg := strings.TrimSpace(res.Stderr)
		if msg == "" {
			msg = strings.TrimSpace(res.Stdout)
		}
		return fmt.Errorf("git meldet: %s", msg)
	}
	return nil
}

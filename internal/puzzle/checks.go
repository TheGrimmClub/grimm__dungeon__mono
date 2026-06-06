package puzzle

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/runner"
)

// --- artifact: proof of a real-world action ---

type artifactCheck struct {
	file    string
	command string
	expect  string
}

func (c artifactCheck) Verify(ctx context.Context, in Input) Result {
	if c.file == "" && c.command == "" {
		return Result{Detail: "Diese Prüfung ist nicht eingerichtet."}
	}

	if c.file != "" {
		path := c.file
		if !filepath.IsAbs(path) {
			path = filepath.Join(in.WorkDir, path)
		}
		if !fileExists(path) {
			return Result{Detail: fmt.Sprintf("Die Datei »%s« fehlt noch.", c.file)}
		}
	}

	if c.command != "" {
		res := runner.Run(ctx, in.WorkDir, "", 0, "sh", "-c", expand(c.command, in))
		switch {
		case res.Err != nil:
			return Result{Detail: "Der Beweis ließ sich nicht prüfen (Befehl nicht ausführbar)."}
		case res.TimedOut:
			return Result{Detail: "Die Prüfung brauchte zu lange."}
		case res.ExitCode != 0:
			return Result{Detail: "Der Beweis hielt der Prüfung nicht stand."}
		case c.expect != "" && !strings.Contains(res.Stdout, c.expect):
			return Result{Detail: "Die Ausgabe war nicht die erwartete."}
		}
	}

	return Result{Passed: true}
}

// --- behavioral: run the student's solution against I/O cases ---

type behavioralCheck struct {
	run   string
	cases []IOCase
}

func (c behavioralCheck) Verify(ctx context.Context, in Input) Result {
	name, args := c.invocation(in)
	for i, cs := range c.cases {
		res := runner.Run(ctx, in.WorkDir, cs.In, 0, name, args...)
		switch {
		case res.Err != nil:
			return Result{Detail: "Deine Lösung ließ sich nicht ausführen."}
		case res.TimedOut:
			return Result{Detail: "Deine Lösung brauchte zu lange."}
		}
		if strings.TrimSpace(res.Stdout) != strings.TrimSpace(cs.Out) {
			return Result{Detail: fmt.Sprintf(
				"Fall %d von %d: erwartet »%s«, bekommen »%s«.",
				i+1, len(c.cases), strings.TrimSpace(cs.Out), strings.TrimSpace(res.Stdout))}
		}
	}
	return Result{Passed: true}
}

// invocation decides how to run the attempt: an explicit `run` template (via the
// shell), or auto-detected from the answer — a .py file with Python, a .go file
// with `go run`, otherwise the answer is treated as a shell one-liner. This is
// what makes the check language-agnostic ("choose your path").
func (c behavioralCheck) invocation(in Input) (string, []string) {
	if c.run != "" {
		return "sh", []string{"-c", expand(c.run, in)}
	}
	a := strings.TrimSpace(in.Answer)
	switch {
	case strings.HasSuffix(a, ".py"):
		py := "python3"
		if !runner.Available(py) {
			py = "python"
		}
		return py, []string{a}
	case strings.HasSuffix(a, ".go"):
		return "go", []string{"run", a}
	default:
		return "sh", []string{"-c", a}
	}
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

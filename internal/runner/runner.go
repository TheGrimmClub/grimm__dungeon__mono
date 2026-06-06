// Package runner executes external commands for puzzle checks: a student's
// Python or Go program, a CLI one-liner, or git. It captures stdout/stderr,
// enforces a timeout, and reports whether a toolchain is even available — so a
// behavioral check can run the student's solution regardless of language
// (decision D003).
package runner

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

// DefaultTimeout bounds any single execution so a runaway loop can't hang grimm.
const DefaultTimeout = 5 * time.Second

// Result is the outcome of running a command.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	TimedOut bool
	Err      error // non-nil for start failures (binary missing, etc.)
}

// Available reports whether a tool is on PATH (e.g. "python3", "go", "task").
func Available(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// Run executes name+args in dir, feeding stdin, with a timeout. dir may be ""
// (current working dir) and stdin may be "".
func Run(ctx context.Context, dir, stdin string, timeout time.Duration, name string, args ...string) Result {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb

	err := cmd.Run()
	res := Result{Stdout: out.String(), Stderr: errb.String()}

	if ctx.Err() == context.DeadlineExceeded {
		res.TimedOut = true
		res.Err = ctx.Err()
		return res
	}
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			res.ExitCode = ee.ExitCode() // ran, but exited non-zero
		} else {
			res.Err = err // failed to start
		}
	}
	return res
}

package session

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/command"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/i18n"
)

// registerLaunchers registers /terminal and /book for /help visibility. Their
// real handling is intercepted in runCommand, because they return an ExecRequest
// (which the write-to-Out command model can't carry). The Run here is never
// called.
func (s *Session) registerLaunchers() {
	noop := func(*command.Context, []string) error { return nil }
	s.reg.Register(&command.Command{Name: "terminal", Summary: i18n.T(i18n.KeyCmdTerminal), Run: noop})
	s.reg.Register(&command.Command{Name: "book", Summary: i18n.T(i18n.KeyCmdBook), Run: noop})
}

// terminalExec builds the request to drop into a shell scoped to the work dir.
func (s *Session) terminalExec() Result {
	if s.workDir == "" {
		return Result{Output: i18n.T(i18n.KeyAlchemistNoDir)}
	}
	name, args := resolveShell()
	return Result{
		Output: i18n.T(i18n.KeyTerminalEnter),
		Exec:   &ExecRequest{Name: name, Args: args, Dir: s.workDir, After: i18n.T(i18n.KeyTerminalReturn)},
	}
}

// bookExec builds the request to edit a file in the work dir.
func (s *Session) bookExec(args []string) Result {
	if s.workDir == "" {
		return Result{Output: i18n.T(i18n.KeyAlchemistNoDir)}
	}
	if len(args) == 0 {
		return Result{Output: i18n.T(i18n.KeyBookUsage)}
	}
	editor := resolveEditor()
	if editor == "" {
		return Result{Output: i18n.T(i18n.KeyNoEditor)}
	}
	// Keep the target inside the work dir (no absolute paths or escaping).
	name := filepath.Base(filepath.Clean(args[0]))
	return Result{
		Output: i18n.T(i18n.KeyBookOpen, name),
		Exec:   &ExecRequest{Name: editor, Args: []string{name}, Dir: s.workDir, After: i18n.T(i18n.KeyBookClosed)},
	}
}

// resolveShell picks the user's shell, falling back sensibly per OS.
func resolveShell() (string, []string) {
	if sh := os.Getenv("SHELL"); sh != "" {
		return sh, nil
	}
	for _, c := range []string{"bash", "zsh", "sh"} {
		if p, err := exec.LookPath(c); err == nil {
			return p, nil
		}
	}
	if runtime.GOOS == "windows" {
		return "cmd", nil
	}
	return "sh", nil
}

// resolveEditor honours $VISUAL/$EDITOR, else the first available friendly editor.
func resolveEditor() string {
	for _, e := range []string{os.Getenv("VISUAL"), os.Getenv("EDITOR")} {
		if e != "" {
			return e
		}
	}
	cands := []string{"micro", "nano", "vim", "vi"}
	if runtime.GOOS == "windows" {
		cands = []string{"notepad"}
	}
	for _, c := range cands {
		if p, err := exec.LookPath(c); err == nil {
			return p
		}
	}
	return ""
}

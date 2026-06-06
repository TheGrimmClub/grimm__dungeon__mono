// Package repl is grimm's input loop. It reads a line, classifies it as either
// a "/command" or free text, and dispatches accordingly. This same loop is the
// seed that later grows into the Claude-Code-shaped shell (req R007).
package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/command"
)

// TextHandler handles free text (anything not starting with "/").
type TextHandler func(ctx *command.Context, text string) error

// REPL ties an input stream, an output stream and a command registry together.
type REPL struct {
	in     *bufio.Scanner
	out    io.Writer
	reg    *command.Registry
	prompt string
	onText TextHandler
}

// New builds a REPL.
func New(in io.Reader, out io.Writer, reg *command.Registry, prompt string, onText TextHandler) *REPL {
	return &REPL{
		in:     bufio.NewScanner(in),
		out:    out,
		reg:    reg,
		prompt: prompt,
		onText: onText,
	}
}

// Run loops until EOF, /quit (command.ErrQuit), or a read error. A handler
// returning command.ErrQuit stops the loop cleanly (no error returned).
func (r *REPL) Run() error {
	ctx := &command.Context{Out: r.out}
	for {
		fmt.Fprint(r.out, r.prompt)
		if !r.in.Scan() {
			// EOF (e.g. Ctrl-D) ends the session like a clean quit.
			fmt.Fprintln(r.out)
			return r.in.Err()
		}
		line := strings.TrimSpace(r.in.Text())
		if line == "" {
			continue
		}
		if err := r.dispatch(ctx, line); err != nil {
			if errors.Is(err, command.ErrQuit) {
				return nil
			}
			return err
		}
	}
}

// dispatch routes a single non-empty line.
func (r *REPL) dispatch(ctx *command.Context, line string) error {
	if strings.HasPrefix(line, "/") {
		name, args := parseCommand(line)
		cmd, ok := r.reg.Get(name)
		if !ok {
			fmt.Fprintln(r.out, unknownCommand(name))
			return nil
		}
		return cmd.Run(ctx, args)
	}
	if r.onText != nil {
		return r.onText(ctx, line)
	}
	return nil
}

// parseCommand splits "/name arg1 arg2" into ("name", ["arg1","arg2"]).
func parseCommand(line string) (name string, args []string) {
	fields := strings.Fields(strings.TrimPrefix(line, "/"))
	if len(fields) == 0 {
		return "", nil
	}
	return fields[0], fields[1:]
}

// unknownCommand is a hook the app overrides via the registry's miss message;
// kept here as a tiny default so the package is usable on its own.
var unknownCommand = func(name string) string {
	return fmt.Sprintf("Unknown command: /%s", name)
}

// SetUnknownCommandMessage lets the app provide a localized "unknown command"
// message (so the German narrative lives in the i18n catalog, not here).
func SetUnknownCommandMessage(f func(name string) string) { unknownCommand = f }

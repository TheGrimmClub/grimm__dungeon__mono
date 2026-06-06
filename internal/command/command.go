// Package command holds the slash-command registry that powers grimm's
// "/command" surface. As the game evolves from a Zork-like REPL toward a
// Claude-Code-shaped shell (decision D006/req R007), new commands are simply
// registered here.
package command

import (
	"errors"
	"io"
	"sort"
)

// ErrQuit is returned by a handler to ask the REPL to stop the loop.
var ErrQuit = errors.New("quit")

// Context carries everything a command handler may need. It will grow to hold
// game state, the player, the runner, etc. as later phases land.
type Context struct {
	Out io.Writer
}

// Handler runs a command with the words that followed its name.
type Handler func(ctx *Context, args []string) error

// Command is a single slash command (the leading "/" is not part of Name).
type Command struct {
	Name    string  // e.g. "help" for the "/help" command
	Summary string  // short, German one-liner shown in /help
	Hidden  bool    // hidden commands (easter eggs) are kept out of /help
	Run     Handler // the behaviour
}

// Registry stores commands by name and remembers registration order.
type Registry struct {
	commands map[string]*Command
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{commands: make(map[string]*Command)}
}

// Register adds (or replaces) a command. A nil command or empty name is ignored.
func (r *Registry) Register(c *Command) {
	if c == nil || c.Name == "" {
		return
	}
	r.commands[c.Name] = c
}

// Get looks up a command by name (without the leading "/").
func (r *Registry) Get(name string) (*Command, bool) {
	c, ok := r.commands[name]
	return c, ok
}

// Visible returns the non-hidden commands sorted by name, for /help.
func (r *Registry) Visible() []*Command {
	out := make([]*Command, 0, len(r.commands))
	for _, c := range r.commands {
		if !c.Hidden {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

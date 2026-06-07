// Package session is the pure, UI-independent core of a grimm play session: it
// turns one submitted line into output text (and a quit signal). The Bubble Tea
// layer (package tui) is a thin view on top, so all dispatch logic stays
// testable without a terminal.
package session

import (
	"bytes"
	"errors"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/command"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/engine"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/i18n"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/voice"
)

// Result is the outcome of submitting one line.
type Result struct {
	Output string
	Quit   bool
}

// Session owns the game, the slash-command registry and the save location.
type Session struct {
	game     *engine.Game
	reg      *command.Registry
	savePath string // "" disables /save

	player  voice.Player // text-to-speech backend (Noop by default)
	voiceOn bool         // whether narration is currently enabled
}

// New builds a session around a game. savePath may be "" to disable saving.
func New(game *engine.Game, savePath string) *Session {
	s := &Session{
		game:     game,
		reg:      command.NewRegistry(),
		savePath: savePath,
		player:   voice.Noop(),
	}
	s.registerBuiltins()
	return s
}

// SetVoice installs a text-to-speech backend (the app injects the OS voice).
func (s *Session) SetVoice(p voice.Player) {
	if p != nil {
		s.player = p
	}
}

// Game exposes the underlying engine (the TUI reads Title/Lit from it).
func (s *Session) Game() *engine.Game { return s.game }

// Submit processes one line of input and returns what to show.
func (s *Session) Submit(line string) Result {
	line = strings.TrimSpace(line)
	switch {
	case line == "":
		return Result{Output: i18n.T(i18n.KeyEmptyInfo)}
	case strings.HasPrefix(line, "/"):
		return s.runCommand(line)
	case normalize(line) == "import antigravity":
		return Result{Output: i18n.T(i18n.KeyEasterEgg) + antigravityArt}
	default:
		out := s.game.Do(line)
		s.narrate(out) // read the room/answer aloud when voice is on
		return Result{Output: out}
	}
}

// narrate speaks engine output aloud when narration is enabled and available.
func (s *Session) narrate(text string) {
	if s.voiceOn && s.player.Available() {
		s.player.Speak(text)
	}
}

// runCommand dispatches a "/name args" line through the registry, capturing the
// command's output and translating command.ErrQuit into Result.Quit.
func (s *Session) runCommand(line string) Result {
	name, args := parseCommand(line)
	cmd, ok := s.reg.Get(name)
	if !ok {
		return Result{Output: i18n.T(i18n.KeyUnknownCommand, "/"+name, "/help")}
	}
	var buf bytes.Buffer
	err := cmd.Run(&command.Context{Out: &buf}, args)
	return Result{
		Output: strings.TrimRight(buf.String(), "\n"),
		Quit:   errors.Is(err, command.ErrQuit),
	}
}

// parseCommand splits "/name a b" into ("name", ["a","b"]).
func parseCommand(line string) (name string, args []string) {
	fields := strings.Fields(strings.TrimPrefix(line, "/"))
	if len(fields) == 0 {
		return "", nil
	}
	return fields[0], fields[1:]
}

func normalize(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

const antigravityArt = `
        .   *        o   .       *
   *        .---.        .    o
       .   /     \   *        .   *
     o     | o o |     .   *
        *   \  ^  /  .        o
            '---'      ~ schwerelos ~
`

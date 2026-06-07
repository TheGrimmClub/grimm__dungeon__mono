// Package session is the pure, UI-independent core of a grimm play session: it
// turns one submitted line into output text (and a quit signal). The Bubble Tea
// layer (package tui) is a thin view on top, so all dispatch logic stays
// testable without a terminal.
package session

import (
	"bytes"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/alchemist"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/command"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/engine"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/i18n"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/voice"
)

// Result is the outcome of submitting one line.
type Result struct {
	Output string
	Quit   bool
	// Exec, when set, asks the UI to suspend and run an external interactive
	// process (a shell or editor), then resume. Kept as plain data so the
	// session stays UI-independent and testable.
	Exec *ExecRequest
}

// ExecRequest describes an external process for the UI to run (via Bubble Tea's
// tea.ExecProcess), e.g. /terminal or /book.
type ExecRequest struct {
	Name  string   // binary to run
	Args  []string // its arguments
	Dir   string   // working directory
	After string   // message to show once it returns
}

// Session owns the game, the slash-command registry and the save location.
type Session struct {
	game     *engine.Game
	reg      *command.Registry
	savePath string // "" disables /save

	player  voice.Player // text-to-speech backend (Noop by default)
	voiceOn bool         // whether narration is currently enabled

	alch    *alchemist.Alchemist // git-as-potions in the student work dir
	workDir string               // student's working directory ("" => unset)
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

// SetWorkDir points /alchemist at the student's working directory, creating it
// if needed so the first potion can be brewed there.
func (s *Session) SetWorkDir(dir string) {
	if dir == "" {
		return
	}
	_ = os.MkdirAll(dir, 0o755)
	s.workDir = dir
	s.alch = alchemist.New(dir)
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
// It reads the prose, not the UI chrome (exits, item numbers, lock footnotes,
// the "(north)" direction hints) — those sound noisy spoken.
func (s *Session) narrate(text string) {
	if s.voiceOn && s.player.Available() {
		s.player.Speak(narratable(text))
	}
}

var (
	reItemNumber = regexp.MustCompile(`^\s*\[\S+\]`) // "[1] …", "[0] …"
	reDirHint    = regexp.MustCompile(`\s*\((?:north|south|east|west|up|down)\)`)
)

// narratable reduces engine output to the speakable prose.
func narratable(text string) string {
	keep := make([]string, 0)
	for _, line := range strings.Split(text, "\n") {
		t := strings.TrimSpace(line)
		switch {
		case t == "",
			reItemNumber.MatchString(t),
			strings.HasPrefix(t, "*"),
			strings.HasPrefix(t, "Ausgänge:"), strings.HasPrefix(t, "Exits:"),
			strings.HasPrefix(t, "Hier liegt:"), strings.HasPrefix(t, "Here lies:"),
			strings.HasPrefix(t, "Du trägst bei dir:"), strings.HasPrefix(t, "You are carrying:"):
			continue
		}
		keep = append(keep, t)
	}
	return reDirHint.ReplaceAllString(strings.Join(keep, " "), "")
}

// runCommand dispatches a "/name args" line through the registry, capturing the
// command's output and translating command.ErrQuit into Result.Quit.
func (s *Session) runCommand(line string) Result {
	name, args := parseCommand(line)

	// /terminal and /book return an ExecRequest, so they bypass the
	// write-to-Out command model.
	switch name {
	case "terminal":
		return s.terminalExec()
	case "book":
		return s.bookExec(args)
	}

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

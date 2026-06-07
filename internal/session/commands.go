package session

import (
	"fmt"
	"io"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/alchemist"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/command"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/engine"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/state"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/i18n"
)

// gameVerb is one row in the "game verbs" section of /help.
type gameVerb struct{ name, descKey string }

// gameVerbs documents the engine vocabulary so /help can teach it. (The engine
// accepts more aliases; these are the canonical ones we advertise.)
var gameVerbs = []gameVerb{
	{"look", i18n.KeyVLook},
	{"go <dir>", i18n.KeyVGo},
	{"take <n|name>", i18n.KeyVTake},
	{"inspect <n|name>", i18n.KeyVInspect},
	{"wear <n|name>", i18n.KeyVWear},
	{"solve <solution>", i18n.KeyVSolve},
	{"inventory", i18n.KeyVInventory},
}

// registerBuiltins installs the slash commands.
func (s *Session) registerBuiltins() {
	s.reg.Register(&command.Command{
		Name:    "help",
		Summary: i18n.T(i18n.KeyCmdHelp),
		Run: func(ctx *command.Context, _ []string) error {
			s.printHelp(ctx.Out)
			return nil
		},
	})

	s.reg.Register(&command.Command{
		Name:    "save",
		Summary: i18n.T(i18n.KeyCmdSave),
		Run: func(ctx *command.Context, _ []string) error {
			if s.savePath == "" {
				fmt.Fprintln(ctx.Out, i18n.T(i18n.KeySaveDisabled))
				return nil
			}
			if err := state.Save(s.savePath, s.game.Snapshot()); err != nil {
				fmt.Fprintln(ctx.Out, i18n.T(i18n.KeySaveFailed))
				return nil
			}
			fmt.Fprintln(ctx.Out, i18n.T(i18n.KeySaved))
			return nil
		},
	})

	s.reg.Register(&command.Command{
		Name:    "class",
		Summary: i18n.T(i18n.KeyCmdClass),
		Run: func(ctx *command.Context, args []string) error {
			s.chooseClass(ctx.Out, args)
			return nil
		},
	})

	s.reg.Register(&command.Command{
		Name:    "alchemist",
		Summary: i18n.T(i18n.KeyCmdAlchemist),
		Run: func(ctx *command.Context, args []string) error {
			s.runAlchemist(ctx.Out, args)
			return nil
		},
	})

	s.reg.Register(&command.Command{
		Name:    "voice",
		Summary: i18n.T(i18n.KeyCmdVoice),
		Run: func(ctx *command.Context, args []string) error {
			s.toggleVoice(ctx.Out, args)
			return nil
		},
	})

	s.reg.Register(&command.Command{
		Name:    "quit",
		Summary: i18n.T(i18n.KeyCmdQuit),
		Run: func(ctx *command.Context, _ []string) error {
			fmt.Fprintln(ctx.Out, i18n.T(i18n.KeyGoodbye))
			return command.ErrQuit
		},
	})

	// Hidden wow-effect: the antigravity easter egg (req R008).
	s.reg.Register(&command.Command{
		Name:   "antigravity",
		Hidden: true,
		Run: func(ctx *command.Context, _ []string) error {
			fmt.Fprint(ctx.Out, i18n.T(i18n.KeyEasterEgg), antigravityArt)
			return nil
		},
	})
}

// runAlchemist drives the in-game potion/git tool in the student work dir.
func (s *Session) runAlchemist(out io.Writer, args []string) {
	if s.alch == nil {
		fmt.Fprintln(out, i18n.T(i18n.KeyAlchemistNoDir))
		return
	}
	msg, err := alchemist.Dispatch(s.alch, args)
	if err != nil {
		fmt.Fprintln(out, err.Error())
		return
	}
	fmt.Fprintln(out, msg)
}

// toggleVoice turns narration on/off ("/voice", "/voice on", "/voice off").
func (s *Session) toggleVoice(out io.Writer, args []string) {
	want := !s.voiceOn // bare /voice toggles
	if len(args) > 0 {
		switch strings.ToLower(args[0]) {
		case "on", "an", "ein":
			want = true
		case "off", "aus":
			want = false
		}
	}
	if want && !s.player.Available() {
		fmt.Fprintln(out, i18n.T(i18n.KeyVoiceUnavailable))
		return
	}
	s.voiceOn = want
	if s.voiceOn {
		fmt.Fprintln(out, i18n.T(i18n.KeyVoiceOn))
	} else {
		s.player.Stop()
		fmt.Fprintln(out, i18n.T(i18n.KeyVoiceOff))
	}
}

// chooseClass lists the paths, or sets one if the player named it.
func (s *Session) chooseClass(out io.Writer, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(out, i18n.T(i18n.KeyClassHeader))
		for _, c := range engine.Classes() {
			fmt.Fprintf(out, "  %-10s — %s\n", c.ID, c.Blurb)
		}
		fmt.Fprintln(out)
		fmt.Fprintln(out, i18n.T(i18n.KeyClassChoose))
		return
	}
	c, ok := s.game.ChooseClass(args[0])
	if !ok {
		fmt.Fprintln(out, i18n.T(i18n.KeyClassUnknown, args[0]))
		return
	}
	fmt.Fprintln(out, i18n.T(i18n.KeyClassChosen, c.Title))
}

// printHelp renders both the slash commands and the game verbs, aligned.
func (s *Session) printHelp(out io.Writer) {
	fmt.Fprintln(out, i18n.T(i18n.KeyHelpCmdHeader))
	cmds := s.reg.Visible()
	width := 0
	for _, c := range cmds {
		if n := len(c.Name) + 1; n > width {
			width = n
		}
	}
	for _, c := range cmds {
		fmt.Fprintf(out, "  %-*s  %s\n", width, "/"+c.Name, c.Summary)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, i18n.T(i18n.KeyHelpVerbHeader))
	vw := 0
	for _, v := range gameVerbs {
		if len(v.name) > vw {
			vw = len(v.name)
		}
	}
	for _, v := range gameVerbs {
		fmt.Fprintf(out, "  %-*s  %s\n", vw, v.name, i18n.T(v.descKey))
	}
}

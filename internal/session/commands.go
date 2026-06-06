package session

import (
	"fmt"
	"io"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/command"
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

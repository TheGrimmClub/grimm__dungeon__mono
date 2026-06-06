// Package app wires grimm together: it loads the world, builds the command
// registry and the verb engine, prints the banner, and runs the REPL.
// cmd/grimm/main.go stays a thin entry point that just calls Run.
package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/assets"
	"github.com/TheGrimmClub/grimm__dungeon__mono/content"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/command"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/engine"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/state"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/i18n"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/repl"
)

// Options configures a grimm session.
type Options struct {
	// SavePath is where progress is loaded from and saved to. An empty string
	// disables persistence (used by tests for hermeticity).
	SavePath string
}

// Run starts grimm with the default save location (~/.grimm/save.yaml).
func Run(in io.Reader, out io.Writer) error {
	path, _ := state.DefaultPath() // empty on error -> persistence disabled
	return RunWith(in, out, Options{SavePath: path})
}

// RunWith starts grimm with explicit options and blocks until the session ends
// (EOF or /quit). It returns a non-nil error only on real failures; a clean quit
// returns nil.
func RunWith(in io.Reader, out io.Writer, opts Options) error {
	w, err := world.Load(content.FS, content.WorldGlob)
	if err != nil {
		return fmt.Errorf("app: loading world: %w", err)
	}
	game := engine.New(w)

	continued := loadProgress(game, opts.SavePath)

	reg := command.NewRegistry()
	registerBuiltins(reg, game, opts.SavePath)
	repl.SetUnknownCommandMessage(func(name string) string {
		return i18n.T(i18n.KeyUnknownCommand, "/"+name, "/help")
	})

	printIntro(out, game, continued)

	r := repl.New(in, out, reg, i18n.T(i18n.KeyPrompt), textHandler(game))
	return r.Run()
}

// loadProgress restores a save if one exists, reporting whether it did.
func loadProgress(game *engine.Game, path string) bool {
	if path == "" || !state.Exists(path) {
		return false
	}
	snap, err := state.Load(path)
	if err != nil {
		return false // a corrupt save shouldn't block a fresh start
	}
	game.Restore(snap)
	return true
}

// printIntro shows the banner, welcome, optional "continued" note and the
// current room, plus a one-time hint about how to talk to the dungeon.
func printIntro(out io.Writer, game *engine.Game, continued bool) {
	fmt.Fprint(out, assets.Banner)
	fmt.Fprintln(out, "  "+i18n.T(i18n.KeyBannerSubtitle))
	fmt.Fprintln(out)
	fmt.Fprintln(out, i18n.T(i18n.KeyWelcome, "/help", "/quit"))
	fmt.Fprintln(out)
	if continued {
		fmt.Fprintln(out, i18n.T(i18n.KeyContinued))
		fmt.Fprintln(out)
	}
	fmt.Fprintln(out, game.Intro())
	fmt.Fprintln(out)
	fmt.Fprintln(out, i18n.T(i18n.KeyVerbHint))
	fmt.Fprintln(out)
}

// registerBuiltins installs the Phase 0/1 command set.
func registerBuiltins(reg *command.Registry, game *engine.Game, savePath string) {
	reg.Register(&command.Command{
		Name:    "help",
		Summary: i18n.T(i18n.KeyCmdHelp),
		Run: func(ctx *command.Context, _ []string) error {
			printHelp(ctx, reg)
			return nil
		},
	})

	reg.Register(&command.Command{
		Name:    "save",
		Summary: i18n.T(i18n.KeyCmdSave),
		Run: func(ctx *command.Context, _ []string) error {
			if savePath == "" {
				fmt.Fprintln(ctx.Out, i18n.T(i18n.KeySaveDisabled))
				return nil
			}
			if err := state.Save(savePath, game.Snapshot()); err != nil {
				fmt.Fprintln(ctx.Out, i18n.T(i18n.KeySaveFailed))
				return nil
			}
			fmt.Fprintln(ctx.Out, i18n.T(i18n.KeySaved))
			return nil
		},
	})

	reg.Register(&command.Command{
		Name:    "quit",
		Summary: i18n.T(i18n.KeyCmdQuit),
		Run: func(ctx *command.Context, _ []string) error {
			fmt.Fprintln(ctx.Out, i18n.T(i18n.KeyGoodbye))
			return command.ErrQuit
		},
	})

	// Hidden wow-effect: the antigravity easter egg (req R008).
	reg.Register(&command.Command{
		Name:   "antigravity",
		Hidden: true,
		Run: func(ctx *command.Context, _ []string) error {
			showAntigravity(ctx.Out)
			return nil
		},
	})
}

// printHelp renders the scroll of (visible) commands, aligned.
func printHelp(ctx *command.Context, reg *command.Registry) {
	fmt.Fprintln(ctx.Out, i18n.T(i18n.KeyHelpHeader))
	cmds := reg.Visible()
	width := 0
	for _, c := range cmds {
		if n := len(c.Name) + 1; n > width { // +1 for the leading "/"
			width = n
		}
	}
	for _, c := range cmds {
		fmt.Fprintf(ctx.Out, "  %-*s  %s\n", width, "/"+c.Name, c.Summary)
	}
	fmt.Fprintln(ctx.Out)
	fmt.Fprintln(ctx.Out, i18n.T(i18n.KeyHelpHint))
}

// textHandler routes free text: the easter egg first, otherwise the verb engine.
func textHandler(game *engine.Game) repl.TextHandler {
	return func(ctx *command.Context, text string) error {
		if normalize(text) == "import antigravity" {
			showAntigravity(ctx.Out)
			return nil
		}
		fmt.Fprintln(ctx.Out, game.Do(text))
		return nil
	}
}

// showAntigravity prints the easter-egg payload.
func showAntigravity(out io.Writer) {
	fmt.Fprint(out, i18n.T(i18n.KeyEasterEgg))
	fmt.Fprint(out, antigravityArt)
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

// Package app wires grimm together: it builds the command registry, prints the
// banner, and runs the REPL. cmd/grimm/main.go stays a thin entry point that
// just calls Run.
package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/assets"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/command"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/i18n"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/repl"
)

// Run starts grimm against the given input/output streams and blocks until the
// session ends (EOF or /quit). It returns a non-nil error only on real I/O
// failures — a clean quit returns nil.
func Run(in io.Reader, out io.Writer) error {
	reg := command.NewRegistry()
	registerBuiltins(reg)

	// Keep the "unknown command" wording in the German i18n catalog.
	repl.SetUnknownCommandMessage(func(name string) string {
		return i18n.T(i18n.KeyUnknownCommand, "/"+name, "/help")
	})

	printIntro(out)

	r := repl.New(in, out, reg, i18n.T(i18n.KeyPrompt), onText)
	return r.Run()
}

// printIntro shows the banner, subtitle and welcome text.
func printIntro(out io.Writer) {
	fmt.Fprint(out, assets.Banner)
	fmt.Fprintln(out, "  "+i18n.T(i18n.KeyBannerSubtitle))
	fmt.Fprintln(out)
	fmt.Fprintln(out, i18n.T(i18n.KeyWelcome, "/help", "/quit"))
	fmt.Fprintln(out)
}

// registerBuiltins installs the Phase 0 command set.
func registerBuiltins(reg *command.Registry) {
	reg.Register(&command.Command{
		Name:    "help",
		Summary: i18n.T(i18n.KeyCmdHelp),
		Run: func(ctx *command.Context, _ []string) error {
			printHelp(ctx, reg)
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

	// Hidden wow-effect: the antigravity easter egg (req R008). Triggered by the
	// secret "/antigravity" command or by whispering "import antigravity".
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

// onText handles free text (non-command input). For now it recognises the
// "import antigravity" easter egg and otherwise gives a gentle placeholder —
// the Zork verb engine arrives in Phase 1.
func onText(ctx *command.Context, text string) error {
	if normalize(text) == "import antigravity" {
		showAntigravity(ctx.Out)
		return nil
	}
	fmt.Fprintln(ctx.Out, i18n.T(i18n.KeyUnknownInput, text))
	return nil
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

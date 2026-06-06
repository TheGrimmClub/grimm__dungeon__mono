// Package alchemist is the shared library behind the "git as potion-brewing"
// tool. It is exposed two ways (decision D006): as the standalone cmd/alchemist
// binary, and later as grimm's in-game "/alchemist" command. Real git wiring
// lands in Phase 3; this Phase 0 stub establishes the grimoire (verb -> git
// mapping) that both surfaces share.
package alchemist

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Verb is one potion action and the git reality it teaches.
type Verb struct {
	Name    string // the potion verb the student types
	Git     string // the real git operation it maps to
	Summary string // short German flavour text
}

// Grimoire is the canonical potion -> git mapping (req R006).
var Grimoire = []Verb{
	{"add", "git add", "lege eine Zutat in den Kessel (stage)"},
	{"brew", "git add + git commit", "braue den Trank — halte den Stand fest"},
	{"bottle", "git push", "fülle den Trank ab und schicke ihn fort"},
	{"discard", "git restore", "schütte die offenen Änderungen weg"},
	{"clean", "git clean", "fege unversionierte Reste vom Tisch"},
	{"look", "git status", "betrachte den Kessel — was brodelt gerade?"},
}

// Run is the standalone tool's entry. For Phase 0 it prints the grimoire; the
// brewing magic (actual git calls) arrives in Phase 3.
func Run(out io.Writer, _ []string) error {
	fmt.Fprintln(out, "Das Grimoire des Alchemisten — Zaubertränke statt Git-Befehle:")
	fmt.Fprintln(out)
	tw := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "  TRANK\tGIT\tWIRKUNG")
	for _, v := range Grimoire {
		fmt.Fprintf(tw, "  %s\t%s\t%s\n", v.Name, v.Git, v.Summary)
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "(Das Brauen selbst erwacht in Phase 3.)")
	return nil
}

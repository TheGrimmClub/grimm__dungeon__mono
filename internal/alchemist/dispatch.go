package alchemist

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

// Dispatch runs a potion verb against an Alchemist and returns the message to
// show. It is shared by the standalone binary and the in-game /alchemist command
// so both behave identically. An empty verb (or "help") prints the grimoire.
func Dispatch(a *Alchemist, args []string) (string, error) {
	if !Available() {
		return "", fmt.Errorf("git ist nicht installiert — ohne git kann der Alchemist nicht brauen")
	}
	if len(args) == 0 || args[0] == "help" {
		return grimoire(), nil
	}

	verb, rest := strings.ToLower(args[0]), args[1:]
	switch verb {
	case "init":
		return a.Init()
	case "add":
		return a.Add(rest)
	case "brew":
		return a.Brew(strings.Join(rest, " "))
	case "bottle":
		return a.Bottle()
	case "discard":
		return a.Discard()
	case "clean":
		return a.Clean()
	case "look", "status":
		return a.Look()
	default:
		return "", fmt.Errorf("»%s« ist kein Trank, den ich kenne. Tippe »help« für das Grimoire", verb)
	}
}

func grimoire() string {
	var b strings.Builder
	b.WriteString("Das Grimoire des Alchemisten — Zaubertränke statt Git-Befehle:\n\n")
	tw := tabwriter.NewWriter(&b, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "  TRANK\tGIT\tWIRKUNG")
	for _, v := range Grimoire {
		fmt.Fprintf(tw, "  %s\t%s\t%s\n", v.Name, v.Git, v.Summary)
	}
	_ = tw.Flush()
	return strings.TrimRight(b.String(), "\n")
}

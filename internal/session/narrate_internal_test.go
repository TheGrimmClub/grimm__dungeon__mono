package session

import (
	"strings"
	"testing"
)

func TestNarratableStripsChrome(t *testing.T) {
	in := strings.Join([]string{
		"Das verwunschene Tor",
		"Nach Norden (north) öffnet sich ein dunkler Gang.",
		"Hier liegt:",
		"  [1] Helm mit Stirnlampe",
		"Ausgänge: north, east, south, west*",
		"* ein Rätsel versperrt diesen Weg — versuche dort zu »go« und dann »solve«.",
	}, "\n")

	got := narratable(in)

	// Prose is kept, direction hints removed.
	if !strings.Contains(got, "Das verwunschene Tor") {
		t.Errorf("title dropped: %q", got)
	}
	if !strings.Contains(got, "Nach Norden öffnet sich ein dunkler Gang") {
		t.Errorf("prose missing or hint not stripped: %q", got)
	}
	// Chrome is gone.
	for _, bad := range []string{"(north)", "[1]", "Ausgänge:", "Hier liegt:", "* ein Rätsel"} {
		if strings.Contains(got, bad) {
			t.Errorf("chrome %q leaked into narration: %q", bad, got)
		}
	}
}

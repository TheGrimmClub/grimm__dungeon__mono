package app

import (
	"strings"
	"testing"
)

// runScript feeds the given input lines through a fresh grimm session and
// returns everything it printed.
func runScript(t *testing.T, input string) string {
	t.Helper()
	var out strings.Builder
	if err := Run(strings.NewReader(input), &out); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	return out.String()
}

func TestIntroAndHelp(t *testing.T) {
	out := runScript(t, "/help\n/quit\n")

	for _, want := range []string{
		`/ ___|`,                    // a stable fragment of the ASCII banner
		"Schriftrolle der Befehle:", // help header (German)
		"/help",                     // listed command
		"/quit",
		"Fackeln verlöschen", // goodbye on /quit
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n---\n%s", want, out)
		}
	}
	// Hidden commands must not leak into /help.
	if strings.Contains(out, "/antigravity") {
		t.Errorf("hidden command /antigravity leaked into help output")
	}
}

func TestAntigravityEasterEggViaText(t *testing.T) {
	out := runScript(t, "import antigravity\n/quit\n")
	if !strings.Contains(out, "schwerelos") {
		t.Errorf("free-text easter egg did not fire\n---\n%s", out)
	}
}

func TestAntigravityEasterEggViaCommand(t *testing.T) {
	out := runScript(t, "/antigravity\n/quit\n")
	if !strings.Contains(out, "Schwerkraft hat heute frei") {
		t.Errorf("/antigravity command did not fire\n---\n%s", out)
	}
}

func TestUnknownCommand(t *testing.T) {
	out := runScript(t, "/fliegen\n/quit\n")
	if !strings.Contains(out, "/fliegen") || !strings.Contains(out, "/help") {
		t.Errorf("unknown-command message missing the offending name or the /help hint\n---\n%s", out)
	}
}

func TestEOFEndsSessionCleanly(t *testing.T) {
	// No /quit — input simply ends. Run must return without error.
	_ = runScript(t, "import antigravity\n")
}

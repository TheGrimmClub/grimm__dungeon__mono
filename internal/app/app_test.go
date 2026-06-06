package app

import (
	"path/filepath"
	"strings"
	"testing"
)

// runScript feeds the given input lines through a fresh grimm session with
// persistence disabled (hermetic) and returns everything it printed.
func runScript(t *testing.T, input string) string {
	t.Helper()
	return runScriptWith(t, input, Options{})
}

func runScriptWith(t *testing.T, input string, opts Options) string {
	t.Helper()
	var out strings.Builder
	if err := RunWith(strings.NewReader(input), &out, opts); err != nil {
		t.Fatalf("RunWith returned error: %v", err)
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

func TestIntroShowsStartRoomAndVerbHint(t *testing.T) {
	out := runScript(t, "/quit\n")
	if !strings.Contains(out, "Das verwunschene Tor") {
		t.Errorf("intro did not describe the start room:\n%s", out)
	}
	if !strings.Contains(out, "Sprich mit dem Verlies") {
		t.Errorf("intro did not show the verb hint:\n%s", out)
	}
}

func TestFreeTextRoutesToEngine(t *testing.T) {
	out := runScript(t, "gehe norden\n/quit\n")
	if !strings.Contains(out, "Halle der schlafenden Maschinen") {
		t.Errorf("free-text verb did not reach the engine:\n%s", out)
	}
}

func TestSaveThenContinue(t *testing.T) {
	opts := Options{SavePath: filepath.Join(t.TempDir(), "save.yaml")}

	// First session: walk north, pick up the lamp, save.
	first := runScriptWith(t, "gehe norden\n/save\n/quit\n", opts)
	if !strings.Contains(first, "versiegelt") {
		t.Fatalf("save did not confirm:\n%s", first)
	}

	// Second session with the same save path: should resume in the hall.
	second := runScriptWith(t, "schau\n/quit\n", opts)
	if !strings.Contains(second, "Du nimmst deinen Weg wieder auf") {
		t.Errorf("second session did not report continuing:\n%s", second)
	}
	if !strings.Contains(second, "Halle der schlafenden Maschinen") {
		t.Errorf("second session did not resume in the saved room:\n%s", second)
	}
}

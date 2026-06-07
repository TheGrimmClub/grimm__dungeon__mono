package voice

import (
	"strings"
	"testing"
)

func TestCleanStripsAnsiAndCollapses(t *testing.T) {
	in := "\x1b[1;36mHallo\x1b[0m   Welt\n\nzweite Zeile"
	if got, want := clean(in), "Hallo Welt zweite Zeile"; got != want {
		t.Errorf("clean = %q, want %q", got, want)
	}
}

func TestCleanCapsLength(t *testing.T) {
	in := strings.Repeat("a", maxSpeechRunes+50)
	if got := clean(in); len([]rune(got)) != maxSpeechRunes {
		t.Errorf("clean length = %d, want %d", len([]rune(got)), maxSpeechRunes)
	}
}

func TestNoopIsUnavailable(t *testing.T) {
	if Noop().Available() {
		t.Error("Noop().Available() = true, want false")
	}
}

func TestParseMacVoicesPrefersBaseGerman(t *testing.T) {
	sample := `Albert              en_US    # Hello! My name is Albert.
Alice               it_IT    # Ciao! Mi chiamo Alice.
Eddy (German (Germany)) de_DE    # Hallo! Ich heiße Eddy.
Anna                de_DE    # Hallo! Ich heiße Anna.
Samantha            en_US    # Hi, I'm Samantha.`
	if got := parseMacVoices(sample); got != "Anna" {
		t.Errorf("parseMacVoices = %q, want Anna (clean base German voice)", got)
	}
}

func TestParseMacVoicesFallsBackToAnyGerman(t *testing.T) {
	sample := "Albert              en_US    # x\nEddy (German (Germany)) de_DE    # y"
	if got := parseMacVoices(sample); got != "Eddy (German (Germany))" {
		t.Errorf("parseMacVoices = %q, want the only German voice", got)
	}
}

func TestParseMacVoicesNoGerman(t *testing.T) {
	if got := parseMacVoices("Albert  en_US  # x\nAlice  it_IT  # y"); got != "" {
		t.Errorf("parseMacVoices = %q, want empty (no German voice)", got)
	}
}

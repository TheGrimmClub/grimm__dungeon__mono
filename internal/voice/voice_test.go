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

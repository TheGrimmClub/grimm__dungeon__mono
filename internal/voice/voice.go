// Package voice reads the game's German narration aloud (text-to-speech). It
// shells out to the operating system's built-in speech (macOS `say`, Linux
// espeak/spd-say, Windows SAPI via PowerShell) behind a Player interface, with
// a no-op fallback so it is safe on machines without a voice and in tests.
//
// (Speech-to-text — voice *input* — is a deliberate later step; see decisions.)
package voice

import (
	"context"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// maxSpeechRunes caps how much we narrate at once, so a huge dump can't make the
// voice ramble for minutes.
const maxSpeechRunes = 600

// Player narrates text aloud. Implementations must be safe for concurrent use.
type Player interface {
	Speak(text string) // narrate; cancels any in-progress narration
	Stop()             // stop any in-progress narration
	Available() bool   // whether a real voice backend exists
}

// Noop is a Player that says nothing (the default and the test/CI backend).
func Noop() Player { return noop{} }

type noop struct{}

func (noop) Speak(string)    {}
func (noop) Stop()           {}
func (noop) Available() bool { return false }

// New returns the best available OS voice backend, or Noop if none is found.
func New() Player {
	switch runtime.GOOS {
	case "darwin":
		if has("say") {
			return &osPlayer{name: "say", build: sayArgs}
		}
	case "linux":
		if has("spd-say") {
			return &osPlayer{name: "spd-say", build: func(t string) []string { return []string{"-l", "de", "-w", t} }}
		}
		for _, c := range []string{"espeak-ng", "espeak"} {
			if has(c) {
				name := c
				return &osPlayer{name: name, build: func(t string) []string { return []string{"-v", "de", t} }}
			}
		}
	case "windows":
		if has("powershell") {
			return &osPlayer{name: "powershell", build: powershellArgs}
		}
	}
	return Noop()
}

func has(tool string) bool { _, err := exec.LookPath(tool); return err == nil }

func sayArgs(text string) []string {
	if v := os.Getenv("GRIMM_VOICE"); v != "" { // e.g. GRIMM_VOICE=Anna for German
		return []string{"-v", v, text}
	}
	return []string{text}
}

func powershellArgs(text string) []string {
	esc := strings.ReplaceAll(text, "'", "''")
	script := "Add-Type -AssemblyName System.Speech; " +
		"$s = New-Object System.Speech.Synthesis.SpeechSynthesizer; $s.Speak('" + esc + "')"
	return []string{"-NoProfile", "-Command", script}
}

// osPlayer narrates by running an OS speech command. Each Speak cancels the
// previous one so narration tracks the latest room rather than piling up.
type osPlayer struct {
	name  string
	build func(text string) []string

	mu     sync.Mutex
	cancel context.CancelFunc
}

func (p *osPlayer) Available() bool { return true }

func (p *osPlayer) Speak(text string) {
	text = clean(text)
	if text == "" {
		return
	}
	p.mu.Lock()
	if p.cancel != nil {
		p.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.mu.Unlock()

	cmd := exec.CommandContext(ctx, p.name, p.build(text)...)
	go func() { _ = cmd.Run() }() // fire-and-forget; cancellation stops it
}

func (p *osPlayer) Stop() {
	p.mu.Lock()
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}
	p.mu.Unlock()
}

var ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// clean strips ANSI styling and collapses whitespace, then caps the length, so
// the spoken text is the words a player would read — not escape codes.
func clean(s string) string {
	s = ansiRE.ReplaceAllString(s, "")
	s = strings.Join(strings.Fields(s), " ")
	if r := []rune(s); len(r) > maxSpeechRunes {
		s = string(r[:maxSpeechRunes])
	}
	return s
}

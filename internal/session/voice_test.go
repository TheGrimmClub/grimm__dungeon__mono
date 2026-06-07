package session_test

import (
	"strings"
	"testing"
)

// fakeVoice records what would be spoken.
type fakeVoice struct {
	spoken []string
	avail  bool
}

func (f *fakeVoice) Speak(text string) { f.spoken = append(f.spoken, text) }
func (f *fakeVoice) Stop()             {}
func (f *fakeVoice) Available() bool   { return f.avail }

func TestVoiceNarratesOnlyWhenOnAndForGameText(t *testing.T) {
	s := newSession(t, "")
	fv := &fakeVoice{avail: true}
	s.SetVoice(fv)

	// Default off: nothing is narrated.
	s.Submit("look")
	if len(fv.spoken) != 0 {
		t.Fatalf("narrated while voice was off: %v", fv.spoken)
	}

	// Turn it on.
	if out := s.Submit("/voice on").Output; !strings.Contains(out, "erwacht") {
		t.Errorf("/voice on did not confirm:\n%s", out)
	}

	// A game verb is now narrated.
	s.Submit("look")
	if len(fv.spoken) == 0 || !strings.Contains(fv.spoken[len(fv.spoken)-1], "Tor") {
		t.Errorf("room text was not narrated: %v", fv.spoken)
	}

	// Slash commands are not narrated.
	before := len(fv.spoken)
	s.Submit("/help")
	if len(fv.spoken) != before {
		t.Error("a slash command should not be narrated")
	}

	// And it can be turned back off.
	s.Submit("/voice off")
	before = len(fv.spoken)
	s.Submit("look")
	if len(fv.spoken) != before {
		t.Error("narration should stop after /voice off")
	}
}

func TestVoiceUnavailableMessage(t *testing.T) {
	s := newSession(t, "") // default Noop player => Available() == false
	if out := s.Submit("/voice on").Output; !strings.Contains(out, "keine Stimme") {
		t.Errorf("expected unavailable message:\n%s", out)
	}
}

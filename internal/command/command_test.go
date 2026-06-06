package command

import "testing"

func TestVisibleExcludesHiddenAndSorts(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&Command{Name: "quit"})
	reg.Register(&Command{Name: "help"})
	reg.Register(&Command{Name: "antigravity", Hidden: true})

	got := reg.Visible()
	if len(got) != 2 {
		t.Fatalf("Visible() returned %d commands, want 2 (hidden excluded)", len(got))
	}
	if got[0].Name != "help" || got[1].Name != "quit" {
		t.Errorf("Visible() not sorted by name: got %q, %q", got[0].Name, got[1].Name)
	}
}

func TestGet(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&Command{Name: "help"})

	if _, ok := reg.Get("help"); !ok {
		t.Error("Get(\"help\") = !ok, want ok")
	}
	if _, ok := reg.Get("missing"); ok {
		t.Error("Get(\"missing\") = ok, want !ok")
	}
}

func TestRegisterIgnoresInvalid(t *testing.T) {
	reg := NewRegistry()
	reg.Register(nil)
	reg.Register(&Command{Name: ""})
	if n := len(reg.Visible()); n != 0 {
		t.Errorf("registry has %d commands after invalid registrations, want 0", n)
	}
}

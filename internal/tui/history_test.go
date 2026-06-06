package tui

import "testing"

func TestHistoryUpDown(t *testing.T) {
	var h history
	h.add("look")
	h.add("go north")
	h.add("inventory")

	// Up walks backward through entries.
	if v, ok := h.prev(); !ok || v != "inventory" {
		t.Fatalf("prev() = %q,%v want inventory", v, ok)
	}
	if v, ok := h.prev(); !ok || v != "go north" {
		t.Fatalf("prev() = %q,%v want 'go north'", v, ok)
	}
	if v, ok := h.prev(); !ok || v != "look" {
		t.Fatalf("prev() = %q,%v want look", v, ok)
	}
	// At the oldest entry, prev stops.
	if _, ok := h.prev(); ok {
		t.Fatal("prev() past oldest should report no move")
	}
	// Down walks forward; past the newest it returns the fresh empty line.
	if v, ok := h.next(); !ok || v != "go north" {
		t.Fatalf("next() = %q,%v want 'go north'", v, ok)
	}
	if v, ok := h.next(); !ok || v != "inventory" {
		t.Fatalf("next() = %q,%v want inventory", v, ok)
	}
	if v, ok := h.next(); !ok || v != "" {
		t.Fatalf("next() = %q,%v want fresh empty line", v, ok)
	}
	if _, ok := h.next(); ok {
		t.Fatal("next() past fresh line should report no move")
	}
}

func TestHistoryIgnoresEmptyAndDupes(t *testing.T) {
	var h history
	h.add("")
	h.add("look")
	h.add("look") // exact repeat ignored
	if len(h.lines) != 1 {
		t.Errorf("history kept %d lines, want 1: %v", len(h.lines), h.lines)
	}
}

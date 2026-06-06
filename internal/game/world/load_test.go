package world

import (
	"testing"
	"testing/fstest"
)

const sampleDoc = `
kind: meta
start: tor
---
kind: room
id: tor
title: Das Tor
description: Ein rostiges Tor.
exits:
  north: { to: halle }
items: [schluessel]
---
kind: room
id: halle
title: Die Halle
description: Eine weite Halle.
exits:
  south: { to: tor }
---
kind: item
id: schluessel
name: Schlüssel
description: Ein alter Schlüssel.
takeable: true
`

func loadSample(t *testing.T) *World {
	t.Helper()
	fsys := fstest.MapFS{"world/dungeon.yaml": {Data: []byte(sampleDoc)}}
	w, err := Load(fsys, "world/*.yaml")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return w
}

func TestLoadBuildsGraph(t *testing.T) {
	w := loadSample(t)

	if w.Start != "tor" {
		t.Errorf("Start = %q, want tor", w.Start)
	}
	if len(w.Rooms) != 2 {
		t.Errorf("loaded %d rooms, want 2", len(w.Rooms))
	}
	tor := w.Room("tor")
	if tor == nil || tor.Title != "Das Tor" {
		t.Fatalf("room tor not loaded correctly: %+v", tor)
	}
	if ex, ok := tor.Exits["north"]; !ok || ex.To != "halle" {
		t.Errorf("tor north exit = %+v, want -> halle", ex)
	}
	if it := w.Item("schluessel"); it == nil || !it.Takeable {
		t.Errorf("item schluessel not loaded as takeable: %+v", it)
	}
}

func TestLoadRejectsDanglingExit(t *testing.T) {
	doc := `
kind: meta
start: a
---
kind: room
id: a
title: A
description: A
exits:
  north: { to: nirgendwo }
`
	fsys := fstest.MapFS{"world/x.yaml": {Data: []byte(doc)}}
	if _, err := Load(fsys, "world/*.yaml"); err == nil {
		t.Fatal("expected error for dangling exit, got nil")
	}
}

func TestLoadRejectsMissingStart(t *testing.T) {
	doc := `
kind: room
id: a
title: A
description: A
`
	fsys := fstest.MapFS{"world/x.yaml": {Data: []byte(doc)}}
	if _, err := Load(fsys, "world/*.yaml"); err == nil {
		t.Fatal("expected error for missing start room, got nil")
	}
}

func TestNormalizeDirection(t *testing.T) {
	cases := map[string]string{
		"north": "north", "n": "north", "norden": "north", // German alias accepted
		"south": "south", "s": "south",
		"e": "east", "west": "west",
		"up": "up", "down": "down", "hoch": "up", "runter": "down",
	}
	for in, want := range cases {
		if got, ok := NormalizeDirection(in); !ok || got != want {
			t.Errorf("NormalizeDirection(%q) = %q,%v want %q", in, got, ok, want)
		}
	}
	if _, ok := NormalizeDirection("sideways"); ok {
		t.Error("NormalizeDirection(sideways) = ok, want !ok")
	}
}

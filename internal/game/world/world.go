// Package world holds the static dungeon: rooms, exits and items, plus the
// loader that reads them from embedded multi-document YAML (decision D004).
// Dynamic per-player state (position, inventory) lives in package entity.
package world

import (
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/puzzle"
)

// Exit is a directed connection from one room to another. If Puzzle is set, the
// exit is locked until that puzzle is solved.
type Exit struct {
	To     string `yaml:"to"`
	Puzzle string `yaml:"puzzle,omitempty"`
}

// Puzzle is an authored challenge: a German prompt, an acceptance check
// (puzzle.Spec), and the text shown on success plus an optional hint.
type Puzzle struct {
	ID      string      `yaml:"id"`
	Prompt  string      `yaml:"prompt"`
	Success string      `yaml:"success"`
	Hint    string      `yaml:"hint"`
	Check   puzzle.Spec `yaml:"check"`
}

// Room is a single location in the dungeon. Exits are keyed by canonical
// direction (norden/sueden/osten/westen/oben/unten). Items lists item IDs
// present in the room at game start.
type Room struct {
	ID          string          `yaml:"id"`
	Title       string          `yaml:"title"`
	Description string          `yaml:"description"`
	Exits       map[string]Exit `yaml:"exits"`
	Items       []string        `yaml:"items"`
	// Details gives authored flavour for scenery words a player might inspect
	// (keyword -> German description), e.g. "wendeltreppe" -> "…führt nach oben".
	Details map[string]string `yaml:"details"`
}

// Item is a thing the player can examine and possibly carry, wear or use.
type Item struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Takeable    bool   `yaml:"takeable"`
	Wearable    bool   `yaml:"wearable"` // can be worn with `wear`
	Light       bool   `yaml:"light"`    // when worn, lights the dungeon (color on)
	Hud         bool   `yaml:"hud"`      // when worn, enables the head-up display
}

// World is the assembled dungeon graph.
type World struct {
	Rooms   map[string]*Room
	Items   map[string]*Item
	Puzzles map[string]*Puzzle
	Start   string
}

// Room returns the room by id, or nil.
func (w *World) Room(id string) *Room { return w.Rooms[id] }

// Item returns the item by id, or nil.
func (w *World) Item(id string) *Item { return w.Items[id] }

// Puzzle returns the puzzle by id, or nil.
func (w *World) Puzzle(id string) *Puzzle { return w.Puzzles[id] }

// canonicalDirections maps the input a player might type to the canonical
// direction key used in room exits. Commands are English (north/south/...); a
// few German words are accepted as forgiving aliases so a curious kid isn't
// punished for typing "norden".
var canonicalDirections = map[string]string{
	"north": "north", "n": "north", "norden": "north", "nord": "north",
	"south": "south", "s": "south", "sueden": "south", "süden": "south",
	"east": "east", "e": "east", "osten": "east", "ost": "east",
	"west": "west", "w": "west", "westen": "west",
	"up": "up", "u": "up", "oben": "up", "hoch": "up",
	"down": "down", "d": "down", "unten": "down", "runter": "down",
}

// NormalizeDirection maps a typed direction word to its canonical exit key.
// Returns ("", false) if the word is not a recognised direction.
func NormalizeDirection(word string) (string, bool) {
	dir, ok := canonicalDirections[strings.ToLower(strings.TrimSpace(word))]
	return dir, ok
}

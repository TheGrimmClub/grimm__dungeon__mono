package engine

import (
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/entity"
)

// Class is a path the player can choose. For now it only changes how the world
// addresses them (the prompt); later phases hang abilities off the chosen class.
type Class struct {
	ID    string // the word the player types (English)
	Title string // how they're addressed afterwards
	Blurb string // German one-line description
}

// classes are the paths offered. "Human" is the class-less default (D015).
var classes = []Class{
	{"alchemist", "Alchemist", "Verwandelt Code in Tränke — Meister des Brauens und der Automatisierung."},
	{"hunter", "Jäger", "Spürt Fehler und Geheimnisse auf — schnell, neugierig, unerschrocken."},
	{"tinkerer", "Tüftler", "Baut aus Nanostaub und Logik, was es noch nicht gibt."},
}

// Classes returns the choosable paths.
func Classes() []Class { return classes }

// HasChosenClass reports whether the player has left the default "Human".
func (g *Game) HasChosenClass() bool { return g.player.Title != entity.DefaultTitle }

// ChooseClass sets the player's class by id, returning the chosen Class.
func (g *Game) ChooseClass(id string) (Class, bool) {
	for _, c := range classes {
		if c.ID == strings.ToLower(strings.TrimSpace(id)) {
			g.player.Title = c.Title
			return c, true
		}
	}
	return Class{}, false
}

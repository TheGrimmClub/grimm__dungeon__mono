package engine

import (
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
)

// German verb vocabularies. Kept generous so students can phrase things a few
// natural ways without a parser fighting them.
var (
	verbLook = set(
		"schau", "schaue", "schauen", "umschau", "umschauen", "umsehen",
		"umsieh", "sieh", "blick", "blicke", "blicken",
	)
	verbGo = set(
		"gehe", "geh", "gehen", "lauf", "laufe", "laufen", "wandere", "betritt",
	)
	verbTake = set(
		"nimm", "nimm!", "nehmen", "nehme", "heb", "hebe", "aufheben",
		"schnapp", "schnappe", "greif", "greife",
	)
	verbExamine = set(
		"untersuche", "untersuchen", "betrachte", "betrachten", "pruefe",
		"prüfe", "prüfen", "pruefen", "inspiziere", "lies", "lesen", "öffne",
		"oeffne",
	)
	verbInventory = set(
		"inventar", "inventur", "inv", "tasche", "taschen", "beutel", "i",
	)
)

// fillers are little German words that carry no meaning for the parser.
var fillers = set(
	"nach", "zum", "zur", "in", "im", "ins", "auf", "an", "den", "die", "das",
	"der", "ein", "eine", "einen", "einem", "mir", "mich", "dem", "richtung",
	"das", "mal", "bitte", "und",
)

func set(words ...string) map[string]bool {
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}

// filterFillers drops filler words, leaving the meaningful tokens.
func filterFillers(words []string) []string {
	out := make([]string, 0, len(words))
	for _, w := range words {
		if !fillers[w] {
			out = append(out, w)
		}
	}
	return out
}

// matchItem finds, among the candidate item ids, the one the query refers to.
// It matches on the item id (exact token) or a forgiving name overlap, so
// "nimm lampe" finds the "Taschenlampe". Returns "" if nothing matches.
func matchItem(w *world.World, candidates []string, query []string) string {
	if len(query) == 0 {
		return ""
	}
	q := strings.Join(query, " ")

	for _, id := range candidates {
		it := w.Item(id)
		if it == nil {
			continue
		}
		name := strings.ToLower(it.Name)
		switch {
		case strings.Contains(name, q), strings.Contains(q, name):
			return id
		}
		for _, tok := range query {
			if tok == strings.ToLower(it.ID) || strings.Contains(name, tok) {
				return id
			}
		}
	}
	return ""
}

// exitDisplay is the stable order and labels for listing a room's exits.
var exitDisplay = []struct{ key, label string }{
	{"norden", "Norden"},
	{"osten", "Osten"},
	{"sueden", "Süden"},
	{"westen", "Westen"},
	{"oben", "oben"},
	{"unten", "unten"},
}

// exitList renders a room's exits in a stable, readable order.
func (g *Game) exitList(r *world.Room) string {
	labels := make([]string, 0, len(r.Exits))
	for _, d := range exitDisplay {
		if _, ok := r.Exits[d.key]; ok {
			labels = append(labels, d.label)
		}
	}
	if len(labels) == 0 {
		return "—"
	}
	return joinList(labels)
}

// joinList joins names as "a, b und c" (German list style).
func joinList(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	default:
		return strings.Join(items[:len(items)-1], ", ") + " und " + items[len(items)-1]
	}
}

package engine

import (
	"strconv"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
)

// English verb vocabularies (commands are English by decision D001). Kept
// generous so students can phrase things a few natural ways.
var (
	verbLook = set(
		"look", "l", "around",
	)
	verbGo = set(
		"go", "move", "walk", "enter", "head",
	)
	verbTake = set(
		"take", "get", "grab", "pick", "pickup",
	)
	verbInspect = set(
		"inspect", "examine", "x", "read", "study", "check",
	)
	verbWear = set(
		"wear", "equip", "don", "put",
	)
	verbInventory = set(
		"inventory", "inv", "i", "items", "bag",
	)
	verbSolve = set(
		"solve", "answer", "antworte", "sage", "loese", "try",
	)
)

// fillers are little words that carry no meaning for the parser (English plus a
// few German ones, since the world speaks German). Direction words are NOT
// fillers — "up"/"down" must survive so "go up" works.
var fillers = set(
	"the", "a", "an", "at", "on", "to", "into", "with", "my",
	"den", "die", "das", "der", "ein", "eine", "einen", "auf", "an", "nach",
)

func set(words ...string) map[string]bool {
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}

// filterFillers drops filler words, leaving the meaningful tokens. Note: "up"
// is a filler for phrases like "pick up", but a bare direction "up" is handled
// before fillers are stripped (see Do/move).
func filterFillers(words []string) []string {
	out := make([]string, 0, len(words))
	for _, w := range words {
		if !fillers[w] {
			out = append(out, w)
		}
	}
	return out
}

// singleInt reports whether the query is a single integer token, and its value.
func singleInt(query []string) (int, bool) {
	if len(query) != 1 {
		return 0, false
	}
	n, err := strconv.Atoi(query[0])
	if err != nil {
		return 0, false
	}
	return n, true
}

// pick resolves a query against an ordered list of item ids: a single integer
// selects by 1-based position; otherwise it falls back to a forgiving name
// match. Returns "" if nothing matches.
func pick(w *world.World, list []string, query []string) string {
	if n, ok := singleInt(query); ok {
		if n >= 1 && n <= len(list) {
			return list[n-1]
		}
		return ""
	}
	return matchItem(w, list, query)
}

// matchItem finds, among the candidate item ids, the one the query refers to.
// It matches on the item id (exact token) or a forgiving name overlap, so
// "wear helm" finds the "Helm mit Stirnlampe". Returns "" if nothing matches.
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
		if strings.Contains(name, q) || strings.Contains(q, name) {
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

// matchDetail returns authored scenery flavour if the query names one of a
// room's detail keywords (case-insensitive, substring-friendly).
func matchDetail(r *world.Room, query []string) string {
	if len(r.Details) == 0 {
		return ""
	}
	q := strings.Join(query, " ")
	for key, text := range r.Details {
		k := strings.ToLower(key)
		if strings.Contains(q, k) || strings.Contains(k, q) {
			return text
		}
	}
	return ""
}

// wordInDescription reports whether a meaningful query token (>=3 letters)
// appears in the room's description text — the cue for a generic "nothing
// special" reply rather than "no such thing".
func wordInDescription(r *world.Room, query []string) bool {
	desc := strings.ToLower(r.Description)
	for _, tok := range query {
		if len(tok) >= 3 && strings.Contains(desc, tok) {
			return true
		}
	}
	return false
}

// exitOrder is the stable order for listing a room's exits. Tokens are English
// (what the player types).
var exitOrder = []string{"north", "east", "south", "west", "up", "down"}

// exitList renders a room's exits in a stable order, styling them when lit and
// marking any door still sealed by an unsolved puzzle with a trailing "*".
func (g *Game) exitList(r *world.Room, lit bool) string {
	tokens := make([]string, 0, len(r.Exits))
	for _, dir := range exitOrder {
		ex, ok := r.Exits[dir]
		if !ok {
			continue
		}
		label := dir
		if ex.Puzzle != "" && !g.player.HasSolved(ex.Puzzle) {
			label += "*"
		}
		if lit {
			label = g.palette.Exit.Render(label)
		}
		tokens = append(tokens, label)
	}
	if len(tokens) == 0 {
		return "—"
	}
	return strings.Join(tokens, ", ")
}

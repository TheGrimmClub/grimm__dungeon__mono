// Package engine is the Zork-like core: it owns the world and the player and
// interprets the German verbs the player types (schau, gehe, nimm, untersuche,
// inventar). Room prose comes from authored YAML; verb feedback comes from the
// i18n catalog.
package engine

import (
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/entity"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/i18n"
)

// Game ties the static world to the dynamic player.
type Game struct {
	world  *world.World
	player *entity.Player
}

// New starts a fresh game with the player in the world's start room.
func New(w *world.World) *Game {
	return &Game{world: w, player: entity.NewPlayer(w.Start)}
}

// Snapshot is the persistable slice of game state (see package state).
type Snapshot struct {
	Location  string   `yaml:"location"`
	Inventory []string `yaml:"inventory"`
	Visited   []string `yaml:"visited"`
}

// Snapshot captures the player's progress.
func (g *Game) Snapshot() Snapshot {
	visited := make([]string, 0, len(g.player.Visited))
	for id := range g.player.Visited {
		visited = append(visited, id)
	}
	return Snapshot{
		Location:  g.player.Location,
		Inventory: append([]string(nil), g.player.Inventory...),
		Visited:   visited,
	}
}

// Restore applies a snapshot, ignoring ids that no longer exist in the world so
// a content change can't corrupt a save.
func (g *Game) Restore(s Snapshot) {
	if g.world.Room(s.Location) != nil {
		g.player.Location = s.Location
	}
	g.player.Inventory = g.player.Inventory[:0]
	for _, id := range s.Inventory {
		if g.world.Item(id) != nil {
			g.player.Take(id)
		}
	}
	g.player.Visited = map[string]bool{g.player.Location: true}
	for _, id := range s.Visited {
		if g.world.Room(id) != nil {
			g.player.Visited[id] = true
		}
	}
}

// Intro returns the description of the starting room (shown once at launch).
func (g *Game) Intro() string { return g.look() }

// Do interprets one line of free-text input and returns the response.
func (g *Game) Do(input string) string {
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(input)))
	if len(fields) == 0 {
		return ""
	}
	verb, rest := fields[0], fields[1:]

	switch {
	case verbLook[verb]:
		if len(filterFillers(rest)) == 0 {
			return g.look()
		}
		return g.examine(rest) // e.g. "schau buch an"
	case verbGo[verb]:
		return g.move(rest)
	case verbTake[verb]:
		return g.take(rest)
	case verbExamine[verb]:
		return g.examine(rest)
	case verbInventory[verb]:
		return g.inventory()
	default:
		// Allow a bare direction word ("norden") to mean "go there".
		if dir, ok := world.NormalizeDirection(verb); ok {
			return g.moveDir(dir)
		}
		return i18n.T(i18n.KeyUnknownVerb)
	}
}

// look fully describes the current room.
func (g *Game) look() string {
	r := g.world.Room(g.player.Location)
	var b strings.Builder
	b.WriteString(r.Title)
	b.WriteString("\n")
	b.WriteString(strings.TrimRight(r.Description, "\n"))

	if loot := g.takeableHere(r); len(loot) > 0 {
		names := make([]string, len(loot))
		for i, id := range loot {
			names[i] = g.world.Item(id).Name
		}
		b.WriteString("\n")
		b.WriteString(i18n.T(i18n.KeyItemsHere, joinList(names)))
	}

	b.WriteString("\n")
	b.WriteString(i18n.T(i18n.KeyExits, g.exitList(r)))
	return b.String()
}

// move resolves a direction from the words after a "go" verb.
func (g *Game) move(rest []string) string {
	words := filterFillers(rest)
	if len(words) == 0 {
		return i18n.T(i18n.KeyWhichDirection)
	}
	dir, ok := world.NormalizeDirection(words[0])
	if !ok {
		return i18n.T(i18n.KeyUnknownDir, words[0])
	}
	return g.moveDir(dir)
}

// moveDir walks through an exit if one exists.
func (g *Game) moveDir(dir string) string {
	r := g.world.Room(g.player.Location)
	ex, ok := r.Exits[dir]
	if !ok {
		return i18n.T(i18n.KeyNoExit)
	}
	g.player.Location = ex.To
	g.player.Visit(ex.To)
	return g.look()
}

// take picks up a takeable item present in the room.
func (g *Game) take(rest []string) string {
	query := filterFillers(rest)
	if len(query) == 0 {
		return i18n.T(i18n.KeyTakeWhat)
	}
	r := g.world.Room(g.player.Location)
	id := matchItem(g.world, g.presentItems(r), query)
	if id == "" {
		return i18n.T(i18n.KeyNotHere)
	}
	it := g.world.Item(id)
	if !it.Takeable {
		return i18n.T(i18n.KeyCannotTake, it.Name)
	}
	g.player.Take(id)
	return i18n.T(i18n.KeyTaken, it.Name)
}

// examine describes an item in the room or in the inventory.
func (g *Game) examine(rest []string) string {
	query := filterFillers(rest)
	if len(query) == 0 {
		return i18n.T(i18n.KeyExamineWhat)
	}
	r := g.world.Room(g.player.Location)
	candidates := append(g.presentItems(r), g.player.Inventory...)
	id := matchItem(g.world, candidates, query)
	if id == "" {
		return i18n.T(i18n.KeyDontSee)
	}
	return strings.TrimRight(g.world.Item(id).Description, "\n")
}

// inventory lists what the player carries.
func (g *Game) inventory() string {
	if len(g.player.Inventory) == 0 {
		return i18n.T(i18n.KeyInventoryEmpty)
	}
	var b strings.Builder
	b.WriteString(i18n.T(i18n.KeyInventoryHead))
	for _, id := range g.player.Inventory {
		b.WriteString("\n  - ")
		b.WriteString(g.world.Item(id).Name)
	}
	return b.String()
}

// presentItems returns room item ids the player has not yet carried off.
func (g *Game) presentItems(r *world.Room) []string {
	out := make([]string, 0, len(r.Items))
	for _, id := range r.Items {
		if !g.player.Has(id) {
			out = append(out, id)
		}
	}
	return out
}

// takeableHere returns the present, takeable item ids (for room descriptions).
func (g *Game) takeableHere(r *world.Room) []string {
	out := make([]string, 0, len(r.Items))
	for _, id := range g.presentItems(r) {
		if it := g.world.Item(id); it != nil && it.Takeable {
			out = append(out, id)
		}
	}
	return out
}

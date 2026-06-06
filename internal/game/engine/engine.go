// Package engine is the Zork-like core: it owns the world and the player and
// interprets the English verbs the player types (look, go, take, inspect,
// inventory, wear). The rule is "type English, the world answers in German":
// commands/directions are English; room prose comes from authored German YAML.
package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/entity"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/i18n"
)

// Game ties the static world to the dynamic player.
type Game struct {
	world   *world.World
	player  *entity.Player
	palette Palette // zero value = plain (tests); app injects ColorPalette
}

// New starts a fresh game with the player in the world's start room. The
// default palette is plain; the app injects ColorPalette for the lit dungeon.
func New(w *world.World) *Game {
	return &Game{world: w, player: entity.NewPlayer(w.Start), palette: PlainPalette()}
}

// SetPalette swaps the rendering palette (the app uses ColorPalette).
func (g *Game) SetPalette(p Palette) { g.palette = p }

// Title is how the world currently addresses the player (e.g. "Human").
func (g *Game) Title() string { return g.player.Title }

// Lit reports whether the player is wearing something that lights the dungeon.
func (g *Game) Lit() bool {
	for _, id := range g.player.Worn {
		if it := g.world.Item(id); it != nil && it.Light {
			return true
		}
	}
	return false
}

// Snapshot is the persistable slice of game state (see package state).
type Snapshot struct {
	Title     string   `yaml:"title"`
	Location  string   `yaml:"location"`
	Inventory []string `yaml:"inventory"`
	Worn      []string `yaml:"worn"`
	Visited   []string `yaml:"visited"`
}

// Snapshot captures the player's progress.
func (g *Game) Snapshot() Snapshot {
	visited := make([]string, 0, len(g.player.Visited))
	for id := range g.player.Visited {
		visited = append(visited, id)
	}
	return Snapshot{
		Title:     g.player.Title,
		Location:  g.player.Location,
		Inventory: append([]string(nil), g.player.Inventory...),
		Worn:      append([]string(nil), g.player.Worn...),
		Visited:   visited,
	}
}

// Restore applies a snapshot, ignoring ids that no longer exist in the world so
// a content change can't corrupt a save.
func (g *Game) Restore(s Snapshot) {
	if s.Title != "" {
		g.player.Title = s.Title
	}
	if g.world.Room(s.Location) != nil {
		g.player.Location = s.Location
	}
	g.player.Inventory = g.player.Inventory[:0]
	for _, id := range s.Inventory {
		if g.world.Item(id) != nil {
			g.player.Take(id)
		}
	}
	g.player.Worn = g.player.Worn[:0]
	for _, id := range s.Worn {
		if g.world.Item(id) != nil {
			g.player.Wear(id)
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
		return g.inspect(rest) // "look at <thing>"
	case verbGo[verb]:
		return g.move(rest)
	case verbTake[verb]:
		return g.take(rest)
	case verbInspect[verb]:
		return g.inspect(rest)
	case verbWear[verb]:
		return g.wear(rest)
	case verbInventory[verb]:
		return g.inventory()
	default:
		// A bare direction ("north") means "go there"; a bare number inspects
		// that item in the room.
		if dir, ok := world.NormalizeDirection(verb); ok && len(rest) == 0 {
			return g.moveDir(dir)
		}
		if _, ok := singleInt(fields); ok {
			return g.inspect(fields)
		}
		return i18n.T(i18n.KeyUnknownVerb)
	}
}

// look fully describes the current room, numbering every item in it.
func (g *Game) look() string {
	r := g.world.Room(g.player.Location)
	lit := g.Lit()

	var b strings.Builder
	b.WriteString(g.styleTitle(r.Title, lit))
	b.WriteString("\n")
	b.WriteString(strings.TrimRight(r.Description, "\n"))

	if items := g.presentItems(r); len(items) > 0 {
		b.WriteString("\n")
		b.WriteString(i18n.T(i18n.KeyItemsHere))
		for i, id := range items {
			label := fmt.Sprintf("[%d] %s", i+1, g.world.Item(id).Name)
			b.WriteString("\n  ")
			b.WriteString(g.styleItem(label, lit))
		}
	}

	b.WriteString("\n")
	b.WriteString(i18n.T(i18n.KeyExits, g.exitList(r, lit)))

	out := b.String()
	if !lit {
		out = g.palette.Dim.Render(out) // the dungeon is dark
	}
	return out
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

// take picks up a takeable item present in the room (by [number] or name).
func (g *Game) take(rest []string) string {
	query := filterFillers(rest)
	if len(query) == 0 {
		return i18n.T(i18n.KeyTakeWhat)
	}
	r := g.world.Room(g.player.Location)
	id := pick(g.world, g.presentItems(r), query)
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

// inspect describes an item in the room or in the inventory (by number or name).
func (g *Game) inspect(rest []string) string {
	query := filterFillers(rest)
	if len(query) == 0 {
		return i18n.T(i18n.KeyExamineWhat)
	}
	r := g.world.Room(g.player.Location)
	id := pick(g.world, g.presentItems(r), query)
	if id == "" {
		id = pick(g.world, g.player.Inventory, query)
	}
	if id == "" {
		return i18n.T(i18n.KeyDontSee)
	}
	return strings.TrimRight(g.world.Item(id).Description, "\n")
}

// wear equips a wearable item; wearing a light source floods the dungeon with
// colour — the first big "aha".
func (g *Game) wear(rest []string) string {
	query := filterFillers(rest)
	if len(query) == 0 {
		return i18n.T(i18n.KeyWearWhat)
	}
	r := g.world.Room(g.player.Location)
	id := pick(g.world, g.presentItems(r), query)
	if id == "" {
		id = pick(g.world, g.player.Inventory, query)
	}
	if id == "" {
		return i18n.T(i18n.KeyDontSee)
	}
	it := g.world.Item(id)
	if !it.Wearable {
		return i18n.T(i18n.KeyCannotWear, it.Name)
	}
	if g.player.Wears(id) {
		return i18n.T(i18n.KeyAlreadyWorn, it.Name)
	}

	wasLit := g.Lit()
	g.player.Take(id) // auto-pick-up if it was lying in the room
	g.player.Wear(id)

	if it.Light && !wasLit {
		// The dramatic moment: announce, then re-render the now-lit room.
		return i18n.T(i18n.KeyHeadlampOn) + "\n\n" + g.look()
	}
	return i18n.T(i18n.KeyWorn, it.Name)
}

// inventory lists what the player carries as a 10-slot hotbar.
func (g *Game) inventory() string {
	if len(g.player.Inventory) == 0 {
		return i18n.T(i18n.KeyInventoryEmpty)
	}
	lit := g.Lit()
	var b strings.Builder
	b.WriteString(i18n.T(i18n.KeyInventoryHead))
	for i, id := range g.player.Inventory {
		name := g.world.Item(id).Name
		if g.player.Wears(id) {
			name += " " + i18n.T(i18n.KeyWornTag)
		}
		line := fmt.Sprintf("[%s] %s", hotbarSlot(i), name)
		b.WriteString("\n  ")
		b.WriteString(g.styleItem(line, lit))
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

// --- styling helpers (no-ops under the zero palette / in tests) ---

func (g *Game) styleTitle(s string, lit bool) string {
	if lit {
		return g.palette.Title.Render(s)
	}
	return s
}

func (g *Game) styleItem(s string, lit bool) string {
	if lit {
		return g.palette.Item.Render(s)
	}
	return s
}

func hotbarSlot(i int) string {
	switch {
	case i < 9:
		return strconv.Itoa(i + 1) // slots 1..9
	case i == 9:
		return "0" // the 10th slot
	default:
		return "·" // beyond the hotbar
	}
}

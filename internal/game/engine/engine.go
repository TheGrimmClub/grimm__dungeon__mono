// Package engine is the Zork-like core: it owns the world and the player and
// interprets the English verbs the player types (look, go, take, inspect,
// inventory, wear). The rule is "type English, the world answers in German":
// commands/directions are English; room prose comes from authored German YAML.
package engine

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/entity"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/i18n"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/puzzle"
)

// Game ties the static world to the dynamic player.
type Game struct {
	world   *world.World
	player  *entity.Player
	palette Palette // zero value = plain (tests); app injects ColorPalette

	active  string                  // id of the puzzle currently blocking the player
	checks  map[string]puzzle.Check // lazily-built checks, keyed by puzzle id
	workDir string                  // student's working directory (artifact/behavioral)
}

// New starts a fresh game with the player in the world's start room. The
// default palette is plain; the app injects ColorPalette for the lit dungeon.
func New(w *world.World) *Game {
	return &Game{
		world:   w,
		player:  entity.NewPlayer(w.Start),
		palette: PlainPalette(),
		checks:  make(map[string]puzzle.Check),
	}
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
	Solved    []string `yaml:"solved"`
}

// Snapshot captures the player's progress.
func (g *Game) Snapshot() Snapshot {
	return Snapshot{
		Title:     g.player.Title,
		Location:  g.player.Location,
		Inventory: append([]string(nil), g.player.Inventory...),
		Worn:      append([]string(nil), g.player.Worn...),
		Visited:   keys(g.player.Visited),
		Solved:    keys(g.player.Solved),
	}
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
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
	g.player.Solved = map[string]bool{}
	for _, id := range s.Solved {
		if g.world.Puzzle(id) != nil {
			g.player.Solve(id)
		}
	}
}

// SetWorkDir sets the student's working directory used by artifact/behavioral
// checks (Phase 3 points this at the alchemist repo; default "").
func (g *Game) SetWorkDir(dir string) { g.workDir = dir }

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
	case verbSolve[verb]:
		return g.solve(rest)
	default:
		// A bare direction ("north") means "go there"; a bare number inspects
		// that item in the room.
		if dir, ok := world.NormalizeDirection(verb); ok && len(rest) == 0 {
			return g.moveDir(dir)
		}
		if _, ok := singleInt(fields); ok {
			return g.inspect(fields)
		}
		// While a puzzle blocks the way, treat free text as an answer attempt.
		if g.active != "" {
			return g.solve(fields)
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
	if g.hasLockedExit(r) {
		b.WriteString("\n")
		b.WriteString(i18n.T(i18n.KeyLockedFootnote))
	}

	out := b.String()
	if !lit {
		out = g.palette.Dim.Render(out) // the dungeon is dark
	}
	return out
}

// hasLockedExit reports whether any of a room's exits is still sealed.
func (g *Game) hasLockedExit(r *world.Room) bool {
	for _, ex := range r.Exits {
		if ex.Puzzle != "" && !g.player.HasSolved(ex.Puzzle) {
			return true
		}
	}
	return false
}

// move resolves a direction from the words after a "go" verb. It scans for the
// first recognised direction, so "go up", "go to north" and "go the north" all
// work.
func (g *Game) move(rest []string) string {
	for _, w := range rest {
		if dir, ok := world.NormalizeDirection(w); ok {
			return g.moveDir(dir)
		}
	}
	words := filterFillers(rest)
	if len(words) == 0 {
		return i18n.T(i18n.KeyWhichDirection)
	}
	return i18n.T(i18n.KeyUnknownDir, words[0])
}

// moveDir walks through an exit if one exists and is not locked. A locked exit
// presents its puzzle instead and becomes the active puzzle.
func (g *Game) moveDir(dir string) string {
	r := g.world.Room(g.player.Location)
	ex, ok := r.Exits[dir]
	if !ok {
		return i18n.T(i18n.KeyNoExit)
	}
	if ex.Puzzle != "" && !g.player.HasSolved(ex.Puzzle) {
		g.active = ex.Puzzle
		return g.puzzlePrompt(ex.Puzzle)
	}
	g.player.Location = ex.To
	g.player.Visit(ex.To)
	return g.look()
}

// puzzlePrompt shows a blocked puzzle's prompt and how to answer it.
func (g *Game) puzzlePrompt(id string) string {
	p := g.world.Puzzle(id)
	var b strings.Builder
	b.WriteString(strings.TrimRight(p.Prompt, "\n"))
	b.WriteString("\n\n")
	b.WriteString(i18n.T(i18n.KeySolveHint))
	return b.String()
}

// solve attempts the active puzzle with the rest of the player's input.
func (g *Game) solve(rest []string) string {
	if g.active == "" {
		return i18n.T(i18n.KeyNoActivePuzzle)
	}
	answer := strings.TrimSpace(strings.Join(rest, " "))
	p := g.world.Puzzle(g.active)
	// Riddles need a typed answer; artifact/behavioral checks read the world, so
	// a bare `solve` (after doing the real-world action) is allowed for them.
	if answer == "" && p.Check.Kind == "answer" {
		return i18n.T(i18n.KeySolveWhat)
	}
	check, err := g.checkFor(g.active)
	if err != nil {
		return i18n.T(i18n.KeyPuzzleBroken)
	}

	res := check.Verify(context.Background(), puzzle.Input{Answer: answer, WorkDir: g.workDir})
	if res.Passed {
		g.player.Solve(g.active)
		g.active = ""
		if p.Success != "" {
			return strings.TrimRight(p.Success, "\n")
		}
		return i18n.T(i18n.KeyPuzzleSolved)
	}

	msg := i18n.T(i18n.KeyPuzzleWrong)
	if res.Detail != "" {
		msg += " " + res.Detail
	}
	if p.Hint != "" {
		msg += "\n" + i18n.T(i18n.KeyHintLabel) + " " + p.Hint
	}
	return msg
}

// checkFor lazily builds and caches the Check for a puzzle.
func (g *Game) checkFor(id string) (puzzle.Check, error) {
	if c, ok := g.checks[id]; ok {
		return c, nil
	}
	c, err := puzzle.Build(g.world.Puzzle(id).Check)
	if err != nil {
		return nil, err
	}
	g.checks[id] = c
	return c, nil
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

// inspect describes an item in the room or inventory (by number or name). If no
// item matches, it falls back to authored scenery details, then to a gentle
// generic reply for any word that appears in the room text — so reading the
// description is always rewarded.
func (g *Game) inspect(rest []string) string {
	query := filterFillers(rest)
	if len(query) == 0 {
		return i18n.T(i18n.KeyExamineWhat)
	}
	r := g.world.Room(g.player.Location)

	if id := pick(g.world, g.presentItems(r), query); id != "" {
		return strings.TrimRight(g.world.Item(id).Description, "\n")
	}
	if id := pick(g.world, g.player.Inventory, query); id != "" {
		return strings.TrimRight(g.world.Item(id).Description, "\n")
	}
	if d := matchDetail(r, query); d != "" {
		return strings.TrimRight(d, "\n")
	}
	if wordInDescription(r, query) {
		return i18n.T(i18n.KeyNothingSpecial)
	}
	return i18n.T(i18n.KeyDontSee)
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

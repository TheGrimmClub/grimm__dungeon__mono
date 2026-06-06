// Package entity holds the dynamic, per-player game state (position, inventory,
// what has been visited) — as opposed to package world, which is the static
// authored dungeon.
package entity

// Player is the dynamic state of the person playing.
type Player struct {
	Location  string          // current room id
	Inventory []string        // item ids carried
	Visited   map[string]bool // room ids already seen
}

// NewPlayer starts a player in the given room.
func NewPlayer(start string) *Player {
	return &Player{
		Location:  start,
		Inventory: []string{},
		Visited:   map[string]bool{start: true},
	}
}

// Has reports whether the player carries the item.
func (p *Player) Has(id string) bool {
	for _, it := range p.Inventory {
		if it == id {
			return true
		}
	}
	return false
}

// Take adds an item to the inventory (no-op if already held).
func (p *Player) Take(id string) {
	if !p.Has(id) {
		p.Inventory = append(p.Inventory, id)
	}
}

// Visit marks a room as seen and reports whether it was new.
func (p *Player) Visit(id string) (firstTime bool) {
	if p.Visited == nil {
		p.Visited = map[string]bool{}
	}
	if p.Visited[id] {
		return false
	}
	p.Visited[id] = true
	return true
}

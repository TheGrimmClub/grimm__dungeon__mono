package engine

// ExitInfo is one move option for the HUD map.
type ExitInfo struct {
	Dir    string
	Locked bool
}

// ItemInfo is one inventory slot for the HUD.
type ItemInfo struct {
	Slot string // hotbar label: "1".."9","0", or "·"
	Name string
	Worn bool
}

// HUDActive reports whether the player wears something that provides a head-up
// display (the helmet). Like the headlamp's colour, the HUD is a helmet feature.
func (g *Game) HUDActive() bool {
	for _, id := range g.player.Worn {
		if it := g.world.Item(id); it != nil && it.Hud {
			return true
		}
	}
	return false
}

// RoomTitle is the title of the room the player is in.
func (g *Game) RoomTitle() string {
	if r := g.world.Room(g.player.Location); r != nil {
		return r.Title
	}
	return ""
}

// ExitsView lists the current room's move options in stable order, flagging any
// still sealed by an unsolved puzzle.
func (g *Game) ExitsView() []ExitInfo {
	r := g.world.Room(g.player.Location)
	if r == nil {
		return nil
	}
	out := make([]ExitInfo, 0, len(r.Exits))
	for _, dir := range exitOrder {
		ex, ok := r.Exits[dir]
		if !ok {
			continue
		}
		out = append(out, ExitInfo{
			Dir:    dir,
			Locked: ex.Puzzle != "" && !g.player.HasSolved(ex.Puzzle),
		})
	}
	return out
}

// InventoryView lists carried items as hotbar slots for the HUD.
func (g *Game) InventoryView() []ItemInfo {
	out := make([]ItemInfo, 0, len(g.player.Inventory))
	for i, id := range g.player.Inventory {
		it := g.world.Item(id)
		if it == nil {
			continue
		}
		out = append(out, ItemInfo{
			Slot: hotbarSlot(i),
			Name: it.Name,
			Worn: g.player.Wears(id),
		})
	}
	return out
}

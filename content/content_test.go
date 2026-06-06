package content_test

import (
	"testing"

	"github.com/TheGrimmClub/grimm__dungeon__mono/content"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
)

// TestEmbeddedDungeonLoads guards the authored content: it must always parse and
// pass referential validation, so a typo in a YAML file fails CI, not a player.
func TestEmbeddedDungeonLoads(t *testing.T) {
	w, err := world.Load(content.FS, content.WorldGlob)
	if err != nil {
		t.Fatalf("embedded dungeon failed to load: %v", err)
	}
	if w.Start != "tor" {
		t.Errorf("start room = %q, want tor", w.Start)
	}
	if len(w.Rooms) != 6 {
		t.Errorf("rooms = %d, want 6", len(w.Rooms))
	}
	if len(w.Puzzles) != 2 {
		t.Errorf("puzzles = %d, want 2", len(w.Puzzles))
	}
	// Spot-check the reward loot in the final room.
	if it := w.Item("kristallkern"); it == nil || !it.Takeable {
		t.Errorf("kristallkern missing or not takeable: %+v", it)
	}
}

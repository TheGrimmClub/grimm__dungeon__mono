package engine

import (
	"strings"
	"testing"

	"github.com/TheGrimmClub/grimm__dungeon__mono/content"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
)

func newGame(t *testing.T) *Game {
	t.Helper()
	w, err := world.Load(content.FS, content.WorldGlob)
	if err != nil {
		t.Fatalf("load world: %v", err)
	}
	return New(w)
}

func TestIntroDescribesStartRoom(t *testing.T) {
	g := newGame(t)
	intro := g.Intro()
	if !strings.Contains(intro, "Das verwunschene Tor") {
		t.Errorf("intro missing start room title:\n%s", intro)
	}
	if !strings.Contains(intro, "Ausgänge: Norden") {
		t.Errorf("intro missing exits line:\n%s", intro)
	}
	if !strings.Contains(intro, "Taschenlampe") {
		t.Errorf("intro should mention the loot in the room:\n%s", intro)
	}
}

func TestMovementAndBareDirection(t *testing.T) {
	g := newGame(t)
	// "gehe norden" -> the hall.
	if got := g.Do("gehe norden"); !strings.Contains(got, "Halle der schlafenden Maschinen") {
		t.Errorf("gehe norden did not reach the hall:\n%s", got)
	}
	// Bare direction "sueden" -> back to the gate.
	if got := g.Do("sueden"); !strings.Contains(got, "Das verwunschene Tor") {
		t.Errorf("bare 'sueden' did not return to the gate:\n%s", got)
	}
	// No exit east from the gate.
	if got := g.Do("gehe osten"); !strings.Contains(got, "kein Weg") {
		t.Errorf("expected 'no path' going east from gate, got:\n%s", got)
	}
}

func TestTakeAndInventory(t *testing.T) {
	g := newGame(t)
	if got := g.Do("nimm lampe"); !strings.Contains(got, "Du nimmst") {
		t.Errorf("could not take the lamp by partial name:\n%s", got)
	}
	if got := g.Do("inventar"); !strings.Contains(got, "Taschenlampe") {
		t.Errorf("lamp not in inventory:\n%s", got)
	}
	// Once taken, the room no longer lists it.
	if got := g.Do("schau"); strings.Contains(got, "Hier liegt: Taschenlampe") {
		t.Errorf("room still advertises the taken lamp:\n%s", got)
	}
}

func TestCannotTakeScenery(t *testing.T) {
	g := newGame(t)
	g.Do("gehe norden") // hall, contains the non-takeable androide
	if got := g.Do("nimm androide"); !strings.Contains(got, "lässt sich nicht nehmen") {
		t.Errorf("expected scenery to be untakeable, got:\n%s", got)
	}
}

func TestExamineRoomAndInventoryItems(t *testing.T) {
	g := newGame(t)
	if got := g.Do("untersuche lampe"); !strings.Contains(got, "Glühwürmchen") {
		t.Errorf("examine of room item failed:\n%s", got)
	}
	g.Do("nimm lampe")
	if got := g.Do("untersuche lampe"); !strings.Contains(got, "Glühwürmchen") {
		t.Errorf("examine of carried item failed:\n%s", got)
	}
	if got := g.Do("untersuche drache"); !strings.Contains(got, "siehst du hier nicht") {
		t.Errorf("examine of absent item should fail gracefully:\n%s", got)
	}
}

func TestUnknownVerb(t *testing.T) {
	g := newGame(t)
	if got := g.Do("tanze"); !strings.Contains(got, "verstehe ich nicht") {
		t.Errorf("unknown verb not handled:\n%s", got)
	}
}

func TestSnapshotRoundTrip(t *testing.T) {
	g := newGame(t)
	g.Do("nimm lampe")
	g.Do("gehe norden")
	snap := g.Snapshot()

	g2 := newGame(t)
	g2.Restore(snap)
	if g2.player.Location != "halle" {
		t.Errorf("restored location = %q, want halle", g2.player.Location)
	}
	if !g2.player.Has("taschenlampe") {
		t.Error("restored player should carry the lamp")
	}
	if !g2.player.Visited["tor"] || !g2.player.Visited["halle"] {
		t.Error("restored visited set incomplete")
	}
}

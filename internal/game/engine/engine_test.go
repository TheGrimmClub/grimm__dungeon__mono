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
	for _, want := range []string{
		"Das verwunschene Tor",    // room title
		"Ausgänge: north",         // exits use English tokens
		"[1] Helm mit Stirnlampe", // numbered item
	} {
		if !strings.Contains(intro, want) {
			t.Errorf("intro missing %q:\n%s", want, intro)
		}
	}
}

func TestMovementAndBareDirection(t *testing.T) {
	g := newGame(t)
	if got := g.Do("go north"); !strings.Contains(got, "Halle der schlafenden Maschinen") {
		t.Errorf("go north did not reach the hall:\n%s", got)
	}
	if got := g.Do("south"); !strings.Contains(got, "Das verwunschene Tor") {
		t.Errorf("bare 'south' did not return to the gate:\n%s", got)
	}
	if got := g.Do("go east"); !strings.Contains(got, "kein Weg") {
		t.Errorf("expected 'no path' east from gate:\n%s", got)
	}
}

func TestTakeByNameAndNumber(t *testing.T) {
	g := newGame(t)
	if got := g.Do("take helm"); !strings.Contains(got, "Du nimmst") {
		t.Errorf("take by name failed:\n%s", got)
	}
	if got := g.Do("inventory"); !strings.Contains(got, "Helm mit Stirnlampe") {
		t.Errorf("helm not in inventory:\n%s", got)
	}
	// Once carried, the room no longer lists it.
	if got := g.Do("look"); strings.Contains(got, "Helm mit Stirnlampe") {
		t.Errorf("room still lists the taken helm:\n%s", got)
	}

	// Number selection on a fresh game.
	g2 := newGame(t)
	if got := g2.Do("take 1"); !strings.Contains(got, "Du nimmst") {
		t.Errorf("take by number failed:\n%s", got)
	}
}

func TestWearHeadlampLightsUp(t *testing.T) {
	g := newGame(t)
	if g.Lit() {
		t.Fatal("dungeon should start dark")
	}
	// `wear 1` auto-takes the helm from the room and lights up.
	got := g.Do("wear 1")
	if !strings.Contains(got, "Farbe") {
		t.Errorf("wearing the headlamp did not announce the light:\n%s", got)
	}
	if !g.Lit() {
		t.Error("Lit() should be true after wearing the headlamp")
	}
	// Re-wearing reports it is already worn.
	if got := g.Do("wear helm"); !strings.Contains(got, "bereits") {
		t.Errorf("expected already-worn message:\n%s", got)
	}
}

func TestCannotWearScenery(t *testing.T) {
	g := newGame(t)
	g.Do("go north") // hall with the non-wearable androide
	if got := g.Do("wear androide"); !strings.Contains(got, "nicht anlegen") {
		t.Errorf("expected androide to be unwearable:\n%s", got)
	}
}

func TestInspectRoomAndInventory(t *testing.T) {
	g := newGame(t)
	if got := g.Do("inspect 1"); !strings.Contains(got, "Glühfaden") {
		t.Errorf("inspect room item by number failed:\n%s", got)
	}
	g.Do("take helm")
	if got := g.Do("inspect helm"); !strings.Contains(got, "Glühfaden") {
		t.Errorf("inspect carried item by name failed:\n%s", got)
	}
	if got := g.Do("inspect dragon"); !strings.Contains(got, "siehst du hier nicht") {
		t.Errorf("inspect of absent item should fail gracefully:\n%s", got)
	}
}

func TestInventoryHotbar(t *testing.T) {
	g := newGame(t)
	g.Do("wear 1") // take + wear the helm
	got := g.Do("inventory")
	if !strings.Contains(got, "[1]") || !strings.Contains(got, "Helm mit Stirnlampe") {
		t.Errorf("hotbar missing slot/name:\n%s", got)
	}
	if !strings.Contains(got, "(angelegt)") {
		t.Errorf("worn item not tagged in hotbar:\n%s", got)
	}
}

func TestUnknownVerb(t *testing.T) {
	g := newGame(t)
	if got := g.Do("dance"); !strings.Contains(got, "verstehe ich nicht") {
		t.Errorf("unknown verb not handled:\n%s", got)
	}
}

func TestSnapshotRoundTrip(t *testing.T) {
	g := newGame(t)
	g.Do("wear 1") // take + wear helm
	g.Do("go north")
	snap := g.Snapshot()

	g2 := newGame(t)
	g2.Restore(snap)
	if g2.player.Location != "halle" {
		t.Errorf("restored location = %q, want halle", g2.player.Location)
	}
	if !g2.player.Has("helm") || !g2.player.Wears("helm") {
		t.Error("restored player should carry and wear the helm")
	}
	if !g2.Lit() {
		t.Error("restored game should still be lit")
	}
	if g2.player.Title != "Human" {
		t.Errorf("restored title = %q, want Human", g2.player.Title)
	}
}

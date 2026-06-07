package engine

import "testing"

func TestHUDActivatesWithHelmet(t *testing.T) {
	g := newGame(t)
	if g.HUDActive() {
		t.Fatal("HUD should be off before wearing the helmet")
	}
	g.Do("wear 1") // take + wear the helmet (hud:true)
	if !g.HUDActive() {
		t.Error("HUD should be active once the helmet is worn")
	}
}

func TestHUDViews(t *testing.T) {
	g := newGame(t)
	g.Do("wear 1") // helmet now carried + worn
	g.Do("go north")
	g.Do("go north") // werkstatt: exits south, up(locked), plus items

	if g.RoomTitle() != "Die Werkstatt des Nanoschmieds" {
		t.Errorf("RoomTitle = %q", g.RoomTitle())
	}

	exits := g.ExitsView()
	var sawLockedUp bool
	for _, e := range exits {
		if e.Dir == "up" && e.Locked {
			sawLockedUp = true
		}
	}
	if !sawLockedUp {
		t.Errorf("expected a locked 'up' exit in %+v", exits)
	}

	inv := g.InventoryView()
	if len(inv) == 0 || inv[0].Slot != "1" {
		t.Fatalf("inventory view wrong: %+v", inv)
	}
	if inv[0].Name != "Helm mit Stirnlampe" || !inv[0].Worn {
		t.Errorf("first slot should be the worn helmet: %+v", inv[0])
	}
}

package state

import (
	"errors"
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/engine"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "save.yaml")

	want := engine.Snapshot{
		Location:  "halle",
		Inventory: []string{"taschenlampe", "maerchenbuch"},
		Visited:   []string{"tor", "halle"},
	}
	if err := Save(path, want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if !Exists(path) {
		t.Fatal("Exists() = false after Save")
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("round trip mismatch:\n got=%+v\nwant=%+v", got, want)
	}
}

func TestLoadMissingReturnsNotExist(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "absent.yaml"))
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Load(missing) err = %v, want fs.ErrNotExist", err)
	}
}

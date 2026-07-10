package content_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TheGrimmClub/grimm__dungeon__mono/content"
)

// TestSeedWorkspace guards the package bridge: after seeding, a work dir must
// hold an importable `grimm` package so behavioral solutions can use Actor.
func TestSeedWorkspace(t *testing.T) {
	dir := t.TempDir()
	if err := content.SeedWorkspace(dir); err != nil {
		t.Fatalf("SeedWorkspace: %v", err)
	}

	for _, name := range []string{"grimm/__init__.py", "grimm/actor.py"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("missing seeded file %s: %v", name, err)
		}
	}

	// The Actor class must be present so `from grimm import Actor` resolves.
	data, err := os.ReadFile(filepath.Join(dir, "grimm", "actor.py"))
	if err != nil {
		t.Fatalf("read actor.py: %v", err)
	}
	if !strings.Contains(string(data), "class Actor") {
		t.Errorf("seeded actor.py does not define Actor")
	}

	// Idempotent: seeding again must not error.
	if err := content.SeedWorkspace(dir); err != nil {
		t.Fatalf("SeedWorkspace (second call): %v", err)
	}
}

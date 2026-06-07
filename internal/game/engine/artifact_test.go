package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestArtifactGateOpensWithRepo drives the repo-tor door: it stays shut until a
// .git directory exists in the work dir, then opens. (We create .git/HEAD
// directly, so the check is exercised without needing git in the test.)
func TestArtifactGateOpensWithRepo(t *testing.T) {
	g := newGame(t)
	dir := t.TempDir()
	g.SetWorkDir(dir)

	g.Do("go north") // halle, whose west door is the artifact gate

	// Going west presents the puzzle; a bare solve fails (no repo yet).
	g.Do("go west")
	if got := g.Do("solve"); !strings.Contains(got, "war nicht die Lösung") {
		t.Errorf("artifact gate should reject before the repo exists:\n%s", got)
	}

	// "Create the repository" — what /alchemist init would produce.
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".git", "HEAD"), []byte("ref: refs/heads/main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if got := g.Do("solve"); !strings.Contains(got, "Lichtfäden") {
		t.Errorf("artifact gate should open once the repo exists:\n%s", got)
	}
	if got := g.Do("go west"); !strings.Contains(got, "Archiv der Versionen") {
		t.Errorf("archive should be reachable after solving:\n%s", got)
	}
}

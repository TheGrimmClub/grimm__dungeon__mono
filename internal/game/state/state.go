// Package state persists a player's progress to disk as YAML. Game progress is
// deliberately kept separate from the student's own work, which lives in their
// git repo via alchemist (decision D005/D011).
package state

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/engine"
)

// saveVersion lets future formats migrate old saves instead of breaking them.
const saveVersion = 1

// file is the on-disk envelope around a game snapshot.
type file struct {
	Version int             `yaml:"version"`
	Game    engine.Snapshot `yaml:"game"`
}

// DefaultPath is where grimm keeps its save: ~/.grimm/save.yaml.
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("state: locating home dir: %w", err)
	}
	return filepath.Join(home, ".grimm", "save.yaml"), nil
}

// Exists reports whether a save file is present at path.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Save writes the snapshot to path, creating the parent directory as needed.
func Save(path string, snap engine.Snapshot) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("state: creating save dir: %w", err)
	}
	data, err := yaml.Marshal(file{Version: saveVersion, Game: snap})
	if err != nil {
		return fmt.Errorf("state: marshaling save: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("state: writing save: %w", err)
	}
	return nil
}

// Load reads a snapshot from path. A missing file returns (zero, fs.ErrNotExist)
// so callers can start a fresh game.
func Load(path string) (engine.Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return engine.Snapshot{}, fs.ErrNotExist
		}
		return engine.Snapshot{}, fmt.Errorf("state: reading save: %w", err)
	}
	var f file
	if err := yaml.Unmarshal(data, &f); err != nil {
		return engine.Snapshot{}, fmt.Errorf("state: parsing save: %w", err)
	}
	return f.Game, nil
}

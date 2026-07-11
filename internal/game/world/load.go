package world

import (
	"errors"
	"fmt"
	"io/fs"

	syon "github.com/object-notation-environment/safe-yaml-object-notation/syon-go"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/puzzle"
)

// fileDoc is the top-level shape of a content .syon file (decision D013). Any of
// the sections may be present; files are split by concern (dungeon rooms, items,
// puzzles, the curriculum wing) but the loader merges them into one World.
type fileDoc struct {
	Start   string   `yaml:"start"`
	Rooms   []Room   `yaml:"rooms"`
	Items   []Item   `yaml:"items"`
	Puzzles []Puzzle `yaml:"puzzles"`
}

// Load reads every file matching glob from fsys, parses each as a single SYON
// document, and assembles them into a World. SYON is safe YAML — no implicit
// typing, no tags/anchors/flow (D013) — parsed by the native Go implementation,
// so the single-binary build (D002) keeps working with no cgo.
//
// Decoupling from the embed.FS (it lives in package content) keeps Load testable
// with an in-memory fs.FS.
func Load(fsys fs.FS, glob string) (*World, error) {
	files, err := fs.Glob(fsys, glob)
	if err != nil {
		return nil, fmt.Errorf("world: glob %q: %w", glob, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("world: no content files match %q", glob)
	}

	w := &World{
		Rooms:   make(map[string]*Room),
		Items:   make(map[string]*Item),
		Puzzles: make(map[string]*Puzzle),
	}
	for _, name := range files {
		if err := loadFile(fsys, name, w); err != nil {
			return nil, err
		}
	}
	if err := w.validate(); err != nil {
		return nil, err
	}
	return w, nil
}

func loadFile(fsys fs.FS, name string, w *World) error {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return fmt.Errorf("world: open %s: %w", name, err)
	}

	var doc fileDoc
	if err := syon.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("world: %s: %w", name, err)
	}

	if doc.Start != "" {
		w.Start = doc.Start
	}
	for i := range doc.Rooms {
		r := doc.Rooms[i]
		if r.ID == "" {
			return fmt.Errorf("world: %s: a room is missing an id", name)
		}
		w.Rooms[r.ID] = &r
	}
	for i := range doc.Items {
		it := doc.Items[i]
		if it.ID == "" {
			return fmt.Errorf("world: %s: an item is missing an id", name)
		}
		w.Items[it.ID] = &it
	}
	for i := range doc.Puzzles {
		p := doc.Puzzles[i]
		if p.ID == "" {
			return fmt.Errorf("world: %s: a puzzle is missing an id", name)
		}
		w.Puzzles[p.ID] = &p
	}
	return nil
}

// validate checks referential integrity: a start room exists and every exit
// and room item points at something real.
func (w *World) validate() error {
	if w.Start == "" {
		return errors.New("world: no start room declared (add `start:` to a content file)")
	}
	if w.Rooms[w.Start] == nil {
		return fmt.Errorf("world: start room %q does not exist", w.Start)
	}
	for id, r := range w.Rooms {
		for dir, ex := range r.Exits {
			if w.Rooms[ex.To] == nil {
				return fmt.Errorf("world: room %q exit %q points at unknown room %q", id, dir, ex.To)
			}
			if ex.Puzzle != "" && w.Puzzles[ex.Puzzle] == nil {
				return fmt.Errorf("world: room %q exit %q references unknown puzzle %q", id, dir, ex.Puzzle)
			}
		}
		for _, itemID := range r.Items {
			if w.Items[itemID] == nil {
				return fmt.Errorf("world: room %q references unknown item %q", id, itemID)
			}
		}
	}
	// Every puzzle's check spec must build, so authoring errors fail at load.
	for id, p := range w.Puzzles {
		if _, err := puzzle.Build(p.Check); err != nil {
			return fmt.Errorf("world: puzzle %q: %w", id, err)
		}
	}
	return nil
}

package world

import (
	"errors"
	"fmt"
	"io"
	"io/fs"

	"gopkg.in/yaml.v3"
)

// Load reads every file matching glob from fsys, parses each as a stream of
// YAML documents, and assembles them into a World. Each document is routed by
// its `kind` field: "room", "item" or "meta" (which names the start room).
//
// Decoupling from the embed.FS (it lives in package content) keeps Load
// testable with an in-memory fs.FS.
func Load(fsys fs.FS, glob string) (*World, error) {
	files, err := fs.Glob(fsys, glob)
	if err != nil {
		return nil, fmt.Errorf("world: glob %q: %w", glob, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("world: no content files match %q", glob)
	}

	w := &World{Rooms: make(map[string]*Room), Items: make(map[string]*Item)}
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
	f, err := fsys.Open(name)
	if err != nil {
		return fmt.Errorf("world: open %s: %w", name, err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	for i := 0; ; i++ {
		var node yaml.Node
		if err := dec.Decode(&node); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("world: %s doc %d: %w", name, i, err)
		}
		if node.Kind == 0 { // empty document (e.g. trailing "---")
			continue
		}
		if err := routeDoc(&node, name, w); err != nil {
			return err
		}
	}
}

// routeDoc decodes a single document into the right type based on its kind.
func routeDoc(node *yaml.Node, name string, w *World) error {
	var head struct {
		Kind  string `yaml:"kind"`
		Start string `yaml:"start"`
	}
	if err := node.Decode(&head); err != nil {
		return fmt.Errorf("world: %s: reading kind: %w", name, err)
	}

	switch head.Kind {
	case "room":
		var r Room
		if err := node.Decode(&r); err != nil {
			return fmt.Errorf("world: %s: room: %w", name, err)
		}
		if r.ID == "" {
			return fmt.Errorf("world: %s: room is missing an id", name)
		}
		w.Rooms[r.ID] = &r
	case "item":
		var it Item
		if err := node.Decode(&it); err != nil {
			return fmt.Errorf("world: %s: item: %w", name, err)
		}
		if it.ID == "" {
			return fmt.Errorf("world: %s: item is missing an id", name)
		}
		w.Items[it.ID] = &it
	case "meta":
		w.Start = head.Start
	default:
		return fmt.Errorf("world: %s: unknown kind %q", name, head.Kind)
	}
	return nil
}

// validate checks referential integrity: a start room exists and every exit
// and room item points at something real.
func (w *World) validate() error {
	if w.Start == "" {
		return errors.New("world: no start room declared (add a `kind: meta` doc with `start:`)")
	}
	if w.Rooms[w.Start] == nil {
		return fmt.Errorf("world: start room %q does not exist", w.Start)
	}
	for id, r := range w.Rooms {
		for dir, ex := range r.Exits {
			if w.Rooms[ex.To] == nil {
				return fmt.Errorf("world: room %q exit %q points at unknown room %q", id, dir, ex.To)
			}
		}
		for _, itemID := range r.Items {
			if w.Items[itemID] == nil {
				return fmt.Errorf("world: room %q references unknown item %q", id, itemID)
			}
		}
	}
	return nil
}

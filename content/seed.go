package content

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

// pkgFS holds the `grimm` Python package (the Actor-focused teaching subset of
// grimm__python__zero) that we write into the student's work dir, so behavioral
// puzzle solutions can `from grimm import Actor`. Vendored snapshot; the source
// of truth is https://github.com/TheGrimmClub/grimm__python__zero.
//
//go:embed pkg/grimm/*.py
var pkgFS embed.FS

// SeedWorkspace writes the embedded `grimm` package into dir/grimm so that a
// Python solution run from dir (where `python3 loesung.py` puts dir on sys.path)
// can import it. It is idempotent — safe to call on every launch — and rewrites
// the files so a member always gets the current package.
func SeedWorkspace(dir string) error {
	return fs.WalkDir(pkgFS, "pkg/grimm", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Map pkg/grimm/... -> <dir>/grimm/...
		rel, err := filepath.Rel("pkg", p)
		if err != nil {
			return err
		}
		target := filepath.Join(dir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := pkgFS.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

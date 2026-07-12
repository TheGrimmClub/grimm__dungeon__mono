# Alchemist

Git as potion-brewing: a thin, German-flavoured wrapper over git that teaches
students version control as crafting. Each potion verb maps to real git.

```
init     git init                 erschaffe deinen Kessel (ein neues Repository)
add      git add                  lege eine Zutat in den Kessel (stage)
brew     git add -A + git commit  braue den Trank — halte den Stand fest
bottle   git push                 fülle den Trank ab und schicke ihn fort
discard  git restore              schütte die offenen Änderungen weg
clean    git clean                fege unversionierte Reste vom Tisch
look     git status               betrachte den Kessel — was brodelt gerade?
```

## Home & history

This Go module is the single source of truth for alchemist (decision D006 in
grimm__dungeon__mono, revised 2026-07-04). It is imported by the dungeon for the
in-game `/alchemist` command and built here as the standalone student binary.
It supersedes the original Python prototype (archived in Z_Archive
2026-07-04_beriah_premerge).

## Usage

- `task build` — binary into ./bin/alchemist
- `task test` / `task check`
- `task run -- look`

## Packages

- `.` — the alchemist library (verbs, dispatch)
- `runner/` — timeout-bounded external command execution (also used by the
  dungeon's puzzle checks)
- `cmd/alchemist/` — standalone CLI

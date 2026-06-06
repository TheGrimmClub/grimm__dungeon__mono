# grimm — TheGrimmClub Dungeon

The lesson tutor for **TheGrimmClub**: an interactive, command-line learning
dungeon that teaches 9th graders the *concepts* behind software engineering —
code reading over writing, how humans and machines think — wrapped in a Brothers
Grimm fairytale world laced with nanobots, androids and real AI.

`grimm` starts life as a Zork-like text adventure and grows, phase by phase, into
a Claude-Code-shaped shell (free-text input + `/commands`, a terminal, an editor,
and `/alchemist` — version control reimagined as potion brewing).

## Status

**Phase 0 — Foundation.** Runnable REPL with a German banner, `/help`, `/quit`
and a hidden `import antigravity` easter egg. The full roadmap lives in the
approved plan; see also [`docs/design`](docs/design).

## Quick start

Tasks are run with [Task](https://taskfile.dev) (`brew install go-task` or see
the [install docs](https://taskfile.dev/installation/)). Run `task` to list them.

```sh
task run          # build & run grimm
task check        # go vet + tests
task ci           # full gate: fmt check + vet + test + build
task build        # binaries into ./bin (grimm, alchemist)
```

Inside grimm: type `/help` for the scroll of commands, `/quit` to leave. Whisper
the old Python words to the dungeon for a surprise.

> Why Task? The `Taskfile.yaml` is also a lesson: a *ritual* (task) bundles many
> steps under one name. In the dungeon, automation is taught as exactly that —
> see [`docs/design/automation-task.md`](docs/design/automation-task.md).

## Layout

| Path | What |
|------|------|
| `cmd/grimm`, `cmd/alchemist` | thin binary entry points |
| `internal/tui` | Bubble Tea UI: scrollback, history, headlamp colour |
| `internal/session`, `internal/command` | pure dispatch (`/commands` + game verbs) |
| `internal/game/{world,entity,engine,state}` | the dungeon + verb engine |
| `internal/i18n` | German narrative text (commands/verbs stay English) |
| `internal/alchemist` | git-as-potion-brewing library (real wiring: Phase 3) |
| `content/` | multi-document YAML rooms/puzzles (Phase 1+) |
| `requirements.yaml`, `decisions.yaml` | living source of truth |

## Conventions

- **Source of truth:** `requirements.yaml` and `decisions.yaml` are updated first
  when scope or decisions change.
- **Git:** before each change, commit the open files, then branch.
- **Language:** narrative in German; commands, code and terminal in English.

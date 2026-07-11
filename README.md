# grimm — TheGrimmClub Dungeon

The lesson tutor for **TheGrimmClub**: an interactive, command-line learning
dungeon that teaches 9th graders the *concepts* behind software engineering —
code reading over writing, how humans and machines think — wrapped in a Brothers
Grimm fairytale world laced with nanobots, androids and real AI.

`grimm` starts life as a Zork-like text adventure and grows, phase by phase, into
a Claude-Code-shaped shell (free-text input + `/commands`, a terminal, an editor,
and `/alchemist` — version control reimagined as potion brewing).

## Status

**Through Phase 2 — Puzzles.** A Bubble Tea dungeon you walk in English verbs
(`look`, `go north`, `take 1`, `wear helm`, `solve echo`) while the world answers
in German. Wearing the headlamp floods the terminal with colour; doors are gated
by puzzles — a riddle (`answer`) and a "choose your path" `behavioral` check
(Python/Go/CLI), backed by `internal/puzzle` + `internal/runner`. Pick a class
with `/class`. Wearing the helmet also raises a **HUD** (inventory + map on the
right) and a **voice** that reads the German aloud (`/voice`; set
`GRIMM_VOICE=Anna` on macOS for a German voice). Version
your work with **`/alchemist`** (git as potion-brewing: `init`/`brew`/`bottle`/…)
in `~/.grimm/work` — one door (the Archiv) opens only once you've created your
own repository. Write code with **`/book <file>`** (your editor) and drop into a
real shell with **`/terminal`**, then `solve loesung.py` — any language. Save/
`/save`, history, hidden `import antigravity`. Roadmap: the approved plan and
[`docs/design`](docs/design).

## Quick start

Tasks are run with [Task](https://taskfile.dev) (`brew install go-task` or see
the [install docs](https://taskfile.dev/installation/)). Run `task` to list them.

```sh
task run          # build & run grimm
task check        # go vet + tests
task ci           # full gate: fmt check + vet + test + build
task build        # binary into ./bin (grimm)
```

Inside grimm: type `/help` for the scroll of commands, `/quit` to leave. Whisper
the old Python words to the dungeon for a surprise.

> Why Task? The `Taskfile.yaml` is also a lesson: a *ritual* (task) bundles many
> steps under one name. In the dungeon, automation is taught as exactly that —
> see [`docs/design/automation-task.md`](docs/design/automation-task.md).

## Layout

| Path | What |
|------|------|
| `cmd/grimm` | thin binary entry point |
| `internal/tui` | Bubble Tea UI: scrollback, history, headlamp colour |
| `internal/session`, `internal/command` | pure dispatch (`/commands` + game verbs) |
| `internal/game/{world,entity,engine,state}` | the dungeon + verb engine |
| `internal/i18n` | German narrative text (commands/verbs stay English) |
| (external) `grimm__toolbox__mono/tools/alchemist` | git-as-potion-brewing library + standalone CLI (imported via go.mod replace) |
| `content/` | multi-document YAML rooms/puzzles (Phase 1+) |
| `requirements.yaml`, `decisions.yaml` | living source of truth |

## Conventions

- **Source of truth:** `requirements.yaml` and `decisions.yaml` are updated first
  when scope or decisions change.
- **Git:** before each change, commit the open files, then branch.
- **Language:** narrative in German; commands, code and terminal in English.

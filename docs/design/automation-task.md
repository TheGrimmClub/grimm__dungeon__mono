# Automation as "Rituals" — teaching Task (taskfile.dev) in grimm

> Decision: [D012]. Requirement: [R013].

## Why Task

Students should meet one **real, cross-platform automation tool** rather than a
toy. [Task](https://taskfile.dev) is a single binary, its `Taskfile.yaml` is
plain YAML (no tab traps like Make), and a task is easy to read aloud. We use it
two ways at once:

1. **For real** — it is this monorepo's canonical task runner (`task build`,
   `task test`, `task ci`). CI runs the very same `task ci` developers run
   locally, so "works on my machine" and "passes CI" become one idea.
2. **In the fiction** — automation is taught as writing **rituals**: a *ritual*
   bundles many steps under one name so you never repeat them by hand. That is
   exactly a function/abstraction, met before any language syntax.

## The in-fiction mapping

| In the dungeon            | In reality                          |
|---------------------------|-------------------------------------|
| a **ritual**              | a `task` (named list of `cmds`)     |
| the **Ritualbuch**        | a `Taskfile.yaml`                   |
| "sprich »task brauen«"    | `task brauen`                       |
| writing a ritual once     | defining a task to avoid repetition |

## What exists now (Phase 1)

- The repo `Taskfile.yaml` is the canonical runner; the Makefile is gone.
- A **readable** in-world artifact: the `ritualbuch` item in the *Werkstatt*
  contains a tiny, real Taskfile and explains the idea. This is code *reading*
  (R001) — no engine work required, pure content.

## Planned: the Phase 2 automation puzzle

Once the puzzle engine and its three check kinds exist (D003), add a room whose
exit is gated by an **automation puzzle** that exercises both the artifact and
behavioral checks:

1. The room presents repetitive busywork (e.g. "greet three sleeping androids,
   each by name") that is tedious to do by hand.
2. The student authors a `Taskfile.yaml` in their work repo with a task that
   produces the required output.
3. Validation, reusing the Phase 2 machinery:
   - **artifact check** — a `Taskfile.yaml` exists and defines the expected task
     name;
   - **behavioral check** — running `task <name>` produces the expected stdout
     (language-agnostic: the check runs the *tool*, not a specific language).
4. Brewing & bottling the Taskfile with `alchemist` (D006) is what finally opens
   the door — tying automation, version control and the artifact check together.

### Open questions for Phase 2

- Do we **bundle** the `task` binary / guide its install as part of "Quest 0"
  (D002), alongside git/Python/Go?
- Should the behavioral check invoke `task` directly, or go through grimm's
  `/terminal` (D009) so the student sees the real command run?

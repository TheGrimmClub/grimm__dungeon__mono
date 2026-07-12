# Alchemist

A task-first git workflow, wrapped in a friendly CLI for students.

Alchemist nudges you to **decide what you're doing before you do it**: name the
task first, write your code, then commit with a message that's already half
written for you. The commands follow a light potion theme.

```
recipe   Scaffold a new project from a template
start    Name the task you're about to work on
brew     Commit your work with a detailed, editable message
bottle   Finalize: optional version tag and push
discard  Throw away uncommitted changes to tracked files
stash    Pause work and set it aside to resume later
clean    Remove untracked files (build artifacts, stray files)
```

## The workflow

```
recipe        # (once) scaffold a project and git init it
  |
start         # name the task -> this becomes your commit title
  |
... write code ...
  |
brew          # stage, summarize, edit the message, commit
  |
bottle        # optional version tag, then push
```

`discard`, `stash`, and `clean` are there for when things go sideways:

- **discard** — undo changes to tracked files (reverts to last commit)
- **stash** — pause and set work aside; `alchemist stash --resume` brings it back
- **clean** — delete untracked junk like build artifacts

## Install

Requires [Go](https://go.dev/dl/) 1.22+ and `git` on your PATH.

```
git clone https://github.com/TheGrimmClub/alchemist.git
cd alchemist
cd source
go mod tidy
go mod download github.com/spf13/cobra
go get github.com/spf13/cobra@v1.8.1
go mod download github.com/cpuguy83/go-md2man/v2
go install .
```

`go install` puts the `alchemist` binary in your Go bin directory. Make sure
that directory is on your PATH (on Windows it's usually `%USERPROFILE%\go\bin`).

If you use [Task](https://taskfile.dev), `task install` does the same thing.

## Editing commit messages

`brew` opens the draft message in an editor so you can refine it. It looks for
[micro](https://micro-editor.github.io/) first, then `$EDITOR`, then falls back
to `notepad` on Windows or `nano` elsewhere.

## How task state is stored

`start` saves the current task to `.alchemist/current.json` in your repo, and
`bottle` clears it. Scaffolded projects ignore that folder automatically; if you
add Alchemist to an existing project, add this line to your `.gitignore`:

```
.alchemist/
```

## For the Grimm Club

- [The Alchemist's Grimoire](./help/GRIMMOIRE/alchemist.grimmoire.md) — a themed reference for every spell,
  with the plain-git equivalent of each so you learn the real craft.
- [Lesson 1: Your First Brew](docs/1__first_brew.lesson.md) — a ~30 minute guided
  first session.

## License

MIT — see [LICENSE](LICENSE).

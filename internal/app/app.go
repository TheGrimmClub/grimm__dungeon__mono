// Package app wires grimm together: it loads the world, builds the colour-aware
// engine, restores any save, assembles the play session, and launches the
// Bubble Tea UI. cmd/grimm/main.go stays a thin entry point that calls Run.
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheGrimmClub/grimm__dungeon__mono/assets"
	"github.com/TheGrimmClub/grimm__dungeon__mono/content"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/engine"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/state"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/world"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/i18n"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/session"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/tui"
	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/voice"
)

// Run starts grimm with the default save location (~/.grimm/save.syon).
func Run() error {
	path, _ := state.DefaultPath() // empty on error -> persistence disabled
	return RunWith(path)
}

// RunWith starts grimm with an explicit save path ("" disables saving) and
// launches the Bubble Tea UI, blocking until the player quits.
func RunWith(savePath string) error {
	sess, intro, err := NewSession(savePath)
	if err != nil {
		return err
	}
	return tui.Run(sess, intro)
}

// NewSession builds a ready-to-play session and its intro text. It is separated
// from the UI so tests can exercise everything without a terminal.
func NewSession(savePath string) (*session.Session, string, error) {
	w, err := world.Load(content.FS, content.WorldGlob)
	if err != nil {
		return nil, "", fmt.Errorf("app: loading world: %w", err)
	}
	game := engine.New(w)
	game.SetPalette(engine.ColorPalette())

	continued := false
	if savePath != "" && state.Exists(savePath) {
		if snap, err := state.Load(savePath); err == nil {
			game.Restore(snap)
			continued = true
		}
	}

	if dir := workDir(savePath); dir != "" {
		game.SetWorkDir(dir) // artifact/behavioral checks see the student repo
	}

	sess := session.New(game, savePath)
	sess.SetVoice(voice.New()) // OS text-to-speech; Noop where unavailable
	if dir := workDir(savePath); dir != "" {
		sess.SetWorkDir(dir) // /alchemist brews here
	}
	return sess, intro(game, continued), nil
}

// workDir returns the student's working directory, next to the save file
// (~/.grimm/work), creating it if needed. Returns "" if it can't be created.
func workDir(savePath string) string {
	if savePath == "" {
		return ""
	}
	dir := filepath.Join(filepath.Dir(savePath), "work")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ""
	}
	// Seed the `grimm` package so behavioral solutions can `from grimm import
	// Actor`. Best-effort: a failure here shouldn't stop the student playing.
	_ = content.SeedWorkspace(dir)
	return dir
}

// intro composes the banner, welcome, optional "continued" note, the starting
// room and a one-time hint about how to talk to the dungeon.
func intro(game *engine.Game, continued bool) string {
	var b strings.Builder
	b.WriteString(assets.Banner)
	b.WriteString("  " + i18n.T(i18n.KeyBannerSubtitle) + "\n\n")
	b.WriteString(i18n.T(i18n.KeyWelcome, "/help", "/quit") + "\n")
	if continued {
		b.WriteString("\n" + i18n.T(i18n.KeyContinued) + "\n")
	}
	b.WriteString("\n" + game.Intro() + "\n")
	b.WriteString("\n" + i18n.T(i18n.KeyVerbHint))
	return b.String()
}

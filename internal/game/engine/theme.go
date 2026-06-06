package engine

import "github.com/charmbracelet/lipgloss"

// Palette styles the room view. The zero Palette renders everything plain, which
// keeps the engine tests deterministic; the app injects ColorPalette so the
// headlamp can flood the dungeon with colour (the first "aha", req R008).
type Palette struct {
	Title lipgloss.Style // room titles
	Exit  lipgloss.Style // direction tokens
	Item  lipgloss.Style // item names / [n] labels
	Dim   lipgloss.Style // whole-room wrap while unlit (the dark)
}

// PlainPalette renders everything unstyled. It is the engine default so logic
// and tests never depend on ANSI; the app swaps in ColorPalette.
func PlainPalette() Palette {
	s := lipgloss.NewStyle()
	return Palette{Title: s, Exit: s, Item: s, Dim: s}
}

// ColorPalette is the lit, in-colour look used by the TUI.
func ColorPalette() Palette {
	return Palette{
		Title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14")), // bright cyan
		Exit:  lipgloss.NewStyle().Foreground(lipgloss.Color("10")),            // green
		Item:  lipgloss.NewStyle().Foreground(lipgloss.Color("11")),            // yellow
		Dim:   lipgloss.NewStyle().Faint(true),                                 // faint grey
	}
}

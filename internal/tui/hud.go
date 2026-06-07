package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/game/engine"
)

const (
	sidebarWidth   = 28 // total columns of the HUD sidebar (incl. borders)
	minWidthForHUD = 64 // below this, the HUD is hidden and the view is full-width
)

var (
	hudBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("14")).
		Padding(0, 1).
		Width(sidebarWidth - 2) // minus the border columns

	hudTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13"))
	hudItem  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	hudExit  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	hudDim   = lipgloss.NewStyle().Faint(true)
)

// renderHUD builds the right-hand sidebar (inventory box over a map box) sized
// to the given inner height. It assumes the helmet is on (caller checks).
func renderHUD(g *engine.Game, height int) string {
	inv := renderInventory(g)
	karte := renderMap(g)

	sidebar := lipgloss.JoinVertical(lipgloss.Left, inv, karte)
	// Pad/trim to the available height so the horizontal join stays aligned.
	return lipgloss.NewStyle().Width(sidebarWidth).Height(height).Render(sidebar)
}

func renderInventory(g *engine.Game) string {
	var b strings.Builder
	b.WriteString(hudTitle.Render("INVENTAR"))
	b.WriteString("\n")
	items := g.InventoryView()
	if len(items) == 0 {
		b.WriteString(hudDim.Render("— leer —"))
	}
	for _, it := range items {
		line := fmt.Sprintf("[%s] %s", it.Slot, it.Name)
		if it.Worn {
			line += " " + hudDim.Render("(an)")
		}
		b.WriteString(hudItem.Render(truncate(line, sidebarWidth-4)))
		b.WriteString("\n")
	}
	return hudBox.Render(strings.TrimRight(b.String(), "\n"))
}

func renderMap(g *engine.Game) string {
	var b strings.Builder
	b.WriteString(hudTitle.Render("KARTE"))
	b.WriteString("\n")
	exits := g.ExitsView()
	if len(exits) == 0 {
		b.WriteString(hudDim.Render("— kein Ausgang —"))
	}
	for _, e := range exits {
		label := e.Dir
		if e.Locked {
			label += " *"
			b.WriteString(hudDim.Render("» " + label))
		} else {
			b.WriteString(hudExit.Render("» " + label))
		}
		b.WriteString("\n")
	}
	return hudBox.Render(strings.TrimRight(b.String(), "\n"))
}

// truncate clips a string to n runes (rough; HUD lines are short ASCII-ish).
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n || n <= 1 {
		return s
	}
	return string(r[:n-1]) + "…"
}

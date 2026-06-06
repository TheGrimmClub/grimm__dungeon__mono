// Package tui is grimm's Bubble Tea view/input layer: a scrollback viewport, a
// text input with up/down history, and a prompt that lights up in colour once
// the player wears the headlamp. All game logic lives in package session, which
// this layer merely drives.
package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/session"
)

// Run launches the full-screen Bubble Tea program. intro is the text shown
// before the first prompt (banner, welcome, the starting room).
func Run(sess *session.Session, intro string) error {
	p := tea.NewProgram(newModel(sess, intro), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type model struct {
	sess       *session.Session
	input      textinput.Model
	vp         viewport.Model
	hist       history
	transcript []string
	ready      bool
}

func newModel(sess *session.Session, intro string) model {
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = ""   // we render our own styled prompt label
	ti.CharLimit = 0 // no limit
	ti.Placeholder = "look · go north · inspect 1 · /help"

	m := model{sess: sess, input: ti}
	m.transcript = []string{intro}
	return m
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w, h := clampSize(msg.Width, msg.Height)
		if !m.ready {
			m.vp = viewport.New(w, h)
			m.ready = true
		} else {
			m.vp.Width, m.vp.Height = w, h
		}
		m.refresh()
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			return m.submit()
		case tea.KeyUp:
			if v, ok := m.hist.prev(); ok {
				m.input.SetValue(v)
				m.input.CursorEnd()
			}
			return m, nil
		case tea.KeyDown:
			if v, ok := m.hist.next(); ok {
				m.input.SetValue(v)
				m.input.CursorEnd()
			}
			return m, nil
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	m.vp, cmd = m.vp.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// submit handles Enter: echo the line, run it through the session, append the
// output, record history, and quit if asked.
func (m model) submit() (tea.Model, tea.Cmd) {
	line := strings.TrimSpace(m.input.Value())
	m.input.Reset()

	// Echo what was typed (the input box is now cleared).
	m.transcript = append(m.transcript, m.promptLabel()+line)
	res := m.sess.Submit(line)
	if res.Output != "" {
		m.transcript = append(m.transcript, res.Output)
	}
	m.hist.add(line)
	m.refresh()
	if res.Quit {
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "\n  initialisiere das Verlies…"
	}
	return m.vp.View() + "\n" + m.promptLabel() + m.input.View()
}

// refresh rebuilds the viewport content from the transcript and re-styles the
// prompt for the current light state, then scrolls to the bottom. It is safe to
// call before the first WindowSizeMsg (the viewport simply isn't drawn yet).
func (m *model) refresh() {
	m.input.PromptStyle = m.promptStyle()
	if !m.ready || m.vp.Height < 1 {
		return
	}
	m.vp.SetContent(strings.Join(m.transcript, "\n\n"))
	m.vp.GotoBottom()
}

// clampSize keeps the viewport dimensions positive (a pty can report 0x0).
func clampSize(width, height int) (w, h int) {
	w, h = width, height-2 // reserve a line for the prompt + spacing
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	return w, h
}

// promptLabel is the "Human> " prefix, styled for the current light state.
func (m model) promptLabel() string {
	return m.promptStyle().Render(m.sess.Game().Title() + "> ")
}

// promptStyle is bright magenta once the dungeon is lit, faint while dark.
func (m model) promptStyle() lipgloss.Style {
	if m.sess.Game().Lit() {
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13"))
	}
	return lipgloss.NewStyle().Faint(true)
}

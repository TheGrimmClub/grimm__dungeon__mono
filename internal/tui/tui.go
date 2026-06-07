// Package tui is grimm's Bubble Tea view/input layer: a scrollback viewport, a
// text input with up/down history, a prompt that lights up once the headlamp is
// worn, and a right-hand HUD (inventory + map) that appears with the helmet.
// All game logic lives in package session, which this layer merely drives.
package tui

import (
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/TheGrimmClub/grimm__dungeon__mono/internal/session"
)

// execDoneMsg is delivered after a suspended external process (/terminal,
// /book) returns control to the TUI.
type execDoneMsg struct {
	after string
	err   error
}

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
	width      int
	height     int
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
		m.width, m.height = msg.Width, msg.Height
		if !m.ready {
			m.vp = viewport.New(1, 1)
			m.ready = true
		}
		m.relayout()
		return m, nil

	case execDoneMsg:
		note := msg.after
		if msg.err != nil {
			note = "(" + msg.err.Error() + ")"
		}
		if note != "" {
			m.transcript = append(m.transcript, note)
		}
		m.relayout()
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
// output, record history, re-layout (the HUD may have just toggled), and quit
// if asked.
func (m model) submit() (tea.Model, tea.Cmd) {
	line := strings.TrimSpace(m.input.Value())
	m.input.Reset()

	m.transcript = append(m.transcript, m.promptLabel()+line)
	res := m.sess.Submit(line)
	if res.Output != "" {
		m.transcript = append(m.transcript, res.Output)
	}
	m.hist.add(line)
	m.relayout()

	// /terminal or /book: suspend the TUI, run the real shell/editor, resume.
	if res.Exec != nil {
		c := exec.Command(res.Exec.Name, res.Exec.Args...)
		c.Dir = res.Exec.Dir
		after := res.Exec.After
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			return execDoneMsg{after: after, err: err}
		})
	}
	if res.Quit {
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "\n  initialisiere das Verlies…"
	}
	_, vpHeight, showHUD := m.geometry()
	main := m.vp.View()
	if showHUD {
		main = lipgloss.JoinHorizontal(lipgloss.Top, main, " ", renderHUD(m.sess.Game(), vpHeight))
	}
	return main + "\n" + m.promptLabel() + m.input.View()
}

// geometry computes the viewport size and whether the HUD is shown, given the
// current terminal size and whether the helmet is on.
func (m model) geometry() (vpWidth, vpHeight int, showHUD bool) {
	vpHeight = m.height - 2 // reserve a line for the prompt + spacing
	if vpHeight < 1 {
		vpHeight = 1
	}
	vpWidth = m.width
	showHUD = m.sess.Game().HUDActive() && m.width >= minWidthForHUD
	if showHUD {
		vpWidth = m.width - sidebarWidth - 1 // sidebar + a spacer column
	}
	if vpWidth < 1 {
		vpWidth = 1
	}
	return vpWidth, vpHeight, showHUD
}

// relayout resizes the viewport for the current geometry and refreshes content.
func (m *model) relayout() {
	if !m.ready {
		return
	}
	w, h, _ := m.geometry()
	m.vp.Width, m.vp.Height = w, h
	m.refresh()
}

// refresh rebuilds the viewport content from the transcript and re-styles the
// prompt for the current light state, then scrolls to the bottom.
func (m *model) refresh() {
	m.input.PromptStyle = m.promptStyle()
	if !m.ready || m.vp.Height < 1 {
		return
	}
	m.vp.SetContent(strings.Join(m.transcript, "\n\n"))
	m.vp.GotoBottom()
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

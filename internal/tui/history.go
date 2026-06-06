package tui

// history is an up/down recall buffer for submitted input lines. It is a small,
// pure type so the arrow-key behaviour can be unit-tested without Bubble Tea.
type history struct {
	lines []string
	// idx is the cursor: 0..len(lines)-1 points at a past line; len(lines)
	// means "the fresh line the user is currently typing".
	idx int
}

// add records a submitted line (ignoring empties and exact repeats of the most
// recent entry) and resets the cursor to the fresh position.
func (h *history) add(line string) {
	if line != "" && (len(h.lines) == 0 || h.lines[len(h.lines)-1] != line) {
		h.lines = append(h.lines, line)
	}
	h.idx = len(h.lines)
}

// prev moves toward older entries (the Up key). It returns the recalled line and
// whether the cursor moved.
func (h *history) prev() (string, bool) {
	if h.idx == 0 {
		return "", false
	}
	h.idx--
	return h.lines[h.idx], true
}

// next moves toward newer entries (the Down key). Past the newest entry it
// returns the empty "fresh" line.
func (h *history) next() (string, bool) {
	if h.idx >= len(h.lines) {
		return "", false
	}
	h.idx++
	if h.idx == len(h.lines) {
		return "", true // back to the fresh line
	}
	return h.lines[h.idx], true
}

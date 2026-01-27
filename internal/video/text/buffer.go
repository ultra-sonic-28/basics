package text

// TextBuffer est un buffer texte générique basé sur une grille.
type TextBuffer struct {
	Cols int
	Rows int

	CursorX int
	CursorY int

	DefaultFG int
	DefaultBG int

	Cells []Cell
}

func NewTextBuffer(cols, rows int, fg, bg int) *TextBuffer {
	tb := &TextBuffer{
		Cols:      cols,
		Rows:      rows,
		DefaultFG: fg,
		DefaultBG: bg,
		Cells:     make([]Cell, cols*rows),
	}

	tb.Clear()
	return tb
}

// --- Effacement ---
func (t *TextBuffer) Clear() {
	for i := range t.Cells {
		t.Cells[i] = Cell{
			Glyph: ' ',
			FG:    t.DefaultFG,
			BG:    t.DefaultBG,
		}
	}
	t.CursorX = 0
	t.CursorY = 0
}

// Index calcule l'index linéaire dans le tableau Cells
func (t *TextBuffer) index(x, y int) int {
	return y*t.Cols + x
}

func (t *TextBuffer) SetCell(x, y int, glyph rune, fg, bg int) {
	if x < 0 || y < 0 || x >= t.Cols || y >= t.Rows {
		return
	}
	t.Cells[t.index(x, y)] = Cell{
		Glyph: glyph,
		FG:    fg,
		BG:    bg,
	}
}

// --- Scroll ---
func (t *TextBuffer) ScrollUp() {
	copy(
		t.Cells[0:],
		t.Cells[t.Cols:],
	)

	// dernière ligne vidée
	start := (t.Rows - 1) * t.Cols
	for i := 0; i < t.Cols; i++ {
		t.Cells[start+i] = Cell{
			Glyph: ' ',
			FG:    t.DefaultFG,
			BG:    t.DefaultBG,
		}
	}

	if t.CursorY > 0 {
		t.CursorY--
	}
}

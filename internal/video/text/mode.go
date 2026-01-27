package text

import (
	"basics/internal/video"
)

// TextMode implémente un mode texte générique (Apple II / CPC / Oric)
type TextMode struct {
	Buffer   *TextBuffer
	Renderer video.Renderer

	CellW int // largeur glyph (ex: 8)
	CellH int // hauteur glyph (ex: 8)

	FG int
	BG int
}

func NewTextMode(
	renderer video.Renderer,
	cols, rows int,
	cellW, cellH int,
	fg, bg int,
) *TextMode {

	tb := NewTextBuffer(cols, rows, fg, bg)

	return &TextMode{
		Buffer:   tb,
		Renderer: renderer,
		CellW:    cellW,
		CellH:    cellH,
		FG:       fg,
		BG:       bg,
	}
}

func (t *TextMode) CursorX() int {
	return t.Buffer.CursorX
}

func (t *TextMode) CursorY() int {
	return t.Buffer.CursorY
}

func (t *TextMode) SetCursor(x, y int) {
	//logger.Debug(fmt.Sprintf("x: %d, y: %d, cols: %d, rows: %d", x, y, t.Buffer.Cols, t.Buffer.Rows))
	if x >= 0 && x < t.Buffer.Cols {
		t.Buffer.CursorX = x
	}
	if y >= 0 && y < t.Buffer.Rows {
		t.Buffer.CursorY = y
	}
}

func (t *TextMode) Home() {
	t.Buffer.Clear()
}

func (t *TextMode) HTab(x int) {
	if x >= 0 && x < t.Buffer.Cols {
		t.Buffer.CursorX = x
	}
}

func (t *TextMode) VTab(y int) {
	if y >= 0 && y < t.Buffer.Rows {
		t.Buffer.CursorY = y
	}
}

func (t *TextMode) PutChar(r rune) {
	switch r {
	case '\n':
		t.NewLine()
	case '\r':
		t.Buffer.CursorX = 0
	default:
		t.putGlyph(r)
	}
}

func (t *TextMode) Print(s string) {
	for _, r := range s {
		t.PutChar(r)
	}
}

func (t *TextMode) putGlyph(r rune) {
	x := t.Buffer.CursorX
	y := t.Buffer.CursorY

	t.Buffer.SetCell(x, y, r, t.FG, t.BG)

	t.Buffer.CursorX++
	if t.Buffer.CursorX >= t.Buffer.Cols {
		t.NewLine()
	}
}

func (t *TextMode) NewLine() {
	t.Buffer.CursorX = 0
	t.Buffer.CursorY++

	if t.Buffer.CursorY >= t.Buffer.Rows {
		t.Buffer.ScrollUp()
		t.Buffer.CursorY = t.Buffer.Rows - 1
	}
}

func (t *TextMode) Backspace() {
	// début de ligne → rien à faire
	if t.Buffer.CursorX == 0 {
		return
	}

	// reculer le curseur
	t.Buffer.CursorX--

	// effacer le caractère (espace avec couleurs courantes)
	t.PutChar(' ')

	// remettre le curseur au bon endroit
	t.SetCursor(t.Buffer.CursorX-1, t.Buffer.CursorY)
}

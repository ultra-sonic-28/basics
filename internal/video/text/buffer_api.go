package text

// --- Dimensions ---

func (b *TextBuffer) Width() int {
	return b.Cols
}

func (b *TextBuffer) Height() int {
	return b.Rows
}

// --- Curseur ---

func (b *TextBuffer) Cursor() (int, int) {
	return b.CursorX, b.CursorY
}

func (b *TextBuffer) SetCursor(x, y int) {
	if x >= 0 && x < b.Cols {
		b.CursorX = x
	}
	if y >= 0 && y < b.Rows {
		b.CursorY = y
	}
}

// --- Cellules ---

func (b *TextBuffer) CellAt(x, y int) Cell {
	if x < 0 || x >= b.Cols || y < 0 || y >= b.Rows {
		return Cell{
			Glyph: ' ',
			FG:    b.DefaultFG,
			BG:    b.DefaultBG,
		}
	}
	return b.Cells[b.index(x, y)]
}

package font

// BitmapFont représente une police bitmap mono-plan.
type BitmapFont struct {
	Width  int
	Height int

	// Glyphs[index][row] = bits de la ligne
	Glyphs map[rune][]byte
}

func (f *BitmapFont) Glyph(r rune) []byte {
	if g, ok := f.Glyphs[r]; ok {
		return g
	}
	// fallback : caractère inconnu → espace
	return f.Glyphs[' ']
}

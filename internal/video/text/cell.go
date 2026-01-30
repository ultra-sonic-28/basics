package text

// Cell représente une cellule texte.
type Cell struct {
	Glyph rune
	FG    int
	BG    int
	Flash bool
}

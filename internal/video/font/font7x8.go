package font

// Font7x8 : vue Apple II (7x8) de la fonte 8x8 existante
// - même glyphes
// - MSB ignoré au rendu
var Font7x8 = &BitmapFont{
	Width:  7,
	Height: 8,
	Glyphs: Truncate8x8To7x8(Font8x8.Glyphs),
}

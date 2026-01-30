package font

import (
	"basics/testutils"
	"testing"
)

func TestBitmapFont_Glyph(t *testing.T) {
	// Création d'une police factice pour les tests
	f := &BitmapFont{
		Width:  8,
		Height: 8,
		Glyphs: map[rune][]byte{
			'A': {0x18, 0x24, 0x42, 0x7E, 0x42, 0x42, 0x42, 0x00},
			'B': {0x7C, 0x42, 0x42, 0x7C, 0x42, 0x42, 0x7C, 0x00},
			' ': {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // fallback
		},
	}

	tests := []struct {
		name   string
		input  rune
		expect []byte
	}{
		{"Known glyph A", 'A', f.Glyphs['A']},
		{"Known glyph B", 'B', f.Glyphs['B']},
		{"Unknown glyph X fallback to space", 'X', f.Glyphs[' ']},
		{"Space glyph itself", ' ', f.Glyphs[' ']},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.Glyph(tt.input)
			testutils.DeepEqual(t, "Glyph bytes", got, tt.expect)
		})
	}
}

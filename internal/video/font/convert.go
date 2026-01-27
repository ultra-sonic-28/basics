package font

func Truncate8x8To7x8(src map[rune][]byte) map[rune][]byte {
	dst := make(map[rune][]byte, len(src))

	for r, glyph := range src {
		if len(glyph) != 8 {
			continue
		}

		out := make([]byte, 8)
		for i := 0; i < 8; i++ {
			// supprime le MSB (bit 7)
			out[i] = glyph[i] & 0x7F
		}

		dst[r] = out
	}

	return dst
}

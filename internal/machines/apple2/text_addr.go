package apple2

const (
	textBase = 0x400
	textCols = 40
	textRows = 24
)

// textAddress calcule l'adresse mémoire Apple II
// correspondant à une position texte (col, row).
func textAddress(col, row int) int {
	if col < 0 || col >= textCols || row < 0 || row >= textRows {
		return textBase
	}

	return textBase +
		(row&0x07)*0x80 +
		(row>>3)*0x28 +
		col
}

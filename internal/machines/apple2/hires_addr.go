package apple2

const (
	hiresBase   = 0x2000
	hiresWidth  = 280
	hiresHeight = 192
)

// hiResAddress retourne l'adresse mémoire et le bit
// correspondant à un pixel (x, y) en hi-res Apple II.
func hiResAddress(x, y int) (addr int, bit int) {
	if x < 0 || x >= hiresWidth || y < 0 || y >= hiresHeight {
		return hiresBase, 0
	}

	byteX := x / 7
	bit = x % 7

	addr = hiresBase +
		(y&0x07)*0x400 +
		(y>>3)*0x80 +
		byteX

	return addr, bit
}

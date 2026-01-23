package video

// Renderer affiche un buffer déjà interprété par la machine.
type Renderer interface {
	RenderText(lines []string)
	RenderChar(x, y int, r rune)
}

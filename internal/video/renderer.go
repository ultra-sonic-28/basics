package video

// Renderer est l'unique abstraction de rendu graphique.
// Il ne contient aucune logique métier.
type Renderer interface {
	// Taille logique du framebuffer courant
	Width() int
	Height() int

	// Nettoyage complet
	Clear()

	// Dessin pixel (mode graphique)
	DrawPixel(x, y int, color int)

	// Dessin caractère (mode texte)
	DrawGlyph(
		x, y int, // position en pixels
		glyph rune, // caractère à dessiner
		fg, bg int, // couleurs (index palette)
	)
}

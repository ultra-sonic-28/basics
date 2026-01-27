package video

// ModeID identifie de manière unique un mode vidéo.
type ModeID string

// ModeInfo décrit un mode vidéo (texte ou graphique).
type ModeInfo struct {
	ID     ModeID // identifiant stable (ex: "apple2.text40")
	Name   string // nom humain
	Width  int    // largeur en pixels ou colonnes
	Height int    // hauteur en pixels ou lignes
	Text   bool   // true = mode texte
}

// Mode est l’interface commune à tous les modes vidéo.
type Mode interface {
	// Info retourne la description statique du mode.
	Info() ModeInfo

	// Reset remet le mode dans son état initial.
	Reset()
}

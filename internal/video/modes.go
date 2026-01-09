package video

// VideoModeID identifie un mode vidéo de manière abstraite.
type VideoModeID int

// ModeInfo décrit les caractéristiques d’un mode vidéo.
type ModeInfo struct {
	Name   string
	Width  int
	Height int
	Text   bool
}

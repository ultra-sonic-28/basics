package video

// Provider expose les capacités vidéo d’une machine.
type Provider interface {
	Modes() []ModeID
	ModeInfo(id ModeID) (ModeInfo, bool)
	DefaultMode() ModeID
}

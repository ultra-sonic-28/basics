package video

// Provider expose les capacités vidéo d’une machine.
type Provider interface {
	Modes() []VideoModeID
	ModeInfo(id VideoModeID) (ModeInfo, bool)
	DefaultMode() VideoModeID
}

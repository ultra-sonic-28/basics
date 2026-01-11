package video

// Device représente un périphérique vidéo logique.
type Device interface {
	//SetMode(VideoModeID) error

	PrintChar(r rune)
	PrintString(s string)

	Plot(x, y int)
	Clear()

	//Width() int
	//Height() int

	SetCursorX(x int)
	SetCursorY(y int)

	Render()
}

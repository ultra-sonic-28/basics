package video

import "io"

// Device représente un périphérique vidéo logique.
type Device interface {
	//SetMode(VideoModeID) error
	NeedsNewLineAfterInput() bool
	SetInput(io.Reader)
	SetOutput(io.Writer)

	PrintChar(r rune)
	PrintString(s string)

	Plot(x, y int)
	Clear()

	//Width() int
	//Height() int

	SetCursorX(x int)
	SetCursorY(y int)

	Render()
	RenderChar(x, y int, r rune)

	ReadLine() (string, error)
}

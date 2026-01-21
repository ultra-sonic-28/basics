package video

import "io"

// Device représente un périphérique vidéo logique.
type Device interface {
	//SetMode(VideoModeID) error
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
	ReadLine() (string, error)
}

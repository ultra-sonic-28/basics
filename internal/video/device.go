package video

import "io"

// Device représente l’interface vidéo vue par le runtime BASIC
type Device interface {
	// --- API BASIC ---
	Clear()
	PrintChar(r rune)
	PrintString(s string)

	SetCursorX(x int)
	SetCursorY(y int)

	Plot(x, y int)

	ReadLine() (string, error)
	GetChar() (rune, error)

	// --- I/O ---
	SetOutput(w io.Writer)
	DisableKeyboard()

	// --- Rendu ---
	Render()
}

package video

import "io"

// Device représente l’interface vidéo vue par le runtime BASIC
type Device interface {
	// --- API BASIC ---
	Clear()
	SetInverse(v bool)
	SetFlash(v bool)
	PrintChar(r rune)
	PrintString(s string)

	SetCursorX(x int)
	SetCursorY(y int)
	CursorX() int
	CursorY() int

	SwitchMode(slot int)

	Plot(x, y int)
	SetColor(c int)

	ReadLine() (string, error)
	GetChar() (rune, error)

	// --- I/O ---
	SetOutput(w io.Writer)
	DisableKeyboard()

	// --- Rendu ---
	Render()
}

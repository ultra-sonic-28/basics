package app

import (
	"errors"

	"basics/internal/machines/apple2"
	"basics/internal/video"

	"github.com/hajimehoshi/ebiten/v2"
)

// EbitenApp implémente ebiten.Game
type EbitenApp struct {
	*BasicEbitenApp
	started  bool
	prevKeys map[ebiten.Key]bool
}

// NewEbitenApp crée une application Ebiten
func NewEbitenApp(basic *BasicEbitenApp) *EbitenApp {
	return &EbitenApp{
		BasicEbitenApp: basic,
		prevKeys:       make(map[ebiten.Key]bool),
	}
}

// Run démarre Ebiten
func (a *EbitenApp) Run() error {

	// Vérification que le device supporte Ebiten
	if _, ok := a.Runtime.Video.(video.EbitenDevice); !ok {
		return errors.New("video device does not support Ebiten")
	}

	ebiten.SetWindowTitle("BASIC – Apple II")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return ebiten.RunGame(a)
}

// ==========================
// ebiten.Game interface
// ==========================

func (a *EbitenApp) Update() error {

	// Lancer l'interpréteur UNE SEULE FOIS
	if !a.started {
		a.started = true

		go a.Interpreter.Run(a.Program)
	}

	if t, ok := a.Runtime.Video.(*apple2.Text40); ok {
		t.Update()
	}

	a.handleInput()

	return nil
}

func (a *EbitenApp) handleInput() {
	t, ok := a.Runtime.Video.(*apple2.Text40)
	if !ok {
		return
	}

	// caractères imprimables
	/* for _, r := range ebiten.InputChars() {
		if r == '\r' || r == '\n' {
			continue
		}
		if r >= 32 && r <= 126 {
			t.InputRune(r)
		} else {
			continue
		}
	} */
	for _, r := range ebiten.InputChars() {
		if r < 32 || r > 126 {
			continue
		}

		// 🔴 MODE GET
		if t.IsGetActive() {
			t.PushGetRune(r)
			return // 👈 STOP : 1 touche suffit
		}

		// 🟢 MODE INPUT
		t.InputRune(r)
	}

	// BACKSPACE uniquement en INPUT
	if !t.IsGetActive() && a.keyJustPressed(ebiten.KeyBackspace) {
		t.Backspace()
	}

	// ENTER uniquement en INPUT
	if !t.IsGetActive() &&
		(a.keyJustPressed(ebiten.KeyEnter) ||
			a.keyJustPressed(ebiten.KeyNumpadEnter)) {
		t.Enter()
	}
}

func (a *EbitenApp) keyJustPressed(k ebiten.Key) bool {
	pressed := ebiten.IsKeyPressed(k)
	prev := a.prevKeys[k]
	a.prevKeys[k] = pressed
	return pressed && !prev
}

package app

import (
	"fmt"

	"basics/internal/video"

	"github.com/hajimehoshi/ebiten/v2"
)

// EbitenApp implémente ebiten.Game
type EbitenApp struct {
	*BasicEbitenApp
	started  bool
	prevKeys map[ebiten.Key]bool

	lastW, lastH int
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
		return fmt.Errorf("video device (%T) does not support Ebiten", a.Runtime.Video)
	}

	if dev, ok := a.Runtime.Video.(interface {
		Width() int
		Height() int
	}); ok {
		a.lastW, a.lastH = dev.Width(), dev.Height()
		ebiten.SetWindowSize(a.lastW, a.lastH)
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

	// Détection du changement de mode vidéo (pour resize fenêtre)
	if dev, ok := a.Runtime.Video.(interface {
		Width() int
		Height() int
	}); ok {
		w, h := dev.Width(), dev.Height()
		if w != a.lastW || h != a.lastH {
			a.lastW, a.lastH = w, h
			ebiten.SetWindowSize(w, h)
		}
	}

	if t, ok := a.Runtime.Video.(interface{ Update() error }); ok {
		t.Update()
	}

	a.handleInput()

	return nil
}

func (a *EbitenApp) handleInput() {
	type inputDevice interface {
		IsGetActive() bool
		PushGetRune(rune)
		InputRune(rune)
		Backspace()
		Enter()
	}

	t, ok := a.Runtime.Video.(inputDevice)
	if !ok {
		return
	}

	for _, r := range ebiten.InputChars() {
		// Autoriser [32–126] et [161–255]
		// De SPACE (32) à ~ (126) ou de ¡ (161) à ÿ (255)
		if !((r >= 32 && r <= 126) || (r >= 161 && r <= 255)) {
			continue
		}

		// MODE GET
		if t.IsGetActive() {
			t.PushGetRune(r)
			return // 👈 STOP : 1 touche suffit
		}

		// MODE INPUT
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

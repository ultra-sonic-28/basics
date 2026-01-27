package app

import (
	"basics/internal/interpreter"
	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/internal/video"

	"github.com/hajimehoshi/ebiten/v2"
)

// BasicEbitenApp contient l'état commun pour une app BASIC graphique
type BasicEbitenApp struct {
	Runtime     *runtime.Runtime
	Interpreter *interpreter.Interpreter
	Program     *parser.Program
}

// NewBasicEbitenApp crée une app graphique BASIC
func NewBasicEbitenApp(
	rt *runtime.Runtime,
	interp *interpreter.Interpreter,
	prog *parser.Program,
) *BasicEbitenApp {
	return &BasicEbitenApp{
		Runtime:     rt,
		Interpreter: interp,
		Program:     prog,
	}
}

func (a *BasicEbitenApp) Draw(screen *ebiten.Image) {
	if dev, ok := a.Runtime.Video.(video.EbitenDevice); ok {
		dev.Draw(screen)
	}
}

func (a *BasicEbitenApp) Layout(w, h int) (int, int) {
	if dev, ok := a.Runtime.Video.(video.EbitenDevice); ok {
		return dev.Layout(w, h)
	}
	return w, h
}

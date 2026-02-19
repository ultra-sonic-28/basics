package machines

import (
	"bufio"
	"os"

	"basics/internal/constants"
	"basics/internal/logger"
	"basics/internal/machines/apple2"
	"basics/internal/machines/tty"
	"basics/internal/runtime"
	ebitenrenderer "basics/internal/video/ebiten"
	"basics/internal/video/font"
)

func NewRuntime(basicType byte, mode bool) (*runtime.Runtime, error) {
	var x, y int
	var video *apple2.Text40
	scale := 2

	switch basicType {

	case constants.BASIC_APPLE:
		if mode {
			// --- Apple II Text 80 ---
			x = 560
			y = 192
		} else {
			// --- Apple II Text 40 ---
			x = 280 // résolution Apple II en mode HGR2
			y = 192
		}

		renderer := ebitenrenderer.New(
			scale*x, scale*y,
			scale, // scale
			apple2.Palette(),
			font.DefaultFontForMode(basicType),
		)

		if mode {
			video = apple2.NewText80(renderer)
		} else {
			video = apple2.NewText40(renderer)
		}
		logger.Info("Instanciate Ebiten renderer")

		return runtime.New(video), nil

	case constants.BASIC_TTY:
		in := bufio.NewReader(os.Stdin)
		out := os.Stdout
		video := tty.New(in, out)
		logger.Info("Instanciate TTY renderer")

		return runtime.New(video), nil

	default:
		return nil, ErrUnsupportedMachine
	}
}

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

func NewRuntime(basicType byte) (*runtime.Runtime, error) {

	switch basicType {

	case constants.BASIC_APPLE:
		// --- Apple II Text 40 ---
		renderer := ebitenrenderer.New(
			280, 192, // résolution Apple II en mode HGR2
			2, // scale
			apple2.Palette(),
			font.DefaultFontForMode(basicType),
		)

		video := apple2.NewText40(renderer)
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

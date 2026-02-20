package machines

import (
	"bufio"
	"fmt"
	"os"

	"basics/internal/constants"
	"basics/internal/logger"
	"basics/internal/machines/apple2"
	"basics/internal/machines/tty"
	"basics/internal/runtime"
	"basics/internal/video"
)

func NewRuntime(basicType byte, mode bool) (*runtime.Runtime, error) {
	var defaultMode video.ModeID

	switch basicType {

	case constants.BASIC_APPLE:

		defaultMode = apple2.ModeAppleText40

		display := apple2.NewAppleDisplay(2, defaultMode)
		mode, _ := display.ModeInfo(display.DefaultMode())
		logger.Info(fmt.Sprintf("Instanciate Ebiten renderer using %s mode (%s)", mode.Name, defaultMode))

		return runtime.New(display), nil

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

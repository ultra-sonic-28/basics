package machines

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"basics/internal/constants"
	"basics/internal/machines/apple2"
	"basics/internal/machines/tty"
	"basics/internal/runtime"
)

func NewRuntime(basicType byte) (*runtime.Runtime, error) {

	switch basicType {

	case constants.BASIC_APPLE:
		provider := apple2.NewVideoProvider()
		renderer := apple2.NewConsoleRenderer()
		vram := apple2.NewVRAM()
		video := apple2.NewVideoDevice(provider, renderer, vram)
		return runtime.New(video), nil

	case constants.BASIC_TTY:
		in := bufio.NewReader(os.Stdin)
		out := &bytes.Buffer{}
		video := tty.New(in, out)
		return runtime.New(video), nil

	default:
		return nil, fmt.Errorf("unsupported BASIC type")
	}
}

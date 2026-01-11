package runtime

import (
	"basics/internal/video"
)

type Runtime struct {
	Video video.Device
	Env   *Environment
}

func New(video video.Device) *Runtime {
	return &Runtime{
		Video: video,
		Env:   NewEnvironment(),
	}
}

func (rt *Runtime) ExecPrint(value string) {
	rt.Video.PrintString(value)
	rt.Video.Render()
}

func (rt *Runtime) ExecPlot(x, y int) {
	rt.Video.Plot(x, y)
	rt.Video.Render()
}

func (rt *Runtime) ExecHTab(x int) {
	rt.Video.SetCursorX(x - 1) // BASIC = 1-based
}

func (rt *Runtime) ExecVTab(y int) {
	rt.Video.SetCursorY(y - 1)
}

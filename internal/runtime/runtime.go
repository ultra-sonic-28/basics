package runtime

import (
	"basics/internal/video"
	"io"
)

type Runtime struct {
	Video  video.Device
	Env    *Environment
	halted bool
}

func New(video video.Device) *Runtime {
	return &Runtime{
		Video: video,
		Env:   NewEnvironment(),
	}
}

func (r *Runtime) SetInput(in io.Reader) {
	r.Video.SetInput(in)
}

func (r *Runtime) SetOutput(out io.Writer) {
	r.Video.SetOutput(out)
}

// runtime/runtime.go
func (rt *Runtime) ExecError(err error) {
	rt.Video.PrintString(err.Error())
	rt.Video.PrintString("\n")
	rt.Video.Render()
}

func (r *Runtime) Halt() {
	r.halted = true
}

func (r *Runtime) IsHalted() bool {
	return r.halted
}

func (rt *Runtime) ExecInput() (string, error) {
	return rt.Video.ReadLine()
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

func (rt *Runtime) ExecHome() {
	rt.Video.Clear()
	rt.Video.SetCursorX(0)
	rt.Video.SetCursorY(0)
}

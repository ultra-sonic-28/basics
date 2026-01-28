package runtime

import (
	"basics/internal/input"
	"basics/internal/video"
	"io"
)

type InputDevice interface {
	GetChar() (rune, error)
}

type Runtime struct {
	Video  video.Device
	Input  input.Device
	Env    *Environment
	halted bool
}

func New(video video.Device) *Runtime {
	return &Runtime{
		Video: video,
		Env:   NewEnvironment(),
	}
}

func (r *Runtime) SetOutput(out io.Writer) {
	r.Video.SetOutput(out)
}

func (rt *Runtime) ExecError(err error) {
	rt.Video.PrintString(err.Error())
	rt.Video.PrintString("\n")
	rt.Video.Render()
}

func (rt *Runtime) Halt() {
	rt.halted = true
}

func (rt *Runtime) IsHalted() bool {
	return rt.halted
}

func (rt *Runtime) ExecInput() (string, error) {
	if rt.Input != nil {
		return rt.Input.ReadLine()
	}

	return rt.Video.ReadLine()
}

func (rt *Runtime) ExecGet() (rune, error) {
	if rt.Input != nil {
		return rt.Input.GetChar()
	}
	return rt.Video.GetChar()
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

func (rt *Runtime) DisableKeyboard() {
	rt.Video.DisableKeyboard()
}

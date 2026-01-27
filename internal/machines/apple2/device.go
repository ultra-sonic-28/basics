package apple2

import (
	"basics/internal/video/text"
	"io"
)

type Device struct {
	Text *text.TextMode

	in  io.Reader
	out io.Writer
}

func NewDevice(textMode *text.TextMode) *Device {
	return &Device{
		Text: textMode,
	}
}

func (d *Device) Clear() {
	d.Text.Home()
}

func (d *Device) PrintString(s string) {
	d.Text.Print(s)
}

func (d *Device) SetCursorX(x int) {
	d.Text.HTab(x)
}

func (d *Device) SetCursorY(y int) {
	d.Text.VTab(y)
}

func (d *Device) Plot(x, y int) {
	// Apple II text mode: no-op
}

func (d *Device) ReadLine() (string, error) {
	return "", nil
}

func (d *Device) SetInput(r io.Reader) {
	d.in = r
}

func (d *Device) SetOutput(w io.Writer) {
	d.out = w
}

func (d *Device) Render() {
	d.Text.Render()
}

package apple2

import (
	"basics/internal/video"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type VideoDevice struct {
	mode     video.VideoModeID
	vram     video.VRAM
	provider video.Provider
	renderer video.Renderer
	cursorX  int
	cursorY  int
	in       *bufio.Reader
	out      io.Writer
}

func NewVideoDevice(
	provider video.Provider,
	renderer video.Renderer,
	vram video.VRAM,
) *VideoDevice {
	return &VideoDevice{
		mode:     provider.DefaultMode(),
		provider: provider,
		renderer: renderer,
		vram:     vram,
		cursorX:  0,
		cursorY:  0,
		in:       bufio.NewReader(os.Stdin),
		out:      &bytes.Buffer{},
	}
}

func (t *VideoDevice) SetInput(r io.Reader) {
	t.in = bufio.NewReader(r)
}

func (t *VideoDevice) SetOutput(w io.Writer) {
	t.out = w
}

func (v *VideoDevice) VRAM() video.VRAM {
	return v.vram
}

func (v *VideoDevice) PrintString(s string) {
	for _, r := range s {
		if r == '\n' {
			v.cursorX = 0
			v.cursorY++
			continue
		}
		v.PrintChar(r)
	}
}

func (v *VideoDevice) PrintChar(r rune) {
	if v.mode != Text40 {
		return
	}

	addr := textAddress(v.cursorX, v.cursorY)
	v.vram.Write(addr, byte(r))

	v.cursorX++
	if v.cursorX >= 40 {
		v.cursorX = 0
		v.cursorY++
	}
}

func (v *VideoDevice) Plot(x, y int) {
	if v.mode != HiRes {
		return
	}

	addr, bit := hiResAddress(x, y)
	value := v.vram.Read(addr)
	value |= (1 << bit)
	v.vram.Write(addr, value)
}

func (v *VideoDevice) SetCursorX(x int) {
	if x < 0 {
		x = 0
	}
	if x >= 40 {
		x = 39
	}
	v.cursorX = x
}

func (v *VideoDevice) SetCursorY(y int) {
	if y < 0 {
		y = 0
	}
	if y >= 24 {
		y = 23
	}
	v.cursorY = y
}

func (v *VideoDevice) MoveTerminalCursor() {
	// +1 car ANSI est 1-based
	fmt.Printf("\033[%d;%dH", v.cursorY+1, v.cursorX+1)
}

func (v *VideoDevice) Clear() {
	v.vram.Clear()
	v.cursorX = 0
	v.cursorY = 0
}

func (v *VideoDevice) Render() {
	if v.mode != Text40 {
		return
	}

	lines := make([]string, 24)

	for y := 0; y < 24; y++ {
		row := make([]rune, 40)
		for x := 0; x < 40; x++ {
			addr := textAddress(x, y)
			ch := v.vram.Read(addr)
			if ch == 0 {
				ch = ' '
			}
			row[x] = rune(ch)
		}
		lines[y] = string(row)
	}

	v.renderer.RenderText(lines)
}

func (v *VideoDevice) ReadLine() (string, error) {
	// synchroniser le curseur terminal avec le curseur Apple II
	v.MoveTerminalCursor()

	line, err := v.in.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")

	// recopier la saisie dans la VRAM
	for _, r := range line {
		v.PrintChar(r)
	}
	v.cursorX = 0
	v.cursorY++

	v.Render()
	return line, err
}

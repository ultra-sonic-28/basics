package apple2

import "basics/internal/video"

type VideoDevice struct {
	mode     video.VideoModeID
	vram     video.VRAM
	provider video.Provider
	renderer video.Renderer
	cursorX  int
	cursorY  int
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
	}
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

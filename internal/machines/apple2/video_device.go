package apple2

import (
	"basics/internal/logger"
	"basics/internal/video"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
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

	cursorVisible bool
	cursorChar    byte
	cursorMu      sync.Mutex
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

func (v *VideoDevice) NeedsNewLineAfterInput() bool {
	return true
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
			v.cursorMu.Lock()
			v.cursorX = 0
			v.cursorY++
			if v.cursorY >= 24 {
				v.cursorY = 23
			}
			v.MoveTerminalCursor()
			v.cursorMu.Unlock()
			logger.Debug(fmt.Sprintf("PrintChar '%c' at %d,%d", byte(r), v.cursorX, v.cursorY))
			continue
		}
		v.PrintChar(r)
	}
}

func (v *VideoDevice) PrintChar(r rune) {
	if v.mode != Text40 {
		return
	}

	v.cursorMu.Lock()
	defer v.cursorMu.Unlock()

	if v.cursorY >= 24 {
		return
	}

	logger.Debug(fmt.Sprintf("PrintChar '%c' at %d,%d", byte(r), v.cursorX, v.cursorY))
	addr := textAddress(v.cursorX, v.cursorY)
	v.vram.Write(addr, byte(r))

	v.cursorX++
	if v.cursorX >= 40 {
		v.cursorX = 0
		v.cursorY++
		if v.cursorY >= 24 {
			v.cursorY = 23
		}
	}

	v.MoveTerminalCursor()
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
	v.cursorMu.Lock()
	defer v.cursorMu.Unlock()

	if x < 0 {
		x = 0
	} else if x >= 40 {
		x = 39
	}

	v.cursorX = x
}

func (v *VideoDevice) SetCursorY(y int) {
	v.cursorMu.Lock()
	defer v.cursorMu.Unlock()

	if y < 0 {
		y = 0
	} else if y >= 24 {
		y = 23
	}

	v.cursorY = y
}

func (v *VideoDevice) MoveTerminalCursor() {
	// +1 car ANSI est 1-based
	logger.Debug(fmt.Sprintf("MoveTerminalCursor at %d,%d", v.cursorX, v.cursorY))
	fmt.Printf("\033[%d;%dH", v.cursorX+1, v.cursorY+1)
}

func (v *VideoDevice) moveCursorBack() {
	v.cursorMu.Lock()
	defer v.cursorMu.Unlock()

	if v.cursorX > 0 {
		v.cursorX--
		return
	}

	// Retour en fin de ligne précédente
	if v.cursorY > 0 {
		v.cursorY--
		v.cursorX = 39
	}
}

func (v *VideoDevice) Clear() {
	v.cursorMu.Lock()
	defer v.cursorMu.Unlock()

	v.vram.Clear()
	v.cursorX = 0
	v.cursorY = 0
	v.cursorVisible = false
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

func (v *VideoDevice) RenderChar(x, y int, r rune) {
	v.renderer.RenderChar(x, y, r)
}

func (v *VideoDevice) ReadLine() (string, error) {
	var buf []rune
	started := false

	stop := make(chan struct{})
	go v.blinkCursor(stop)

	for {
		ch, _, err := v.in.ReadRune()
		if err != nil {
			close(stop)
			return "", err
		}

		// IGNORER les \n AVANT la première saisie
		if !started {
			if ch == '\n' || ch == '\r' {
				continue
			}
			started = true
		}

		switch ch {
		case '\n', '\r':
			close(stop)
			return string(buf), nil

		case 0x08, 0x7F: // BACKSPACE
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				v.moveCursorBack()
				v.PrintChar(' ')
				v.moveCursorBack()
				v.RenderChar(v.cursorX, v.cursorX, ' ')
			}

		default:
			v.eraseCursor()

			buf = append(buf, ch)

			//x, y := v.cursorX, v.cursorY
			v.PrintChar(ch)
			//v.RenderChar(x, y, ch)
			v.renderer.RenderChar(v.cursorX-1, v.cursorY, ch)

			v.drawCursor()
		}
	}
}

func (v *VideoDevice) drawCursor() {
	v.cursorMu.Lock()
	defer v.cursorMu.Unlock()

	logger.Debug(fmt.Sprintf("drawCursor at %d,%d", v.cursorX, v.cursorY))
	addr := textAddress(v.cursorX, v.cursorY)
	v.cursorChar = v.vram.Read(addr)
	v.vram.Write(addr, '_')
	v.RenderChar(v.cursorX, v.cursorY, '_')
}

func (v *VideoDevice) eraseCursor() {
	v.cursorMu.Lock()
	defer v.cursorMu.Unlock()

	logger.Debug(fmt.Sprintf("eraseCursor at %d,%d", v.cursorX, v.cursorY))
	addr := textAddress(v.cursorX, v.cursorY)
	v.vram.Write(addr, v.cursorChar)
	v.RenderChar(v.cursorX, v.cursorY, rune(v.cursorChar))
}

func (v *VideoDevice) blinkCursor(stop <-chan struct{}) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			v.eraseCursor()
			return

		case <-ticker.C:
			if v.cursorVisible {
				v.eraseCursor()
			} else {
				v.drawCursor()
			}
			v.cursorVisible = !v.cursorVisible
		}
	}
}

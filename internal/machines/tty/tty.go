package tty

import (
	"basics/internal/video"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type TTYDevice struct {
	buffer  []rune
	in      *bufio.Reader
	out     io.Writer
	inverse bool
	flash   bool
	cursorX int
}

func New(in io.Reader, out io.Writer) video.Device {
	return &TTYDevice{
		in:      bufio.NewReader(in),
		out:     out,
		inverse: false,
		flash:   false,
		cursorX: 0,
	}
}

func (t *TTYDevice) SetInput(r io.Reader) {
	t.in = bufio.NewReader(r)
}

func (t *TTYDevice) SetOutput(w io.Writer) {
	t.out = w
}

func (t *TTYDevice) PrintString(s string) {
	for _, r := range s {
		t.PrintChar(r)
	}
}

func (t *TTYDevice) PrintChar(r rune) {
	t.buffer = append(t.buffer, r)
	if r == '\n' {
		t.cursorX = 0
	} else {
		t.cursorX++
	}
}

func (t *TTYDevice) Plot(x, y int) {}

func (t *TTYDevice) SetCursorX(x int) {
	if x > t.cursorX {
		padding := strings.Repeat(" ", x-t.cursorX)
		t.buffer = append(t.buffer, []rune(padding)...)
	}
	t.cursorX = x
}

func (t *TTYDevice) SetCursorY(y int) {}

func (t *TTYDevice) CursorX() int { return t.cursorX }
func (t *TTYDevice) CursorY() int { return 0 }

func (t *TTYDevice) SwitchMode(slot int) {}

func (t *TTYDevice) Clear() {
	t.buffer = nil
	fmt.Print("\033[2J\033[H")
}

func (t *TTYDevice) Render() {
	if t.flash {
		fmt.Fprintf(t.out, "\x1b[5m%s", string(t.buffer)) // ANSI blink
		t.buffer = nil
		return
	}

	if t.inverse {
		fmt.Fprintf(t.out, "\x1b[0m\033[7m%s", string(t.buffer)) // inverse video
	} else {
		fmt.Fprintf(t.out, "\x1b[0m\033[27m%s", string(t.buffer)) // normal
	}

	t.buffer = nil
}

func (t *TTYDevice) ReadLine() (string, error) {
	line, err := t.in.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimRight(line, "\r\n"), err
}

func (t *TTYDevice) GetChar() (rune, error) {
	reader := bufio.NewReader(os.Stdin)
	r, _, err := reader.ReadRune()
	return r, err
}

func (t *TTYDevice) DisableKeyboard() {}

func (t *TTYDevice) SetInverse(v bool) {
	t.inverse = v
}

func (t *TTYDevice) SetFlash(v bool) {
	t.flash = v
}

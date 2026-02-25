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
	cursorY int
}

func New(in io.Reader, out io.Writer) video.Device {
	return &TTYDevice{
		in:      bufio.NewReader(in),
		out:     out,
		inverse: false,
		flash:   false,
		cursorX: 0,
		cursorY: 0,
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
		t.cursorY++
	} else {
		t.cursorX++
	}
}

func (t *TTYDevice) Plot(x, y int) {}

func (t *TTYDevice) SetCursorX(x int) {
	// Utilisation de l'escape code ANSI pour le positionnement horizontal absolu (1-based)
	t.buffer = append(t.buffer, []rune(fmt.Sprintf("\033[%dG", x+1))...)
	t.cursorX = x
}

func (t *TTYDevice) SetCursorY(y int) {
	// CUP (Cursor Position) positionne à la fois la ligne et la colonne (1-based)
	// On conserve la colonne actuelle pour ne changer que la ligne (comportement VTAB)
	t.buffer = append(t.buffer, []rune(fmt.Sprintf("\033[%d;%dH", y+1, t.cursorX+1))...)
	t.cursorY = y
}

func (t *TTYDevice) CursorX() int { return t.cursorX }
func (t *TTYDevice) CursorY() int { return t.cursorY }

func (t *TTYDevice) SwitchMode(slot int) {}

func (t *TTYDevice) Clear() {
	t.buffer = nil
	fmt.Fprintf(t.out, "\033[2J\033[H")
	t.cursorX = 0
	t.cursorY = 0
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

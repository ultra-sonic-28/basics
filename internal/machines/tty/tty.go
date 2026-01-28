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
	buffer []rune
	in     *bufio.Reader
	out    io.Writer
}

func New(in io.Reader, out io.Writer) video.Device {
	return &TTYDevice{
		in:  bufio.NewReader(in),
		out: out,
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
}

func (t *TTYDevice) Plot(x, y int) {}

func (t *TTYDevice) SetCursorX(x int) {}

func (t *TTYDevice) SetCursorY(y int) {}

func (t *TTYDevice) Clear() {
	t.buffer = nil
	fmt.Print("\033[2J\033[H")
}

func (t *TTYDevice) Render() {
	fmt.Fprint(t.out, string(t.buffer))
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

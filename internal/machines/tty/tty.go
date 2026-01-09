package tty

import (
	"basics/internal/video"
	"fmt"
)

type TTYDevice struct {
	buffer []rune
}

func New() video.Device {
	return &TTYDevice{}
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

func (t *TTYDevice) Clear() {
	t.buffer = nil
	fmt.Print("\033[2J\033[H")
}

func (t *TTYDevice) Render() {
	fmt.Print(string(t.buffer))
	t.buffer = nil
}

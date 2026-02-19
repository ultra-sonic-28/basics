package apple2

import (
	"basics/internal/video"
	"basics/internal/video/text"
	"bufio"
	"io"
	"strings"
)

func NewText80(renderer video.Renderer) *Text40 {
	mode := text.NewTextMode(
		renderer,
		80, 24, // ← 80 colonnes
		7, 8,
		1, 0,
	)

	return &Text40{
		Mode:        mode,
		renderer:    renderer,
		in:          bufio.NewReader(strings.NewReader("")),
		out:         io.Discard,
		inputBuffer: make([]rune, 0, 64),
		lineReady:   false,
		allowInput:  false,
		inverse:     false,
		flash:       false,
	}
}

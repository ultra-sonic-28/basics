package renderer

import "basics/internal/video/text"

type Renderer interface {
	Clear()
	DrawTextBuffer(buf text.TextBuffer)
	Present()
}

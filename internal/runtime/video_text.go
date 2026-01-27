package runtime

import "basics/internal/video/text"

// TextVideo est une implémentation runtime basée sur un TextMode
type TextVideo struct {
	Text *text.TextMode
}

func (v *TextVideo) Home() {
	v.Text.Home()
}

func (v *TextVideo) Print(s string) {
	v.Text.Print(s)
}

func (v *TextVideo) HTab(x int) {
	v.Text.HTab(x)
}

func (v *TextVideo) VTab(y int) {
	v.Text.VTab(y)
}

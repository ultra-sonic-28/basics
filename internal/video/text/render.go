package text

import ebitenrenderer "basics/internal/video/ebiten"

func (t *TextMode) Render() {
	for y := 0; y < t.Buffer.Rows; y++ {
		for x := 0; x < t.Buffer.Cols; x++ {
			cell := t.Buffer.Cells[y*t.Buffer.Cols+x]

			px := x * t.CellW
			py := y * t.CellH

			fg := cell.FG
			bg := cell.BG

			if cell.Flash && (t.Renderer.(*ebitenrenderer.Renderer).FrameCounter/30)%2 == 0 {
				fg, bg = bg, fg
			}
			t.Renderer.DrawGlyph(
				px,
				py,
				cell.Glyph,
				fg,
				bg,
			)
		}
	}
}

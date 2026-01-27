package text

func (t *TextMode) Render() {
	for y := 0; y < t.Buffer.Rows; y++ {
		for x := 0; x < t.Buffer.Cols; x++ {
			cell := t.Buffer.Cells[y*t.Buffer.Cols+x]

			px := x * t.CellW
			py := y * t.CellH

			t.Renderer.DrawGlyph(
				px,
				py,
				cell.Glyph,
				cell.FG,
				cell.BG,
			)
		}
	}
}

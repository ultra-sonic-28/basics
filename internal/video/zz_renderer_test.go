package video

import (
	"testing"

	"basics/testutils"
)

// MockRenderer implémente l'interface Renderer pour tests
type MockRenderer struct {
	cleared     bool
	pixelsDrawn []struct {
		x, y, color int
	}
	glyphsDrawn []struct {
		x, y   int
		glyph  rune
		fg, bg int
	}
	width, height int
}

func (m *MockRenderer) Width() int  { return m.width }
func (m *MockRenderer) Height() int { return m.height }
func (m *MockRenderer) Clear()      { m.cleared = true }
func (m *MockRenderer) DrawPixel(x, y int, color int) {
	m.pixelsDrawn = append(m.pixelsDrawn, struct{ x, y, color int }{x, y, color})
}
func (m *MockRenderer) DrawGlyph(x, y int, glyph rune, fg, bg int) {
	m.glyphsDrawn = append(m.glyphsDrawn, struct {
		x, y   int
		glyph  rune
		fg, bg int
	}{x, y, glyph, fg, bg})
}

func TestMockRenderer(t *testing.T) {
	mock := &MockRenderer{width: 80, height: 25}

	// Test Width / Height
	testutils.Equal(t, "Width", mock.Width(), 80)
	testutils.Equal(t, "Height", mock.Height(), 25)

	// Test Clear
	mock.Clear()
	testutils.True(t, "Clear called", mock.cleared)

	// Test DrawPixel
	mock.DrawPixel(1, 2, 42)
	testutils.Equal(t, "DrawPixel count", len(mock.pixelsDrawn), 1)
	testutils.Equal(t, "DrawPixel x", mock.pixelsDrawn[0].x, 1)
	testutils.Equal(t, "DrawPixel y", mock.pixelsDrawn[0].y, 2)
	testutils.Equal(t, "DrawPixel color", mock.pixelsDrawn[0].color, 42)

	// Test DrawGlyph
	mock.DrawGlyph(5, 6, 'A', 1, 0)
	testutils.Equal(t, "DrawGlyph count", len(mock.glyphsDrawn), 1)
	testutils.Equal(t, "DrawGlyph x", mock.glyphsDrawn[0].x, 5)
	testutils.Equal(t, "DrawGlyph y", mock.glyphsDrawn[0].y, 6)
	testutils.Equal(t, "DrawGlyph rune", mock.glyphsDrawn[0].glyph, 'A')
	testutils.Equal(t, "DrawGlyph fg", mock.glyphsDrawn[0].fg, 1)
	testutils.Equal(t, "DrawGlyph bg", mock.glyphsDrawn[0].bg, 0)
}

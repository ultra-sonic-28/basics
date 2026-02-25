package video

import (
	"basics/testutils"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// MockEbitenDevice implémente EbitenDevice pour les tests
type MockEbitenDevice struct {
	DrawCalled                bool
	LayoutCalled              bool
	LayoutW, LayoutH          int
	DeviceWidth, DeviceHeight int
}

func (m *MockEbitenDevice) Draw(screen *ebiten.Image) {
	m.DrawCalled = true
}

func (m *MockEbitenDevice) Layout(w, h int) (int, int) {
	m.LayoutCalled = true
	m.LayoutW = w
	m.LayoutH = h
	return m.DeviceWidth, m.DeviceHeight
}

// Implémentation minimale de Device (interface parent) pour le mock
func (m *MockEbitenDevice) Width() int                                 { return m.DeviceWidth }
func (m *MockEbitenDevice) Height() int                                { return m.DeviceHeight }
func (m *MockEbitenDevice) CursorX() int                               { return 0 }
func (m *MockEbitenDevice) CursorY() int                               { return 0 }
func (m *MockEbitenDevice) SetCursorX(x int)                           {}
func (m *MockEbitenDevice) SetCursorY(y int)                           {}
func (m *MockEbitenDevice) Clear()                                     {}
func (m *MockEbitenDevice) DrawPixel(x, y int, color int)              {}
func (m *MockEbitenDevice) DrawGlyph(x, y int, glyph rune, fg, bg int) {}

func TestMockEbitenDevice_Draw(t *testing.T) {
	mock := &MockEbitenDevice{}
	screen := ebiten.NewImage(10, 10)

	testutils.False(t, "DrawCalled initially", mock.DrawCalled)
	mock.Draw(screen)
	testutils.True(t, "DrawCalled after Draw()", mock.DrawCalled)
}

func TestMockEbitenDevice_Layout(t *testing.T) {
	mock := &MockEbitenDevice{
		DeviceWidth:  320,
		DeviceHeight: 240,
	}

	testutils.False(t, "LayoutCalled initially", mock.LayoutCalled)
	w, h := mock.Layout(800, 600)
	testutils.True(t, "LayoutCalled after Layout()", mock.LayoutCalled)
	testutils.Equal(t, "Layout returned width", w, 320)
	testutils.Equal(t, "Layout returned height", h, 240)
	testutils.Equal(t, "Layout input w recorded", mock.LayoutW, 800)
	testutils.Equal(t, "Layout input h recorded", mock.LayoutH, 600)
}

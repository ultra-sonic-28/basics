package video

import (
	"basics/testutils"
	"bytes"
	"errors"
	"io"
	"testing"
)

// MockDevice implémente l'interface Device pour les tests
type MockDevice struct {
	ClearCalled      bool
	InverseSet       bool
	FlashSet         bool
	PrintedChars     []rune
	PrintedStrings   []string
	CurX, CurY       int
	PlottedPoints    [][2]int
	ReadLineOutput   string
	ReadLineErr      error
	GetCharOutput    rune
	GetCharErr       error
	Output           *bytes.Buffer
	KeyboardDisabled bool
	RenderCalled     bool
}

func (m *MockDevice) Clear()                    { m.ClearCalled = true }
func (m *MockDevice) SetInverse(v bool)         { m.InverseSet = v }
func (m *MockDevice) SetFlash(v bool)           { m.FlashSet = v }
func (m *MockDevice) PrintChar(r rune)          { m.PrintedChars = append(m.PrintedChars, r) }
func (m *MockDevice) PrintString(s string)      { m.PrintedStrings = append(m.PrintedStrings, s) }
func (m *MockDevice) SetCursorX(x int)          { m.CurX = x }
func (m *MockDevice) SetCursorY(y int)          { m.CurY = y }
func (m *MockDevice) CursorX() int              { return m.CurX }
func (m *MockDevice) CursorY() int              { return m.CurY }
func (m *MockDevice) Plot(x, y int)             { m.PlottedPoints = append(m.PlottedPoints, [2]int{x, y}) }
func (m *MockDevice) ReadLine() (string, error) { return m.ReadLineOutput, m.ReadLineErr }
func (m *MockDevice) GetChar() (rune, error)    { return m.GetCharOutput, m.GetCharErr }
func (m *MockDevice) SetOutput(w io.Writer)     { m.Output = w.(*bytes.Buffer) }
func (m *MockDevice) DisableKeyboard()          { m.KeyboardDisabled = true }
func (m *MockDevice) Render()                   { m.RenderCalled = true }

func TestMockDevice_ClearInverseFlash(t *testing.T) {
	dev := &MockDevice{}

	testutils.False(t, "ClearCalled initially", dev.ClearCalled)
	dev.Clear()
	testutils.True(t, "ClearCalled after Clear()", dev.ClearCalled)

	dev.SetInverse(true)
	testutils.True(t, "InverseSet after SetInverse(true)", dev.InverseSet)

	dev.SetFlash(true)
	testutils.True(t, "FlashSet after SetFlash(true)", dev.FlashSet)
}

func TestMockDevice_Printing(t *testing.T) {
	dev := &MockDevice{}

	dev.PrintChar('A')
	dev.PrintChar('B')
	dev.PrintString("Hello")

	testutils.DeepEqual(t, "PrintedChars", dev.PrintedChars, []rune{'A', 'B'})
	testutils.DeepEqual(t, "PrintedStrings", dev.PrintedStrings, []string{"Hello"})
}

func TestMockDevice_CursorAndPlot(t *testing.T) {
	dev := &MockDevice{}

	dev.SetCursorX(5)
	dev.SetCursorY(10)
	testutils.Equal(t, "CursorX", dev.CurX, 5)
	testutils.Equal(t, "CursorY", dev.CurY, 10)

	dev.Plot(1, 2)
	dev.Plot(3, 4)
	testutils.DeepEqual(t, "PlottedPoints", dev.PlottedPoints, [][2]int{{1, 2}, {3, 4}})
}

func TestMockDevice_IO(t *testing.T) {
	dev := &MockDevice{}
	buf := &bytes.Buffer{}
	dev.SetOutput(buf)

	testutils.Equal(t, "Output buffer set", dev.Output, buf)

	dev.DisableKeyboard()
	testutils.True(t, "KeyboardDisabled", dev.KeyboardDisabled)

	dev.ReadLineOutput = "input line"
	dev.ReadLineErr = errors.New("fail")
	line, err := dev.ReadLine()
	testutils.Equal(t, "ReadLine output", line, "input line")
	testutils.Equal(t, "ReadLine error", err, dev.ReadLineErr)

	dev.GetCharOutput = 'X'
	dev.GetCharErr = errors.New("err")
	ch, err := dev.GetChar()
	testutils.Equal(t, "GetChar output", ch, 'X')
	testutils.Equal(t, "GetChar error", err, dev.GetCharErr)
}

func TestMockDevice_Render(t *testing.T) {
	dev := &MockDevice{}

	testutils.False(t, "RenderCalled initially", dev.RenderCalled)
	dev.Render()
	testutils.True(t, "RenderCalled after Render()", dev.RenderCalled)
}

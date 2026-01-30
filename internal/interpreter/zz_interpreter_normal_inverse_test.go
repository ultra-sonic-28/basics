package interpreter

import (
	"basics/internal/lexer"
	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
	"io"
	"testing"
)

// FakeVideoDevice permet de capturer les appels à SetInverse
type FakeVideoDevice struct {
	InverseCalls []bool
}

func NewFakeVideoDevice() *FakeVideoDevice {
	return &FakeVideoDevice{}
}

// ---- API BASIC ----

func (f *FakeVideoDevice) Clear() {}

func (f *FakeVideoDevice) SetInverse(v bool) {
	f.InverseCalls = append(f.InverseCalls, v)
}

func (f *FakeVideoDevice) SetFlash(v bool)      {}
func (f *FakeVideoDevice) PrintChar(r rune)     {}
func (f *FakeVideoDevice) PrintString(s string) {}
func (f *FakeVideoDevice) SetCursorX(x int)     {}
func (f *FakeVideoDevice) SetCursorY(y int)     {}
func (f *FakeVideoDevice) Plot(x, y int)        {}
func (f *FakeVideoDevice) ReadLine() (string, error) {
	return "", nil
}
func (f *FakeVideoDevice) GetChar() (rune, error) {
	return 0, nil
}

// ---- I/O ----

func (f *FakeVideoDevice) SetOutput(w io.Writer) {}
func (f *FakeVideoDevice) DisableKeyboard()      {}

// ---- Render ----

func (f *FakeVideoDevice) Render() {}

func TestInterpreter_NORMAL_INVERSE(t *testing.T) {
	source := `
10 NORMAL
20 INVERSE
30 NORMAL
`

	// Lexer + Parser
	tokens := lexer.Lex(source)
	p := parser.New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)

	// Fake video
	fakeVideo := NewFakeVideoDevice()

	// Runtime
	rt := &runtime.Runtime{
		Video: fakeVideo,
		Env:   runtime.NewEnvironment(),
	}

	// Interpreter
	i := New(rt)
	i.Run(prog)

	// Vérification des appels
	expected := []bool{false, true, false}

	testutils.Equal(t, "inverse call count",
		len(fakeVideo.InverseCalls),
		len(expected),
	)

	for idx, v := range expected {
		testutils.Equal(
			t,
			"SetInverse call",
			fakeVideo.InverseCalls[idx],
			v,
		)
	}
}

func TestInterpreter_NORMAL_INVERSE_SingleLine(t *testing.T) {
	source := `
10 NORMAL:INVERSE:NORMAL
`

	tokens := lexer.Lex(source)
	p := parser.New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)

	fakeVideo := NewFakeVideoDevice()

	rt := &runtime.Runtime{
		Video: fakeVideo,
		Env:   runtime.NewEnvironment(),
	}

	i := New(rt)
	i.Run(prog)

	expected := []bool{false, true, false}

	testutils.Equal(t, "inverse call count",
		len(fakeVideo.InverseCalls),
		len(expected),
	)

	for i, v := range expected {
		testutils.Equal(t, "inverse value", fakeVideo.InverseCalls[i], v)
	}
}

func TestInterpreter_NORMAL_INVERSE_WithPrint(t *testing.T) {
	source := `
10 NORMAL
20 PRINT "HELLO"
30 INVERSE
40 PRINT "WORLD"
`

	tokens := lexer.Lex(source)
	p := parser.New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)

	fakeVideo := NewFakeVideoDevice()

	rt := &runtime.Runtime{
		Video: fakeVideo,
		Env:   runtime.NewEnvironment(),
	}

	i := New(rt)
	i.Run(prog)

	testutils.DeepEqual(t, "inverse calls", fakeVideo.InverseCalls, []bool{
		false,
		true,
	})
}

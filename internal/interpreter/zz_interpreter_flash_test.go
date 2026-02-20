package interpreter

import (
	"basics/internal/lexer"
	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
	"io"
	"testing"
)

// FakeVideoDeviceFlash permet de capturer les appels à SetFlash
type FakeVideoDeviceFlash struct {
	FlashCalls []bool
}

func NewFakeVideoDeviceFlash() *FakeVideoDeviceFlash {
	return &FakeVideoDeviceFlash{}
}

// ---- API BASIC ----

func (f *FakeVideoDeviceFlash) Clear() {}

func (f *FakeVideoDeviceFlash) SetFlash(v bool) {
	f.FlashCalls = append(f.FlashCalls, v)
}

func (f *FakeVideoDeviceFlash) SetInverse(v bool)    {}
func (f *FakeVideoDeviceFlash) PrintChar(r rune)     {}
func (f *FakeVideoDeviceFlash) PrintString(s string) {}
func (f *FakeVideoDeviceFlash) SetCursorX(x int)     {}
func (f *FakeVideoDeviceFlash) SetCursorY(y int)     {}
func (f *FakeVideoDeviceFlash) SwitchMode(slot int)  {}
func (f *FakeVideoDeviceFlash) Plot(x, y int)        {}
func (f *FakeVideoDeviceFlash) ReadLine() (string, error) {
	return "", nil
}
func (f *FakeVideoDeviceFlash) GetChar() (rune, error) {
	return 0, nil
}

// ---- I/O ----

func (f *FakeVideoDeviceFlash) SetOutput(w io.Writer) {}
func (f *FakeVideoDeviceFlash) DisableKeyboard()      {}

// ---- Render ----

func (f *FakeVideoDeviceFlash) Render() {}

func TestInterpreter_FLASH(t *testing.T) {
	source := `
10 NORMAL
20 FLASH
30 NORMAL
`

	// Lexer + Parser
	tokens := lexer.Lex(source)
	p := parser.New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)

	// Fake video
	fakeVideo := NewFakeVideoDeviceFlash()

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

	testutils.Equal(t, "flash call count",
		len(fakeVideo.FlashCalls),
		len(expected),
	)

	for idx, v := range expected {
		testutils.Equal(
			t,
			"SetFlash call",
			fakeVideo.FlashCalls[idx],
			v,
		)
	}
}

func TestInterpreter_FLASH_SingleLine(t *testing.T) {
	source := `
10 NORMAL:FLASH:NORMAL
`

	tokens := lexer.Lex(source)
	p := parser.New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)

	fakeVideo := NewFakeVideoDeviceFlash()

	rt := &runtime.Runtime{
		Video: fakeVideo,
		Env:   runtime.NewEnvironment(),
	}

	i := New(rt)
	i.Run(prog)

	expected := []bool{false, true, false}

	testutils.Equal(t, "flash call count",
		len(fakeVideo.FlashCalls),
		len(expected),
	)

	for i, v := range expected {
		testutils.Equal(t, "flash value", fakeVideo.FlashCalls[i], v)
	}
}

func TestInterpreter_FLASH_WithPrint(t *testing.T) {
	source := `
10 NORMAL
20 PRINT "HELLO"
30 FLASH
40 PRINT "WORLD"
`

	tokens := lexer.Lex(source)
	p := parser.New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)

	fakeVideo := NewFakeVideoDeviceFlash()

	rt := &runtime.Runtime{
		Video: fakeVideo,
		Env:   runtime.NewEnvironment(),
	}

	i := New(rt)
	i.Run(prog)

	testutils.DeepEqual(t, "flash calls", fakeVideo.FlashCalls, []bool{
		false,
		true,
	})
}

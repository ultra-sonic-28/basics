package interpreter

import (
	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
	"io"
	"testing"
)

type MockVideo struct{}

func (m *MockVideo) Clear()               {}
func (m *MockVideo) SetInverse(v bool)    {}
func (m *MockVideo) SetFlash(v bool)      {}
func (m *MockVideo) PrintChar(r rune)     {}
func (m *MockVideo) PrintString(s string) {}
func (m *MockVideo) SetCursorX(x int)     {}
func (m *MockVideo) SetCursorY(y int)     {}
func (m *MockVideo) SwitchMode(slot int)  {}
func (m *MockVideo) Plot(x, y int)        {}
func (m *MockVideo) ReadLine() (string, error) {
	return "", nil
}
func (m *MockVideo) GetChar() (rune, error) {
	return 0, nil
}
func (m *MockVideo) SetOutput(w io.Writer) {}
func (m *MockVideo) DisableKeyboard()      {}
func (m *MockVideo) Render()               {}

func TestInterpreter_ValAssignment(t *testing.T) {
	// 60 C% = VAL("3.5")
	
	p := &parser.Program{
		Lines: []*parser.Line{
			{
				Number: 60,
				Stmts: []parser.Statement{
					&parser.LetStmt{
						Name: "C%",
						Value: &parser.ValExpr{
							Expr: &parser.StringLiteral{Value: "3.5"},
						},
					},
				},
			},
		},
	}

	rt := runtime.New(&MockVideo{})
	interp := New(rt)
	
	// Exec the program
	interp.Run(p)
	
	val, ok := rt.Env.Get("C%")
	testutils.True(t, "variable C% should exist", ok)
	testutils.Equal(t, "C% should be INTEGER", val.Type, runtime.INTEGER)
	testutils.Equal(t, "C% should be 3", val.Int, 3)
}

package parser

import (
	"basics/internal/lexer"
	"basics/testutils"
	"testing"
)

func TestParse_CLEAR_Only(t *testing.T) {
	source := `10 CLEAR`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 1)

	line10 := prog.Lines[0]
	testutils.Equal(t, "stmt count", len(line10.Stmts), 1)

	_, ok := line10.Stmts[0].(*ClearStmt)
	testutils.True(t, "stmt is ClearStmt", ok)
}

func TestParse_CLEAR_InProgram(t *testing.T) {
	source := `
10 LET A = 10
20 CLEAR
30 PRINT A
`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 3)

	// ---- Line 20 : CLEAR ----
	line20 := prog.Lines[1]
	testutils.Equal(t, "line 20 stmt count", len(line20.Stmts), 1)

	_, ok := line20.Stmts[0].(*ClearStmt)
	testutils.True(t, "stmt is ClearStmt", ok)
}

func TestParse_CLEAR_WithColon(t *testing.T) {
	source := `10 CLEAR:PRINT "OK"`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 1)

	line10 := prog.Lines[0]
	testutils.Equal(t, "stmt count", len(line10.Stmts), 2)

	_, ok := line10.Stmts[0].(*ClearStmt)
	testutils.True(t, "first stmt is ClearStmt", ok)

	_, ok = line10.Stmts[1].(*PrintStmt)
	testutils.True(t, "second stmt is PrintStmt", ok)
}

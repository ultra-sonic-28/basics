package parser

import (
	"basics/internal/lexer"
	"basics/testutils"
	"testing"
)

func TestParse_NORMAL_INVERSE_Multiline(t *testing.T) {
	source := `
10 NORMAL
20 PRINT "HELLO"
30 INVERSE
40 PRINT "WORLD"
`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 4)

	// ---- Line 10 : NORMAL ----
	line10 := prog.Lines[0]
	testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 1)

	_, ok := line10.Stmts[0].(*NormalStmt)
	testutils.True(t, "stmt is NormalStmt", ok)

	// ---- Line 20 : PRINT ----
	line20 := prog.Lines[1]
	testutils.Equal(t, "line 20 stmt count", len(line20.Stmts), 1)

	_, ok = line20.Stmts[0].(*PrintStmt)
	testutils.True(t, "stmt is PrintStmt", ok)

	// ---- Line 30 : INVERSE ----
	line30 := prog.Lines[2]
	testutils.Equal(t, "line 30 stmt count", len(line30.Stmts), 1)

	_, ok = line30.Stmts[0].(*InverseStmt)
	testutils.True(t, "stmt is InverseStmt", ok)

	// ---- Line 40 : PRINT ----
	line40 := prog.Lines[3]
	testutils.Equal(t, "line 40 stmt count", len(line40.Stmts), 1)

	_, ok = line40.Stmts[0].(*PrintStmt)
	testutils.True(t, "stmt is PrintStmt", ok)
}

func TestParse_NORMAL_INVERSE_WithPrintChain(t *testing.T) {
	source := `
10 NORMAL
20 PRINT "HELLO ";
30 INVERSE
40 PRINT "WORLD";
50 NORMAL
60 PRINT "!"
`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 6)

	_, ok := prog.Lines[0].Stmts[0].(*NormalStmt)
	testutils.True(t, "line 10 NORMAL", ok)

	_, ok = prog.Lines[1].Stmts[0].(*PrintStmt)
	testutils.True(t, "line 20 PRINT", ok)

	_, ok = prog.Lines[2].Stmts[0].(*InverseStmt)
	testutils.True(t, "line 30 INVERSE", ok)

	_, ok = prog.Lines[3].Stmts[0].(*PrintStmt)
	testutils.True(t, "line 40 PRINT", ok)

	_, ok = prog.Lines[4].Stmts[0].(*NormalStmt)
	testutils.True(t, "line 50 NORMAL", ok)

	_, ok = prog.Lines[5].Stmts[0].(*PrintStmt)
	testutils.True(t, "line 60 PRINT", ok)
}

func TestParse_NORMAL_INVERSE_SingleLine(t *testing.T) {
	source := `
10 NORMAL:PRINT "HELLO ";:INVERSE:PRINT "WORLD";:NORMAL:PRINT "!"
`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 1)

	line := prog.Lines[0]
	testutils.Equal(t, "stmt count", len(line.Stmts), 6)

	_, ok := line.Stmts[0].(*NormalStmt)
	testutils.True(t, "stmt 0 NORMAL", ok)

	_, ok = line.Stmts[1].(*PrintStmt)
	testutils.True(t, "stmt 1 PRINT", ok)

	_, ok = line.Stmts[2].(*InverseStmt)
	testutils.True(t, "stmt 2 INVERSE", ok)

	_, ok = line.Stmts[3].(*PrintStmt)
	testutils.True(t, "stmt 3 PRINT", ok)

	_, ok = line.Stmts[4].(*NormalStmt)
	testutils.True(t, "stmt 4 NORMAL", ok)

	_, ok = line.Stmts[5].(*PrintStmt)
	testutils.True(t, "stmt 5 PRINT", ok)
}

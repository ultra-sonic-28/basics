package parser

import (
	"basics/internal/lexer"
	"basics/testutils"
	"testing"
)

func TestParse_GET_Simple(t *testing.T) {
	source := `10 GET A$`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 1)

	// ---- Line 10 : GET A$ ----
	line10 := prog.Lines[0]
	testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 1)

	getStmt, ok := line10.Stmts[0].(*GetStmt)
	testutils.True(t, "stmt is GetStmt", ok)

	testutils.Equal(t, "var name", getStmt.Var.Name, "A$")
}

func TestParse_GET_NumericVar(t *testing.T) {
	source := `
10 GET A
`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 1)

	line10 := prog.Lines[0]
	testutils.Equal(t, "stmt count", len(line10.Stmts), 1)

	getStmt, ok := line10.Stmts[0].(*GetStmt)
	testutils.True(t, "stmt is GetStmt", ok)

	testutils.Equal(t, "var name", getStmt.Var.Name, "A")
}

func TestParse_GET_InProgram(t *testing.T) {
	source := `
10 PRINT "Appuyez sur une touche pour continuer"
20 GET A$
30 PRINT "Merci"
`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 3)

	// ---- Line 20 : GET A$ ----
	line20 := prog.Lines[1]
	testutils.Equal(t, "line 20 stmt count", len(line20.Stmts), 1)

	getStmt, ok := line20.Stmts[0].(*GetStmt)
	testutils.True(t, "stmt is GetStmt", ok)

	testutils.Equal(t, "var name", getStmt.Var.Name, "A$")
}

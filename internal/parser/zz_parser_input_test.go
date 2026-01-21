package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_INPUT_Simple(t *testing.T) {
	source := `
10 PRINT "Entrez votre nom : ";
20 INPUT N$
30 PRINT "Entrez votre age : ";
40 INPUT A
50 PRINT:PRINT N$;", vous avez ";A;" ans"
`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 5)

	// ---- Line 20 : INPUT N$ ----
	line20 := prog.Lines[1]
	testutils.Equal(t, "line 20 stmt count", len(line20.Stmts), 1)

	inputStmt, ok := line20.Stmts[0].(*InputStmt)
	testutils.True(t, "stmt is InputStmt", ok)

	testutils.True(t, "no prompt", inputStmt.Prompt == nil)
	testutils.Equal(t, "var count", len(inputStmt.Vars), 1)
	testutils.Equal(t, "var name", inputStmt.Vars[0].Name, "N$")

	// ---- Line 40 : INPUT A ----
	line40 := prog.Lines[3]
	testutils.Equal(t, "line 40 stmt count", len(line40.Stmts), 1)

	inputStmt, ok = line40.Stmts[0].(*InputStmt)
	testutils.True(t, "stmt is InputStmt", ok)

	testutils.True(t, "no prompt", inputStmt.Prompt == nil)
	testutils.Equal(t, "var count", len(inputStmt.Vars), 1)
	testutils.Equal(t, "var name", inputStmt.Vars[0].Name, "A")
}

func TestParse_INPUT_WithPromptAndMultipleVars(t *testing.T) {
	source := `
10 PRINT "Multiply 2 numbers"
20 INPUT "Enter 2 values: ";A,B
30 PRINT "A*B is ";A*B
`

	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()

	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "line count", len(prog.Lines), 3)

	// ---- Line 20 ----
	line20 := prog.Lines[1]
	testutils.Equal(t, "line 20 stmt count", len(line20.Stmts), 1)

	inputStmt, ok := line20.Stmts[0].(*InputStmt)
	testutils.True(t, "stmt is InputStmt", ok)

	// Prompt
	testutils.True(t, "prompt present", inputStmt.Prompt != nil)
	testutils.Equal(
		t,
		"prompt value",
		inputStmt.Prompt.Value,
		"Enter 2 values: ",
	)

	// Variables
	testutils.Equal(t, "var count", len(inputStmt.Vars), 2)
	testutils.Equal(t, "var 0", inputStmt.Vars[0].Name, "A")
	testutils.Equal(t, "var 1", inputStmt.Vars[1].Name, "B")
}

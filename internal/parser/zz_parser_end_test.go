package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_END_Statement(t *testing.T) {
	source := `
10 END
`

	// Lexing
	tokens := lexer.Lex(source)

	// Parsing
	p := New(tokens)
	prog, errs := p.ParseProgram()

	// --- Assertions ---
	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "one line parsed", len(prog.Lines), 1)

	line := prog.Lines[0]
	testutils.Equal(t, "line number", line.Number, 10)
	testutils.Equal(t, "one statement", len(line.Stmts), 1)

	stmt := line.Stmts[0]

	// Type assertion
	_, ok := stmt.(*EndStmt)
	testutils.Equal(t, "statement is EndStmt", ok, true)
}

func TestParse_PRINT_COLON_END(t *testing.T) {
	source := `
10 PRINT "HELLO": END
`

	// Lexing
	tokens := lexer.Lex(source)

	// Parsing
	p := New(tokens)
	prog, errs := p.ParseProgram()

	// --- Assertions ---
	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "one line parsed", len(prog.Lines), 1)

	line := prog.Lines[0]
	testutils.Equal(t, "line number", line.Number, 10)
	testutils.Equal(t, "two statements", len(line.Stmts), 2)

	// --- Statement 1 : PRINT ---
	_, ok := line.Stmts[0].(*PrintStmt)
	testutils.Equal(t, "stmt[0] is PrintStmt", ok, true)

	// --- Statement 2 : END ---
	_, ok = line.Stmts[1].(*EndStmt)
	testutils.Equal(t, "stmt[1] is EndStmt", ok, true)
}

func TestParse_PRINT_END_PRINT_SameLine(t *testing.T) {
	source := `
10 PRINT "A": END: PRINT "B"
`

	// Lexing
	tokens := lexer.Lex(source)

	// Parsing
	p := New(tokens)
	prog, errs := p.ParseProgram()

	// --- Assertions globales ---
	testutils.Equal(t, "no parser errors", len(errs), 0)
	testutils.Equal(t, "one line parsed", len(prog.Lines), 1)

	line := prog.Lines[0]
	testutils.Equal(t, "line number", line.Number, 10)
	testutils.Equal(t, "three statements", len(line.Stmts), 3)

	// --- Statement 1 : PRINT "A" ---
	ps1, ok := line.Stmts[0].(*PrintStmt)
	testutils.Equal(t, "stmt[0] is PrintStmt", ok, true)
	testutils.Equal(t, "stmt[0] expr count", len(ps1.Exprs), 1)

	sl1, ok := ps1.Exprs[0].(*StringLiteral)
	testutils.Equal(t, "stmt[0] expr is StringLiteral", ok, true)
	testutils.Equal(t, "stmt[0] string value", sl1.Value, "A")

	// --- Statement 2 : END ---
	_, ok = line.Stmts[1].(*EndStmt)
	testutils.Equal(t, "stmt[1] is EndStmt", ok, true)

	// --- Statement 3 : PRINT "B" ---
	ps2, ok := line.Stmts[2].(*PrintStmt)
	testutils.Equal(t, "stmt[2] is PrintStmt", ok, true)
	testutils.Equal(t, "stmt[2] expr count", len(ps2.Exprs), 1)

	sl2, ok := ps2.Exprs[0].(*StringLiteral)
	testutils.Equal(t, "stmt[2] expr is StringLiteral", ok, true)
	testutils.Equal(t, "stmt[2] string value", sl2.Value, "B")
}

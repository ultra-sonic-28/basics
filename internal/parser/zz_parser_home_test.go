package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_HOME_Statement(t *testing.T) {
	source := `
10 HOME
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
	testutils.Equal(t, "one statement", len(line.Stmts), 1)

	stmt := line.Stmts[0]

	// --- Type assertion ---
	home, ok := stmt.(*HomeStmt)
	testutils.Equal(t, "statement is HomeStmt", ok, true)

	// --- Bonus : vérification position ---
	testutils.Equal(t, "line stored", home.Line, 2)
}

func TestParse_IF_THEN_HOME(t *testing.T) {
	source := `
10 IF A=1 THEN HOME
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
	testutils.Equal(t, "one statement", len(line.Stmts), 1)

	// --- IF statement ---
	ifStmt, ok := line.Stmts[0].(*IfStmt)
	testutils.Equal(t, "stmt is IfStmt", ok, true)

	// --- THEN block ---
	testutils.Equal(t, "then block size", len(ifStmt.Then), 1)

	_, ok = ifStmt.Then[0].(*HomeStmt)
	testutils.Equal(t, "then stmt is HomeStmt", ok, true)

	// --- ELSE block absent ---
	testutils.Equal(t, "else block empty", len(ifStmt.Else), 0)
}

func TestParse_HOME_COLON_PRINT(t *testing.T) {
	source := `
10 HOME: PRINT "OK"
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
	testutils.Equal(t, "two statements", len(line.Stmts), 2)

	// --- Statement 1 : HOME ---
	_, ok := line.Stmts[0].(*HomeStmt)
	testutils.Equal(t, "stmt[0] is HomeStmt", ok, true)

	// --- Statement 2 : PRINT ---
	ps, ok := line.Stmts[1].(*PrintStmt)
	testutils.Equal(t, "stmt[1] is PrintStmt", ok, true)

	testutils.Equal(t, "print expr count", len(ps.Exprs), 1)
	sl, ok := ps.Exprs[0].(*StringLiteral)
	testutils.Equal(t, "print expr is StringLiteral", ok, true)
	testutils.Equal(t, "print value", sl.Value, "OK")
}

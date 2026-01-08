package parser

import (
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestParser_AST_Deep(t *testing.T) {
	tokens := []token.Token{
		// 10 LET A = 1 + 2 * 3
		{Type: token.LINENUM, Literal: "10"},
		{Type: token.KEYWORD, Literal: "LET"},
		{Type: token.IDENT, Literal: "A"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.NUMBER, Literal: "1"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.NUMBER, Literal: "2"},
		{Type: token.ASTERISK, Literal: "*"},
		{Type: token.NUMBER, Literal: "3"},
		{Type: token.EOL},

		// 20 PRINT A; "!"
		{Type: token.LINENUM, Literal: "20"},
		{Type: token.KEYWORD, Literal: "PRINT"},
		{Type: token.IDENT, Literal: "A"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.STRING, Literal: "!"},
		{Type: token.EOL},

		{Type: token.EOF},
	}

	p := New(tokens)
	prog, errs := p.ParseProgram()

	// --- erreurs ---
	testutils.Equal(t, "no parser errors", len(errs), 0)

	// --- programme ---
	testutils.Equal(t, "line count", len(prog.Lines), 2)

	// =====================================================
	// 🔹 Ligne 10
	// =====================================================
	line10 := prog.Lines[0]
	testutils.Equal(t, "line 10 number", line10.Number, 10)
	testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 1)

	let, ok := line10.Stmts[0].(*LetStmt)
	testutils.True(t, "line 10 is LetStmt", ok)
	testutils.Equal(t, "LET variable", let.Name, "A")

	// LET value = 1 + 2 * 3
	plus, ok := let.Value.(*InfixExpr)
	testutils.True(t, "LET value is InfixExpr", ok)
	testutils.Equal(t, "outer op", plus.Op, "+")

	// gauche: 1
	leftNum, ok := plus.Left.(*NumberLiteral)
	testutils.True(t, "left is NumberLiteral", ok)
	testutils.Equal(t, "left value", leftNum.Value, 1.0)

	// droite: 2 * 3
	mul, ok := plus.Right.(*InfixExpr)
	testutils.True(t, "right is InfixExpr", ok)
	testutils.Equal(t, "inner op", mul.Op, "*")

	n2, ok := mul.Left.(*NumberLiteral)
	testutils.True(t, "mul left number", ok)
	testutils.Equal(t, "mul left value", n2.Value, 2.0)

	n3, ok := mul.Right.(*NumberLiteral)
	testutils.True(t, "mul right number", ok)
	testutils.Equal(t, "mul right value", n3.Value, 3.0)

	// =====================================================
	// 🔹 Ligne 20
	// =====================================================
	line20 := prog.Lines[1]
	testutils.Equal(t, "line 20 number", line20.Number, 20)
	testutils.Equal(t, "line 20 stmt count", len(line20.Stmts), 1)

	printStmt, ok := line20.Stmts[0].(*PrintStmt)
	testutils.True(t, "line 20 is PrintStmt", ok)

	testutils.Equal(t, "print expr count", len(printStmt.Exprs), 2)
	testutils.Equal(t, "print separator count", len(printStmt.Separators), 1)
	testutils.Equal(t, "separator", printStmt.Separators[0], ';')

	// PRINT A
	id, ok := printStmt.Exprs[0].(*Identifier)
	testutils.True(t, "first expr is Identifier", ok)
	testutils.Equal(t, "identifier name", id.Name, "A")

	// PRINT "!"
	str, ok := printStmt.Exprs[1].(*StringLiteral)
	testutils.True(t, "second expr is StringLiteral", ok)
	testutils.Equal(t, "string literal", str.Value, "!")
}

package parser

import (
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestParser_Golden_AST_Markdown(t *testing.T) {
	tokens := []token.Token{
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
	testutils.Equal(t, "no errors", len(errs), 0)

	got := ASTToMarkdownTable(prog)

	want := `| Path | Type | Value |
|------|------|-------|
| Program/Line[10] | Line |  |
| Program/Line[10]/Stmt[0] | LetStmt | A |
| Program/Line[10]/Stmt[0]/Value | InfixExpr | + |
| Program/Line[10]/Stmt[0]/Value/Left | NumberLiteral | 1 |
| Program/Line[10]/Stmt[0]/Value/Right | InfixExpr | * |
| Program/Line[10]/Stmt[0]/Value/Right/Left | NumberLiteral | 2 |
| Program/Line[10]/Stmt[0]/Value/Right/Right | NumberLiteral | 3 |
| Program/Line[20] | Line |  |
| Program/Line[20]/Stmt[0] | PrintStmt |  |
| Program/Line[20]/Stmt[0]/Expr[0] | Identifier | A |
| Program/Line[20]/Stmt[0]/Expr[1] | StringLiteral | ! |
`

	testutils.Equal(t, "AST markdown", got, want)
}

package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParser_Golden_Positioning(t *testing.T) {
	source := `10 HTAB 10
20 VTAB 5
30 HTAB A + 1: VTAB B * 2
40 PRINT TAB(20); "HELLO"
50 PRINT SPC(10); "A"
60 PRINT SPC(2 * A); "A"
70 PRINT RND(1)
80 PRINT SIN(0)
`
	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()
	testutils.Equal(t, "no errors", len(errs), 0)

	got := ASTToMarkdownTable(prog)

	want := `| Path | Type | Value |
|------|------|-------|
| Program/Line[10] | Line |  |
| Program/Line[10]/Stmt[0] | HTabStmt |  |
| Program/Line[10]/Stmt[0]/Expr | NumberLiteral | 10 |
| Program/Line[20] | Line |  |
| Program/Line[20]/Stmt[0] | VTabStmt |  |
| Program/Line[20]/Stmt[0]/Expr | NumberLiteral | 5 |
| Program/Line[30] | Line |  |
| Program/Line[30]/Stmt[0] | HTabStmt |  |
| Program/Line[30]/Stmt[0]/Expr | InfixExpr | + |
| Program/Line[30]/Stmt[0]/Expr/Left | Identifier | A |
| Program/Line[30]/Stmt[0]/Expr/Right | NumberLiteral | 1 |
| Program/Line[30]/Stmt[1] | VTabStmt |  |
| Program/Line[30]/Stmt[1]/Expr | InfixExpr | * |
| Program/Line[30]/Stmt[1]/Expr/Left | Identifier | B |
| Program/Line[30]/Stmt[1]/Expr/Right | NumberLiteral | 2 |
| Program/Line[40] | Line |  |
| Program/Line[40]/Stmt[0] | PrintStmt |  |
| Program/Line[40]/Stmt[0]/Expr[0] | TabExpr |  |
| Program/Line[40]/Stmt[0]/Expr[0]/Expr | NumberLiteral | 20 |
| Program/Line[40]/Stmt[0]/Expr[1] | StringLiteral | HELLO |
| Program/Line[50] | Line |  |
| Program/Line[50]/Stmt[0] | PrintStmt |  |
| Program/Line[50]/Stmt[0]/Expr[0] | SpcExpr |  |
| Program/Line[50]/Stmt[0]/Expr[0]/Expr | NumberLiteral | 10 |
| Program/Line[50]/Stmt[0]/Expr[1] | StringLiteral | A |
| Program/Line[60] | Line |  |
| Program/Line[60]/Stmt[0] | PrintStmt |  |
| Program/Line[60]/Stmt[0]/Expr[0] | SpcExpr |  |
| Program/Line[60]/Stmt[0]/Expr[0]/Expr | InfixExpr | * |
| Program/Line[60]/Stmt[0]/Expr[0]/Expr/Left | NumberLiteral | 2 |
| Program/Line[60]/Stmt[0]/Expr[0]/Expr/Right | Identifier | A |
| Program/Line[60]/Stmt[0]/Expr[1] | StringLiteral | A |
| Program/Line[70] | Line |  |
| Program/Line[70]/Stmt[0] | PrintStmt |  |
| Program/Line[70]/Stmt[0]/Expr[0] | RndExpr |  |
| Program/Line[70]/Stmt[0]/Expr[0]/Expr | NumberLiteral | 1 |
| Program/Line[80] | Line |  |
| Program/Line[80]/Stmt[0] | PrintStmt |  |
| Program/Line[80]/Stmt[0]/Expr[0] | SinExpr |  |
| Program/Line[80]/Stmt[0]/Expr[0]/Expr | NumberLiteral | 0 |
`

	testutils.Equal(t, "AST markdown for positioning statements", got, want)
}

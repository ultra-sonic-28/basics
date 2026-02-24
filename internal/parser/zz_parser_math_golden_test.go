package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParser_Golden_Math(t *testing.T) {
	source := `10 PRINT INT(1.5); ABS(-5); SGN(-10); SQR(16)
`
	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()
	testutils.Equal(t, "no errors", len(errs), 0)

	got := ASTToMarkdownTable(prog)

	want := `| Path | Type | Value |
|------|------|-------|
| Program/Line[10] | Line |  |
| Program/Line[10]/Stmt[0] | PrintStmt |  |
| Program/Line[10]/Stmt[0]/Expr[0] | IntExpr |  |
| Program/Line[10]/Stmt[0]/Expr[0]/Expr | NumberLiteral | 1.5 |
| Program/Line[10]/Stmt[0]/Expr[1] | AbsExpr |  |
| Program/Line[10]/Stmt[0]/Expr[1]/Expr | PrefixExpr | - |
| Program/Line[10]/Stmt[0]/Expr[1]/Expr/Right | NumberLiteral | 5 |
| Program/Line[10]/Stmt[0]/Expr[2] | SgnExpr |  |
| Program/Line[10]/Stmt[0]/Expr[2]/Expr | PrefixExpr | - |
| Program/Line[10]/Stmt[0]/Expr[2]/Expr/Right | NumberLiteral | 10 |
| Program/Line[10]/Stmt[0]/Expr[3] | SqrExpr |  |
| Program/Line[10]/Stmt[0]/Expr[3]/Expr | NumberLiteral | 16 |
`

	testutils.Equal(t, "AST markdown for Math functions", got, want)
}

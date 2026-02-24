package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParser_Golden_Strings(t *testing.T) {
	source := `10 PRINT LEN("APPLE"); ASC("A")
20 PRINT CHR$(65); STR$(123); VAL("456")
30 PRINT LEFT$("ABC", 1); RIGHT$("ABC", 1); MID$("ABC", 1, 1)
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
| Program/Line[10]/Stmt[0]/Expr[0] | LenExpr |  |
| Program/Line[10]/Stmt[0]/Expr[0]/Expr | StringLiteral | APPLE |
| Program/Line[10]/Stmt[0]/Expr[1] | AscExpr |  |
| Program/Line[10]/Stmt[0]/Expr[1]/Expr | StringLiteral | A |
| Program/Line[20] | Line |  |
| Program/Line[20]/Stmt[0] | PrintStmt |  |
| Program/Line[20]/Stmt[0]/Expr[0] | ChrExpr |  |
| Program/Line[20]/Stmt[0]/Expr[0]/Expr | NumberLiteral | 65 |
| Program/Line[20]/Stmt[0]/Expr[1] | StrExpr |  |
| Program/Line[20]/Stmt[0]/Expr[1]/Expr | NumberLiteral | 123 |
| Program/Line[20]/Stmt[0]/Expr[2] | ValExpr |  |
| Program/Line[20]/Stmt[0]/Expr[2]/Expr | StringLiteral | 456 |
| Program/Line[30] | Line |  |
| Program/Line[30]/Stmt[0] | PrintStmt |  |
| Program/Line[30]/Stmt[0]/Expr[0] | LeftExpr |  |
| Program/Line[30]/Stmt[0]/Expr[0]/Str | StringLiteral | ABC |
| Program/Line[30]/Stmt[0]/Expr[0]/Len | NumberLiteral | 1 |
| Program/Line[30]/Stmt[0]/Expr[1] | RightExpr |  |
| Program/Line[30]/Stmt[0]/Expr[1]/Str | StringLiteral | ABC |
| Program/Line[30]/Stmt[0]/Expr[1]/Len | NumberLiteral | 1 |
| Program/Line[30]/Stmt[0]/Expr[2] | MidExpr |  |
| Program/Line[30]/Stmt[0]/Expr[2]/Str | StringLiteral | ABC |
| Program/Line[30]/Stmt[0]/Expr[2]/Start | NumberLiteral | 1 |
| Program/Line[30]/Stmt[0]/Expr[2]/Len | NumberLiteral | 1 |
`

	testutils.Equal(t, "AST markdown for String functions", got, want)
}

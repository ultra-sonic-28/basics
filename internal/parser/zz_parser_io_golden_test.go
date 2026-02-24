package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParser_Golden_IO(t *testing.T) {
	source := `10 PRINT
20 PRINT "HELLO"
30 PRINT "A="; A, "B="; B
40 PRINT TAB(10); "WORLD";
50 INPUT A
60 INPUT "VALUE? "; B, C
70 GET K$
80 GET N
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
| Program/Line[20] | Line |  |
| Program/Line[20]/Stmt[0] | PrintStmt |  |
| Program/Line[20]/Stmt[0]/Expr[0] | StringLiteral | HELLO |
| Program/Line[30] | Line |  |
| Program/Line[30]/Stmt[0] | PrintStmt |  |
| Program/Line[30]/Stmt[0]/Expr[0] | StringLiteral | A= |
| Program/Line[30]/Stmt[0]/Expr[1] | Identifier | A |
| Program/Line[30]/Stmt[0]/Expr[2] | StringLiteral | B= |
| Program/Line[30]/Stmt[0]/Expr[3] | Identifier | B |
| Program/Line[40] | Line |  |
| Program/Line[40]/Stmt[0] | PrintStmt |  |
| Program/Line[40]/Stmt[0]/Expr[0] | TabExpr |  |
| Program/Line[40]/Stmt[0]/Expr[0]/Expr | NumberLiteral | 10 |
| Program/Line[40]/Stmt[0]/Expr[1] | StringLiteral | WORLD |
| Program/Line[50] | Line |  |
| Program/Line[50]/Stmt[0] | InputStmt |  |
| Program/Line[50]/Stmt[0]/Var[0] | Identifier | A |
| Program/Line[60] | Line |  |
| Program/Line[60]/Stmt[0] | InputStmt |  |
| Program/Line[60]/Stmt[0]/Prompt | StringLiteral | VALUE?  |
| Program/Line[60]/Stmt[0]/Var[0] | Identifier | B |
| Program/Line[60]/Stmt[0]/Var[1] | Identifier | C |
| Program/Line[70] | Line |  |
| Program/Line[70]/Stmt[0] | GetStmt |  |
| Program/Line[70]/Stmt[0]/Var | Identifier | K$ |
| Program/Line[80] | Line |  |
| Program/Line[80]/Stmt[0] | GetStmt |  |
| Program/Line[80]/Stmt[0]/Var | Identifier | N |
`

	testutils.Equal(t, "AST markdown for IO statements", got, want)
}

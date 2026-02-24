package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParser_Golden_Style(t *testing.T) {
	source := `10 INVERSE
20 PRINT "INVERSE TEXT"
30 NORMAL
40 PRINT "NORMAL TEXT"
50 FLASH
60 PRINT "FLASHING TEXT"
70 NORMAL: PRINT "BACK TO NORMAL"
`
	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()
	testutils.Equal(t, "no errors", len(errs), 0)

	got := ASTToMarkdownTable(prog)

	want := `| Path | Type | Value |
|------|------|-------|
| Program/Line[10] | Line |  |
| Program/Line[10]/Stmt[0] | InverseStmt |  |
| Program/Line[20] | Line |  |
| Program/Line[20]/Stmt[0] | PrintStmt |  |
| Program/Line[20]/Stmt[0]/Expr[0] | StringLiteral | INVERSE TEXT |
| Program/Line[30] | Line |  |
| Program/Line[30]/Stmt[0] | NormalStmt |  |
| Program/Line[40] | Line |  |
| Program/Line[40]/Stmt[0] | PrintStmt |  |
| Program/Line[40]/Stmt[0]/Expr[0] | StringLiteral | NORMAL TEXT |
| Program/Line[50] | Line |  |
| Program/Line[50]/Stmt[0] | FlashStmt |  |
| Program/Line[60] | Line |  |
| Program/Line[60]/Stmt[0] | PrintStmt |  |
| Program/Line[60]/Stmt[0]/Expr[0] | StringLiteral | FLASHING TEXT |
| Program/Line[70] | Line |  |
| Program/Line[70]/Stmt[0] | NormalStmt |  |
| Program/Line[70]/Stmt[1] | PrintStmt |  |
| Program/Line[70]/Stmt[1]/Expr[0] | StringLiteral | BACK TO NORMAL |
`

	testutils.Equal(t, "AST markdown for style statements", got, want)
}

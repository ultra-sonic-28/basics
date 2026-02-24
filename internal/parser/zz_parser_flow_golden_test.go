package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParser_Golden_Flow(t *testing.T) {
	source := `10 FOR I = 1 TO 10 STEP 2
20 PRINT I
30 NEXT I
40 IF A = 1 THEN GOTO 100
50 IF B = 2 THEN PRINT "B IS 2": GOTO 200 ELSE END
60 GOTO 300
70 GOSUB 400
80 RETURN
90 GOTO A + 10
100 GOSUB B * 5
`
	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()
	testutils.Equal(t, "no errors", len(errs), 0)

	got := ASTToMarkdownTable(prog)

	want := `| Path | Type | Value |
|------|------|-------|
| Program/Line[10] | Line |  |
| Program/Line[10]/Stmt[0] | ForStmt | I |
| Program/Line[10]/Stmt[0]/Start | NumberLiteral | 1 |
| Program/Line[10]/Stmt[0]/End | NumberLiteral | 10 |
| Program/Line[10]/Stmt[0]/Step | NumberLiteral | 2 |
| Program/Line[20] | Line |  |
| Program/Line[20]/Stmt[0] | PrintStmt |  |
| Program/Line[20]/Stmt[0]/Expr[0] | Identifier | I |
| Program/Line[30] | Line |  |
| Program/Line[30]/Stmt[0] | NextStmt | I |
| Program/Line[40] | Line |  |
| Program/Line[40]/Stmt[0] | IfStmt |  |
| Program/Line[40]/Stmt[0]/Cond | InfixExpr | = |
| Program/Line[40]/Stmt[0]/Cond/Left | Identifier | A |
| Program/Line[40]/Stmt[0]/Cond/Right | NumberLiteral | 1 |
| Program/Line[40]/Stmt[0]/Then[0] | GotoStmt |  |
| Program/Line[40]/Stmt[0]/Then[0]/Target | NumberLiteral | 100 |
| Program/Line[50] | Line |  |
| Program/Line[50]/Stmt[0] | IfStmt |  |
| Program/Line[50]/Stmt[0]/Cond | InfixExpr | = |
| Program/Line[50]/Stmt[0]/Cond/Left | Identifier | B |
| Program/Line[50]/Stmt[0]/Cond/Right | NumberLiteral | 2 |
| Program/Line[50]/Stmt[0]/Then[0] | PrintStmt |  |
| Program/Line[50]/Stmt[0]/Then[0]/Expr[0] | StringLiteral | B IS 2 |
| Program/Line[50]/Stmt[0]/Then[1] | GotoStmt |  |
| Program/Line[50]/Stmt[0]/Then[1]/Target | NumberLiteral | 200 |
| Program/Line[50]/Stmt[0]/Else[0] | EndStmt |  |
| Program/Line[60] | Line |  |
| Program/Line[60]/Stmt[0] | GotoStmt |  |
| Program/Line[60]/Stmt[0]/Target | NumberLiteral | 300 |
| Program/Line[70] | Line |  |
| Program/Line[70]/Stmt[0] | GosubStmt |  |
| Program/Line[70]/Stmt[0]/Target | NumberLiteral | 400 |
| Program/Line[80] | Line |  |
| Program/Line[80]/Stmt[0] | ReturnStmt |  |
| Program/Line[90] | Line |  |
| Program/Line[90]/Stmt[0] | GotoStmt |  |
| Program/Line[90]/Stmt[0]/Target | InfixExpr | + |
| Program/Line[90]/Stmt[0]/Target/Left | Identifier | A |
| Program/Line[90]/Stmt[0]/Target/Right | NumberLiteral | 10 |
| Program/Line[100] | Line |  |
| Program/Line[100]/Stmt[0] | GosubStmt |  |
| Program/Line[100]/Stmt[0]/Target | InfixExpr | * |
| Program/Line[100]/Stmt[0]/Target/Left | Identifier | B |
| Program/Line[100]/Stmt[0]/Target/Right | NumberLiteral | 5 |
`

	testutils.Equal(t, "AST markdown for flow control statements", got, want)
}

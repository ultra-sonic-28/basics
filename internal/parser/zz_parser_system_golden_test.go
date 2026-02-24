package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParser_Golden_System(t *testing.T) {
	source := `10 REM THIS IS A COMMENT
20 PR# 0
30 PR# 3
40 PR# A + 1
50 END
`
	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()
	testutils.Equal(t, "no errors", len(errs), 0)

	got := ASTToMarkdownTable(prog)

	want := `| Path | Type | Value |
|------|------|-------|
| Program/Line[10] | Line |  |
| Program/Line[10]/Stmt[0] | RemStmt | THIS IS A COMMENT |
| Program/Line[20] | Line |  |
| Program/Line[20]/Stmt[0] | PrStmt |  |
| Program/Line[20]/Stmt[0]/Slot | NumberLiteral | 0 |
| Program/Line[30] | Line |  |
| Program/Line[30]/Stmt[0] | PrStmt |  |
| Program/Line[30]/Stmt[0]/Slot | NumberLiteral | 3 |
| Program/Line[40] | Line |  |
| Program/Line[40]/Stmt[0] | PrStmt |  |
| Program/Line[40]/Stmt[0]/Slot | InfixExpr | + |
| Program/Line[40]/Stmt[0]/Slot/Left | Identifier | A |
| Program/Line[40]/Stmt[0]/Slot/Right | NumberLiteral | 1 |
| Program/Line[50] | Line |  |
| Program/Line[50]/Stmt[0] | EndStmt |  |
`

	testutils.Equal(t, "AST markdown for system statements", got, want)
}

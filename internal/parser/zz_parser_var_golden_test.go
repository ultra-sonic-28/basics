package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParser_Golden_Var(t *testing.T) {
	source := `10 LET A = 1
20 B$ = "HELLO"
30 DIM C(10), D$(5, 5)
40 LET C(1) = 2
50 D$(1, 2) = "WORLD"
60 CLEAR
`
	tokens := lexer.Lex(source)
	p := New(tokens)
	prog, errs := p.ParseProgram()
	testutils.Equal(t, "no errors", len(errs), 0)

	got := ASTToMarkdownTable(prog)

	want := `| Path | Type | Value |
|------|------|-------|
| Program/Line[10] | Line |  |
| Program/Line[10]/Stmt[0] | LetStmt | A |
| Program/Line[10]/Stmt[0]/Value | NumberLiteral | 1 |
| Program/Line[20] | Line |  |
| Program/Line[20]/Stmt[0] | LetStmt | B$ |
| Program/Line[20]/Stmt[0]/Value | StringLiteral | HELLO |
| Program/Line[30] | Line |  |
| Program/Line[30]/Stmt[0] | DimStmt |  |
| Program/Line[30]/Stmt[0]/Array[0] | DimDecl | C |
| Program/Line[30]/Stmt[0]/Array[0]/Dim[0] | NumberLiteral | 10 |
| Program/Line[30]/Stmt[0]/Array[1] | DimDecl | D$ |
| Program/Line[30]/Stmt[0]/Array[1]/Dim[0] | NumberLiteral | 5 |
| Program/Line[30]/Stmt[0]/Array[1]/Dim[1] | NumberLiteral | 5 |
| Program/Line[40] | Line |  |
| Program/Line[40]/Stmt[0] | LetStmt | C |
| Program/Line[40]/Stmt[0]/Index[0] | NumberLiteral | 1 |
| Program/Line[40]/Stmt[0]/Value | NumberLiteral | 2 |
| Program/Line[50] | Line |  |
| Program/Line[50]/Stmt[0] | LetStmt | D$ |
| Program/Line[50]/Stmt[0]/Index[0] | NumberLiteral | 1 |
| Program/Line[50]/Stmt[0]/Index[1] | NumberLiteral | 2 |
| Program/Line[50]/Stmt[0]/Value | StringLiteral | WORLD |
| Program/Line[60] | Line |  |
| Program/Line[60]/Stmt[0] | ClearStmt |  |
`

	testutils.Equal(t, "AST markdown for variables statements", got, want)
}

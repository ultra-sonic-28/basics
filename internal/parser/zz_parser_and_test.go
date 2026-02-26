package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_AND_Expression(t *testing.T) {
	tests := []struct {
		name   string
		source string
		assert func(t *testing.T, stmt Statement)
	}{
		{
			name:   "Simple bitwise AND",
			source: "10 PRINT 8 AND 4",
			assert: func(t *testing.T, stmt Statement) {
				ps, ok := stmt.(*PrintStmt)
				testutils.True(t, "is PrintStmt", ok)
				testutils.Equal(t, "expr count", len(ps.Exprs), 1)
				
				infix, ok := ps.Exprs[0].(*InfixExpr)
				testutils.True(t, "is InfixExpr", ok)
				testutils.Equal(t, "op", infix.Op, "AND")
				
				left := infix.Left.(*NumberLiteral)
				testutils.Equal(t, "left", left.Value, 8.0)
				
				right := infix.Right.(*NumberLiteral)
				testutils.Equal(t, "right", right.Value, 4.0)
			},
		},
		{
			name:   "Logical AND with comparisons",
			source: "10 IF A > 0 AND B < 5 THEN PRINT",
			assert: func(t *testing.T, stmt Statement) {
				ifStmt, ok := stmt.(*IfStmt)
				testutils.True(t, "is IfStmt", ok)
				
				infix, ok := ifStmt.Cond.(*InfixExpr)
				testutils.True(t, "is InfixExpr", ok)
				testutils.Equal(t, "op", infix.Op, "AND")
				
				left, ok := infix.Left.(*InfixExpr)
				testutils.True(t, "left is comparison", ok)
				testutils.Equal(t, "left op", left.Op, ">")
				
				right, ok := infix.Right.(*InfixExpr)
				testutils.True(t, "right is comparison", ok)
				testutils.Equal(t, "right op", right.Op, "<")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexer.Lex(tt.source)
			p := New(tokens)
			prog, errs := p.ParseProgram()
			testutils.Equal(t, "no errors", len(errs), 0)
			testutils.Equal(t, "one line", len(prog.Lines), 1)
			
			tt.assert(t, prog.Lines[0].Stmts[0])
		})
	}
}

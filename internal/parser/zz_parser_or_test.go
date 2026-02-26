package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_OR_Expression(t *testing.T) {
	tests := []struct {
		name   string
		source string
		assert func(t *testing.T, stmt Statement)
	}{
		{
			name:   "Simple bitwise OR",
			source: "10 PRINT 8 OR 4",
			assert: func(t *testing.T, stmt Statement) {
				ps, ok := stmt.(*PrintStmt)
				testutils.True(t, "is PrintStmt", ok)
				testutils.Equal(t, "expr count", len(ps.Exprs), 1)
				
				infix, ok := ps.Exprs[0].(*InfixExpr)
				testutils.True(t, "is InfixExpr", ok)
				testutils.Equal(t, "op", infix.Op, "OR")
				
				left := infix.Left.(*NumberLiteral)
				testutils.Equal(t, "left", left.Value, 8.0)
				
				right := infix.Right.(*NumberLiteral)
				testutils.Equal(t, "right", right.Value, 4.0)
			},
		},
		{
			name:   "Logical OR with comparisons",
			source: "10 IF A > 0 OR B < 5 THEN PRINT",
			assert: func(t *testing.T, stmt Statement) {
				ifStmt, ok := stmt.(*IfStmt)
				testutils.True(t, "is IfStmt", ok)
				
				infix, ok := ifStmt.Cond.(*InfixExpr)
				testutils.True(t, "is InfixExpr", ok)
				testutils.Equal(t, "op", infix.Op, "OR")
				
				left, ok := infix.Left.(*InfixExpr)
				testutils.True(t, "left is comparison", ok)
				testutils.Equal(t, "left op", left.Op, ">")
				
				right, ok := infix.Right.(*InfixExpr)
				testutils.True(t, "right is comparison", ok)
				testutils.Equal(t, "right op", right.Op, "<")
			},
		},
		{
			name:   "Precedence: AND higher than OR",
			source: "10 IF A OR B AND C THEN PRINT",
			assert: func(t *testing.T, stmt Statement) {
				// Expected: A OR (B AND C)
				ifStmt := stmt.(*IfStmt)
				orExpr := ifStmt.Cond.(*InfixExpr)
				testutils.Equal(t, "root op is OR", orExpr.Op, "OR")
				
				leftIdent := orExpr.Left.(*Identifier)
				testutils.Equal(t, "left ident is A", leftIdent.Name, "A")
				
				andExpr := orExpr.Right.(*InfixExpr)
				testutils.Equal(t, "right expr is AND", andExpr.Op, "AND")
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

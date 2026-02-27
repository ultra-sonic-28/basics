package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_NOT_Expression(t *testing.T) {
	tests := []struct {
		name   string
		source string
		assert func(t *testing.T, stmt Statement)
	}{
		{
			name:   "Simple NOT",
			source: "10 PRINT NOT 0",
			assert: func(t *testing.T, stmt Statement) {
				ps := stmt.(*PrintStmt)
				prefix := ps.Exprs[0].(*PrefixExpr)
				testutils.Equal(t, "op", prefix.Op, "NOT")
				testutils.Equal(t, "val", prefix.Right.(*NumberLiteral).Value, 0.0)
			},
		},
		{
			name:   "NOT precedence over SUM",
			source: "10 PRINT NOT 4 - 4",
			assert: func(t *testing.T, stmt Statement) {
				// (NOT 4) - 4
				ps := stmt.(*PrintStmt)
				infix := ps.Exprs[0].(*InfixExpr)
				testutils.Equal(t, "root op", infix.Op, "-")
				
				prefix := infix.Left.(*PrefixExpr)
				testutils.Equal(t, "left prefix op", prefix.Op, "NOT")
				testutils.Equal(t, "prefix operand", prefix.Right.(*NumberLiteral).Value, 4.0)
				
				right := infix.Right.(*NumberLiteral)
				testutils.Equal(t, "right operand", right.Value, 4.0)
			},
		},
		{
			name:   "NOT with parentheses",
			source: "10 PRINT NOT (4 - 4)",
			assert: func(t *testing.T, stmt Statement) {
				ps := stmt.(*PrintStmt)
				prefix := ps.Exprs[0].(*PrefixExpr)
				testutils.Equal(t, "op", prefix.Op, "NOT")
				
				infix := prefix.Right.(*InfixExpr)
				testutils.Equal(t, "inner op", infix.Op, "-")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexer.Lex(tt.source)
			p := New(tokens)
			prog, errs := p.ParseProgram()
			testutils.Equal(t, "no errors", len(errs), 0)
			tt.assert(t, prog.Lines[0].Stmts[0])
		})
	}
}

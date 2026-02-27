package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_WAIT_Statement(t *testing.T) {
	tests := []struct {
		name   string
		source string
		assert func(t *testing.T, stmt Statement)
	}{
		{
			name:   "Simple WAIT",
			source: "10 WAIT 1000",
			assert: func(t *testing.T, stmt Statement) {
				ws, ok := stmt.(*WaitStmt)
				testutils.True(t, "is WaitStmt", ok)
				num, ok := ws.Expr.(*NumberLiteral)
				testutils.True(t, "expr is NumberLiteral", ok)
				testutils.Equal(t, "value", num.Value, 1000.0)
			},
		},
		{
			name:   "WAIT with expression",
			source: "10 WAIT DELAI * 5",
			assert: func(t *testing.T, stmt Statement) {
				ws, ok := stmt.(*WaitStmt)
				testutils.True(t, "is WaitStmt", ok)
				infix, ok := ws.Expr.(*InfixExpr)
				testutils.True(t, "expr is InfixExpr", ok)
				testutils.Equal(t, "op", infix.Op, "*")
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

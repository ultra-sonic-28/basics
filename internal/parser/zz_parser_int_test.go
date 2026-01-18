package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_INT_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "INT with positive float",
			source:    `10 PRINT INT(1.75)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt, ok := prog.Lines[0].Stmts[0].(*PrintStmt)
				testutils.True(t, "is PrintStmt", ok)
				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				intExpr, ok := printStmt.Exprs[0].(*IntExpr)
				testutils.True(t, "expression is IntExpr", ok)

				num, ok := intExpr.Expr.(*NumberLiteral)
				testutils.True(t, "INT argument is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 1.75)
			},
		},
		{
			name:      "INT with negative float",
			source:    `10 PRINT INT(-1.32)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				intExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*IntExpr)

				prefix, ok := intExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				num, ok := prefix.Right.(*NumberLiteral)
				testutils.True(t, "prefix contains NumberLiteral", ok)
				testutils.Equal(t, "value", num.Value, 1.32)
			},
		},
		{
			name:      "INT with float variable",
			source:    `10 PRINT INT(A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				intExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*IntExpr)

				_, ok := intExpr.Expr.(*Identifier)
				testutils.True(t, "INT arg is Identifier", ok)
			},
		},
		{
			name:      "INT with expression A * 3.74",
			source:    `10 PRINT INT(A * 3.74)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				intExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*IntExpr)

				infix, ok := intExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*Identifier)
				testutils.True(t, "left is Identifier", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right is NumberLiteral", ok)
			},
		},
		{
			name:      "INT with integer variable",
			source:    `10 PRINT INT(I%)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				intExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*IntExpr)

				id, ok := intExpr.Expr.(*Identifier)
				testutils.True(t, "arg is Identifier", ok)
				testutils.Equal(t, "identifier name", id.Name, "I%")
			},
		},
		{
			name:      "INT with mixed expression I% * A",
			source:    `10 PRINT INT(I% * A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				intExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*IntExpr)

				infix, ok := intExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*Identifier)
				testutils.True(t, "left is Identifier", ok)

				_, ok = infix.Right.(*Identifier)
				testutils.True(t, "right is Identifier", ok)
			},
		},
		{
			name:      "INT with nested expression -(A + 3.2)",
			source:    `10 PRINT INT(-(A + 3.2))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				intExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*IntExpr)

				prefix, ok := intExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				infix, ok := prefix.Right.(*InfixExpr)
				testutils.True(t, "prefix contains InfixExpr", ok)

				_, ok = infix.Left.(*Identifier)
				testutils.True(t, "left is Identifier", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right is NumberLiteral", ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexer.Lex(tt.source)
			p := New(tokens)
			prog, errs := p.ParseProgram()

			testutils.Equal(t, "no parser errors", len(errs), 0)
			testutils.Equal(t, "line count", len(prog.Lines), tt.lineCount)

			tt.assertFn(t, prog)
		})
	}
}

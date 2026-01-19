package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_SGN_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "SGN with positive integer",
			source:    `10 PRINT SGN(3)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt, ok := prog.Lines[0].Stmts[0].(*PrintStmt)
				testutils.True(t, "is PrintStmt", ok)
				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				sgnExpr, ok := printStmt.Exprs[0].(*SgnExpr)
				testutils.True(t, "expression is SgnExpr", ok)

				num, ok := sgnExpr.Expr.(*NumberLiteral)
				testutils.True(t, "SGN argument is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 3)
			},
		},
		{
			name:      "SGN with negative integer",
			source:    `10 PRINT SGN(-3)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				sgnExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SgnExpr)

				prefix, ok := sgnExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				num, ok := prefix.Right.(*NumberLiteral)
				testutils.True(t, "prefix contains NumberLiteral", ok)
				testutils.Equal(t, "value", num.Value, 3)
			},
		},
		{
			name:      "SGN with positive float",
			source:    `10 PRINT SGN(1.75)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt, ok := prog.Lines[0].Stmts[0].(*PrintStmt)
				testutils.True(t, "is PrintStmt", ok)
				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				sgnExpr, ok := printStmt.Exprs[0].(*SgnExpr)
				testutils.True(t, "expression is SgnExpr", ok)

				num, ok := sgnExpr.Expr.(*NumberLiteral)
				testutils.True(t, "SGN argument is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 1.75)
			},
		},
		{
			name:      "SGN with negative float",
			source:    `10 PRINT SGN(-1.75)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				sgnExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SgnExpr)

				prefix, ok := sgnExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				num, ok := prefix.Right.(*NumberLiteral)
				testutils.True(t, "prefix contains NumberLiteral", ok)
				testutils.Equal(t, "value", num.Value, 1.75)
			},
		},
		{
			name:      "SGN with float variable",
			source:    `10 PRINT SGN(A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				sgnExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SgnExpr)

				_, ok := sgnExpr.Expr.(*Identifier)
				testutils.True(t, "SGN arg is Identifier", ok)
			},
		},
		{
			name:      "SGN with expression A * 3.74",
			source:    `10 PRINT SGN(A * 3.74)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				sgnExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SgnExpr)

				infix, ok := sgnExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*Identifier)
				testutils.True(t, "left is Identifier", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right is NumberLiteral", ok)
			},
		},
		{
			name:      "SGN with integer variable",
			source:    `10 PRINT SGN(I%)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				sgnExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SgnExpr)

				id, ok := sgnExpr.Expr.(*Identifier)
				testutils.True(t, "arg is Identifier", ok)
				testutils.Equal(t, "identifier name", id.Name, "I%")
			},
		},
		{
			name:      "SGN with negative integer variable",
			source:    `10 PRINT SGN(-I%)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				sgnExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SgnExpr)

				prefix, ok := sgnExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				id, ok := prefix.Right.(*Identifier)
				testutils.True(t, "prefix contains Identifier", ok)
				testutils.Equal(t, "identifier name", id.Name, "I%")
			},
		},
		{
			name:      "SGN with mixed expression I% * A",
			source:    `10 PRINT SGN(I% * A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				sgnExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SgnExpr)

				infix, ok := sgnExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*Identifier)
				testutils.True(t, "left is Identifier", ok)

				_, ok = infix.Right.(*Identifier)
				testutils.True(t, "right is Identifier", ok)
			},
		},
		{
			name:      "SGN with nested expression -(I% * A)",
			source:    `10 PRINT SGN(-(I% * A))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				sgnExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SgnExpr)

				prefix, ok := sgnExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				infix, ok := prefix.Right.(*InfixExpr)
				testutils.True(t, "prefix contains InfixExpr", ok)

				_, ok = infix.Left.(*Identifier)
				testutils.True(t, "left is Identifier", ok)

				_, ok = infix.Right.(*Identifier)
				testutils.True(t, "right is Identifier", ok)
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

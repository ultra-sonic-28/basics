package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_SPC_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "SPC with positive integer",
			source:    `10 PRINT SPC(25)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt, ok := prog.Lines[0].Stmts[0].(*PrintStmt)
				testutils.True(t, "is PrintStmt", ok)
				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				spcExpr, ok := printStmt.Exprs[0].(*SpcExpr)
				testutils.True(t, "expression is SpcExpr", ok)

				num, ok := spcExpr.Expr.(*NumberLiteral)
				testutils.True(t, "SPC argument is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 25)
			},
		},
		{
			name:      "SPC with variable",
			source:    `10 PRINT SPC(A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				spcExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SpcExpr)

				id, ok := spcExpr.Expr.(*Identifier)
				testutils.True(t, "SPC arg is Identifier", ok)
				testutils.Equal(t, "identifier name", id.Name, "A")
			},
		},
		{
			name:      "SPC with expression A*2",
			source:    `10 PRINT SPC(A*2)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				spcExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SpcExpr)

				infix, ok := spcExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*Identifier)
				testutils.True(t, "left is Identifier", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right is NumberLiteral", ok)
			},
		},
		{
			name:      "SPC with negative value",
			source:    `10 PRINT SPC(-3)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				spcExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SpcExpr)

				prefix, ok := spcExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				num, ok := prefix.Right.(*NumberLiteral)
				testutils.True(t, "prefix contains NumberLiteral", ok)
				testutils.Equal(t, "value", num.Value, 3)
			},
		},
		{
			name:      "Multiple SPC in same PRINT",
			source:    `10 PRINT SPC(10);SPC(20);SPC(30)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "three expressions", len(printStmt.Exprs), 3)

				for i := 0; i < 3; i++ {
					_, ok := printStmt.Exprs[i].(*SpcExpr)
					testutils.True(t, "expr is SpcExpr", ok)
				}
			},
		},
		{
			name:      "SPC mixed with string",
			source:    `10 PRINT SPC(15);"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				_, ok := printStmt.Exprs[0].(*SpcExpr)
				testutils.True(t, "first expr is SpcExpr", ok)

				_, ok = printStmt.Exprs[1].(*StringLiteral)
				testutils.True(t, "second expr is StringLiteral", ok)
			},
		},
		{
			name:      "SPC with nested expression -(A*2)",
			source:    `10 PRINT SPC(-(A*2))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				spcExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*SpcExpr)

				prefix, ok := spcExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				_, ok = prefix.Right.(*InfixExpr)
				testutils.True(t, "prefix contains InfixExpr", ok)
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

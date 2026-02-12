package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_TAB_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "TAB with positive integer",
			source:    `10 PRINT TAB(25)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt, ok := prog.Lines[0].Stmts[0].(*PrintStmt)
				testutils.True(t, "is PrintStmt", ok)
				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				tabExpr, ok := printStmt.Exprs[0].(*TabExpr)
				testutils.True(t, "expression is TabExpr", ok)

				num, ok := tabExpr.Expr.(*NumberLiteral)
				testutils.True(t, "TAB argument is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 25)
			},
		},
		{
			name:      "TAB with variable",
			source:    `10 PRINT TAB(A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				tabExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*TabExpr)

				id, ok := tabExpr.Expr.(*Identifier)
				testutils.True(t, "TAB arg is Identifier", ok)
				testutils.Equal(t, "identifier name", id.Name, "A")
			},
		},
		{
			name:      "TAB with expression A*2",
			source:    `10 PRINT TAB(A*2)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				tabExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*TabExpr)

				infix, ok := tabExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*Identifier)
				testutils.True(t, "left is Identifier", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right is NumberLiteral", ok)
			},
		},
		{
			name:      "TAB with negative value",
			source:    `10 PRINT TAB(-3)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				tabExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*TabExpr)

				prefix, ok := tabExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				num, ok := prefix.Right.(*NumberLiteral)
				testutils.True(t, "prefix contains NumberLiteral", ok)
				testutils.Equal(t, "value", num.Value, 3)
			},
		},
		{
			name:      "Multiple TAB in same PRINT",
			source:    `10 PRINT TAB(10);TAB(20);TAB(30)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "three expressions", len(printStmt.Exprs), 3)

				for i := 0; i < 3; i++ {
					_, ok := printStmt.Exprs[i].(*TabExpr)
					testutils.True(t, "expr is TabExpr", ok)
				}
			},
		},
		{
			name:      "TAB mixed with string",
			source:    `10 PRINT TAB(15);"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				_, ok := printStmt.Exprs[0].(*TabExpr)
				testutils.True(t, "first expr is TabExpr", ok)

				_, ok = printStmt.Exprs[1].(*StringLiteral)
				testutils.True(t, "second expr is StringLiteral", ok)
			},
		},
		{
			name:      "TAB with nested expression -(A*2)",
			source:    `10 PRINT TAB(-(A*2))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				tabExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*TabExpr)

				prefix, ok := tabExpr.Expr.(*PrefixExpr)
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

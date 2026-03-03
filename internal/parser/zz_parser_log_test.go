package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_LOG_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "LOG with positive integer",
			source:    `10 PRINT LOG(25)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt, ok := prog.Lines[0].Stmts[0].(*PrintStmt)
				testutils.True(t, "is PrintStmt", ok)
				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				logExpr, ok := printStmt.Exprs[0].(*LogExpr)
				testutils.True(t, "expression is logExpr", ok)

				num, ok := logExpr.Expr.(*NumberLiteral)
				testutils.True(t, "LOG argument is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 25)
			},
		},
		{
			name:      "LOG with variable",
			source:    `10 PRINT LOG(A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				logExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*LogExpr)

				id, ok := logExpr.Expr.(*Identifier)
				testutils.True(t, "LOG arg is Identifier", ok)
				testutils.Equal(t, "identifier name", id.Name, "A")
			},
		},
		{
			name:      "LOG with expression A*2",
			source:    `10 PRINT LOG(A*2)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				logExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*LogExpr)

				infix, ok := logExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*Identifier)
				testutils.True(t, "left is Identifier", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right is NumberLiteral", ok)
			},
		},
		{
			name:      "LOG with negative value",
			source:    `10 PRINT LOG(-3)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				logExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*LogExpr)

				prefix, ok := logExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				num, ok := prefix.Right.(*NumberLiteral)
				testutils.True(t, "prefix contains NumberLiteral", ok)
				testutils.Equal(t, "value", num.Value, 3)
			},
		},
		{
			name:      "Multiple LOG in same PRINT",
			source:    `10 PRINT LOG(10);LOG(20);LOG(30)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "three expressions", len(printStmt.Exprs), 3)

				for i := 0; i < 3; i++ {
					_, ok := printStmt.Exprs[i].(*LogExpr)
					testutils.True(t, "expr is logExpr", ok)
				}
			},
		},
		{
			name:      "LOG mixed with string",
			source:    `10 PRINT LOG(15);"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				_, ok := printStmt.Exprs[0].(*LogExpr)
				testutils.True(t, "first expr is logExpr", ok)

				_, ok = printStmt.Exprs[1].(*StringLiteral)
				testutils.True(t, "second expr is StringLiteral", ok)
			},
		},
		{
			name:      "LOG with nested expression -(A*2)",
			source:    `10 PRINT LOG(-(A*2))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				logExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*LogExpr)

				prefix, ok := logExpr.Expr.(*PrefixExpr)
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

package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_RND_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "RND with positive integer",
			source:    `10 PRINT RND(25)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt, ok := prog.Lines[0].Stmts[0].(*PrintStmt)
				testutils.True(t, "is PrintStmt", ok)
				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				rndExpr, ok := printStmt.Exprs[0].(*RndExpr)
				testutils.True(t, "expression is rndExpr", ok)

				num, ok := rndExpr.Expr.(*NumberLiteral)
				testutils.True(t, "RND argument is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 25)
			},
		},
		{
			name:      "RND with variable",
			source:    `10 PRINT RND(A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				rndExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*RndExpr)

				id, ok := rndExpr.Expr.(*Identifier)
				testutils.True(t, "RND arg is Identifier", ok)
				testutils.Equal(t, "identifier name", id.Name, "A")
			},
		},
		{
			name:      "RND with expression A*2",
			source:    `10 PRINT RND(A*2)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				rndExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*RndExpr)

				infix, ok := rndExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*Identifier)
				testutils.True(t, "left is Identifier", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right is NumberLiteral", ok)
			},
		},
		{
			name:      "RND with negative value",
			source:    `10 PRINT RND(-3)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				rndExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*RndExpr)

				prefix, ok := rndExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				num, ok := prefix.Right.(*NumberLiteral)
				testutils.True(t, "prefix contains NumberLiteral", ok)
				testutils.Equal(t, "value", num.Value, 3)
			},
		},
		{
			name:      "Multiple RND in same PRINT",
			source:    `10 PRINT RND(10);RND(20);RND(30)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "three expressions", len(printStmt.Exprs), 3)

				for i := 0; i < 3; i++ {
					_, ok := printStmt.Exprs[i].(*RndExpr)
					testutils.True(t, "expr is rndExpr", ok)
				}
			},
		},
		{
			name:      "RND mixed with string",
			source:    `10 PRINT RND(15);"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				_, ok := printStmt.Exprs[0].(*RndExpr)
				testutils.True(t, "first expr is rndExpr", ok)

				_, ok = printStmt.Exprs[1].(*StringLiteral)
				testutils.True(t, "second expr is StringLiteral", ok)
			},
		},
		{
			name:      "RND with nested expression -(A*2)",
			source:    `10 PRINT RND(-(A*2))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				rndExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*RndExpr)

				prefix, ok := rndExpr.Expr.(*PrefixExpr)
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

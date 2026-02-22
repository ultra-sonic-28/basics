package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_STR_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "STR$ with integer",
			source:    `10 PRINT STR$(5)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				strExpr, ok := printStmt.Exprs[0].(*StrExpr)
				testutils.True(t, "expression is strExpr", ok)

				num, ok := strExpr.Expr.(*NumberLiteral)
				testutils.True(t, "arg is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 5)
			},
		},
		{
			name:      "STR$ assigned to variable",
			source:    `10 A$ = STR$(5)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				assign := prog.Lines[0].Stmts[0].(*LetStmt)

				_, ok := assign.Value.(*StrExpr)
				testutils.True(t, "value is strExpr", ok)
			},
		},
		{
			name:      "STR$ with variable",
			source:    `10 PRINT STR$(A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				strExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*StrExpr)

				_, ok := strExpr.Expr.(*Identifier)
				testutils.True(t, "arg is Identifier", ok)
			},
		},
		{
			name:      "STR$ with arithmetic expression",
			source:    `10 PRINT STR$(A*2+1)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				strExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*StrExpr)

				infix, ok := strExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*InfixExpr)
				testutils.True(t, "left side is InfixExpr", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right side is NumberLiteral", ok)
			},
		},
		{
			name:      "STR$ with nested expression -(A*2)",
			source:    `10 PRINT STR$(-(A*2))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				strExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*StrExpr)

				prefix, ok := strExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				_, ok = prefix.Right.(*InfixExpr)
				testutils.True(t, "prefix contains InfixExpr", ok)
			},
		},
		{
			name:      "Multiple STR$ in same PRINT",
			source:    `10 PRINT STR$(2);STR$(1)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				for i := 0; i < 2; i++ {
					_, ok := printStmt.Exprs[i].(*StrExpr)
					testutils.True(t, "expr is strExpr", ok)
				}
			},
		},
		{
			name:      "STR$ mixed with string",
			source:    `10 PRINT STR$(2);"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				_, ok := printStmt.Exprs[0].(*StrExpr)
				testutils.True(t, "first expr is strExpr", ok)

				_, ok = printStmt.Exprs[1].(*StringLiteral)
				testutils.True(t, "second expr is StringLiteral", ok)
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

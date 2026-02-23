package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_VAL_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "VAL with literal string",
			source:    `10 PRINT VAL("123456")`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				ValExpr, ok := printStmt.Exprs[0].(*ValExpr)
				testutils.True(t, "expression is ValExpr", ok)

				str, ok := ValExpr.Expr.(*StringLiteral)
				testutils.True(t, "arg is StringLiteral", ok)
				testutils.Equal(t, "string value", str.Value, "123456")
			},
		},
		{
			name:      "VAL assigned to variable",
			source:    `10 A = VAL("123456")`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				assign := prog.Lines[0].Stmts[0].(*LetStmt)

				_, ok := assign.Value.(*ValExpr)
				testutils.True(t, "value is ValExpr", ok)
			},
		},
		{
			name:      "VAL with arithmetic expression",
			source:    `10 PRINT VAL("123" + "456")`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				ValExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*ValExpr)

				infix, ok := ValExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*StringLiteral)
				testutils.True(t, "left side is StringLiteral", ok)

				_, ok = infix.Right.(*StringLiteral)
				testutils.True(t, "right side is StringLiteral", ok)
			},
		},
		{
			name:      "Multiple VAL in same PRINT",
			source:    `10 PRINT VAL("ABCDE");VAL("XYZ")`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				for i := 0; i < 2; i++ {
					_, ok := printStmt.Exprs[i].(*ValExpr)
					testutils.True(t, "expr is ValExpr", ok)
				}
			},
		},
		{
			name:      "VAL mixed with string",
			source:    `10 PRINT VAL("123456");"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				_, ok := printStmt.Exprs[0].(*ValExpr)
				testutils.True(t, "first expr is ValExpr", ok)

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

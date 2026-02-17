package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_LEN_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "LEN with literal string",
			source:    `10 PRINT LEN("APPLESOFT")`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				LenExpr, ok := printStmt.Exprs[0].(*LenExpr)
				testutils.True(t, "expression is LenExpr", ok)

				str, ok := LenExpr.Expr.(*StringLiteral)
				testutils.True(t, "arg is StringLiteral", ok)
				testutils.Equal(t, "string value", str.Value, "APPLESOFT")
			},
		},
		{
			name:      "LEN assigned to variable",
			source:    `10 A = LEN("APPLESOFT")`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				assign := prog.Lines[0].Stmts[0].(*LetStmt)

				_, ok := assign.Value.(*LenExpr)
				testutils.True(t, "value is LenExpr", ok)
			},
		},
		{
			name:      "LEN with arithmetic expression",
			source:    `10 PRINT LEN("APPLESOFT" + "WARE")`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				LenExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*LenExpr)

				infix, ok := LenExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*StringLiteral)
				testutils.True(t, "left side is StringLiteral", ok)

				_, ok = infix.Right.(*StringLiteral)
				testutils.True(t, "right side is StringLiteral", ok)
			},
		},
		{
			name:      "Multiple LEN in same PRINT",
			source:    `10 PRINT LEN("ABCDE");LEN("XYZ")`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				for i := 0; i < 2; i++ {
					_, ok := printStmt.Exprs[i].(*LenExpr)
					testutils.True(t, "expr is LenExpr", ok)
				}
			},
		},
		{
			name:      "LEN mixed with string",
			source:    `10 PRINT LEN("ABCDE");"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				_, ok := printStmt.Exprs[0].(*LenExpr)
				testutils.True(t, "first expr is LenExpr", ok)

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

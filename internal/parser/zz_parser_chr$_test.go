package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_CHR_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "CHR$ with integer",
			source:    `10 PRINT CHR$(35)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				chrExpr, ok := printStmt.Exprs[0].(*ChrExpr)
				testutils.True(t, "expression is chrExpr", ok)

				num, ok := chrExpr.Expr.(*NumberLiteral)
				testutils.True(t, "arg is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 35)
			},
		},
		{
			name:      "CHR$ assigned to variable",
			source:    `10 A$ = CHR$(35)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				assign := prog.Lines[0].Stmts[0].(*LetStmt)

				_, ok := assign.Value.(*ChrExpr)
				testutils.True(t, "value is chrExpr", ok)
			},
		},
		{
			name:      "CHR$ with variable",
			source:    `10 PRINT CHR$(A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				chrExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*ChrExpr)

				_, ok := chrExpr.Expr.(*Identifier)
				testutils.True(t, "arg is Identifier", ok)
			},
		},
		{
			name:      "CHR$ with arithmetic expression",
			source:    `10 PRINT CHR$(A*2+1)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				chrExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*ChrExpr)

				infix, ok := chrExpr.Expr.(*InfixExpr)
				testutils.True(t, "arg is InfixExpr", ok)

				_, ok = infix.Left.(*InfixExpr)
				testutils.True(t, "left side is InfixExpr", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right side is NumberLiteral", ok)
			},
		},
		{
			name:      "CHR$ with nested expression -(A*2)",
			source:    `10 PRINT CHR$(-(A*2))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				chrExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*ChrExpr)

				prefix, ok := chrExpr.Expr.(*PrefixExpr)
				testutils.True(t, "arg is PrefixExpr", ok)

				_, ok = prefix.Right.(*InfixExpr)
				testutils.True(t, "prefix contains InfixExpr", ok)
			},
		},
		{
			name:      "Multiple CHR$ in same PRINT",
			source:    `10 PRINT CHR$(2);CHR$(1)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				for i := 0; i < 2; i++ {
					_, ok := printStmt.Exprs[i].(*ChrExpr)
					testutils.True(t, "expr is chrExpr", ok)
				}
			},
		},
		{
			name:      "CHR$ mixed with string",
			source:    `10 PRINT CHR$(2);"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				_, ok := printStmt.Exprs[0].(*ChrExpr)
				testutils.True(t, "first expr is chrExpr", ok)

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

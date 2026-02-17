package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_RIGHT_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "RIGHT$ with literal string and integer",
			source:    `10 PRINT RIGHT$("APPLESOFT",5)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				RightExpr, ok := printStmt.Exprs[0].(*RightExpr)
				testutils.True(t, "expression is RightExpr", ok)

				_, ok = RightExpr.StrExpr.(*StringLiteral)
				testutils.True(t, "first arg is StringLiteral", ok)

				num, ok := RightExpr.LenExpr.(*NumberLiteral)
				testutils.True(t, "second arg is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 5)
			},
		},
		{
			name:      "RIGHT$ assigned to variable",
			source:    `10 A$ = RIGHT$("APPLESOFT",5)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				assign := prog.Lines[0].Stmts[0].(*LetStmt)

				_, ok := assign.Value.(*RightExpr)
				testutils.True(t, "value is RightExpr", ok)
			},
		},
		{
			name:      "RIGHT$ with variable count",
			source:    `10 PRINT RIGHT$("APPLESOFT",A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				RightExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*RightExpr)

				_, ok := RightExpr.LenExpr.(*Identifier)
				testutils.True(t, "second arg is Identifier", ok)
			},
		},
		{
			name:      "RIGHT$ with arithmetic expression",
			source:    `10 PRINT RIGHT$("APPLESOFT",A*2+1)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				RightExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*RightExpr)

				infix, ok := RightExpr.LenExpr.(*InfixExpr)
				testutils.True(t, "second arg is InfixExpr", ok)

				_, ok = infix.Left.(*InfixExpr)
				testutils.True(t, "left side is InfixExpr", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right side is NumberLiteral", ok)
			},
		},
		{
			name:      "RIGHT$ with nested expression -(A*2)",
			source:    `10 PRINT RIGHT$("APPLESOFT",-(A*2))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				RightExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*RightExpr)

				prefix, ok := RightExpr.LenExpr.(*PrefixExpr)
				testutils.True(t, "second arg is PrefixExpr", ok)

				_, ok = prefix.Right.(*InfixExpr)
				testutils.True(t, "prefix contains InfixExpr", ok)
			},
		},
		{
			name:      "Multiple RIGHT$ in same PRINT",
			source:    `10 PRINT RIGHT$("ABCDE",2);RIGHT$("XYZ",1)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				for i := 0; i < 2; i++ {
					_, ok := printStmt.Exprs[i].(*RightExpr)
					testutils.True(t, "expr is RightExpr", ok)
				}
			},
		},
		{
			name:      "RIGHT$ mixed with string",
			source:    `10 PRINT RIGHT$("ABCDE",2);"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				_, ok := printStmt.Exprs[0].(*RightExpr)
				testutils.True(t, "first expr is RightExpr", ok)

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

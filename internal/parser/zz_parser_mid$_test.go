package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_MID_Expressions(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "MID$ with literal string and integer 2 Args",
			source:    `10 PRINT MID$("APPLESOFT",5)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				MidExpr, ok := printStmt.Exprs[0].(*MidExpr)
				testutils.True(t, "expression is MidExpr", ok)

				_, ok = MidExpr.StrExpr.(*StringLiteral)
				testutils.True(t, "first arg is StringLiteral", ok)

				num, ok := MidExpr.Start.(*NumberLiteral)
				testutils.True(t, "second arg is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 5)

				// Pas de paramètre Len, donc Len doit être nil
				testutils.True(t, "len arg is nil", MidExpr.Len == nil)
			},
		},
		{
			name:      "MID$ with literal string and integer 3 Args",
			source:    `10 PRINT MID$("APPLESOFT",5, 2)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "one expression", len(printStmt.Exprs), 1)

				MidExpr, ok := printStmt.Exprs[0].(*MidExpr)
				testutils.True(t, "expression is MidExpr", ok)

				_, ok = MidExpr.StrExpr.(*StringLiteral)
				testutils.True(t, "first arg is StringLiteral", ok)

				num, ok := MidExpr.Start.(*NumberLiteral)
				testutils.True(t, "second arg is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 5)

				num, ok = MidExpr.Len.(*NumberLiteral)
				testutils.True(t, "third arg is NumberLiteral", ok)
				testutils.Equal(t, "number value", num.Value, 2)
			},
		},
		{
			name:      "MID$ 2 Args assigned to variable",
			source:    `10 A$ = MID$("APPLESOFT",5)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				assign := prog.Lines[0].Stmts[0].(*LetStmt)

				MidExpr, ok := assign.Value.(*MidExpr)
				testutils.True(t, "value is MidExpr", ok)

				// Pas de paramètre Len, donc Len doit être nil
				testutils.True(t, "len arg is nil", MidExpr.Len == nil)
			},
		},
		{
			name:      "MID$ 3 Args assigned to variable",
			source:    `10 A$ = MID$("APPLESOFT",5,2)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				assign := prog.Lines[0].Stmts[0].(*LetStmt)

				MidExpr, ok := assign.Value.(*MidExpr)
				testutils.True(t, "value is MidExpr", ok)

				// Len passé en paramètre, donc Len ne doit pas être nil
				testutils.True(t, "len arg is nil", MidExpr.Len != nil)
			},
		},
		{
			name:      "MID$ with variable 2 Args",
			source:    `10 PRINT MID$("APPLESOFT",A)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				MidExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*MidExpr)

				_, ok := MidExpr.Start.(*Identifier)
				testutils.True(t, "second arg is Identifier", ok)

				// Pas de paramètre Len, donc Len doit être nil
				testutils.True(t, "len arg is nil", MidExpr.Len == nil)
			},
		},
		{
			name:      "MID$ with variable 3 Args",
			source:    `10 PRINT MID$("APPLESOFT",A, B)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				MidExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*MidExpr)

				_, ok := MidExpr.Start.(*Identifier)
				testutils.True(t, "second arg is Identifier", ok)

				_, ok = MidExpr.Len.(*Identifier)
				testutils.True(t, "thrid arg is Identifier", ok)
			},
		},
		{
			name:      "MID$ with arithmetic expression 2 Args",
			source:    `10 PRINT MID$("APPLESOFT",A*2+1)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				MidExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*MidExpr)

				infix, ok := MidExpr.Start.(*InfixExpr)
				testutils.True(t, "second arg is InfixExpr", ok)

				_, ok = infix.Left.(*InfixExpr)
				testutils.True(t, "left side is InfixExpr", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right side is NumberLiteral", ok)

				// Pas de paramètre Len, donc Len doit être nil
				testutils.True(t, "len arg is nil", MidExpr.Len == nil)
			},
		},
		{
			name:      "MID$ with arithmetic expression 3 Args",
			source:    `10 PRINT MID$("APPLESOFT",A*2+1, B*2+1)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				MidExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*MidExpr)

				infix, ok := MidExpr.Start.(*InfixExpr)
				testutils.True(t, "second arg is InfixExpr", ok)

				_, ok = infix.Left.(*InfixExpr)
				testutils.True(t, "left side is InfixExpr", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right side is NumberLiteral", ok)

				infix, ok = MidExpr.Len.(*InfixExpr)
				testutils.True(t, "third arg is InfixExpr", ok)

				_, ok = infix.Left.(*InfixExpr)
				testutils.True(t, "left side is InfixExpr", ok)

				_, ok = infix.Right.(*NumberLiteral)
				testutils.True(t, "right side is NumberLiteral", ok)
			},
		},
		{
			name:      "MID$ with nested expression -(A*2) 2 Args",
			source:    `10 PRINT MID$("APPLESOFT",-(A*2))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				MidExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*MidExpr)

				prefix, ok := MidExpr.Start.(*PrefixExpr)
				testutils.True(t, "second arg is PrefixExpr", ok)

				_, ok = prefix.Right.(*InfixExpr)
				testutils.True(t, "prefix contains InfixExpr", ok)

				// Pas de paramètre Len, donc Len doit être nil
				testutils.True(t, "len arg is nil", MidExpr.Len == nil)
			},
		},
		{
			name:      "MID$ with nested expression -(A*2) 3 Args",
			source:    `10 PRINT MID$("APPLESOFT",-(A*2), -(B*2))`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				MidExpr := prog.Lines[0].Stmts[0].(*PrintStmt).Exprs[0].(*MidExpr)

				prefix, ok := MidExpr.Start.(*PrefixExpr)
				testutils.True(t, "second arg is PrefixExpr", ok)

				_, ok = prefix.Right.(*InfixExpr)
				testutils.True(t, "prefix contains InfixExpr", ok)

				prefix, ok = MidExpr.Len.(*PrefixExpr)
				testutils.True(t, "third arg is PrefixExpr", ok)

				_, ok = prefix.Right.(*InfixExpr)
				testutils.True(t, "prefix contains InfixExpr", ok)
			},
		},
		{
			name:      "Multiple MID$ in same PRINT 2 Args",
			source:    `10 PRINT MID$("ABCDE",2);MID$("XYZ",1)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				for i := 0; i < 2; i++ {
					MidExpr, ok := printStmt.Exprs[i].(*MidExpr)
					testutils.True(t, "expr is MidExpr", ok)

					// Pas de paramètre Len, donc Len doit être nil
					testutils.True(t, "len arg is nil", MidExpr.Len == nil)
				}
			},
		},
		{
			name:      "Multiple MID$ in same PRINT 3 Args",
			source:    `10 PRINT MID$("ABCDE",2,1);MID$("XYZ",1,1)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				for i := 0; i < 2; i++ {
					MidExpr, ok := printStmt.Exprs[i].(*MidExpr)
					testutils.True(t, "expr is MidExpr", ok)

					// Len passé en paramètre, donc Len ne doit pas être nil
					testutils.True(t, "len arg is nil", MidExpr.Len != nil)
				}
			},
		},
		{
			name:      "MID$ mixed with string 2 Args",
			source:    `10 PRINT MID$("ABCDE",2);"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				MidExpr, ok := printStmt.Exprs[0].(*MidExpr)
				testutils.True(t, "first expr is MidExpr", ok)

				_, ok = printStmt.Exprs[1].(*StringLiteral)
				testutils.True(t, "second expr is StringLiteral", ok)

				// Pas de paramètre Len, donc Len doit être nil
				testutils.True(t, "len arg is nil", MidExpr.Len == nil)
			},
		},
		{
			name:      "MID$ mixed with string 3 Args",
			source:    `10 PRINT MID$("ABCDE",2,2);"HELLO"`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				printStmt := prog.Lines[0].Stmts[0].(*PrintStmt)

				testutils.Equal(t, "two expressions", len(printStmt.Exprs), 2)

				MidExpr, ok := printStmt.Exprs[0].(*MidExpr)
				testutils.True(t, "first expr is MidExpr", ok)

				_, ok = printStmt.Exprs[1].(*StringLiteral)
				testutils.True(t, "second expr is StringLiteral", ok)

				// Len passé en paramètre, donc Len ne doit pas être nil
				testutils.True(t, "len arg is nil", MidExpr.Len != nil)
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

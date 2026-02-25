package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_COLOR_Statement(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		lineCount  int
		assertFunc func(t *testing.T, prog *Program)
	}{
		{
			name: "COLOR same line with PRINT",
			source: `
10 COLOR = 15: PRINT "HELLO"
`,
			lineCount: 1,
			assertFunc: func(t *testing.T, prog *Program) {
				line := prog.Lines[0]
				testutils.Equal(t, "line number", line.Number, 10)
				testutils.Equal(t, "two statements", len(line.Stmts), 2)

				c, ok := line.Stmts[0].(*ColorStmt)
				testutils.Equal(t, "first stmt is ColorStmt", ok, true)
				testutils.Equal(t, "COLOR value", c.Expr.(*NumberLiteral).Value, 15.0)

				_, ok = line.Stmts[1].(*PrintStmt)
				testutils.Equal(t, "second stmt is PrintStmt", ok, true)
			},
		},
		{
			name: "COLOR on its own line",
			source: `
10 COLOR = 1
30 PRINT "HELLO"
`,
			lineCount: 2,
			assertFunc: func(t *testing.T, prog *Program) {
				line10 := prog.Lines[0]
				testutils.Equal(t, "line 10 number", line10.Number, 10)
				testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 1)

				c, ok := line10.Stmts[0].(*ColorStmt)
				testutils.Equal(t, "line 10 is ColorStmt", ok, true)
				testutils.Equal(t, "COLOR value", c.Expr.(*NumberLiteral).Value, 1.0)

				line30 := prog.Lines[1]
				testutils.Equal(t, "line 30 number", line30.Number, 30)
				testutils.Equal(t, "line 30 stmt count", len(line30.Stmts), 1)

				_, ok = line30.Stmts[0].(*PrintStmt)
				testutils.Equal(t, "line 30 is PrintStmt", ok, true)
			},
		},
		{
			name: "COLOR with expression",
			source: `
10 COLOR = X + 1
`,
			lineCount: 1,
			assertFunc: func(t *testing.T, prog *Program) {
				line := prog.Lines[0]
				c, ok := line.Stmts[0].(*ColorStmt)
				testutils.Equal(t, "stmt is ColorStmt", ok, true)

				_, ok = c.Expr.(*InfixExpr)
				testutils.Equal(t, "Expr is InfixExpr", ok, true)
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

			tt.assertFunc(t, prog)
		})
	}
}

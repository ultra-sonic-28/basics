package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_PLOT_Statement(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		lineCount  int
		assertFunc func(t *testing.T, prog *Program)
	}{
		{
			name: "PLOT same line with PRINT",
			source: `
10 PLOT 10, 20: PRINT "HELLO"
`,
			lineCount: 1,
			assertFunc: func(t *testing.T, prog *Program) {
				line := prog.Lines[0]
				testutils.Equal(t, "line number", line.Number, 10)
				testutils.Equal(t, "two statements", len(line.Stmts), 2)

				p, ok := line.Stmts[0].(*PlotStmt)
				testutils.Equal(t, "first stmt is PlotStmt", ok, true)
				testutils.Equal(t, "X value", p.X.(*NumberLiteral).Value, 10.0)
				testutils.Equal(t, "Y value", p.Y.(*NumberLiteral).Value, 20.0)

				_, ok = line.Stmts[1].(*PrintStmt)
				testutils.Equal(t, "second stmt is PrintStmt", ok, true)
			},
		},
		{
			name: "PLOT on its own line",
			source: `
10 PLOT 5, 5
30 PRINT "HELLO"
`,
			lineCount: 2,
			assertFunc: func(t *testing.T, prog *Program) {
				line10 := prog.Lines[0]
				testutils.Equal(t, "line 10 number", line10.Number, 10)
				testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 1)

				p, ok := line10.Stmts[0].(*PlotStmt)
				testutils.Equal(t, "line 10 is PlotStmt", ok, true)
				testutils.Equal(t, "X value", p.X.(*NumberLiteral).Value, 5.0)
				testutils.Equal(t, "Y value", p.Y.(*NumberLiteral).Value, 5.0)

				line30 := prog.Lines[1]
				testutils.Equal(t, "line 30 number", line30.Number, 30)
				testutils.Equal(t, "line 30 stmt count", len(line30.Stmts), 1)

				_, ok = line30.Stmts[0].(*PrintStmt)
				testutils.Equal(t, "line 30 is PrintStmt", ok, true)
			},
		},
		{
			name: "PLOT with expressions",
			source: `
10 PLOT X + 1, Y * 2
`,
			lineCount: 1,
			assertFunc: func(t *testing.T, prog *Program) {
				line := prog.Lines[0]
				p, ok := line.Stmts[0].(*PlotStmt)
				testutils.Equal(t, "stmt is PlotStmt", ok, true)

				_, ok = p.X.(*InfixExpr)
				testutils.Equal(t, "X is InfixExpr", ok, true)

				_, ok = p.Y.(*InfixExpr)
				testutils.Equal(t, "Y is InfixExpr", ok, true)
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

package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_HTAB_Statement(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		lineCount  int
		assertFunc func(t *testing.T, prog *Program)
	}{
		{
			name: "HTAB same line with PRINT",
			source: `
10 HTAB 10: PRINT "HELLO"
`,
			lineCount: 1,
			assertFunc: func(t *testing.T, prog *Program) {
				line := prog.Lines[0]
				testutils.Equal(t, "line number", line.Number, 10)
				testutils.Equal(t, "two statements", len(line.Stmts), 2)

				_, ok := line.Stmts[0].(*HTabStmt)
				testutils.Equal(t, "first stmt is HTabStmt", ok, true)

				_, ok = line.Stmts[1].(*PrintStmt)
				testutils.Equal(t, "second stmt is PrintStmt", ok, true)
			},
		},
		{
			name: "HTAB on its own line",
			source: `
10 HTAB 10
30 PRINT "HELLO"
`,
			lineCount: 2,
			assertFunc: func(t *testing.T, prog *Program) {
				line10 := prog.Lines[0]
				testutils.Equal(t, "line 10 number", line10.Number, 10)
				testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 1)

				_, ok := line10.Stmts[0].(*HTabStmt)
				testutils.Equal(t, "line 10 is HTabStmt", ok, true)

				line30 := prog.Lines[1]
				testutils.Equal(t, "line 30 number", line30.Number, 30)
				testutils.Equal(t, "line 30 stmt count", len(line30.Stmts), 1)

				_, ok = line30.Stmts[0].(*PrintStmt)
				testutils.Equal(t, "line 30 is PrintStmt", ok, true)
			},
		},
		{
			name: "HTAB inside FOR loop",
			source: `
10 FOR A = 1 TO 5
20 HTAB A * 2: PRINT A
30 NEXT A
`,
			lineCount: 3,
			assertFunc: func(t *testing.T, prog *Program) {
				// FOR
				line10 := prog.Lines[0]
				testutils.Equal(t, "line 10 number", line10.Number, 10)
				testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 1)

				_, ok := line10.Stmts[0].(*ForStmt)
				testutils.Equal(t, "line 10 is ForStmt", ok, true)

				// HTAB + PRINT
				line20 := prog.Lines[1]
				testutils.Equal(t, "line 20 number", line20.Number, 20)
				testutils.Equal(t, "line 20 stmt count", len(line20.Stmts), 2)

				_, ok = line20.Stmts[0].(*HTabStmt)
				testutils.Equal(t, "line 20 first stmt is HTabStmt", ok, true)

				_, ok = line20.Stmts[1].(*PrintStmt)
				testutils.Equal(t, "line 20 second stmt is PrintStmt", ok, true)

				// NEXT
				line30 := prog.Lines[2]
				testutils.Equal(t, "line 30 number", line30.Number, 30)
				testutils.Equal(t, "line 30 stmt count", len(line30.Stmts), 1)

				_, ok = line30.Stmts[0].(*NextStmt)
				testutils.Equal(t, "line 30 is NextStmt", ok, true)
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

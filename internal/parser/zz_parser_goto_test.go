package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_GOTO_Statements(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		lineCount  int
		assertFunc func(t *testing.T, prog *Program)
	}{
		{
			name: "Simple GOTO with constant line",
			source: `
10 PRINT "First line"
20 GOTO 40
30 PRINT "Third line"
40 PRINT "Second line"
`,
			lineCount: 4,
			assertFunc: func(t *testing.T, prog *Program) {
				// 10 PRINT
				line10 := prog.Lines[0]
				testutils.Equal(t, "line 10 number", line10.Number, 10)
				testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 1)
				_, ok := line10.Stmts[0].(*PrintStmt)
				testutils.Equal(t, "line 10 is PrintStmt", ok, true)

				// 20 GOTO
				line20 := prog.Lines[1]
				testutils.Equal(t, "line 20 stmt count", len(line20.Stmts), 1)
				_, ok = line20.Stmts[0].(*GotoStmt)
				testutils.Equal(t, "line 20 is GotoStmt", ok, true)

				// 30 PRINT
				line30 := prog.Lines[2]
				_, ok = line30.Stmts[0].(*PrintStmt)
				testutils.Equal(t, "line 30 is PrintStmt", ok, true)

				// 40 PRINT
				line40 := prog.Lines[3]
				_, ok = line40.Stmts[0].(*PrintStmt)
				testutils.Equal(t, "line 40 is PrintStmt", ok, true)
			},
		},
		{
			name: "GOTO with variable target",
			source: `
10 REM GOTO Example
15 JUMP = 80
20 PRINT "First line"
30 GOTO 60
40 PRINT "Third line"
50 GOTO JUMP
60 PRINT "Second line"
70 GOTO 40
80 PRINT "Last line"
`,
			lineCount: 9,
			assertFunc: func(t *testing.T, prog *Program) {
				// Assignment
				_, ok := prog.Lines[1].Stmts[0].(*LetStmt)
				testutils.Equal(t, "line 15 is AssignStmt", ok, true)

				// GOTO 60
				_, ok = prog.Lines[3].Stmts[0].(*GotoStmt)
				testutils.Equal(t, "line 30 is GotoStmt", ok, true)

				// GOTO JUMP
				_, ok = prog.Lines[5].Stmts[0].(*GotoStmt)
				testutils.Equal(t, "line 50 is GotoStmt", ok, true)

				// GOTO 40
				_, ok = prog.Lines[7].Stmts[0].(*GotoStmt)
				testutils.Equal(t, "line 70 is GotoStmt", ok, true)
			},
		},
		{
			name: "GOTO with expression target",
			source: `
10 REM GOTO Example
15 JUMP = 40
20 PRINT "First line"
30 GOTO 60
40 PRINT "Third line"
50 GOTO JUMP * 2
60 PRINT "Second line"
70 GOTO 40
80 PRINT "Last line"
`,
			lineCount: 9,
			assertFunc: func(t *testing.T, prog *Program) {
				line50 := prog.Lines[5]
				testutils.Equal(t, "line 50 stmt count", len(line50.Stmts), 1)

				gotoStmt, ok := line50.Stmts[0].(*GotoStmt)
				testutils.Equal(t, "line 50 is GotoStmt", ok, true)

				// expression must exist
				testutils.True(t, "goto has expression", gotoStmt.Expr != nil)
			},
		},
		{
			name: "GOTO on same line with PRINT",
			source: `
10 PRINT "First line" : GOTO 30
20 PRINT "Third line"
30 PRINT "Second line"
`,
			lineCount: 3,
			assertFunc: func(t *testing.T, prog *Program) {
				line10 := prog.Lines[0]
				testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 2)

				_, ok := line10.Stmts[0].(*PrintStmt)
				testutils.Equal(t, "first stmt is PrintStmt", ok, true)

				_, ok = line10.Stmts[1].(*GotoStmt)
				testutils.Equal(t, "second stmt is GotoStmt", ok, true)
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

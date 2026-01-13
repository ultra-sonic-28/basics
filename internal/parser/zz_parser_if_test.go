package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_IF_Statements(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		lineCount  int
		assertFunc func(t *testing.T, prog *Program)
	}{
		{
			name: "IF with implicit GOTO",
			source: `
10 LET count = 0
20 PRINT "Count: ", count
30 LET count = count + 1
40 IF count < 10 THEN GOTO 20
50 PRINT "All done!"
`,
			lineCount: 5,
			assertFunc: func(t *testing.T, prog *Program) {
				line40 := prog.Lines[3]
				testutils.Equal(t, "line 40 stmt count", len(line40.Stmts), 1)

				ifStmt, ok := line40.Stmts[0].(*IfStmt)
				testutils.Equal(t, "line 40 is IfStmt", ok, true)

				testutils.True(t, "condition exists", ifStmt.Cond != nil)
				testutils.Equal(t, "then stmt count", len(ifStmt.Then), 1)
				testutils.Equal(t, "else stmt count", len(ifStmt.Else), 0)

				_, ok = ifStmt.Then[0].(*GotoStmt)
				testutils.Equal(t, "THEN is GotoStmt", ok, true)
			},
		},
		{
			name: "IF THEN GOTO without ELSE",
			source: `
10 LET count = 0
20 PRINT "Count: ", count
30 LET count = count + 1
40 IF count > 9 THEN GOTO 60
50 GOTO 20
60 PRINT "All done!"
`,
			lineCount: 6,
			assertFunc: func(t *testing.T, prog *Program) {
				ifStmt := prog.Lines[3].Stmts[0].(*IfStmt)

				testutils.Equal(t, "then stmt count", len(ifStmt.Then), 1)
				testutils.Equal(t, "else stmt count", len(ifStmt.Else), 0)
			},
		},
		{
			name: "IF THEN ELSE with multiple statements",
			source: `
10 LET count = 0
20 PRINT "Count: ", count
30 LET count = count + 1
40 IF count < 10 THEN PRINT "Go to line 20" : GOTO 20 ELSE PRINT "Go to line 60" : GOTO 60
50 END
60 PRINT "All done!"
`,
			lineCount: 6,
			assertFunc: func(t *testing.T, prog *Program) {
				ifStmt := prog.Lines[3].Stmts[0].(*IfStmt)

				testutils.Equal(t, "then stmt count", len(ifStmt.Then), 2)
				testutils.Equal(t, "else stmt count", len(ifStmt.Else), 2)

				_, ok := ifStmt.Then[0].(*PrintStmt)
				testutils.True(t, "then[0] is PrintStmt", ok)

				_, ok = ifStmt.Then[1].(*GotoStmt)
				testutils.True(t, "then[1] is GotoStmt", ok)

				_, ok = ifStmt.Else[0].(*PrintStmt)
				testutils.True(t, "else[0] is PrintStmt", ok)

				_, ok = ifStmt.Else[1].(*GotoStmt)
				testutils.True(t, "else[1] is GotoStmt", ok)
			},
		},
		{
			name: "IF THEN <number> implicit GOTO",
			source: `
10 LET count = 0
20 PRINT "Count: ", count
30 LET count = count + 1
40 IF count > 9 THEN 60
50 GOTO 20
60 PRINT "All done!"
`,
			lineCount: 6,
			assertFunc: func(t *testing.T, prog *Program) {
				ifStmt := prog.Lines[3].Stmts[0].(*IfStmt)

				testutils.Equal(t, "then stmt count", len(ifStmt.Then), 1)
				_, ok := ifStmt.Then[0].(*GotoStmt)
				testutils.True(t, "implicit GOTO created", ok)
			},
		},
		{
			name: "Complex IF in Fibonacci program",
			source: `
10 REM 20 first Fibonacci numbers
15 LET n = 20
20 LET first = 0
30 LET second = 1
40 PRINT first
50 IF n = 1 THEN GOTO 150
60 PRINT second
70 IF n = 2 THEN GOTO 150
80 LET count = 3
90 LET next = first + second
100 PRINT next
110 LET first = second
120 LET second = next
130 LET count = count + 1
140 IF count < (n + 1) THEN GOTO 90
150 PRINT "All done!"
`,
			lineCount: 16,
			assertFunc: func(t *testing.T, prog *Program) {
				line50 := prog.Lines[5]
				_, ok := line50.Stmts[0].(*IfStmt)
				testutils.True(t, "line 50 is IfStmt", ok)

				line70 := prog.Lines[7]
				_, ok = line70.Stmts[0].(*IfStmt)
				testutils.True(t, "line 70 is IfStmt", ok)

				line140 := prog.Lines[14]
				_, ok = line140.Stmts[0].(*IfStmt)
				testutils.True(t, "line 140 is IfStmt", ok)
			},
		},
		{
			name: "IF inside FOR loop (prime numbers)",
			source: `
5 MAX = 50
10 PRINT "NOMBRES PREMIERS JUSQU'A "; MAX
20 FOR N = 2 TO MAX
30 P = 1
40 FOR D = 2 TO N/2
50 T = D
60 IF T >= N THEN 90
70 T = T + D
80 GOTO 60
90 IF T = N THEN P = 0
100 NEXT D
110 IF P = 1 THEN PRINT N
120 NEXT N
130 PRINT "All done!"
`,
			lineCount: 14,
			assertFunc: func(t *testing.T, prog *Program) {
				// IF T >= N
				_, ok := prog.Lines[6].Stmts[0].(*IfStmt)
				testutils.True(t, "line 60 is IfStmt", ok)

				// IF T = N
				_, ok = prog.Lines[9].Stmts[0].(*IfStmt)
				testutils.True(t, "line 90 is IfStmt", ok)

				// IF P = 1
				_, ok = prog.Lines[11].Stmts[0].(*IfStmt)
				testutils.True(t, "line 110 is IfStmt", ok)
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

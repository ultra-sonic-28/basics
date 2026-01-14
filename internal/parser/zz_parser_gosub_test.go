package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_GOSUB_RETURN_Statements(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		lineCount  int
		assertFunc func(t *testing.T, prog *Program)
	}{
		{
			name: "Simple GOSUB / RETURN",
			source: `
10 PRINT "Hello "
20 GOSUB 100
30 PRINT "!!!"
40 END
100 PRINT "World"
110 RETURN
`,
			lineCount: 6,
			assertFunc: func(t *testing.T, prog *Program) {
				// line 20 : GOSUB
				line20 := prog.Lines[1]
				testutils.Equal(t, "line 20 stmt count", len(line20.Stmts), 1)
				_, ok := line20.Stmts[0].(*GosubStmt)
				testutils.True(t, "line 20 is GosubStmt", ok)

				// line 110 : RETURN
				line110 := prog.Lines[5]
				testutils.Equal(t, "line 110 stmt count", len(line110.Stmts), 1)
				_, ok = line110.Stmts[0].(*ReturnStmt)
				testutils.True(t, "line 110 is ReturnStmt", ok)
			},
		},
		{
			name: "Inline GOSUB with PRINT before and after",
			source: `
10 PRINT "Hello " : GOSUB 100 : PRINT "!!!"
30 END
100 PRINT "World" : RETURN
`,
			lineCount: 3,
			assertFunc: func(t *testing.T, prog *Program) {
				line10 := prog.Lines[0]
				testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 3)

				_, ok := line10.Stmts[0].(*PrintStmt)
				testutils.True(t, "stmt 0 is PrintStmt", ok)

				_, ok = line10.Stmts[1].(*GosubStmt)
				testutils.True(t, "stmt 1 is GosubStmt", ok)

				_, ok = line10.Stmts[2].(*PrintStmt)
				testutils.True(t, "stmt 2 is PrintStmt", ok)

				line100 := prog.Lines[2]
				testutils.Equal(t, "line 100 stmt count", len(line100.Stmts), 2)

				_, ok = line100.Stmts[1].(*ReturnStmt)
				testutils.True(t, "line 100 has ReturnStmt", ok)
			},
		},
		{
			name: "GOSUB inside FOR loop",
			source: `
5 REM **** Ce programme affiche la table de 4 ****
10 PRINT "TABLE DE 4 :"
20 FOR I=1 TO 10
25 GOSUB 100
30 PRINT I, V
40 NEXT I
50 END
100 V = I * 4 : RETURN
`,
			lineCount: 8,
			assertFunc: func(t *testing.T, prog *Program) {
				// line 25 : GOSUB
				line25 := prog.Lines[3]
				testutils.Equal(t, "line 25 stmt count", len(line25.Stmts), 1)
				_, ok := line25.Stmts[0].(*GosubStmt)
				testutils.True(t, "line 25 is GosubStmt", ok)

				// line 100 : assignment + RETURN
				line100 := prog.Lines[7]
				testutils.Equal(t, "line 100 stmt count", len(line100.Stmts), 2)

				_, ok = line100.Stmts[0].(*LetStmt)
				testutils.True(t, "line 100 first stmt is LetStmt", ok)

				_, ok = line100.Stmts[1].(*ReturnStmt)
				testutils.True(t, "line 100 second stmt is ReturnStmt", ok)
			},
		},
		{
			name: "GOSUB with expression target",
			source: `
10 PRINT "Hello " : A=50 : GOSUB A*2 : PRINT "!!!"
30 END
100 PRINT "World" : RETURN
`,
			lineCount: 3,
			assertFunc: func(t *testing.T, prog *Program) {
				line10 := prog.Lines[0]
				testutils.Equal(t, "line 10 stmt count", len(line10.Stmts), 4)

				gosub, ok := line10.Stmts[2].(*GosubStmt)
				testutils.True(t, "stmt 2 is GosubStmt", ok)

				testutils.True(t, "GOSUB has expression", gosub.Expr != nil)

				line100 := prog.Lines[2]
				_, ok = line100.Stmts[1].(*ReturnStmt)
				testutils.True(t, "line 100 has ReturnStmt", ok)
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

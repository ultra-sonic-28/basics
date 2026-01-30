package interpreter

import (
	"basics/internal/common"
	"basics/internal/constants"
	"basics/internal/lexer"
	"basics/internal/machines"
	"basics/internal/parser"
	"basics/testutils"
	"bytes"
	"fmt"
	"testing"
)

type clearTestCase struct {
	name    string
	program string
	want    string
}

func TestCLEAR_TableDriven(t *testing.T) {
	tests := []clearTestCase{
		{
			name: "CLEAR float variable",
			program: `
10 LET A = 1.5
20 PRINT "A=";A
30 CLEAR
40 PRINT "A=";A
`,
			want: "A=1.5\nA=0\n",
		},
		{
			name: "CLEAR integer variable",
			program: `
10 LET I% = 42
20 PRINT "I%=";I%
30 CLEAR
40 PRINT "I%=";I%
`,
			want: "I%=42\nI%=0\n",
		},
		{
			name: "CLEAR string variable",
			program: `
10 LET S$ = "Hello"
20 PRINT "S$=";S$
30 CLEAR
40 PRINT "S$=";S$
`,
			want: "S$=Hello\nS$=\n",
		},
		{
			name: "CLEAR all variables",
			program: `
10 LET A = 1.5
20 LET I% = 42
30 LET S$ = "Hello"
40 PRINT "A=";A
50 PRINT "I%=";I%
60 PRINT "S$=";S$
70 CLEAR
80 PRINT "A=";A
90 PRINT "I%=";I%
100 PRINT "S$=";S$
`,
			want: "A=1.5\nI%=42\nS$=Hello\nA=0\nI%=0\nS$=\n",
		},
	}

	rt, _ := machines.NewRuntime(constants.BASIC_TTY)

	for tIndex, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			rt.SetOutput(out)

			i := New(rt)

			tokens := lexer.Lex(tc.program)
			p := parser.New(tokens)
			prog, errs := p.ParseProgram()
			testutils.Equal(t, "no parser errors", len(errs), 0)

			i.Run(prog)

			got := common.StripANSI(out.String())
			testutils.True(
				t,
				fmt.Sprintf(
					"tests[%d]\n--- EXPECTED ---\n%q\n--- GOT ---\n%q\n",
					tIndex,
					tc.want,
					got,
				),
				got == tc.want,
			)
		})
	}
}

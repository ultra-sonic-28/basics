package interpreter

import (
	"basics/internal/constants"
	"basics/internal/input"
	"basics/internal/lexer"
	"basics/internal/machines"
	"basics/internal/parser"
	"basics/testutils"
	"bytes"
	"fmt"
	"testing"
)

type getTestCase struct {
	name    string
	program string
	input   string
	want    string
}

func TestGET_TableDriven(t *testing.T) {
	tests := []getTestCase{
		{
			name: "GET single character",
			program: `
10 GET A$
20 PRINT A$
`,
			input: "x",
			want:  "x\n",
		},
		{
			name: "GET does not echo character",
			program: `
10 PRINT "Press key"
20 GET A$
30 PRINT "Done"
`,
			input: "q",
			want:  "Press key\nDone\n",
		},
		{
			name: "GET stores only first character",
			program: `
10 GET A$
20 PRINT A$
`,
			input: "abcd",
			want:  "a\n",
		},
		{
			name: "GET does not require enter",
			program: `
10 GET A$
20 PRINT A$
`,
			input: "Z\n",
			want:  "Z\n",
		},
		{
			name: "GET before PRINT",
			program: `
10 GET A$
20 PRINT "Key was: ";A$
`,
			input: "K",
			want:  "Key was: K\n",
		},
	}

	rt, _ := machines.NewRuntime(constants.BASIC_TTY)

	for tIndex, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := input.NewFakeInput(tc.input)
			out := &bytes.Buffer{}

			rt.Input = in
			rt.SetOutput(out)

			i := New(rt)

			tokens := lexer.Lex(tc.program)
			p := parser.New(tokens)
			prog, errs := p.ParseProgram()
			testutils.Equal(t, "no parser errors", len(errs), 0)

			i.Run(prog)

			got := out.String()
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

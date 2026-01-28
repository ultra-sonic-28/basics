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

type inputTestCase struct {
	name    string
	program string
	input   string
	want    string
}

func TestINPUT_TableDriven(t *testing.T) {
	tests := []inputTestCase{
		{
			name: "INPUT numeric single",
			program: `
10 INPUT A
20 PRINT A
`,
			input: "42\n",
			want:  "? 42\n",
		},
		{
			name: "INPUT string single",
			program: `
10 INPUT A$
20 PRINT A$
`,
			input: "hello\n",
			want:  "? hello\n",
		},
		{
			name: "INPUT numeric with prompt",
			program: `
10 INPUT "Enter value: ";A
20 PRINT A
`,
			input: "12\n",
			want:  "Enter value: 12\n",
		},
		{
			name: "INPUT string with prompt",
			program: `
10 INPUT "Name: ";N$
20 PRINT N$
`,
			input: "Dom\n",
			want:  "Name: Dom\n",
		},
		{
			name: "INPUT multiple numeric",
			program: `
10 INPUT "Enter 2 values: ";A,B
20 PRINT A*B
`,
			input: "10,52\n",
			want:  "Enter 2 values: 520\n",
		},
		{
			name: "INPUT multiple strings",
			program: `
10 INPUT "Enter 2 strings: ";A$,B$
20 PRINT A$;B$
`,
			input: "foo,bar\n",
			want:  "Enter 2 strings: foobar\n",
		},
		{
			name: "INPUT without prompt prints question mark",
			program: `
10 INPUT A
`,
			input: "5\n",
			want:  "? ",
		},
		/* {
			name: "INPUT REENTER on invalid numeric",
			program: `
		10 INPUT A
		20 PRINT A
		`,
			input: "abc\n",
			want:  "\n?REENTER\n? ",
		}, */
	}

	rt, _ := machines.NewRuntime(constants.BASIC_TTY)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := input.NewFakeInput(tc.input)
			out := &bytes.Buffer{}

			// connecter le runtime au test lors de son exécution
			// objectif : passer la saisie au programme et récupérer la sortie de celui-ci
			rt.Input = in
			rt.SetOutput(out)

			i := New(rt)

			source := tc.program
			tokens := lexer.Lex(source)
			p := parser.New(tokens)
			prog, _ := p.ParseProgram()

			i.Run(prog)

			got := out.String()
			testutils.True(
				t,
				fmt.Sprintf("\n--- EXPECTED ---\n%q\n--- GOT ---\n%q\n", tc.want, got),
				got == tc.want,
			)
		})
	}
}

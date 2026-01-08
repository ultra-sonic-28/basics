package parser

import (
	"basics/internal/lexer"
	"basics/testutils"
	"strings"
	"testing"
)

func TestParser_FOR_NEXT_Golden(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		golden string
	}{
		{
			name: "FOR TO NEXT",
			input: `
10 FOR I = 1 TO 10
20 NEXT I
`,
			golden: `
| Line | Statement | Details |
|------|-----------|---------|
| 10 | FOR | var=I, from=1, to=10, step=1 |
| 20 | NEXT | var=I, for_line=10 |
`,
		},
		{
			name: "FOR TO STEP NEXT",
			input: `
10 FOR I = 0 TO 20 STEP 2
20 NEXT I
`,
			golden: `
| Line | Statement | Details |
|------|-----------|---------|
| 10 | FOR | var=I, from=0, to=20, step=2 |
| 20 | NEXT | var=I, for_line=10 |
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tokens := lexer.Lex(tt.input)
			/* if len(lexErrs) != 0 {
				t.Fatalf("lexer errors: %+v", lexErrs)
			} */

			p := New(tokens)

			prog, errs := p.ParseProgram()
			testutils.Equal(t, "no errors", len(errs), 0)

			got := strings.TrimSpace(programToMarkdown(prog))
			want := strings.TrimSpace(tt.golden)
			testutils.Equal(t, "AST markdown", got, want)
		})
	}
}

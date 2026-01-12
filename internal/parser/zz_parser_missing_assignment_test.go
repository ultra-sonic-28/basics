package parser

import (
	"basics/internal/lexer"
	"basics/testutils"
	"testing"
)

func TestParse_LetMissingEqual(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			"LET without =",
			`
10 REM Test
20 LET A 3
30 PRINT A
`,
		},
		{
			"implicit LET without =",
			`
10 REM Test
20 A 3
30 PRINT A
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexer.Lex(tt.source)
			p := New(tokens)
			_, errs := p.ParseProgram()

			testutils.Equal(t, "one error", len(errs), 2)
			testutils.Equal(t, "error message", errs[0].Msg, "EXPECTED '='")
			testutils.Equal(t, "error line", errs[0].Line, 3)
		})
	}
}

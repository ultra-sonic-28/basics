package parser

import (
	"fmt"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_CHR_Errors(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "CHR$ with empty parentheses",
			source: `10 PRINT CHR$()`,
		},
		{
			name:   "CHR$ with trailing comma",
			source: `10 PRINT CHR$(3,)`,
		},
		{
			name:   "CHR$ with leading comma",
			source: `10 PRINT CHR$(,3)`,
		},
		{
			name:   "CHR$ with missing closing paren",
			source: `10 PRINT CHR$(3.5`,
		},
		{
			name:   "CHR$ with only opening paren",
			source: `10 PRINT CHR$(`,
		},
		{
			name:   "CHR$ without parentheses",
			source: `10 PRINT CHR$ 5`,
		},
		{
			name:   "CHR$ with two arguments",
			source: `10 PRINT CHR$(5,2)`,
		},
		{
			name:   "CHR$ nested missing paren",
			source: `10 PRINT CHR$((5)`,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexer.Lex(tt.source)
			p := New(tokens)

			prog, errs := p.ParseProgram()

			// Le parser DOIT signaler une erreur
			testutils.True(t,
				fmt.Sprintf("tests[%d] - parser should return errors", i),
				len(errs) > 0,
			)

			// Le programme ne doit pas être nil
			testutils.True(t,
				fmt.Sprintf("tests[%d] - program is not nil", i),
				prog != nil,
			)
		})
	}
}

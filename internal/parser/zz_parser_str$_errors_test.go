package parser

import (
	"fmt"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_STR_Errors(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "STR$ with empty parentheses",
			source: `10 PRINT STR$()`,
		},
		{
			name:   "STR$ with trailing comma",
			source: `10 PRINT STR$(3,)`,
		},
		{
			name:   "STR$ with leading comma",
			source: `10 PRINT STR$(,3)`,
		},
		{
			name:   "STR$ with missing closing paren",
			source: `10 PRINT STR$(3.5`,
		},
		{
			name:   "STR$ with only opening paren",
			source: `10 PRINT STR$(`,
		},
		{
			name:   "STR$ without parentheses",
			source: `10 PRINT STR$ 5`,
		},
		{
			name:   "STR$ with two arguments",
			source: `10 PRINT STR$(5,2)`,
		},
		{
			name:   "STR$ nested missing paren",
			source: `10 PRINT STR$((5)`,
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

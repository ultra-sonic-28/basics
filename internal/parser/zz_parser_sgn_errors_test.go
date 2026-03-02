package parser

import (
	"fmt"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_SGN_Errors(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "SGN with empty parentheses",
			source: `10 PRINT SGN()`,
		},
		{
			name:   "SGN with trailing comma",
			source: `10 PRINT SGN(A,)`,
		},
		{
			name:   "SGN with missing closing paren",
			source: `10 PRINT SGN(A`,
		},
		{
			name:   "SGN with only opening paren",
			source: `10 PRINT SGN(`,
		},
		{
			name:   "SGN without parentheses",
			source: `10 PRINT SGN 10`,
		},
		{
			name:   "SGN with multiple arguments",
			source: `10 PRINT SGN(10,20)`,
		},
		{
			name:   "SGN nested missing paren",
			source: `10 PRINT SGN((A+2)`,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexer.Lex(tt.source)
			p := New(tokens)

			prog, errs := p.ParseProgram()

			// Le parser DOIT signaler une erreur
			testutils.True(t, fmt.Sprintf("tests[%d] - parser should return errors", i), len(errs) > 0)

			// Le programme ne doit pas être nil
			testutils.True(t, fmt.Sprintf("tests[%d] - program is not nil", i), prog != nil)
		})
	}
}

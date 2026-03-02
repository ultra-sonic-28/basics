package parser

import (
	"fmt"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_COS_Errors(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "COS with empty parentheses",
			source: `10 PRINT COS()`,
		},
		{
			name:   "COS with trailing comma",
			source: `10 PRINT COS(A,)`,
		},
		{
			name:   "COS with missing closing paren",
			source: `10 PRINT COS(A`,
		},
		{
			name:   "COS with only opening paren",
			source: `10 PRINT COS(`,
		},
		{
			name:   "COS without parentheses",
			source: `10 PRINT COS 10`,
		},
		{
			name:   "COS with multiple arguments",
			source: `10 PRINT COS(10,20)`,
		},
		{
			name:   "COS nested missing paren",
			source: `10 PRINT COS((A+2)`,
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

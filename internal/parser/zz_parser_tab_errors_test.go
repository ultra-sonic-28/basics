package parser

import (
	"fmt"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_TAB_Errors(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "TAB with empty parentheses",
			source: `10 PRINT TAB()`,
		},
		{
			name:   "TAB with trailing comma",
			source: `10 PRINT TAB(A,)`,
		},
		{
			name:   "TAB with missing closing paren",
			source: `10 PRINT TAB(A`,
		},
		{
			name:   "TAB with only opening paren",
			source: `10 PRINT TAB(`,
		},
		{
			name:   "TAB without parentheses",
			source: `10 PRINT TAB 10`,
		},
		{
			name:   "TAB with multiple arguments",
			source: `10 PRINT TAB(10,20)`,
		},
		{
			name:   "TAB nested missing paren",
			source: `10 PRINT TAB((A+2)`,
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

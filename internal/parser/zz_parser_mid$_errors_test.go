package parser

import (
	"fmt"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_MID_Errors(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "MID$ with empty parentheses",
			source: `10 PRINT MID$()`,
		},
		{
			name:   "MID$ with only one argument",
			source: `10 PRINT MID$("ABC")`,
		},
		{
			name:   "MID$ with trailing comma",
			source: `10 PRINT MID$("ABC",)`,
		},
		{
			name:   "MID$ with leading comma",
			source: `10 PRINT MID$(,"ABC")`,
		},
		{
			name:   "MID$ with missing closing paren",
			source: `10 PRINT MID$("ABC",5`,
		},
		{
			name:   "MID$ with only opening paren",
			source: `10 PRINT MID$(`,
		},
		{
			name:   "MID$ without parentheses",
			source: `10 PRINT MID$ "ABC",5`,
		},
		{
			name:   "MID$ with four arguments",
			source: `10 PRINT MID$("ABC",5,2,1)`,
		},
		{
			name:   "MID$ nested missing paren",
			source: `10 PRINT MID$(("ABC",5)`,
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

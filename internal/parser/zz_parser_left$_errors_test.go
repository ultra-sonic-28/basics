package parser

import (
	"fmt"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_LEFT_Errors(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "LEFT$ with empty parentheses",
			source: `10 PRINT LEFT$()`,
		},
		{
			name:   "LEFT$ with only one argument",
			source: `10 PRINT LEFT$("ABC")`,
		},
		{
			name:   "LEFT$ with trailing comma",
			source: `10 PRINT LEFT$("ABC",)`,
		},
		{
			name:   "LEFT$ with leading comma",
			source: `10 PRINT LEFT$(,"ABC")`,
		},
		{
			name:   "LEFT$ with missing closing paren",
			source: `10 PRINT LEFT$("ABC",5`,
		},
		{
			name:   "LEFT$ with only opening paren",
			source: `10 PRINT LEFT$(`,
		},
		{
			name:   "LEFT$ without parentheses",
			source: `10 PRINT LEFT$ "ABC",5`,
		},
		{
			name:   "LEFT$ with three arguments",
			source: `10 PRINT LEFT$("ABC",5,2)`,
		},
		{
			name:   "LEFT$ nested missing paren",
			source: `10 PRINT LEFT$(("ABC",5)`,
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

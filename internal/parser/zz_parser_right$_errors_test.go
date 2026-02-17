package parser

import (
	"fmt"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_RIGHT_Errors(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "RIGHT$ with empty parentheses",
			source: `10 PRINT RIGHT$()`,
		},
		{
			name:   "RIGHT$ with only one argument",
			source: `10 PRINT RIGHT$("ABC")`,
		},
		{
			name:   "RIGHT$ with trailing comma",
			source: `10 PRINT RIGHT$("ABC",)`,
		},
		{
			name:   "RIGHT$ with leading comma",
			source: `10 PRINT RIGHT$(,"ABC")`,
		},
		{
			name:   "RIGHT$ with missing closing paren",
			source: `10 PRINT RIGHT$("ABC",5`,
		},
		{
			name:   "RIGHT$ with only opening paren",
			source: `10 PRINT RIGHT$(`,
		},
		{
			name:   "RIGHT$ without parentheses",
			source: `10 PRINT RIGHT$ "ABC",5`,
		},
		{
			name:   "RIGHT$ with three arguments",
			source: `10 PRINT RIGHT$("ABC",5,2)`,
		},
		{
			name:   "RIGHT$ nested missing paren",
			source: `10 PRINT RIGHT$(("ABC",5)`,
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

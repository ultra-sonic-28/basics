package parser

import (
	"fmt"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_LEN_Errors(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "LEN with empty parentheses",
			source: `10 PRINT LEN()`,
		},
		{
			name:   "LEN with trailing comma",
			source: `10 PRINT LEN("ABC",)`,
		},
		{
			name:   "LEN with leading comma",
			source: `10 PRINT LEN(,"ABC")`,
		},
		{
			name:   "LEN with missing closing paren",
			source: `10 PRINT LEN("ABC"`,
		},
		{
			name:   "LEN with only opening paren",
			source: `10 PRINT LEN(`,
		},
		{
			name:   "LEN without parentheses",
			source: `10 PRINT LEN "ABC"`,
		},
		{
			name:   "LEN with three arguments",
			source: `10 PRINT LEN("ABC",5,2)`,
		},
		{
			name:   "LEN nested missing paren",
			source: `10 PRINT LEN(("ABC",5)`,
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

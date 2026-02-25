package parser

import (
	"strings"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_PLOT_Errors(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedError string
	}{
		{
			name: "PLOT missing comma",
			source: `
10 PLOT 10 20
`,
			expectedError: "EXPECTED , AFTER PLOT X",
		},
		{
			name: "PLOT missing Y",
			source: `
10 PLOT 10,
`,
			expectedError: "INVALID EXPRESSION", // parseExpression(LOWEST) will fail
		},
		{
			name: "PLOT missing X",
			source: `
10 PLOT , 20
`,
			expectedError: "INVALID EXPRESSION",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexer.Lex(tt.source)
			p := New(tokens)
			_, errs := p.ParseProgram()

			testutils.True(t, "parser should report errors", len(errs) > 0)

			found := false
			for _, err := range errs {
				if strings.Contains(err.Error(), tt.expectedError) {
					found = true
					break
				}
			}

			testutils.True(t,
				"error "+tt.expectedError+" should be reported",
				found,
			)
		})
	}
}

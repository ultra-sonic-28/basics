package parser

import (
	"strings"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_COLOR_Errors(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedError string
	}{
		{
			name: "COLOR missing =",
			source: `
10 COLOR 15
`,
			expectedError: "EXPECTED = AFTER COLOR",
		},
		{
			name: "COLOR missing expression",
			source: `
10 COLOR = : PRINT "HELLO"
`,
			expectedError: "INVALID EXPRESSION",
		},
		{
			name: "COLOR end of line after =",
			source: `
10 COLOR =
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

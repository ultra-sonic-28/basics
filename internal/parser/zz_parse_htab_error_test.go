package parser

import (
	"strings"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_HTAB_MissingExpression(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name: "HTAB followed by colon",
			source: `
10 HTAB : PRINT "HELLO"
`,
		},
		{
			name: "HTAB end of line",
			source: `
10 HTAB
30 PRINT "HELLO"
`,
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
				if strings.Contains(err.Error(), "EXPECTED EXPRESSION AFTER HTAB") {
					found = true
					break
				}
			}

			testutils.True(t,
				"error EXPECTED EXPRESSION AFTER HTAB should be reported",
				found,
			)
		})
	}
}

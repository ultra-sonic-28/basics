package parser

import (
	"strings"
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_VTAB_MissingExpression(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name: "VTAB followed by colon",
			source: `
10 VTAB : PRINT "HELLO"
`,
		},
		{
			name: "VTAB end of line",
			source: `
10 VTAB
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
				if strings.Contains(err.Error(), "EXPECTED EXPRESSION AFTER VTAB") {
					found = true
					break
				}
			}

			testutils.True(t,
				"error EXPECTED EXPRESSION AFTER VTAB should be reported",
				found,
			)
		})
	}
}

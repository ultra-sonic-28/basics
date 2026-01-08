package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"testing"
)

func TestDumpTokens_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []token.Token
		expected string
	}{
		{
			name: "single token with literal",
			tokens: []token.Token{
				{
					Type:    token.NUMBER,
					Literal: "42",
					Line:    1,
					Column:  3,
				},
			},
			expected: `[1:3] NUMBER   "42"
`,
		},
		{
			name: "single token without literal (EOF)",
			tokens: []token.Token{
				{
					Type:   token.EOF,
					Line:   2,
					Column: 1,
				},
			},
			expected: `[2:1] EOF     
`,
		},
		{
			name: "mixed tokens with and without literals",
			tokens: []token.Token{
				{
					Type:    token.IDENT,
					Literal: "X",
					Line:    5,
					Column:  10,
				},
				{
					Type:   token.EOL,
					Line:   5,
					Column: 11,
				},
				{
					Type:    token.STRING,
					Literal: "HELLO",
					Line:    6,
					Column:  1,
				},
			},
			expected: `[5:10] IDENT    "X"
[5:11] EOL     
[6:1] STRING   "HELLO"
`,
		},
		{
			name:     "empty token slice",
			tokens:   []token.Token{},
			expected: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := testutils.CaptureStdout(t, func() {
				DumpTokens(tt.tokens)
			})

			testutils.Equal(t, "DumpTokens output mismatch", output, tt.expected)
		})
	}
}

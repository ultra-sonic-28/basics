package lexer

import (
	"testing"

	"basics/internal/common"
	"basics/internal/token"
	"basics/testutils"
)

func TestLex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.TokenType
		literals []string
	}{
		{
			name:  "simple LET",
			input: "10 LET A = 5\n",
			expected: []token.TokenType{
				token.LINENUM,
				token.KEYWORD,
				token.IDENT,
				token.EQUAL,
				token.NUMBER,
				token.EOL,
				token.EOF,
			},
			literals: []string{
				"10", "LET", "A", "=", "5", "\n", "",
			},
		},
		{
			name:  "PRINT string",
			input: "10 PRINT \"HELLO\"\n",
			expected: []token.TokenType{
				token.LINENUM,
				token.KEYWORD,
				token.STRING,
				token.EOL,
				token.EOF,
			},
			literals: []string{
				"10", "PRINT", "HELLO", "\n", "",
			},
		},
		{
			name:  "math expression",
			input: "10 A = 1 + 2 * 3\n",
			expected: []token.TokenType{
				token.LINENUM,
				token.IDENT,
				token.EQUAL,
				token.NUMBER,
				token.PLUS,
				token.NUMBER,
				token.ASTERISK,
				token.NUMBER,
				token.EOL,
				token.EOF,
			},
			literals: []string{
				"10", "A", "=", "1", "+", "2", "*", "3", "\n", "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := Lex(tt.input)

			testutils.Equal(
				t,
				"token count",
				len(tokens),
				len(tt.expected),
			)

			for i, tok := range tokens {
				testutils.Equal(
					t,
					"token type at index "+common.Itoa(i),
					tok.Type,
					tt.expected[i],
				)
				testutils.Equal(
					t,
					"token literal at index "+common.Itoa(i),
					tok.Literal,
					tt.literals[i],
				)
			}
		})
	}
}

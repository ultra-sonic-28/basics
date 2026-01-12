package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLexer_HTAB(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []struct {
			typ     token.TokenType
			literal string
			line    int
			column  int
		}
	}{
		{
			name:  "HTAB same line with PRINT",
			input: `10 HTAB 10: PRINT "HELLO"`,
			tokens: []struct {
				typ     token.TokenType
				literal string
				line    int
				column  int
			}{
				{token.LINENUM, "10", 1, 1},
				{token.KEYWORD, "HTAB", 1, 4},
				{token.NUMBER, "10", 1, 9},
				{token.COLON, ":", 1, 11},
				{token.KEYWORD, "PRINT", 1, 13},
				{token.STRING, "HELLO", 1, 19},
				{token.EOF, "", 1, 26},
			},
		},
		{
			name: "HTAB on its own line",
			input: `10 HTAB 10
30 PRINT "HELLO"`,
			tokens: []struct {
				typ     token.TokenType
				literal string
				line    int
				column  int
			}{
				{token.LINENUM, "10", 1, 1},
				{token.KEYWORD, "HTAB", 1, 4},
				{token.NUMBER, "10", 1, 9},
				{token.EOL, "\n", 2, 0},

				{token.LINENUM, "30", 2, 1},
				{token.KEYWORD, "PRINT", 2, 4},
				{token.STRING, "HELLO", 2, 10},
				{token.EOF, "", 2, 17},
			},
		},
		{
			name: "HTAB with expression inside FOR",
			input: `10 FOR A = 1 TO 5
20 HTAB A * 2: PRINT A
30 NEXT A`,
			tokens: []struct {
				typ     token.TokenType
				literal string
				line    int
				column  int
			}{
				{token.LINENUM, "10", 1, 1},
				{token.KEYWORD, "FOR", 1, 4},
				{token.IDENT, "A", 1, 8},
				{token.EQUAL, "=", 1, 10},
				{token.NUMBER, "1", 1, 12},
				{token.KEYWORD, "TO", 1, 14},
				{token.NUMBER, "5", 1, 17},
				{token.EOL, "\n", 2, 0},

				{token.LINENUM, "20", 2, 1},
				{token.KEYWORD, "HTAB", 2, 4},
				{token.IDENT, "A", 2, 9},
				{token.ASTERISK, "*", 2, 11},
				{token.NUMBER, "2", 2, 13},
				{token.COLON, ":", 2, 14},
				{token.KEYWORD, "PRINT", 2, 16},
				{token.IDENT, "A", 2, 22},
				{token.EOL, "\n", 3, 0},

				{token.LINENUM, "30", 3, 1},
				{token.KEYWORD, "NEXT", 3, 4},
				{token.IDENT, "A", 3, 9},
				{token.EOF, "", 3, 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)

			for i, expected := range tt.tokens {
				tok := l.NextToken()

				testutils.Equal(t,
					fmt.Sprintf("[%s][%d] type", tt.name, i),
					tok.Type, expected.typ)

				testutils.Equal(t,
					fmt.Sprintf("[%s][%d] literal", tt.name, i),
					tok.Literal, expected.literal)

				testutils.Equal(t,
					fmt.Sprintf("[%s][%d] line", tt.name, i),
					tok.Line, expected.line)

				testutils.Equal(t,
					fmt.Sprintf("[%s][%d] column", tt.name, i),
					tok.Column, expected.column)
			}
		})
	}
}

package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLexer_PLOT(t *testing.T) {
	tests := []struct {
		name  string
		input string
		toks  []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}
	}{
		{
			name:  "Simple PLOT",
			input: "10 PLOT 10, 20",
			toks: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LINENUM, "10"},
				{token.KEYWORD, "PLOT"},
				{token.NUMBER, "10"},
				{token.COMMA, ","},
				{token.NUMBER, "20"},
				{token.EOF, ""},
			},
		},
		{
			name:  "PLOT with expressions",
			input: "20 PLOT X + 1, Y * 2",
			toks: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LINENUM, "20"},
				{token.KEYWORD, "PLOT"},
				{token.IDENT, "X"},
				{token.PLUS, "+"},
				{token.NUMBER, "1"},
				{token.COMMA, ","},
				{token.IDENT, "Y"},
				{token.ASTERISK, "*"},
				{token.NUMBER, "2"},
				{token.EOF, ""},
			},
		},
		{
			name:  "PLOT in mixed line",
			input: "30 GR : PLOT 5, 5",
			toks: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LINENUM, "30"},
				{token.KEYWORD, "GR"},
				{token.COLON, ":"},
				{token.KEYWORD, "PLOT"},
				{token.NUMBER, "5"},
				{token.COMMA, ","},
				{token.NUMBER, "5"},
				{token.EOF, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for i, et := range tt.toks {
				tok := l.NextToken()
				msg := fmt.Sprintf("%s - token %d wrong type. got=%v, want=%v", tt.name, i, tok.Type, et.expectedType)
				testutils.True(t, msg, tok.Type == et.expectedType)

				msg = fmt.Sprintf("%s - token %d wrong literal. got=%q, want=%q", tt.name, i, tok.Literal, et.expectedLiteral)
				testutils.Equal(t, msg, tok.Literal, et.expectedLiteral)
			}
		})
	}
}

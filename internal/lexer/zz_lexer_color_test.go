package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLexer_COLOR(t *testing.T) {
	tests := []struct {
		name  string
		input string
		toks  []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}
	}{
		{
			name:  "Simple COLOR assignment",
			input: "10 COLOR = 5",
			toks: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LINENUM, "10"},
				{token.KEYWORD, "COLOR"},
				{token.EQUAL, "="},
				{token.NUMBER, "5"},
				{token.EOF, ""},
			},
		},
		{
			name:  "COLOR with expression",
			input: "20 COLOR=A+1",
			toks: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LINENUM, "20"},
				{token.KEYWORD, "COLOR"},
				{token.EQUAL, "="},
				{token.IDENT, "A"},
				{token.PLUS, "+"},
				{token.NUMBER, "1"},
				{token.EOF, ""},
			},
		},
		{
			name:  "COLOR in mixed line",
			input: "30 GR: COLOR=15: PLOT 0,0",
			toks: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LINENUM, "30"},
				{token.KEYWORD, "GR"},
				{token.COLON, ":"},
				{token.KEYWORD, "COLOR"},
				{token.EQUAL, "="},
				{token.NUMBER, "15"},
				{token.COLON, ":"},
				{token.KEYWORD, "PLOT"},
				{token.NUMBER, "0"},
				{token.COMMA, ","},
				{token.NUMBER, "0"},
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

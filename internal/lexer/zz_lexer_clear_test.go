package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLexer_CLEAR_Only(t *testing.T) {
	input := `10 CLEAR`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "CLEAR", 1, 4},
		{token.EOF, "", 1, 9},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		msg := fmt.Sprintf("tests[%d] - token type wrong. got=%q, want=%q",
			i, tok.Type, tt.expectedType)
		testutils.True(t, msg, tok.Type == tt.expectedType)

		msg = fmt.Sprintf("tests[%d] - literal wrong. got=%q, want=%q",
			i, tok.Literal, tt.expectedLiteral)
		testutils.Equal(t, msg, tok.Literal, tt.expectedLiteral)

		msg = fmt.Sprintf("tests[%d] - line wrong. got=%d, want=%d",
			i, tok.Line, tt.expectedLine)
		testutils.True(t, msg, tok.Line == tt.expectedLine)

		msg = fmt.Sprintf("tests[%d] - column wrong. got=%d, want=%d",
			i, tok.Column, tt.expectedColumn)
		testutils.True(t, msg, tok.Column == tt.expectedColumn)
	}
}

func TestLexer_CLEAR_InProgram(t *testing.T) {
	input := `10 LET A = 10
20 CLEAR
30 PRINT A
`

	tests := []struct {
		typ token.TokenType
		lit string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "LET"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.NUMBER, "10"},
		{token.EOL, "\n"},

		{token.LINENUM, "20"},
		{token.KEYWORD, "CLEAR"},
		{token.EOL, "\n"},

		{token.LINENUM, "30"},
		{token.KEYWORD, "PRINT"},
		{token.IDENT, "A"},
		{token.EOL, "\n"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		msg := fmt.Sprintf("tests[%d] - token type wrong. got=%q, want=%q",
			i, tok.Type, tt.typ)
		testutils.True(t, msg, tok.Type == tt.typ)

		msg = fmt.Sprintf("tests[%d] - literal wrong", i)
		testutils.Equal(t, msg, tok.Literal, tt.lit)
	}
}

func TestLexer_CLEAR_WithColon(t *testing.T) {
	input := `10 CLEAR:PRINT "OK"`

	tests := []struct {
		typ token.TokenType
		lit string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "CLEAR"},
		{token.COLON, ":"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "OK"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		testutils.True(t,
			fmt.Sprintf("tests[%d] - type", i),
			tok.Type == tt.typ)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - literal", i),
			tok.Literal, tt.lit)
	}
}

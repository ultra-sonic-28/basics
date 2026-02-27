package lexer

import (
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_NOT(t *testing.T) {
	input := "10 IF NOT A > 0 THEN PRINT NOT 4 - 4"
	
	expected := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "IF"},
		{token.NOT, "NOT"},
		{token.IDENT, "A"},
		{token.GT, ">"},
		{token.NUMBER, "0"},
		{token.KEYWORD, "THEN"},
		{token.KEYWORD, "PRINT"},
		{token.NOT, "NOT"},
		{token.NUMBER, "4"},
		{token.MINUS, "-"},
		{token.NUMBER, "4"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range expected {
		tok := l.NextToken()
		testutils.Equal(t, "token type at index", tok.Type, tt.expectedType)
		testutils.Equal(t, "token literal at index", tok.Literal, tt.expectedLiteral)
		_ = i
	}
}

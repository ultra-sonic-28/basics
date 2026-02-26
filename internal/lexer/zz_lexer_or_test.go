package lexer

import (
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_OR(t *testing.T) {
	input := "10 IF A > 0 OR B < 10 THEN PRINT 8 OR 4"
	
	expected := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "IF"},
		{token.IDENT, "A"},
		{token.GT, ">"},
		{token.NUMBER, "0"},
		{token.OR, "OR"},
		{token.IDENT, "B"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.KEYWORD, "THEN"},
		{token.KEYWORD, "PRINT"},
		{token.NUMBER, "8"},
		{token.OR, "OR"},
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

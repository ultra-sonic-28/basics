package lexer

import (
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_WAIT(t *testing.T) {
	input := "10 WAIT 1000: WAIT DELAI * 5"
	
	expected := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "WAIT"},
		{token.NUMBER, "1000"},
		{token.COLON, ":"},
		{token.KEYWORD, "WAIT"},
		{token.IDENT, "DELAI"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "5"},
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

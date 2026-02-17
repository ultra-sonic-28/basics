package lexer

import (
	"fmt"
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_LEN_Function(t *testing.T) {
	input := `10 REM LEN Function
20 PRINT LEN("APPLESOFT")
30 A = LEN("APPLESOFT")
40 A% = LEN("APPLESOFT" + "WARE")
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 10 REM LEN Function
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.EOL, "\n"},

		// 20 PRINT LEN("APPLESOFT")
		{token.LINENUM, "20"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "LEN"},
		{token.LPAREN, "("},
		{token.STRING, "APPLESOFT"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 30 A = LEN("APPLESOFT")
		{token.LINENUM, "30"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.KEYWORD, "LEN"},
		{token.LPAREN, "("},
		{token.STRING, "APPLESOFT"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 40 A% = LEN("APPLESOFT" + "WARE")
		{token.LINENUM, "40"},
		{token.IDENT, "A%"},
		{token.EQUAL, "="},
		{token.KEYWORD, "LEN"},
		{token.LPAREN, "("},
		{token.STRING, "APPLESOFT"},
		{token.PLUS, "+"},
		{token.STRING, "WARE"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},
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
	}
}

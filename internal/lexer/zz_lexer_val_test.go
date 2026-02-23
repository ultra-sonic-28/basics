package lexer

import (
	"fmt"
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_VAL_Function(t *testing.T) {
	input := `10 REM VAL Function
20 PRINT VAL("159753")
30 A = VAL("159753")
40 A% = VAL("159" + "753")
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 10 REM VAL Function
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.COMMENT, "VAL Function"},
		{token.EOL, "\n"},

		// 20 PRINT VAL("159753")
		{token.LINENUM, "20"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "VAL"},
		{token.LPAREN, "("},
		{token.STRING, "159753"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 30 A = VAL("159753")
		{token.LINENUM, "30"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.KEYWORD, "VAL"},
		{token.LPAREN, "("},
		{token.STRING, "159753"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 40 A% = VAL("159" + "753")
		{token.LINENUM, "40"},
		{token.IDENT, "A%"},
		{token.EQUAL, "="},
		{token.KEYWORD, "VAL"},
		{token.LPAREN, "("},
		{token.STRING, "159"},
		{token.PLUS, "+"},
		{token.STRING, "753"},
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

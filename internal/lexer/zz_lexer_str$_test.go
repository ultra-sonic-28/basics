package lexer

import (
	"fmt"
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_STR_Function(t *testing.T) {
	input := `10 REM STR$ Function
20 A$ = "APPLESOFT"
30 PRINT STR$(2*LEN(A$))
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 10 REM STR$ Function
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.COMMENT, "STR$ Function"},
		{token.EOL, "\n"},

		// 20 A$ = "APPLESOFT"
		{token.LINENUM, "20"},
		{token.IDENT, "A$"},
		{token.EQUAL, "="},
		{token.STRING, "APPLESOFT"},
		{token.EOL, "\n"},

		// 30 PRINT STR$(2*LEN(A$))
		{token.LINENUM, "30"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "STR$"},
		{token.LPAREN, "("},
		{token.NUMBER, "2"},
		{token.ASTERISK, "*"},
		{token.KEYWORD, "LEN"},
		{token.LPAREN, "("},
		{token.IDENT, "A$"},
		{token.RPAREN, ")"},
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

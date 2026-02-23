package lexer

import (
	"fmt"
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_CHR_Function(t *testing.T) {
	input := `10 REM CHR$ Function
20 A = 23
30 PRINT CHR$(2*A)
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 10 REM CHR$ Function
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.COMMENT, "CHR$ Function"},
		{token.EOL, "\n"},

		// 20 A = "123"
		{token.LINENUM, "20"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.NUMBER, "23"},
		{token.EOL, "\n"},

		// 30 PRINT CHR$(2*35)
		{token.LINENUM, "30"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "CHR$"},
		{token.LPAREN, "("},
		{token.NUMBER, "2"},
		{token.ASTERISK, "*"},
		{token.IDENT, "A"},
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

package lexer

import (
	"fmt"
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_ASC_Function(t *testing.T) {
	input := `10 REM ASC Function
20 PRINT ASC("A")
30 A = ASC("A")
40 A% = ASC("ABC")
50 PRINT ASC(A$)
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 10 REM ASC Function
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.COMMENT, "ASC Function"},
		{token.EOL, "\n"},

		// 20 PRINT ASC("A")
		{token.LINENUM, "20"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ASC"},
		{token.LPAREN, "("},
		{token.STRING, "A"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 30 A = ASC("A")
		{token.LINENUM, "30"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.KEYWORD, "ASC"},
		{token.LPAREN, "("},
		{token.STRING, "A"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 40 A% = ASC("ABC")
		{token.LINENUM, "40"},
		{token.IDENT, "A%"},
		{token.EQUAL, "="},
		{token.KEYWORD, "ASC"},
		{token.LPAREN, "("},
		{token.STRING, "ABC"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 50 PRINT ASC(A$)
		{token.LINENUM, "50"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ASC"},
		{token.LPAREN, "("},
		{token.IDENT, "A$"},
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

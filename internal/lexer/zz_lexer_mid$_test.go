package lexer

import (
	"fmt"
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_MID_Function(t *testing.T) {
	input := `10 REM MID$ Function
20 PRINT MID$("APPLESOFT",5)
30 A$ = MID$("APPLESOFT", 5)
40 A=2
50 A$ = MID$("APPLESOFT", A*2 + 1)
60 PRINT MID$("APPLESOFT", 4, 2)
70 B=1
80 PRINT MID$("APPLESOFT", A*2 + 1, B+1)
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 10 REM MID$ Function
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.EOL, "\n"},

		// 20 PRINT MID$("APPLESOFT",5)
		{token.LINENUM, "20"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "MID$"},
		{token.LPAREN, "("},
		{token.STRING, "APPLESOFT"},
		{token.COMMA, ","},
		{token.NUMBER, "5"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 30 A$ = MID$("APPLESOFT", 5)
		{token.LINENUM, "30"},
		{token.IDENT, "A$"},
		{token.EQUAL, "="},
		{token.KEYWORD, "MID$"},
		{token.LPAREN, "("},
		{token.STRING, "APPLESOFT"},
		{token.COMMA, ","},
		{token.NUMBER, "5"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 40 A=2
		{token.LINENUM, "40"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.NUMBER, "2"},
		{token.EOL, "\n"},

		// 50 A$ = MID$("APPLESOFT", A*2 + 1)
		{token.LINENUM, "50"},
		{token.IDENT, "A$"},
		{token.EQUAL, "="},
		{token.KEYWORD, "MID$"},
		{token.LPAREN, "("},
		{token.STRING, "APPLESOFT"},
		{token.COMMA, ","},
		{token.IDENT, "A"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "2"},
		{token.PLUS, "+"},
		{token.NUMBER, "1"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 60 PRINT MID$("APPLESOFT", 4, 2)
		{token.LINENUM, "60"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "MID$"},
		{token.LPAREN, "("},
		{token.STRING, "APPLESOFT"},
		{token.COMMA, ","},
		{token.NUMBER, "4"},
		{token.COMMA, ","},
		{token.NUMBER, "2"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 70 B=1
		{token.LINENUM, "70"},
		{token.IDENT, "B"},
		{token.EQUAL, "="},
		{token.NUMBER, "1"},
		{token.EOL, "\n"},

		// 80 PRINT MID$("APPLESOFT", A*2 + 1, B+1)
		{token.LINENUM, "80"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "MID$"},
		{token.LPAREN, "("},
		{token.STRING, "APPLESOFT"},
		{token.COMMA, ","},
		{token.IDENT, "A"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "2"},
		{token.PLUS, "+"},
		{token.NUMBER, "1"},
		{token.COMMA, ","},
		{token.IDENT, "B"},
		{token.PLUS, "+"},
		{token.NUMBER, "1"},
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

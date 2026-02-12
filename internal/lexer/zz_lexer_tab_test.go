package lexer

import (
	"fmt"
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_TAB_Function(t *testing.T) {
	input := `10 HOME
20 PRINT TAB(25);"TAB Testing"
30 PRINT TAB(3)
40 PRINT TAB(A)
50 PRINT TAB(A*2)
60 PRINT TAB(10);TAB(20);TAB(30)
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 10 HOME
		{token.LINENUM, "10"},
		{token.KEYWORD, "HOME"},
		{token.EOL, "\n"},

		// 20 PRINT TAB(25);"TAB Testing"
		{token.LINENUM, "20"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "TAB"},
		{token.LPAREN, "("},
		{token.NUMBER, "25"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.STRING, "TAB Testing"},
		{token.EOL, "\n"},

		// 30 PRINT TAB(3)
		{token.LINENUM, "30"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "TAB"},
		{token.LPAREN, "("},
		{token.NUMBER, "3"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 40 PRINT TAB(A)
		{token.LINENUM, "40"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "TAB"},
		{token.LPAREN, "("},
		{token.IDENT, "A"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 50 PRINT TAB(A*2)
		{token.LINENUM, "50"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "TAB"},
		{token.LPAREN, "("},
		{token.IDENT, "A"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "2"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 60 PRINT TAB(10);TAB(20);TAB(30)
		{token.LINENUM, "60"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "TAB"},
		{token.LPAREN, "("},
		{token.NUMBER, "10"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.KEYWORD, "TAB"},
		{token.LPAREN, "("},
		{token.NUMBER, "20"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.KEYWORD, "TAB"},
		{token.LPAREN, "("},
		{token.NUMBER, "30"},
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

func TestLexer_TAB_With_Position(t *testing.T) {
	input := `10 PRINT TAB(12)`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "PRINT", 1, 4},
		{token.KEYWORD, "TAB", 1, 10},
		{token.LPAREN, "(", 1, 13},
		{token.NUMBER, "12", 1, 14},
		{token.RPAREN, ")", 1, 16},
		{token.EOF, "", 1, 17},
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

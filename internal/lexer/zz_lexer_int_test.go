package lexer

import (
	"fmt"
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_INT_Function(t *testing.T) {
	input := `10 REM INT Function
20 PRINT INT(1.75)
30 PRINT INT(1.32)
40 PRINT INT(-1.75)
50 PRINT INT(-1.32)
60 A=2.8746841
70 PRINT INT(A)
80 PRINT INT(A*3.74)
90 I%=5
100 PRINT INT(I%)
110 PRINT INT(I%*A)
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 10 REM INT Function
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.EOL, "\n"},

		// 20 PRINT INT(1.75)
		{token.LINENUM, "20"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "INT"},
		{token.LPAREN, "("},
		{token.NUMBER, "1.75"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 30 PRINT INT(1.32)
		{token.LINENUM, "30"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "INT"},
		{token.LPAREN, "("},
		{token.NUMBER, "1.32"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 40 PRINT INT(-1.75)
		{token.LINENUM, "40"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "INT"},
		{token.LPAREN, "("},
		{token.MINUS, "-"},
		{token.NUMBER, "1.75"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 50 PRINT INT(-1.32)
		{token.LINENUM, "50"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "INT"},
		{token.LPAREN, "("},
		{token.MINUS, "-"},
		{token.NUMBER, "1.32"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 60 A=2.8746841
		{token.LINENUM, "60"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.NUMBER, "2.8746841"},
		{token.EOL, "\n"},

		// 70 PRINT INT(A)
		{token.LINENUM, "70"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "INT"},
		{token.LPAREN, "("},
		{token.IDENT, "A"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 80 PRINT INT(A*3.74)
		{token.LINENUM, "80"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "INT"},
		{token.LPAREN, "("},
		{token.IDENT, "A"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "3.74"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 90 I%=5
		{token.LINENUM, "90"},
		{token.IDENT, "I%"},
		{token.EQUAL, "="},
		{token.NUMBER, "5"},
		{token.EOL, "\n"},

		// 100 PRINT INT(I%)
		{token.LINENUM, "100"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "INT"},
		{token.LPAREN, "("},
		{token.IDENT, "I%"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 110 PRINT INT(I%*A)
		{token.LINENUM, "110"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "INT"},
		{token.LPAREN, "("},
		{token.IDENT, "I%"},
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

func TestLexer_INT_Negative_Expression(t *testing.T) {
	input := `10 PRINT INT(-(A+3.2))`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "PRINT", 1, 4},
		{token.KEYWORD, "INT", 1, 10},
		{token.LPAREN, "(", 1, 13},
		{token.MINUS, "-", 1, 14},
		{token.LPAREN, "(", 1, 15},
		{token.IDENT, "A", 1, 16},
		{token.PLUS, "+", 1, 17},
		{token.NUMBER, "3.2", 1, 18},
		{token.RPAREN, ")", 1, 21},
		{token.RPAREN, ")", 1, 22},
		{token.EOF, "", 1, 23},
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

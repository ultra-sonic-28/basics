package lexer

import (
	"fmt"
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_ABS_Function(t *testing.T) {
	input := `10 REM ABS Function
20 PRINT ABS(1.75)
30 PRINT ABS(-1.75)
40 A=2.8746841
50 PRINT ABS(A)
55 PRINT ABS(-A)
60 PRINT ABS(A*3.74)
65 PRINT ABS(-A*3.74)
70 I%=5
80 PRINT ABS(I%)
90 PRINT ABS(-I%)
100 PRINT ABS(I%*A)
110 PRINT ABS(-(I%*A))
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 10 REM ABS Function
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.EOL, "\n"},

		// 20 PRINT ABS(1.75)
		{token.LINENUM, "20"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ABS"},
		{token.LPAREN, "("},
		{token.NUMBER, "1.75"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 30 PRINT ABS(-1.75)
		{token.LINENUM, "30"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ABS"},
		{token.LPAREN, "("},
		{token.MINUS, "-"},
		{token.NUMBER, "1.75"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 40 A=2.8746841
		{token.LINENUM, "40"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.NUMBER, "2.8746841"},
		{token.EOL, "\n"},

		// 50 PRINT ABS(A)
		{token.LINENUM, "50"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ABS"},
		{token.LPAREN, "("},
		{token.IDENT, "A"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 55 PRINT ABS(-A)
		{token.LINENUM, "55"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ABS"},
		{token.LPAREN, "("},
		{token.MINUS, "-"},
		{token.IDENT, "A"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 60 PRINT ABS(A*3.74)
		{token.LINENUM, "60"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ABS"},
		{token.LPAREN, "("},
		{token.IDENT, "A"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "3.74"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 65 PRINT ABS(-A*3.74)
		{token.LINENUM, "65"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ABS"},
		{token.LPAREN, "("},
		{token.MINUS, "-"},
		{token.IDENT, "A"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "3.74"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 70 I%=5
		{token.LINENUM, "70"},
		{token.IDENT, "I%"},
		{token.EQUAL, "="},
		{token.NUMBER, "5"},
		{token.EOL, "\n"},

		// 80 PRINT ABS(I%)
		{token.LINENUM, "80"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ABS"},
		{token.LPAREN, "("},
		{token.IDENT, "I%"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 90 PRINT ABS(-I%)
		{token.LINENUM, "90"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ABS"},
		{token.LPAREN, "("},
		{token.MINUS, "-"},
		{token.IDENT, "I%"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 100 PRINT ABS(I%*A)
		{token.LINENUM, "100"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ABS"},
		{token.LPAREN, "("},
		{token.IDENT, "I%"},
		{token.ASTERISK, "*"},
		{token.IDENT, "A"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 110 PRINT ABS(-I%*A)
		{token.LINENUM, "110"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "ABS"},
		{token.LPAREN, "("},
		{token.MINUS, "-"},
		{token.LPAREN, "("},
		{token.IDENT, "I%"},
		{token.ASTERISK, "*"},
		{token.IDENT, "A"},
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

func TestLexer_ABS_Negative_Expression(t *testing.T) {
	input := `10 PRINT ABS(-(A+3.2))`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "PRINT", 1, 4},
		{token.KEYWORD, "ABS", 1, 10},
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

package lexer

import (
	"fmt"
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestLexer_LOG_Function(t *testing.T) {
	input := `10 REM LOG Function
20 PRINT LOG(1.75)
30 A=2.8746841
40 PRINT LOG(A)
50 PRINT LOG(A*3.74)
60 I%=5
70 PRINT LOG(I%)
80 PRINT LOG(I%*A)
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// 10 REM LOG Function
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.COMMENT, "LOG Function"},
		{token.EOL, "\n"},

		// 20 PRINT LOG(1.75)
		{token.LINENUM, "20"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "LOG"},
		{token.LPAREN, "("},
		{token.NUMBER, "1.75"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 30 A=2.8746841
		{token.LINENUM, "30"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.NUMBER, "2.8746841"},
		{token.EOL, "\n"},

		// 40 PRINT LOG(A)
		{token.LINENUM, "40"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "LOG"},
		{token.LPAREN, "("},
		{token.IDENT, "A"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 50 PRINT LOG(A*3.74)
		{token.LINENUM, "50"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "LOG"},
		{token.LPAREN, "("},
		{token.IDENT, "A"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "3.74"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 60 I%=5
		{token.LINENUM, "60"},
		{token.IDENT, "I%"},
		{token.EQUAL, "="},
		{token.NUMBER, "5"},
		{token.EOL, "\n"},

		// 70 PRINT LOG(I%)
		{token.LINENUM, "70"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "LOG"},
		{token.LPAREN, "("},
		{token.IDENT, "I%"},
		{token.RPAREN, ")"},
		{token.EOL, "\n"},

		// 80 PRINT LOG(I%*A)
		{token.LINENUM, "80"},
		{token.KEYWORD, "PRINT"},
		{token.KEYWORD, "LOG"},
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

func TestLexer_LOG_Negative_Expression(t *testing.T) {
	input := `10 PRINT LOG(-(A+3.2))`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "PRINT", 1, 4},
		{token.KEYWORD, "LOG", 1, 10},
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

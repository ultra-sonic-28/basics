package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLexer_GET_Simple(t *testing.T) {
	input := `10 GET A$`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "GET", 1, 4},
		{token.IDENT, "A$", 1, 8},
		{token.EOF, "", 1, 10},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - token type wrong", i),
			tok.Type, tt.expectedType)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - literal wrong", i),
			tok.Literal, tt.expectedLiteral)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - line wrong", i),
			tok.Line, tt.expectedLine)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - column wrong", i),
			tok.Column, tt.expectedColumn)
	}
}

func TestLexer_GET_InProgram(t *testing.T) {
	input := `1 REM ***** Exemple de sous-routine *****
10 HOME:PRINT "Ce programme est une demo"
20 GOSUB 100
30 PRINT "Test"
40 END
100 PRINT "Pressez une touche";
110 GET X$
120 HOME
130 RETURN
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		// Line 1
		{token.LINENUM, "1", 1, 1},
		{token.KEYWORD, "REM", 1, 3},
		{token.EOL, "\n", 2, 0},

		// Line 10
		{token.LINENUM, "10", 2, 1},
		{token.KEYWORD, "HOME", 2, 4},
		{token.COLON, ":", 2, 8},
		{token.KEYWORD, "PRINT", 2, 9},
		{token.STRING, "Ce programme est une demo", 2, 15},
		{token.EOL, "\n", 3, 0},

		// Line 20
		{token.LINENUM, "20", 3, 1},
		{token.KEYWORD, "GOSUB", 3, 4},
		{token.NUMBER, "100", 3, 10},
		{token.EOL, "\n", 4, 0},

		// Line 30
		{token.LINENUM, "30", 4, 1},
		{token.KEYWORD, "PRINT", 4, 4},
		{token.STRING, "Test", 4, 10},
		{token.EOL, "\n", 5, 0},

		// Line 40
		{token.LINENUM, "40", 5, 1},
		{token.KEYWORD, "END", 5, 4},
		{token.EOL, "\n", 6, 0},

		// Line 100
		{token.LINENUM, "100", 6, 1},
		{token.KEYWORD, "PRINT", 6, 5},
		{token.STRING, "Pressez une touche", 6, 11},
		{token.SEMICOLON, ";", 6, 31},
		{token.EOL, "\n", 7, 0},

		// Line 110 (GET)
		{token.LINENUM, "110", 7, 1},
		{token.KEYWORD, "GET", 7, 5},
		{token.IDENT, "X$", 7, 9},
		{token.EOL, "\n", 8, 0},

		// Line 120
		{token.LINENUM, "120", 8, 1},
		{token.KEYWORD, "HOME", 8, 5},
		{token.EOL, "\n", 9, 0},

		// Line 130
		{token.LINENUM, "130", 9, 1},
		{token.KEYWORD, "RETURN", 9, 5},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - token type wrong", i),
			tok.Type, tt.expectedType)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - literal wrong", i),
			tok.Literal, tt.expectedLiteral)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - line wrong", i),
			tok.Line, tt.expectedLine)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - column wrong", i),
			tok.Column, tt.expectedColumn)
	}
}

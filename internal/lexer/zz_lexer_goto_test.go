package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLexer_GOTO_Number(t *testing.T) {
	input := `10 PRINT "First line"
20 GOTO 40
30 PRINT "Third line"
40 PRINT "Second line"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		// 10 PRINT "First line"
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "PRINT", 1, 4},
		{token.STRING, "First line", 1, 10},
		{token.EOL, "\n", 2, 0},

		// 20 GOTO 40
		{token.LINENUM, "20", 2, 1},
		{token.KEYWORD, "GOTO", 2, 4},
		{token.NUMBER, "40", 2, 9},
		{token.EOL, "\n", 3, 0},

		// 30 PRINT "Third line"
		{token.LINENUM, "30", 3, 1},
		{token.KEYWORD, "PRINT", 3, 4},
		{token.STRING, "Third line", 3, 10},
		{token.EOL, "\n", 4, 0},

		// 40 PRINT "Second line"
		{token.LINENUM, "40", 4, 1},
		{token.KEYWORD, "PRINT", 4, 4},
		{token.STRING, "Second line", 4, 10},
		{token.EOF, "", 4, 23},
	}

	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		testutils.Equal(t, "type", tok.Type, tt.expectedType)
		testutils.Equal(t, "literal", tok.Literal, tt.expectedLiteral)
		testutils.Equal(t, "line", tok.Line, tt.expectedLine)
		testutils.Equal(t, "column", tok.Column, tt.expectedColumn)
	}
}

func TestLexer_GOTO_Identifier(t *testing.T) {
	input := `10 REM GOTO Example
15 JUMP = 80
20 PRINT "First line"
30 GOTO 60
40 PRINT "Third line"
50 GOTO JUMP
60 PRINT "Second line"
70 GOTO 40
80 PRINT "Last line"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.EOL, "\n"},

		{token.LINENUM, "15"},
		{token.IDENT, "JUMP"},
		{token.EQUAL, "="},
		{token.NUMBER, "80"},
		{token.EOL, "\n"},

		{token.LINENUM, "20"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "First line"},
		{token.EOL, "\n"},

		{token.LINENUM, "30"},
		{token.KEYWORD, "GOTO"},
		{token.NUMBER, "60"},
		{token.EOL, "\n"},

		{token.LINENUM, "40"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "Third line"},
		{token.EOL, "\n"},

		{token.LINENUM, "50"},
		{token.KEYWORD, "GOTO"},
		{token.IDENT, "JUMP"},
		{token.EOL, "\n"},

		{token.LINENUM, "60"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "Second line"},
		{token.EOL, "\n"},

		{token.LINENUM, "70"},
		{token.KEYWORD, "GOTO"},
		{token.NUMBER, "40"},
		{token.EOL, "\n"},

		{token.LINENUM, "80"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "Last line"},
		{token.EOF, ""},
	}

	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()
		testutils.Equal(t, "type", tok.Type, tt.expectedType)
		testutils.Equal(t, "literal", tok.Literal, tt.expectedLiteral)
	}
}

func TestLexer_GOTO_Expression(t *testing.T) {
	input := `10 REM GOTO Example
15 JUMP = 40
50 GOTO JUMP * 2`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "REM"},
		{token.EOL, "\n"},

		{token.LINENUM, "15"},
		{token.IDENT, "JUMP"},
		{token.EQUAL, "="},
		{token.NUMBER, "40"},
		{token.EOL, "\n"},

		{token.LINENUM, "50"},
		{token.KEYWORD, "GOTO"},
		{token.IDENT, "JUMP"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "2"},
		{token.EOF, ""},
	}

	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()
		testutils.Equal(t, "type", tok.Type, tt.expectedType)
		testutils.Equal(t, "literal", tok.Literal, tt.expectedLiteral)
	}
}

func TestLexer_GOTO_SameLine(t *testing.T) {
	input := `10 PRINT "First line" : GOTO 30`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "First line"},
		{token.COLON, ":"},
		{token.KEYWORD, "GOTO"},
		{token.NUMBER, "30"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		testutils.Equal(t, fmt.Sprintf("tests[%d]", i), tok.Type, tt.expectedType)
		testutils.Equal(t, fmt.Sprintf("tests[%d]", i), tok.Literal, tt.expectedLiteral)
	}
}

func TestLexer_GOTO_VariableInline(t *testing.T) {
	input := `5 A = 30
10 PRINT "First line" : GOTO A`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "5"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.NUMBER, "30"},
		{token.EOL, "\n"},

		{token.LINENUM, "10"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "First line"},
		{token.COLON, ":"},
		{token.KEYWORD, "GOTO"},
		{token.IDENT, "A"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		testutils.Equal(t, fmt.Sprintf("tests[%d]", i), tok.Type, tt.expectedType)
		testutils.Equal(t, fmt.Sprintf("tests[%d]", i), tok.Literal, tt.expectedLiteral)
	}
}

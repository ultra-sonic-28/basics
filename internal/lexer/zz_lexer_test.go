package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

// TestNextToken_MainCases teste la plupart des cas avec vérification des positions
func TestNextToken_MainCases(t *testing.T) {
	input := `10 PRINT "Hello"
20 REM This is a comment
30 LET X = 5
40 X = X + 1
50 GOTO 10`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		// Ligne 10
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "PRINT", 1, 4},
		{token.STRING, "Hello", 1, 10},
		{token.EOL, "\n", 2, 0},

		// Ligne 20 (REM)
		{token.LINENUM, "20", 2, 1},
		{token.KEYWORD, "REM", 2, 4},
		{token.EOL, "\n", 3, 0}, // le reste de la ligne est ignoré

		// Ligne 30
		{token.LINENUM, "30", 3, 1},
		{token.KEYWORD, "LET", 3, 4},
		{token.IDENT, "X", 3, 8},
		{token.EQUAL, "=", 3, 10},
		{token.NUMBER, "5", 3, 12},
		{token.EOL, "\n", 4, 0},

		// Ligne 40
		{token.LINENUM, "40", 4, 1},
		{token.IDENT, "X", 4, 4},
		{token.EQUAL, "=", 4, 6},
		{token.IDENT, "X", 4, 8},
		{token.PLUS, "+", 4, 10},
		{token.NUMBER, "1", 4, 12},
		{token.EOL, "\n", 5, 0},

		// Ligne 50
		{token.LINENUM, "50", 5, 1},
		{token.KEYWORD, "GOTO", 5, 4},
		{token.NUMBER, "10", 5, 9},
		{token.EOF, "", 5, 11},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		msg := fmt.Sprintf("tests[%d] - token type wrong. got=%q, want=%q", i, tok.Type, tt.expectedType)
		testutils.True(t, msg, tok.Type == tt.expectedType)

		msg = fmt.Sprintf("tests[%d] - literal wrong. got=%q, want=%q", i, tok.Literal, tt.expectedLiteral)
		testutils.Equal(t, msg, tok.Literal, tt.expectedLiteral)

		msg = fmt.Sprintf("tests[%d] - line wrong. got=%d, want=%d", i, tok.Line, tt.expectedLine)
		testutils.True(t, msg, tok.Line == tt.expectedLine)

		msg = fmt.Sprintf("tests[%d] - column wrong. got=%d, want=%d", i, tok.Column, tt.expectedColumn)
		testutils.True(t, msg, tok.Column == tt.expectedColumn)
	}
}

// TestLexer_OperatorsAndPunctuation avec positions
func TestLexer_OperatorsAndPunctuation(t *testing.T) {
	input := `()+-*/^=,:;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LPAREN, "(", 1, 1},
		{token.RPAREN, ")", 1, 2},
		{token.PLUS, "+", 1, 3},
		{token.MINUS, "-", 1, 4},
		{token.ASTERISK, "*", 1, 5},
		{token.SLASH, "/", 1, 6},
		{token.CARET, "^", 1, 7},
		{token.EQUAL, "=", 1, 8},
		{token.COMMA, ",", 1, 9},
		{token.COLON, ":", 1, 10},
		{token.SEMICOLON, ";", 1, 11},
		{token.EOF, "", 1, 12},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		msg := fmt.Sprintf("tests[%d] - token type wrong. got=%q, want=%q", i, tok.Type, tt.expectedType)
		testutils.True(t, msg, tok.Type == tt.expectedType)

		msg = fmt.Sprintf("tests[%d] - literal wrong. got=%q, want=%q", i, tok.Literal, tt.expectedLiteral)
		testutils.Equal(t, msg, tok.Literal, tt.expectedLiteral)

		msg = fmt.Sprintf("tests[%d] - line wrong. got=%d, want=%d", i, tok.Line, tt.expectedLine)
		testutils.True(t, msg, tok.Line == tt.expectedLine)

		msg = fmt.Sprintf("tests[%d] - column wrong. got=%d, want=%d", i, tok.Column, tt.expectedColumn)
		testutils.True(t, msg, tok.Column == tt.expectedColumn)
	}
}

// TestLexer_IdentifiersAndKeywords avec positions
func TestLexer_IdentifiersAndKeywords(t *testing.T) {
	input := `10 LET X = 42 PRINT REM`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "LET", 1, 4},
		{token.IDENT, "X", 1, 8},
		{token.EQUAL, "=", 1, 10},
		{token.NUMBER, "42", 1, 12},
		{token.KEYWORD, "PRINT", 1, 15},
		{token.KEYWORD, "REM", 1, 21},
		{token.EOF, "", 1, 24},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		msg := fmt.Sprintf("tests[%d] - token type wrong. got=%q, want=%q", i, tok.Type, tt.expectedType)
		testutils.True(t, msg, tok.Type == tt.expectedType)

		msg = fmt.Sprintf("tests[%d] - literal wrong. got=%q, want=%q", i, tok.Literal, tt.expectedLiteral)
		testutils.Equal(t, msg, tok.Literal, tt.expectedLiteral)

		msg = fmt.Sprintf("tests[%d] - line wrong. got=%d, want=%d", i, tok.Line, tt.expectedLine)
		testutils.True(t, msg, tok.Line == tt.expectedLine)

		msg = fmt.Sprintf("tests[%d] - column wrong. got=%d, want=%d", i, tok.Column, tt.expectedColumn)
		testutils.True(t, msg, tok.Column == tt.expectedColumn)
	}
}

// TestLexer_StringLiterals avec positions
func TestLexer_StringLiterals(t *testing.T) {
	input := `"Hello, world" "Another string"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.STRING, "Hello, world", 1, 1},
		{token.STRING, "Another string", 1, 16},
		{token.EOF, "", 1, 32},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		msg := fmt.Sprintf("tests[%d] - token type wrong. got=%q, want=%q", i, tok.Type, tt.expectedType)
		testutils.True(t, msg, tok.Type == tt.expectedType)

		msg = fmt.Sprintf("tests[%d] - literal wrong. got=%q, want=%q", i, tok.Literal, tt.expectedLiteral)
		testutils.Equal(t, msg, tok.Literal, tt.expectedLiteral)

		msg = fmt.Sprintf("tests[%d] - line wrong. got=%d, want=%d", i, tok.Line, tt.expectedLine)
		testutils.True(t, msg, tok.Line == tt.expectedLine)

		msg = fmt.Sprintf("tests[%d] - column wrong. got=%d, want=%d", i, tok.Column, tt.expectedColumn)
		testutils.True(t, msg, tok.Column == tt.expectedColumn)
	}
}

// TestLexer_Numbers avec positions
func TestLexer_Numbers(t *testing.T) {
	input := `10 123 45.67 0.89`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.NUMBER, "123", 1, 4},
		{token.NUMBER, "45.67", 1, 8},
		{token.NUMBER, "0.89", 1, 14},
		{token.EOF, "", 1, 18},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		msg := fmt.Sprintf("tests[%d] - token type wrong. got=%q, want=%q", i, tok.Type, tt.expectedType)
		testutils.True(t, msg, tok.Type == tt.expectedType)

		msg = fmt.Sprintf("tests[%d] - literal wrong. got=%q, want=%q", i, tok.Literal, tt.expectedLiteral)
		testutils.Equal(t, msg, tok.Literal, tt.expectedLiteral)

		msg = fmt.Sprintf("tests[%d] - line wrong. got=%d, want=%d", i, tok.Line, tt.expectedLine)
		testutils.True(t, msg, tok.Line == tt.expectedLine)

		msg = fmt.Sprintf("tests[%d] - column wrong. got=%d, want=%d", i, tok.Column, tt.expectedColumn)
		testutils.True(t, msg, tok.Column == tt.expectedColumn)
	}
}

// TestLexer_EOLandEOF avec positions
func TestLexer_EOLandEOF(t *testing.T) {
	input := "10 LET X=1\n20 PRINT X\n"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "LET", 1, 4},
		{token.IDENT, "X", 1, 8},
		{token.EQUAL, "=", 1, 9},
		{token.NUMBER, "1", 1, 10},
		{token.EOL, "\n", 2, 0},
		{token.LINENUM, "20", 2, 1},
		{token.KEYWORD, "PRINT", 2, 4},
		{token.IDENT, "X", 2, 10},
		{token.EOL, "\n", 3, 0},
		{token.EOF, "", 3, 1},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()

		msg := fmt.Sprintf("tests[%d] - token type wrong. got=%q, want=%q", i, tok.Type, tt.expectedType)
		testutils.True(t, msg, tok.Type == tt.expectedType)

		msg = fmt.Sprintf("tests[%d] - literal wrong. got=%q, want=%q", i, tok.Literal, tt.expectedLiteral)
		testutils.Equal(t, msg, tok.Literal, tt.expectedLiteral)

		msg = fmt.Sprintf("tests[%d] - line wrong. got=%d, want=%d", i, tok.Line, tt.expectedLine)
		testutils.True(t, msg, tok.Line == tt.expectedLine)

		msg = fmt.Sprintf("tests[%d] - column wrong. got=%d, want=%d", i, tok.Column, tt.expectedColumn)
		testutils.True(t, msg, tok.Column == tt.expectedColumn)
	}
}

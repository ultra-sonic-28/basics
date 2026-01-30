package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLexer_FLASH_Multiline(t *testing.T) {
	input := `10 NORMAL
20 PRINT "HELLO"
30 FLASH
40 PRINT "WORLD"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "NORMAL", 1, 4},
		{token.EOL, "\n", 2, 0},

		{token.LINENUM, "20", 2, 1},
		{token.KEYWORD, "PRINT", 2, 4},
		{token.STRING, "HELLO", 2, 10},
		{token.EOL, "\n", 3, 0},

		{token.LINENUM, "30", 3, 1},
		{token.KEYWORD, "FLASH", 3, 4},
		{token.EOL, "\n", 4, 0},

		{token.LINENUM, "40", 4, 1},
		{token.KEYWORD, "PRINT", 4, 4},
		{token.STRING, "WORLD", 4, 10},

		{token.EOF, "", 4, 17},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		testutils.True(t,
			fmt.Sprintf("tests[%d] - token type wrong, want %d got %d", i, tt.expectedType, tok.Type),
			tok.Type == tt.expectedType,
		)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - literal wrong", i),
			tok.Literal,
			tt.expectedLiteral,
		)

		testutils.True(t,
			fmt.Sprintf("tests[%d] - line wrong, want %s got %s", i, tt.expectedLiteral, tok.Literal),
			tok.Line == tt.expectedLine,
		)

		testutils.True(t,
			fmt.Sprintf("tests[%d] - column wrong, want %d got %d", i, tt.expectedColumn, tok.Column),
			tok.Column == tt.expectedColumn,
		)
	}
}

func TestLexer_FLASH_WithSemicolon(t *testing.T) {
	input := `10 NORMAL
20 PRINT "HELLO ";
30 FLASH
40 PRINT "WORLD";
50 NORMAL
60 PRINT "!"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "NORMAL", 1, 4},
		{token.EOL, "\n", 2, 0},

		{token.LINENUM, "20", 2, 1},
		{token.KEYWORD, "PRINT", 2, 4},
		{token.STRING, "HELLO ", 2, 10},
		{token.SEMICOLON, ";", 2, 18},
		{token.EOL, "\n", 3, 0},

		{token.LINENUM, "30", 3, 1},
		{token.KEYWORD, "FLASH", 3, 4},
		{token.EOL, "\n", 4, 0},

		{token.LINENUM, "40", 4, 1},
		{token.KEYWORD, "PRINT", 4, 4},
		{token.STRING, "WORLD", 4, 10},
		{token.SEMICOLON, ";", 4, 17},
		{token.EOL, "\n", 5, 0},

		{token.LINENUM, "50", 5, 1},
		{token.KEYWORD, "NORMAL", 5, 4},
		{token.EOL, "\n", 6, 0},

		{token.LINENUM, "60", 6, 1},
		{token.KEYWORD, "PRINT", 6, 4},
		{token.STRING, "!", 6, 10},

		{token.EOF, "", 6, 13},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		testutils.True(t,
			fmt.Sprintf("tests[%d] - token type wrong, want %d got %d", i, tt.expectedType, tok.Type),
			tok.Type == tt.expectedType,
		)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - literal wrong", i),
			tok.Literal,
			tt.expectedLiteral,
		)

		testutils.True(t,
			fmt.Sprintf("tests[%d] - line wrong, want %s got %s", i, tt.expectedLiteral, tok.Literal),
			tok.Line == tt.expectedLine,
		)

		testutils.True(t,
			fmt.Sprintf("tests[%d] - column wrong, want %d got %d", i, tt.expectedColumn, tok.Column),
			tok.Column == tt.expectedColumn,
		)
	}
}

func TestLexer_FLASH_SingleLineWithColon(t *testing.T) {
	input := `10 NORMAL:PRINT "HELLO ";:FLASH:PRINT "WORLD";:NORMAL:PRINT "!"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "NORMAL", 1, 4},
		{token.COLON, ":", 1, 10},

		{token.KEYWORD, "PRINT", 1, 11},
		{token.STRING, "HELLO ", 1, 17},
		{token.SEMICOLON, ";", 1, 25},
		{token.COLON, ":", 1, 26},

		{token.KEYWORD, "FLASH", 1, 27},
		{token.COLON, ":", 1, 32},

		{token.KEYWORD, "PRINT", 1, 33},
		{token.STRING, "WORLD", 1, 39},
		{token.SEMICOLON, ";", 1, 46},
		{token.COLON, ":", 1, 47},

		{token.KEYWORD, "NORMAL", 1, 48},
		{token.COLON, ":", 1, 54},

		{token.KEYWORD, "PRINT", 1, 55},
		{token.STRING, "!", 1, 61},

		{token.EOF, "", 1, 64},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		testutils.True(t,
			fmt.Sprintf("tests[%d] - token type wrong, want %d got %d", i, tt.expectedType, tok.Type),
			tok.Type == tt.expectedType,
		)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - literal wrong", i),
			tok.Literal,
			tt.expectedLiteral,
		)

		testutils.True(t,
			fmt.Sprintf("tests[%d] - line wrong, want %s got %s", i, tt.expectedLiteral, tok.Literal),
			tok.Line == tt.expectedLine,
		)

		testutils.True(t,
			fmt.Sprintf("tests[%d] - column wrong, want %d got %d", i, tt.expectedColumn, tok.Column),
			tok.Column == tt.expectedColumn,
		)
	}
}

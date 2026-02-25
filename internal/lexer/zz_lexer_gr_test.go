package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLexer_GR_AfterPrint(t *testing.T) {
	input := `10 PRINT "HELLO": GR`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "PRINT", 1, 4},
		{token.STRING, "HELLO", 1, 10},
		{token.COLON, ":", 1, 17},
		{token.KEYWORD, "GR", 1, 19},
		{token.EOF, "", 1, 21},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		testutils.True(t,
			fmt.Sprintf("tests[%d] - token type wrong", i),
			tok.Type == tt.expectedType,
		)

		testutils.Equal(t,
			fmt.Sprintf("tests[%d] - literal wrong", i),
			tok.Literal,
			tt.expectedLiteral,
		)

		testutils.True(t,
			fmt.Sprintf("tests[%d] - line wrong", i),
			tok.Line == tt.expectedLine,
		)

		testutils.True(t,
			fmt.Sprintf("tests[%d] - column wrong", i),
			tok.Column == tt.expectedColumn,
		)
	}
}

func TestLexer_GR_Only(t *testing.T) {
	input := `10 GR`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LINENUM, "10", 1, 1},
		{token.KEYWORD, "GR", 1, 4},
		{token.EOF, "", 1, 6},
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

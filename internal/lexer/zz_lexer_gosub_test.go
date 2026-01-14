package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLexer_GOSUB_RETURN_Simple(t *testing.T) {
	input := `10 PRINT "Hello "
20 GOSUB 100
30 PRINT "!!!"
40 END
100 PRINT "World"
110 RETURN`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "Hello "},
		{token.EOL, "\n"},

		{token.LINENUM, "20"},
		{token.KEYWORD, "GOSUB"},
		{token.NUMBER, "100"},
		{token.EOL, "\n"},

		{token.LINENUM, "30"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "!!!"},
		{token.EOL, "\n"},

		{token.LINENUM, "40"},
		{token.KEYWORD, "END"},
		{token.EOL, "\n"},

		{token.LINENUM, "100"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "World"},
		{token.EOL, "\n"},

		{token.LINENUM, "110"},
		{token.KEYWORD, "RETURN"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		testutils.Equal(t, fmt.Sprintf("tests[%d] type", i), tok.Type, tt.expectedType)
		testutils.Equal(t, fmt.Sprintf("tests[%d] literal", i), tok.Literal, tt.expectedLiteral)
	}
}

func TestLexer_GOSUB_RETURN_Inline(t *testing.T) {
	input := `10 PRINT "Hello " : GOSUB 100 : PRINT "!!!"
30 END
100 PRINT "World" : RETURN`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "Hello "},
		{token.COLON, ":"},
		{token.KEYWORD, "GOSUB"},
		{token.NUMBER, "100"},
		{token.COLON, ":"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "!!!"},
		{token.EOL, "\n"},

		{token.LINENUM, "30"},
		{token.KEYWORD, "END"},
		{token.EOL, "\n"},

		{token.LINENUM, "100"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "World"},
		{token.COLON, ":"},
		{token.KEYWORD, "RETURN"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		testutils.Equal(t, fmt.Sprintf("tests[%d] type", i), tok.Type, tt.expectedType)
		testutils.Equal(t, fmt.Sprintf("tests[%d] literal", i), tok.Literal, tt.expectedLiteral)
	}
}

func TestLexer_GOSUB_RETURN_WithLoop(t *testing.T) {
	input := `5 REM **** Ce programme affiche la table de 4 ****
10 PRINT "TABLE DE 4 :"
20 FOR I=1 TO 10
25 GOSUB 100
30 PRINT I, V
40 NEXT I
50 END
100 V = I * 4 : RETURN`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "5"},
		{token.KEYWORD, "REM"},
		{token.EOL, "\n"},

		{token.LINENUM, "10"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "TABLE DE 4 :"},
		{token.EOL, "\n"},

		{token.LINENUM, "20"},
		{token.KEYWORD, "FOR"},
		{token.IDENT, "I"},
		{token.EQUAL, "="},
		{token.NUMBER, "1"},
		{token.KEYWORD, "TO"},
		{token.NUMBER, "10"},
		{token.EOL, "\n"},

		{token.LINENUM, "25"},
		{token.KEYWORD, "GOSUB"},
		{token.NUMBER, "100"},
		{token.EOL, "\n"},

		{token.LINENUM, "30"},
		{token.KEYWORD, "PRINT"},
		{token.IDENT, "I"},
		{token.COMMA, ","},
		{token.IDENT, "V"},
		{token.EOL, "\n"},

		{token.LINENUM, "40"},
		{token.KEYWORD, "NEXT"},
		{token.IDENT, "I"},
		{token.EOL, "\n"},

		{token.LINENUM, "50"},
		{token.KEYWORD, "END"},
		{token.EOL, "\n"},

		{token.LINENUM, "100"},
		{token.IDENT, "V"},
		{token.EQUAL, "="},
		{token.IDENT, "I"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "4"},
		{token.COLON, ":"},
		{token.KEYWORD, "RETURN"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		testutils.Equal(t, fmt.Sprintf("tests[%d] type", i), tok.Type, tt.expectedType)
		testutils.Equal(t, fmt.Sprintf("tests[%d] literal", i), tok.Literal, tt.expectedLiteral)
	}
}

func TestLexer_GOSUB_Expression(t *testing.T) {
	input := `10 PRINT "Hello ":A=50:GOSUB A*2:PRINT "!!!"
30 END
100 PRINT "World":RETURN`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LINENUM, "10"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "Hello "},
		{token.COLON, ":"},
		{token.IDENT, "A"},
		{token.EQUAL, "="},
		{token.NUMBER, "50"},
		{token.COLON, ":"},
		{token.KEYWORD, "GOSUB"},
		{token.IDENT, "A"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "2"},
		{token.COLON, ":"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "!!!"},
		{token.EOL, "\n"},

		{token.LINENUM, "30"},
		{token.KEYWORD, "END"},
		{token.EOL, "\n"},

		{token.LINENUM, "100"},
		{token.KEYWORD, "PRINT"},
		{token.STRING, "World"},
		{token.COLON, ":"},
		{token.KEYWORD, "RETURN"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		testutils.Equal(t, fmt.Sprintf("tests[%d] type", i), tok.Type, tt.expectedType)
		testutils.Equal(t, fmt.Sprintf("tests[%d] literal", i), tok.Literal, tt.expectedLiteral)
	}
}

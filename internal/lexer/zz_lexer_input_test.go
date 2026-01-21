package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLexer_INPUT_SimpleAndPrompt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name: "INPUT with variables and PRINT",
			input: `10 PRINT "Entrez votre nom : ";
20 INPUT N$
30 PRINT "Entrez votre age : ";
40 INPUT A
50 PRINT:PRINT N$;", vous avez ";A;" ans"`,
			expected: []token.Token{
				{Type: token.LINENUM, Literal: "10"},
				{Type: token.KEYWORD, Literal: "PRINT"},
				{Type: token.STRING, Literal: "Entrez votre nom : "},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.EOL, Literal: "\n"},

				{Type: token.LINENUM, Literal: "20"},
				{Type: token.KEYWORD, Literal: "INPUT"},
				{Type: token.IDENT, Literal: "N$"},
				{Type: token.EOL, Literal: "\n"},

				{Type: token.LINENUM, Literal: "30"},
				{Type: token.KEYWORD, Literal: "PRINT"},
				{Type: token.STRING, Literal: "Entrez votre age : "},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.EOL, Literal: "\n"},

				{Type: token.LINENUM, Literal: "40"},
				{Type: token.KEYWORD, Literal: "INPUT"},
				{Type: token.IDENT, Literal: "A"},
				{Type: token.EOL, Literal: "\n"},

				{Type: token.LINENUM, Literal: "50"},
				{Type: token.KEYWORD, Literal: "PRINT"},
				{Type: token.COLON, Literal: ":"},
				{Type: token.KEYWORD, Literal: "PRINT"},
				{Type: token.IDENT, Literal: "N$"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.STRING, Literal: ", vous avez "},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.IDENT, Literal: "A"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.STRING, Literal: " ans"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := Lex(tc.input)

			testutils.Equal(
				t,
				"token count mismatch",
				len(tokens),
				len(tc.expected),
			)

			for i := range tc.expected {
				testutils.Equal(
					t,
					fmt.Sprintf("token[%d] type mismatch", i),
					tokens[i].Type,
					tc.expected[i].Type,
				)

				testutils.Equal(
					t,
					fmt.Sprintf("token[%d] literal mismatch", i),
					tokens[i].Literal,
					tc.expected[i].Literal,
				)
			}
		})
	}
}

func TestLexer_INPUT_MultipleVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name: "INPUT with prompt and multiple vars",
			input: `10 PRINT "Multiply 2 numbers"
20 INPUT "Enter 2 values: ";A,B
30 PRINT "A*B is ";A*B`,
			expected: []token.Token{
				{Type: token.LINENUM, Literal: "10"},
				{Type: token.KEYWORD, Literal: "PRINT"},
				{Type: token.STRING, Literal: "Multiply 2 numbers"},
				{Type: token.EOL, Literal: "\n"},

				{Type: token.LINENUM, Literal: "20"},
				{Type: token.KEYWORD, Literal: "INPUT"},
				{Type: token.STRING, Literal: "Enter 2 values: "},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.IDENT, Literal: "A"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "B"},
				{Type: token.EOL, Literal: "\n"},

				{Type: token.LINENUM, Literal: "30"},
				{Type: token.KEYWORD, Literal: "PRINT"},
				{Type: token.STRING, Literal: "A*B is "},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.IDENT, Literal: "A"},
				{Type: token.ASTERISK, Literal: "*"},
				{Type: token.IDENT, Literal: "B"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := Lex(tc.input)

			testutils.Equal(
				t,
				"token count mismatch",
				len(tokens),
				len(tc.expected),
			)

			for i := range tc.expected {
				testutils.Equal(
					t,
					fmt.Sprintf("token[%d] type mismatch", i),
					tokens[i].Type,
					tc.expected[i].Type,
				)

				testutils.Equal(
					t,
					fmt.Sprintf("token[%d] literal mismatch", i),
					tokens[i].Literal,
					tc.expected[i].Literal,
				)
			}
		})
	}
}

package token

import (
	"basics/testutils"
	"fmt"
	"testing"
)

func TestToken_TypeName(t *testing.T) {
	tests := []struct {
		tok      Token
		expected string
	}{
		// Spéciaux
		{Token{Type: ILLEGAL}, "ILLEGAL"},
		{Token{Type: EOF}, "EOF"},
		{Token{Type: EOL}, "EOL"},

		// Spéciaux BASIC
		{Token{Type: LINENUM}, "LINENUM"},

		// Littéraux
		{Token{Type: NUMBER}, "NUMBER"},
		{Token{Type: STRING}, "STRING"},
		{Token{Type: IDENT}, "IDENT"},

		// Opérateurs
		{Token{Type: PLUS}, "+"},
		{Token{Type: MINUS}, "-"},
		{Token{Type: ASTERISK}, "*"},
		{Token{Type: SLASH}, "/"},
		{Token{Type: CARET}, "^"},
		{Token{Type: EQUAL}, "="},
		{Token{Type: LT}, "<"},
		{Token{Type: LTE}, "<="},
		{Token{Type: GT}, ">"},
		{Token{Type: GTE}, ">="},
		{Token{Type: NEQ}, "<>"},

		// Délimiteurs
		{Token{Type: LPAREN}, "("},
		{Token{Type: RPAREN}, ")"},
		{Token{Type: COMMA}, ","},
		{Token{Type: COLON}, ":"},
		{Token{Type: SEMICOLON}, ";"},

		// Keyword
		{Token{Type: KEYWORD}, "KEYWORD"},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("tests[%d] - TypeName() wrong:", i)
		testutils.Equal(t, msg, tt.tok.TypeName(), tt.expected)
	}
}

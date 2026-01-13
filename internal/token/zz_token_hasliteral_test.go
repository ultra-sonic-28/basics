package token

import (
	"basics/testutils"
	"fmt"
	"testing"
)

func TestToken_HasLiteral(t *testing.T) {
	tests := []struct {
		tok      Token
		expected bool
	}{
		// Spéciaux
		{Token{Type: ILLEGAL}, true},
		{Token{Type: EOF}, false},
		{Token{Type: EOL}, false},

		// Spéciaux BASIC
		{Token{Type: LINENUM}, true},

		// Littéraux
		{Token{Type: NUMBER}, true},
		{Token{Type: STRING}, true},
		{Token{Type: IDENT}, true},

		// Opérateurs
		{Token{Type: PLUS}, true},
		{Token{Type: MINUS}, true},
		{Token{Type: ASTERISK}, true},
		{Token{Type: SLASH}, true},
		{Token{Type: CARET}, true},
		{Token{Type: EQUAL}, true},
		{Token{Type: LT}, true},
		{Token{Type: LTE}, true},
		{Token{Type: GT}, true},
		{Token{Type: GTE}, true},
		{Token{Type: NEQ}, true},

		// Délimiteurs
		{Token{Type: LPAREN}, true},
		{Token{Type: RPAREN}, true},
		{Token{Type: COMMA}, true},
		{Token{Type: COLON}, true},
		{Token{Type: SEMICOLON}, true},

		// Keyword
		{Token{Type: KEYWORD}, true},
	}

	for i, tt := range tests {
		got := tt.tok.HasLiteral()
		msg := fmt.Sprintf("tests[%d] - HasLiteral() wrong: got=%v, want=%v", i, got, tt.expected)
		testutils.True(t, msg, got == tt.expected)
	}
}

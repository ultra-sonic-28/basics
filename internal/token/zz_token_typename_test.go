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
		{Token{Type: ILLEGAL}, "ILLEGAL"},
		{Token{Type: EOF}, "EOF"},
		{Token{Type: EOL}, "EOL"},
		{Token{Type: LINENUM}, "LINENUM"},
		{Token{Type: NUMBER}, "NUMBER"},
		{Token{Type: STRING}, "STRING"},
		{Token{Type: IDENT}, "IDENT"},
		{Token{Type: PLUS}, "+"},
		{Token{Type: MINUS}, "-"},
		{Token{Type: KEYWORD}, "KEYWORD"},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("tests[%d] - TypeName() wrong:", i)
		testutils.Equal(t, msg, tt.tok.TypeName(), tt.expected)
	}
}

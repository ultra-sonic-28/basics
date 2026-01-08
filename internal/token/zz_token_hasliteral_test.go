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
		{Token{Type: ILLEGAL}, true},
		{Token{Type: EOF}, false},
		{Token{Type: EOL}, false},
		{Token{Type: LINENUM}, true},
		{Token{Type: NUMBER}, true},
		{Token{Type: STRING}, true},
		{Token{Type: IDENT}, true},
		{Token{Type: PLUS}, true},
		{Token{Type: KEYWORD}, true},
	}

	for i, tt := range tests {
		got := tt.tok.HasLiteral()
		msg := fmt.Sprintf("tests[%d] - HasLiteral() wrong: got=%v, want=%v", i, got, tt.expected)
		testutils.True(t, msg, got == tt.expected)
	}
}

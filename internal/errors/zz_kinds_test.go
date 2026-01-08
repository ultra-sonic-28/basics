package errors

import (
	"basics/testutils"
	"fmt"
	"testing"
)

func TestKind_String(t *testing.T) {
	tests := []struct {
		name string
		kind Kind
		want string
	}{
		{
			name: "Lexical error",
			kind: Lexical,
			want: "LEXICAL ERROR",
		},
		{
			name: "Syntax error",
			kind: Syntax,
			want: "SYNTAX ERROR",
		},
		{
			name: "Semantic error",
			kind: Semantic,
			want: "SEMANTIC ERROR",
		},
		{
			name: "Unknown error kind",
			kind: Kind(99),
			want: "ERROR",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := fmt.Sprintf("tests[%d] Kind.String() wrong:", i)
			testutils.Equal(t, msg, tt.kind.String(), tt.want)
		})
	}
}

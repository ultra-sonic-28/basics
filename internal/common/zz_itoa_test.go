package common

import (
	"testing"

	"basics/testutils"
)

func TestItoa(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want string
	}{
		{
			name: "zero",
			in:   0,
			want: "0",
		},
		{
			name: "positive number",
			in:   42,
			want: "42",
		},
		{
			name: "negative number",
			in:   -7,
			want: "-7",
		},
		{
			name: "large number",
			in:   123456789,
			want: "123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Itoa(tt.in)
			testutils.Equal(
				t,
				"Itoa("+tt.name+")",
				got,
				tt.want,
			)
		})
	}
}

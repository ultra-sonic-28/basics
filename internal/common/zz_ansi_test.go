package common

import (
	"basics/testutils"
	"fmt"
	"testing"
)

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no ANSI",
			input: "Hello World",
			want:  "Hello World",
		},
		{
			name:  "simple color",
			input: "\x1b[31mRed\x1b[0m",
			want:  "Red",
		},
		{
			name:  "multiple attributes",
			input: "\x1b[1;32mGreen Bold\x1b[0m",
			want:  "Green Bold",
		},
		{
			name:  "mixed text",
			input: "Error: \x1b[31mFAIL\x1b[0m!",
			want:  "Error: FAIL!",
		},
		{
			name:  "cursor movement",
			input: "\x1b[2K\rLine",
			want:  "\rLine",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripANSI(tt.input)
			str := fmt.Sprintf(
				"StripANSI(%q): got %q, want %q",
				tt.input,
				got,
				tt.want,
			)
			testutils.True(t, str, got == tt.want)
		})
	}
}

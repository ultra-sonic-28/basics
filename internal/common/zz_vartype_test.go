package common

import (
	"fmt"
	"testing"

	"basics/internal/runtime"
	"basics/testutils"
)

func TestVarType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "integer variable",
			input:    "A%",
			expected: "int",
		},
		{
			name:     "string variable",
			input:    "A$",
			expected: "string",
		},
		{
			name:     "float variable default",
			input:    "A",
			expected: "float",
		},
		{
			name:     "float variable with number",
			input:    "A1",
			expected: "float",
		},
		{
			name:     "suffix only percent",
			input:    "%",
			expected: "int",
		},
		{
			name:     "suffix only dollar",
			input:    "$",
			expected: "string",
		},
		{
			name:     "empty name defaults to float",
			input:    "",
			expected: "float",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VarType(tt.input)

			str := fmt.Sprintf(
				"VarType(%q): got %q, want %q",
				tt.input,
				got,
				tt.expected,
			)
			testutils.True(t, str, got == tt.expected)
		})
	}
}

func TestVarTypeAsInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected runtime.ValueType
	}{
		{
			name:     "integer variable",
			input:    "A%",
			expected: runtime.INTEGER,
		},
		{
			name:     "string variable",
			input:    "A$",
			expected: runtime.STRING,
		},
		{
			name:     "number variable default",
			input:    "A",
			expected: runtime.NUMBER,
		},
		{
			name:     "number variable with digits",
			input:    "X42",
			expected: runtime.NUMBER,
		},
		{
			name:     "suffix only percent",
			input:    "%",
			expected: runtime.INTEGER,
		},
		{
			name:     "suffix only dollar",
			input:    "$",
			expected: runtime.STRING,
		},
		{
			name:     "empty name defaults to number",
			input:    "",
			expected: runtime.NUMBER,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VarTypeAsInt(tt.input)

			str := fmt.Sprintf(
				"VarTypeAsInt(%q): got %q, want %q",
				tt.input,
				got,
				tt.expected,
			)
			testutils.True(t, str, got == tt.expected)
		})
	}
}

package interpreter

import (
	"fmt"
	"testing"

	"basics/testutils"
)

func TestVarType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Integer variable",
			input:    "A%",
			expected: "int",
		},
		{
			name:     "String variable",
			input:    "A$",
			expected: "string",
		},
		{
			name:     "Float variable default",
			input:    "A",
			expected: "float",
		},
		{
			name:     "Float with digits",
			input:    "X1",
			expected: "float",
		},
		{
			name:     "Integer with digits",
			input:    "COUNT%",
			expected: "int",
		},
		{
			name:     "String with digits",
			input:    "NAME1$",
			expected: "string",
		},
		{
			name:     "Suffix not at end (%)",
			input:    "A%B",
			expected: "float",
		},
		{
			name:     "Suffix not at end ($)",
			input:    "A$B",
			expected: "float",
		},
		{
			name:     "Empty name",
			input:    "",
			expected: "float",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VarType(tt.input)
			testutils.True(t, fmt.Sprintf("VarType(%q) = %q, want %q", tt.input, got, tt.expected), got == tt.expected)
		})
	}
}

func TestVarType_ErrorsAndEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "float",
		},
		{
			name:     "Only percent sign",
			input:    "%",
			expected: "int",
		},
		{
			name:     "Only dollar sign",
			input:    "$",
			expected: "string",
		},
		{
			name:     "Invalid suffix",
			input:    "A#",
			expected: "float",
		},
		{
			name:     "Whitespace name",
			input:    "   ",
			expected: "float",
		},
		{
			name:     "Trailing whitespace after suffix",
			input:    "A% ",
			expected: "float",
		},
		{
			name:     "Suffix not at end",
			input:    "A%B",
			expected: "float",
		},
		{
			name:     "Multiple suffixes percent",
			input:    "A%%",
			expected: "int",
		},
		{
			name:     "Multiple suffixes dollar",
			input:    "A$$",
			expected: "string",
		},
		{
			name:     "Lowercase variable name",
			input:    "value%",
			expected: "int",
		},
		{
			name:     "Unicode variable name",
			input:    "π$",
			expected: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VarType(tt.input)
			testutils.True(t, fmt.Sprintf("VarType(%q) = %q, want %q", tt.input, got, tt.expected), got == tt.expected)
		})
	}
}

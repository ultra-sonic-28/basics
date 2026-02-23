package common

import (
	"basics/testutils"
	"fmt"
	"testing"
)

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "Zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "Positive integer",
			input:    42,
			expected: "42",
		},
		{
			name:     "Negative integer",
			input:    -7,
			expected: "-7",
		},
		{
			name:     "Float with decimals",
			input:    0.5,
			expected: "0.5",
		},
		{
			name:     "Negative float",
			input:    -0.75,
			expected: "-0.75",
		},
		{
			name:     "Float 3.14",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "Float 3e-5 (decimal)",
			input:    3e-5,
			expected: "0.00003",
		},
		{
			name:     "Float 3e-6 (decimal)",
			input:    3e-6,
			expected: "0.000003",
		},
		{
			name:     "Float 3e-7 (scientific)",
			input:    3e-7,
			expected: "3e-07",
		},
		{
			name:     "Float 3e-10 (scientific)",
			input:    3e-10,
			expected: "3e-10",
		},
		{
			name:     "Float scientific notation large",
			input:    10000000000.5,
			expected: "1.00000000005e+10",
		},
		{
			name:     "Negative float",
			input:    -12.75,
			expected: "-12.75",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatNumber(tt.input)
			testutils.True(t, fmt.Sprintf("formatNumber(%v) = %q, want %q", tt.input, got, tt.expected), got == tt.expected)
		})
	}
}

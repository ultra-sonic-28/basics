package interpreter

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
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "Float without trailing zeros",
			input:    2.5,
			expected: "2.5",
		},
		{
			name:     "Float scientific notation small",
			input:    0.000001,
			expected: "1e-06",
		},
		{
			name:     "Float scientific notation large",
			input:    1000000.5,
			expected: "1.0000005e+06",
		},
		{
			name:     "Negative float",
			input:    -12.75,
			expected: "-12.75",
		},
		{
			name:     "Float very close to integer but not exact",
			input:    1.0000000001,
			expected: "1.0000000001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatNumber(tt.input)
			testutils.True(t, fmt.Sprintf("formatNumber(%v) = %q, want %q", tt.input, got, tt.expected), got == tt.expected)
		})
	}
}

package common

import (
	"basics/testutils"
	"fmt"
	"math"
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
		{
			name:     "Positive infinity",
			input:    math.Inf(1),
			expected: "OVER FLOW ERROR",
		},
		{
			name:     "Negative infinity",
			input:    math.Inf(-1),
			expected: "OVER FLOW ERROR",
		},
		{
			name:     "NaN",
			input:    math.NaN(),
			expected: "OVER FLOW ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatNumber(tt.input)
			testutils.True(t, fmt.Sprintf("formatNumber(%v) = %q, want %q", tt.input, got, tt.expected), got == tt.expected)
		})
	}
}

func TestFormatNumber_Special(t *testing.T) {
	testutils.Equal(t, "Positive infinity", FormatNumber(math.Inf(1)), "OVER FLOW ERROR")
	testutils.Equal(t, "Negative infinity", FormatNumber(math.Inf(-1)), "OVER FLOW ERROR")
	testutils.Equal(t, "NaN", FormatNumber(math.NaN()), "OVER FLOW ERROR")
}

func TestParseFloatApplesoft(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedVal   float64
		expectedError string
	}{
		{
			name:        "Simple integer",
			input:       "123",
			expectedVal: 123,
		},
		{
			name:        "Negative integer",
			input:       "-456",
			expectedVal: -456,
		},
		{
			name:        "Positive signed integer",
			input:       "+789",
			expectedVal: 789,
		},
		{
			name:        "Simple float",
			input:       "3.14",
			expectedVal: 3.14,
		},
		{
			name:        "Float starting with dot",
			input:       ".5",
			expectedVal: 0.5,
		},
		{
			name:        "Float ending with dot",
			input:       "1.",
			expectedVal: 1.0,
		},
		{
			name:        "Scientific notation lowercase",
			input:       "1e2",
			expectedVal: 100,
		},
		{
			name:        "Scientific notation uppercase",
			input:       "1E2",
			expectedVal: 100,
		},
		{
			name:        "Scientific notation with dot",
			input:       "1.2e2",
			expectedVal: 120,
		},
		{
			name:        "Scientific notation negative exponent",
			input:       "1.2e-2",
			expectedVal: 0.012,
		},
		{
			name:        "Scientific notation positive signed exponent",
			input:       "1.2e+2",
			expectedVal: 120,
		},
		{
			name:        "Leading spaces",
			input:       "   123",
			expectedVal: 123,
		},
		{
			name:        "Leading zeros",
			input:       "000123",
			expectedVal: 123,
		},
		{
			name:        "Stops at invalid character (alpha)",
			input:       "123ABC",
			expectedVal: 123,
		},
		{
			name:        "Stops at invalid character (space)",
			input:       "123 456",
			expectedVal: 123,
		},
		{
			name:        "Starts with invalid character",
			input:       "ABC123",
			expectedVal: 0,
		},
		{
			name:        "Stops at second dot",
			input:       "1.2.3",
			expectedVal: 1.2,
		},
		{
			name:        "Stops at second exponent",
			input:       "1e2e3",
			expectedVal: 100,
		},
		{
			name:        "Just a plus",
			input:       "+",
			expectedVal: 0,
		},
		{
			name:        "Just a minus",
			input:       "-",
			expectedVal: 0,
		},
		{
			name:        "Just a dot",
			input:       ".",
			expectedVal: 0,
		},
		{
			name:        "Just an e",
			input:       "e",
			expectedVal: 0,
		},
		{
			name:        "Empty string",
			input:       "",
			expectedVal: 0,
		},
		{
			name:        "Large number",
			input:       "1e38",
			expectedVal: 1e38,
		},
		{
			name:          "Overflow",
			input:         "1e400",
			expectedVal:   0,
			expectedError: "OVER FLOW ERROR",
		},
		{
			name:        "Space after minus stops parsing",
			input:       "- 123",
			expectedVal: 0,
		},
		{
			name:        "Trailing e followed by nothing",
			input:       "1e",
			expectedVal: 0, // Current implementation returns 0 if ParseFloat fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFloatApplesoft(tt.input)
			if tt.expectedError != "" {
				testutils.True(t, "expected error", err != nil)
				testutils.Equal(t, "error message", err.Error(), tt.expectedError)
			} else {
				testutils.True(t, fmt.Sprintf("no error expected for %q, got %v", tt.input, err), err == nil)
				testutils.Equal(t, fmt.Sprintf("ParseFloatApplesoft(%q)", tt.input), got, tt.expectedVal)
			}
		})
	}
}

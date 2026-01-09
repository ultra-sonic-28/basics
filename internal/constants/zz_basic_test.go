package constants

import (
	"basics/testutils"
	"fmt"
	"testing"
)

func TestBasicConstants(t *testing.T) {
	tests := []struct {
		name     string
		constVal byte
		expected byte
	}{
		{"BASIC_TTY", BASIC_TTY, 0},
		{"BASIC_APPLE", BASIC_APPLE, 1},
		{"BASIC_C64", BASIC_C64, 2},
		{"BASIC_AMS", BASIC_AMS, 3},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("tests[%d] - %s wrong value. got=%d, want=%d", i, tt.name, tt.constVal, tt.expected)
		testutils.True(t, msg, tt.constVal == tt.expected)
	}
}

func TestBasicVersionMap(t *testing.T) {
	tests := []struct {
		constVal byte
		expected byte
	}{
		{BASIC_TTY, 10},
		{BASIC_APPLE, 10},
		{BASIC_C64, 10},
		{BASIC_AMS, 10},
	}

	for i, tt := range tests {
		got, ok := BasicVersion[tt.constVal]
		msg := fmt.Sprintf("tests[%d] - BasicVersion key %d missing", i, tt.constVal)
		testutils.True(t, msg, ok)
		if !ok {
			continue
		}
		msg = fmt.Sprintf("tests[%d] - BasicVersion[%d] wrong. got=%d, want=%d", i, tt.constVal, got, tt.expected)
		testutils.True(t, msg, got == tt.expected)
	}
}

func TestBasicNameMap(t *testing.T) {
	tests := []struct {
		constVal byte
		expected string
	}{
		{BASIC_TTY, "TTY"},
		{BASIC_APPLE, "APPLE"},
		{BASIC_C64, "C64"},
		{BASIC_AMS, "AMS6128"},
	}

	for i, tt := range tests {
		got, ok := BasicName[tt.constVal]
		msg := fmt.Sprintf("tests[%d] - BasicName key %d missing", i, tt.constVal)
		testutils.True(t, msg, ok)
		if !ok {
			continue
		}
		msg = fmt.Sprintf("tests[%d] - BasicName[%d] wrong. got=%q, want=%q", i, tt.constVal, got, tt.expected)
		testutils.True(t, msg, got == tt.expected)
	}
}

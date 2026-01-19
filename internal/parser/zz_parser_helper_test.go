package parser

import (
	"basics/testutils"
	"fmt"
	"testing"
)

func TestStmtName(t *testing.T) {
	tests := []struct {
		name     string
		stmt     Statement
		expected string
	}{
		{"HOME", &HomeStmt{}, "HOME"},
		{"PRINT", &PrintStmt{}, "PRINT"},
		{"LET", &LetStmt{}, "LET"},
		{"IF", &IfStmt{}, "IF"},
		{"IFMULTI", &IfJumpStmt{}, "IFMULTI"},
		{"GOTO", &GotoStmt{}, "GOTO"},
		{"GOSUB", &GosubStmt{}, "GOSUB"},
		{"RETURN", &ReturnStmt{}, "RETURN"},
		{"FOR", &ForStmt{}, "FOR"},
		{"NEXT", &NextStmt{}, "NEXT"},
		{"END", &EndStmt{}, "END"},
		{"HTAB", &HTabStmt{}, "HTAB"},
		{"VTAB", &VTabStmt{}, "VTAB"},
		{"UNKNOWN", nil, "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StmtName(tt.stmt)
			testutils.True(t, fmt.Sprintf("StmtName(%T) = %q, expected %q", tt.stmt, got, tt.expected), got == tt.expected)
		})
	}
}

func TestStmtArgs(t *testing.T) {
	tests := []struct {
		name     string
		stmt     Statement
		expected string
	}{
		{
			name:     "PRINT",
			stmt:     &PrintStmt{},
			expected: " ->",
		},
		{
			name: "LET",
			stmt: &LetStmt{
				Name: "A",
			},
			expected: " A ->",
		},
		{
			name:     "IFJUMP",
			stmt:     &IfJumpStmt{},
			expected: " ->",
		},
		{
			name: "FOR",
			stmt: &ForStmt{
				Var: "I",
			},
			expected: " I",
		},
		{
			name: "NEXT",
			stmt: &NextStmt{
				Var: "I",
			},
			expected: " I",
		},
		{
			name:     "NO_ARGS_HOME",
			stmt:     &HomeStmt{},
			expected: "",
		},
		{
			name:     "NO_ARGS_END",
			stmt:     &EndStmt{},
			expected: "",
		},
		{
			name:     "UNKNOWN",
			stmt:     nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StmtArgs(tt.stmt)
			testutils.True(t, fmt.Sprintf("StmtArgs(%T) = %q, expected %q", tt.stmt, got, tt.expected), got == tt.expected)
		})
	}
}

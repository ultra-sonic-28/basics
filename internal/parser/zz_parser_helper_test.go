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
		{"INPUT", &InputStmt{}, "INPUT"},
		{"GET", &GetStmt{}, "GET"},
		{"END", &EndStmt{}, "END"},
		{"HTAB", &HTabStmt{}, "HTAB"},
		{"VTAB", &VTabStmt{}, "VTAB"},
		{"NORMAL", &NormalStmt{}, "NORMAL"},
		{"INVERSE", &InverseStmt{}, "INVERSE"},
		{"FLASH", &FlashStmt{}, "FLASH"},
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
	var vars []*Identifier

	tests := []struct {
		name     string
		stmt     Statement
		expected string
	}{
		// Statements with args
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
			name: "GET",
			stmt: &GetStmt{
				Var: &Identifier{
					Name: "A$",
				},
			},
			expected: " A$ ->",
		},
		{
			name: "INPUT with prompt",
			stmt: &InputStmt{
				Prompt: &StringLiteral{
					Value: "Enter:",
				},
				Vars: append(vars, &Identifier{
					Name: "A$",
				},
				),
			},
			expected: " -> Enter: -> A$",
		},
		{
			name: "INPUT without prompt",
			stmt: &InputStmt{
				Prompt: nil,
				Vars: append(vars, &Identifier{
					Name: "A$",
				},
				),
			},
			expected: " -> A$",
		},

		// Statements without args
		{
			name:     "NO_ARGS_HOME",
			stmt:     &HomeStmt{},
			expected: "",
		},
		{
			name:     "NO_ARGS_RETURN",
			stmt:     &ReturnStmt{},
			expected: "",
		},
		{
			name:     "NO_ARGS_END",
			stmt:     &EndStmt{},
			expected: "",
		},
		{
			name:     "NO_ARGS_NORMAL",
			stmt:     &NormalStmt{},
			expected: "",
		},
		{
			name:     "NO_ARGS_INVERSE",
			stmt:     &InverseStmt{},
			expected: "",
		},
		{
			name:     "NO_ARGS_FLASH",
			stmt:     &FlashStmt{},
			expected: "",
		},
		{
			name:     "NO_ARGS_CLEAR",
			stmt:     &ClearStmt{},
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

func TestStmtExprValue(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected string
	}{
		{
			name: "NumberLiteral integer",
			expr: &NumberLiteral{
				Value: 10,
			},
			expected: "10",
		},
		{
			name: "NumberLiteral float",
			expr: &NumberLiteral{
				Value: 3.14,
			},
			expected: "3.14",
		},
		{
			name: "StringLiteral",
			expr: &StringLiteral{
				Value: "HELLO",
			},
			expected: "\"HELLO\"",
		},
		{
			name: "Identifier",
			expr: &Identifier{
				Name: "A$",
			},
			expected: "A$",
		},
		{
			name:     "Nil expression",
			expr:     nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StmtExprValue(tt.expr)
			testutils.True(
				t,
				fmt.Sprintf("StmtExprValue(%T) = %q, expected %q", tt.expr, got, tt.expected),
				got == tt.expected,
			)
		})
	}
}

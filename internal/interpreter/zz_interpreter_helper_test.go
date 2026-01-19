package interpreter

import (
	"basics/internal/parser"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestLogTrace_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		inst     Instruction
		pc       int
		nextPC   int
		sExpr    string
		expected string
	}{
		{
			name: "HOME",
			inst: Instruction{
				LineNum: 10,
				Stmt:    &parser.HomeStmt{},
			},
			pc:       0,
			nextPC:   1,
			sExpr:    "",
			expected: "Executing line: 10, pc: 0, nextPC: 1 - [HOME] ",
		},
		{
			name: "LET",
			inst: Instruction{
				LineNum: 20,
				Stmt: &parser.LetStmt{
					Name:  "A",
					Value: &parser.NumberLiteral{Value: 42},
				},
			},
			pc:       1,
			nextPC:   2,
			sExpr:    "42",
			expected: "Executing line: 20, pc: 1, nextPC: 2 - [LET] A -> 42",
		},
		{
			name: "PRINT",
			inst: Instruction{
				LineNum: 30,
				Stmt: &parser.PrintStmt{
					Exprs: []parser.Expression{
						&parser.StringLiteral{Value: "HELLO"},
					},
				},
			},
			pc:       2,
			nextPC:   3,
			sExpr:    "\"HELLO\"",
			expected: "Executing line: 30, pc: 2, nextPC: 3 - [PRINT] -> \"HELLO\"",
		},
		{
			name: "FOR",
			inst: Instruction{
				LineNum: 40,
				Stmt: &parser.ForStmt{
					Var: "I",
				},
			},
			pc:       3,
			nextPC:   4,
			sExpr:    "",
			expected: "Executing line: 40, pc: 3, nextPC: 4 - [FOR] I ",
		},
		{
			name: "NEXT",
			inst: Instruction{
				LineNum: 50,
				Stmt: &parser.NextStmt{
					Var: "I",
				},
			},
			pc:       4,
			nextPC:   5,
			sExpr:    "",
			expected: "Executing line: 50, pc: 4, nextPC: 5 - [NEXT] I ",
		},
		{
			name: "GOTO",
			inst: Instruction{
				LineNum: 60,
				Stmt:    &parser.GotoStmt{},
			},
			pc:       5,
			nextPC:   10,
			sExpr:    "100",
			expected: "Executing line: 60, pc: 5, nextPC: 10 - [GOTO] 100",
		},
		{
			name: "END",
			inst: Instruction{
				LineNum: 70,
				Stmt:    &parser.EndStmt{},
			},
			pc:       10,
			nextPC:   -1,
			sExpr:    "",
			expected: "Executing line: 70, pc: 10, nextPC: -1 - [END] ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := LogTrace(tt.inst, tt.pc, tt.nextPC, tt.sExpr)

			testutils.True(t, fmt.Sprintf("\nexpected: %q\ngot:      %q", tt.expected, got), got == tt.expected)
		})
	}
}

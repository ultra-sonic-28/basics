package parser

import (
	"basics/testutils"
	"fmt"
	"testing"
)

// ------------------------
// Tests pour dumpExpr
// ------------------------
func TestDumpExpr(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected string
	}{
		{
			name:     "NumberLiteral",
			expr:     &NumberLiteral{Value: 42, Line: 1, Column: 1, Token: "42"},
			expected: "Number 42\n",
		},
		{
			name:     "StringLiteral",
			expr:     &StringLiteral{Value: "hello", Line: 2, Column: 3, Token: "\"hello\""},
			expected: "String \"hello\"\n",
		},
		{
			name:     "Identifier",
			expr:     &Identifier{Name: "X", Line: 3, Column: 5, Token: "X"},
			expected: "Ident X\n",
		},
		{
			name: "PrefixExpr",
			expr: &PrefixExpr{
				Op:    "-",
				Right: &NumberLiteral{Value: 5, Line: 4, Column: 2, Token: "5"},
			},
			expected: "Prefix -\n  Number 5\n",
		},
		{
			name: "InfixExpr",
			expr: &InfixExpr{
				Left:  &NumberLiteral{Value: 2, Line: 5, Column: 1, Token: "2"},
				Op:    "+",
				Right: &NumberLiteral{Value: 3, Line: 5, Column: 3, Token: "3"},
			},
			expected: "Infix +\n  Number 2\n  Number 3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := testutils.CaptureStdout(t, func() {
				dumpExpr(tt.expr, "")
			})
			msg := fmt.Sprintf("tests[%s] - dumpExpr output mismatch.", tt.name)
			testutils.Equal(t, msg, output, tt.expected)
		})
	}
}

// ------------------------
// Tests pour dumpStatement
// ------------------------
func TestDumpStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     Statement
		expected string
	}{
		{
			name:     "PrintStmt",
			stmt:     &PrintStmt{Exprs: []Expression{&NumberLiteral{Value: 42, Line: 1, Column: 1, Token: "42"}}},
			expected: "PRINT\n  EXPR 0:\n    Number 42\n",
		},
		{
			name:     "LetStmt",
			stmt:     &LetStmt{Name: "X", Value: &NumberLiteral{Value: 10, Line: 2, Column: 1, Token: "10"}},
			expected: "LET X\n  Number 10\n",
		},
		{
			name: "ForStmt without Step",
			stmt: &ForStmt{
				Var:     "I",
				Start:   &NumberLiteral{Value: 1, Line: 3, Column: 1, Token: "1"},
				End:     &NumberLiteral{Value: 10, Line: 3, Column: 3, Token: "10"},
				LineNum: 100,
				Step:    nil,
			},
			expected: "FOR I (Line 100)\n  FROM:\n    Number 1\n  TO:\n    Number 10\n",
		},
		{
			name: "ForStmt with Step",
			stmt: &ForStmt{
				Var:     "J",
				Start:   &NumberLiteral{Value: 0, Line: 4, Column: 1, Token: "0"},
				End:     &NumberLiteral{Value: 5, Line: 4, Column: 3, Token: "5"},
				LineNum: 200,
				Step:    &NumberLiteral{Value: 2, Line: 4, Column: 2, Token: "2"},
			},
			expected: "FOR J (Line 200)\n  STEP:\n    Number 2\n  FROM:\n    Number 0\n  TO:\n    Number 5\n",
		},
		{
			name:     "NextStmt with ForLineNum",
			stmt:     &NextStmt{Var: "I", ForLineNum: 100},
			expected: "NEXT I (FOR Line 100)\n",
		},
		{
			name:     "NextStmt without ForLineNum",
			stmt:     &NextStmt{Var: "J"},
			expected: "NEXT J\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := testutils.CaptureStdout(t, func() {
				dumpStatement(tt.stmt, "")
			})
			msg := fmt.Sprintf("tests[%s] - dumpStatement output mismatch.", tt.name)
			testutils.Equal(t, msg, output, tt.expected)
		})
	}
}

// ------------------------
// Tests pour DumpProgram
// ------------------------
func TestDumpProgram(t *testing.T) {
	p := &Program{
		Lines: []*Line{
			{
				Number: 10,
				Stmts: []Statement{
					&LetStmt{Name: "X", Value: &NumberLiteral{Value: 5, Line: 1, Column: 1, Token: "5"}},
					&PrintStmt{Exprs: []Expression{&Identifier{Name: "X", Line: 1, Column: 2, Token: "X"}}},
				},
			},
			{
				Number: 20,
				Stmts: []Statement{
					&ForStmt{
						Var:     "I",
						Start:   &NumberLiteral{Value: 1, Line: 2, Column: 1, Token: "1"},
						End:     &NumberLiteral{Value: 3, Line: 2, Column: 3, Token: "3"},
						LineNum: 20,
					},
					&NextStmt{Var: "I", ForLineNum: 20},
				},
			},
		},
	}

	expected := "Line 10\n  LET X\n    Number 5\n  PRINT\n    EXPR 0:\n      Ident X\nLine 20\n  FOR I (Line 20)\n    FROM:\n      Number 1\n    TO:\n      Number 3\n  NEXT I (FOR Line 20)\n"

	output := testutils.CaptureStdout(t, func() {
		DumpProgram(p)
	})

	testutils.Equal(t, "DumpProgram output mismatch", output, expected)
}

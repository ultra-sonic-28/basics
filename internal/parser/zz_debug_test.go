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
			name:     "EndStmt",
			stmt:     &EndStmt{},
			expected: "END\n",
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
		{
			name: "GotoStmt",
			stmt: &GotoStmt{
				Expr: &NumberLiteral{Value: 40, Line: 3, Column: 1, Token: "40"},
			},
			expected: "GOTO\n  Number 40\n",
		},
		{
			name: "GotoStmt with expression",
			stmt: &GotoStmt{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "JUMP", Token: "JUMP"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2, Token: "2"},
				},
			},
			expected: "" +
				"GOTO\n" +
				"  Infix *\n" +
				"    Ident JUMP\n" +
				"    Number 2\n",
		},
		{
			name: "HTabStmt",
			stmt: &HTabStmt{
				Expr: &NumberLiteral{Value: 10, Line: 1, Column: 1, Token: "10"},
			},
			expected: "HTAB\n  Number 10\n",
		},
		{
			name: "HTabStmt with expression",
			stmt: &HTabStmt{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A", Token: "A"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2, Token: "2"},
				},
			},
			expected: "" +
				"HTAB\n" +
				"  Infix *\n" +
				"    Ident A\n" +
				"    Number 2\n",
		},
		{
			name: "HTabStmt with parenthesized expression",
			stmt: &HTabStmt{
				Expr: &InfixExpr{
					Left: &InfixExpr{
						Left:  &Identifier{Name: "A", Token: "A"},
						Op:    "+",
						Right: &Identifier{Name: "B", Token: "B"},
					},
					Op: "*",
					Right: &NumberLiteral{
						Value: 2,
						Token: "2",
					},
				},
			},
			expected: "" +
				"HTAB\n" +
				"  Infix *\n" +
				"    Infix +\n" +
				"      Ident A\n" +
				"      Ident B\n" +
				"    Number 2\n",
		},
		{
			name: "VTabStmt",
			stmt: &VTabStmt{
				Expr: &Identifier{Name: "A", Line: 2, Column: 3, Token: "A"},
			},
			expected: "VTAB\n  Ident A\n",
		},
		{
			name: "VTabStmt with expression",
			stmt: &VTabStmt{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "V", Token: "V"},
					Op:    "+",
					Right: &NumberLiteral{Value: 5, Token: "5"},
				},
			},
			expected: "" +
				"VTAB\n" +
				"  Infix +\n" +
				"    Ident V\n" +
				"    Number 5\n",
		},
		{
			name: "VTabStmt with parenthesized expression",
			stmt: &VTabStmt{
				Expr: &InfixExpr{
					Left: &InfixExpr{
						Left:  &Identifier{Name: "A", Token: "A"},
						Op:    "+",
						Right: &Identifier{Name: "B", Token: "B"},
					},
					Op: "*",
					Right: &NumberLiteral{
						Value: 2,
						Token: "2",
					},
				},
			},
			expected: "" +
				"VTAB\n" +
				"  Infix *\n" +
				"    Infix +\n" +
				"      Ident A\n" +
				"      Ident B\n" +
				"    Number 2\n",
		},
		{
			name: "IfStmt without ELSE",
			stmt: &IfStmt{
				Cond: &InfixExpr{
					Left:  &Identifier{Name: "A", Token: "A"},
					Op:    "<",
					Right: &NumberLiteral{Value: 10, Token: "10"},
				},
				Then: []Statement{
					&GotoStmt{
						Expr: &NumberLiteral{Value: 20, Token: "20"},
					},
				},
			},
			expected: "" +
				"IF\n" +
				"  Infix <\n" +
				"    Ident A\n" +
				"    Number 10\n" +
				"THEN\n" +
				"  GOTO\n" +
				"    Number 20\n",
		},
		{
			name: "IfStmt with ELSE",
			stmt: &IfStmt{
				Cond: &Identifier{Name: "X", Token: "X"},
				Then: []Statement{
					&PrintStmt{
						Exprs: []Expression{
							&StringLiteral{Value: "YES", Token: "\"YES\""},
						},
					},
				},
				Else: []Statement{
					&PrintStmt{
						Exprs: []Expression{
							&StringLiteral{Value: "NO", Token: "\"NO\""},
						},
					},
				},
			},
			expected: "" +
				"IF\n" +
				"  Ident X\n" +
				"THEN\n" +
				"  PRINT\n" +
				"    EXPR 0:\n" +
				"      String \"YES\"\n" +
				"ELSE\n" +
				"  PRINT\n" +
				"    EXPR 0:\n" +
				"      String \"NO\"\n",
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

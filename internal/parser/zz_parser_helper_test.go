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
		{"CLEAR", &ClearStmt{}, "CLEAR"},
		{"DIM", &DimStmt{}, "DIM"},
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
			expected: "  ->",
		},
		{
			name: "PRINT with identifier",
			stmt: &PrintStmt{
				Exprs: []Expression{
					&Identifier{Name: "A", Token: "A"},
				},
			},
			expected: " A ->",
		},
		{
			name: "PRINT with number literal",
			stmt: &PrintStmt{
				Exprs: []Expression{
					&NumberLiteral{Value: 42, Token: "42"},
				},
			},
			expected: " 42 ->",
		},
		{
			name: "PRINT with string literal",
			stmt: &PrintStmt{
				Exprs: []Expression{
					&StringLiteral{Value: "HELLO", Token: "\"HELLO\""},
				},
			},
			expected: " \"HELLO\" ->",
		},
		{
			name: "PRINT with infix expression",
			stmt: &PrintStmt{
				Exprs: []Expression{
					&InfixExpr{
						Left:  &Identifier{Name: "A", Token: "A"},
						Op:    "+",
						Right: &NumberLiteral{Value: 1, Token: "1"},
					},
				},
			},
			expected: " A + 1 ->",
		},
		{
			name: "PRINT with multiple expressions",
			stmt: &PrintStmt{
				Exprs: []Expression{
					&Identifier{Name: "A", Token: "A"},
					&Identifier{Name: "B", Token: "B"},
				},
			},
			expected: " AB ->",
		},
		{
			name: "PRINT with mixed expressions",
			stmt: &PrintStmt{
				Exprs: []Expression{
					&Identifier{Name: "A", Token: "A"},
					&InfixExpr{
						Left:  &Identifier{Name: "B", Token: "B"},
						Op:    "*",
						Right: &NumberLiteral{Value: 2, Token: "2"},
					},
					&StringLiteral{Value: "OK", Token: "\"OK\""},
				},
			},
			expected: " AB * 2\"OK\" ->",
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
		{
			name: "LET with single index identifier",
			stmt: &LetStmt{
				Name: "A",
				Indices: []Expression{
					&Identifier{Name: "I"},
				},
			},
			expected: " A(I) ->",
		},
		{
			name: "LET with single numeric index",
			stmt: &LetStmt{
				Name: "A",
				Indices: []Expression{
					&NumberLiteral{Value: 3, Token: "3"},
				},
			},
			expected: " A(3) ->",
		},
		{
			name: "LET with infix index",
			stmt: &LetStmt{
				Name: "A",
				Indices: []Expression{
					&InfixExpr{
						Left:  &Identifier{Name: "I", Token: "I"},
						Op:    "+",
						Right: &NumberLiteral{Value: 1, Token: "1"},
					},
				},
			},
			expected: " A(I + 1) ->",
		},
		{
			name: "LET with multiple indices",
			stmt: &LetStmt{
				Name: "A",
				Indices: []Expression{
					&Identifier{Name: "I", Token: "I"},
					&Identifier{Name: "J", Token: "J"},
				},
			},
			expected: " A(I,J) ->",
		},
		{
			name: "LET with mixed indices expressions",
			stmt: &LetStmt{
				Name: "A",
				Indices: []Expression{
					&Identifier{Name: "I", Token: "I"},
					&InfixExpr{
						Left:  &Identifier{Name: "J", Token: "J"},
						Op:    "+",
						Right: &NumberLiteral{Value: 2, Token: "2"},
					},
				},
			},
			expected: " A(I,J + 2) ->",
		},
		{
			name: "DIM single array single numeric dimension",
			stmt: &DimStmt{
				Arrays: []DimDecl{
					{
						Name: "A",
						Dimensions: []Expression{
							&NumberLiteral{Value: 10, Token: "10"},
						},
					},
				},
			},
			expected: " -> A(10)",
		},
		{
			name: "DIM single array identifier dimension",
			stmt: &DimStmt{
				Arrays: []DimDecl{
					{
						Name: "A%",
						Dimensions: []Expression{
							&Identifier{Name: "SIZE", Token: "SIZE"},
						},
					},
				},
			},
			expected: " -> A%(SIZE)",
		},
		{
			name: "DIM single array multiple numeric dimensions",
			stmt: &DimStmt{
				Arrays: []DimDecl{
					{
						Name: "B",
						Dimensions: []Expression{
							&NumberLiteral{Value: 10, Token: "10"},
							&NumberLiteral{Value: 20, Token: "20"},
						},
					},
				},
			},
			expected: " -> B(10, 20)",
		},
		{
			name: "DIM with infix expression dimension",
			stmt: &DimStmt{
				Arrays: []DimDecl{
					{
						Name: "C",
						Dimensions: []Expression{
							&InfixExpr{
								Left:  &Identifier{Name: "N", Token: "N"},
								Op:    "*",
								Right: &NumberLiteral{Value: 2, Token: "2"},
							},
						},
					},
				},
			},
			expected: " -> C(N * 2)",
		},
		{
			name: "DIM multiple arrays",
			stmt: &DimStmt{
				Arrays: []DimDecl{
					{
						Name: "A",
						Dimensions: []Expression{
							&NumberLiteral{Value: 10, Token: "10"},
						},
					},
					{
						Name: "B$",
						Dimensions: []Expression{
							&Identifier{Name: "SIZE", Token: "SIZE"},
							&NumberLiteral{Value: 5, Token: "5"},
						},
					},
				},
			},
			expected: " -> A(10)B$(SIZE, 5)",
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
			testutils.True(t, fmt.Sprintf("StmtArgs(%T): got %q, expected %q", tt.stmt, got, tt.expected), got == tt.expected)
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
			name: "IndexExpr with single numeric index",
			expr: &IndexExpr{
				Name: "A",
				Indices: []Expression{
					&NumberLiteral{Value: 1},
				},
			},
			expected: "A(1) ",
		},
		{
			name: "IndexExpr with identifier index",
			expr: &IndexExpr{
				Name: "TAB",
				Indices: []Expression{
					&Identifier{Name: "I"},
				},
			},
			expected: "TAB(I) ",
		},
		{
			name: "IndexExpr with infix index",
			expr: &IndexExpr{
				Name: "ARR",
				Indices: []Expression{
					&InfixExpr{
						Left:  &Identifier{Name: "I"},
						Op:    "+",
						Right: &NumberLiteral{Value: 1},
					},
				},
			},
			expected: "ARR(I + 1) ",
		},
		{
			name: "IndexExpr with multiple indices",
			expr: &IndexExpr{
				Name: "M",
				Indices: []Expression{
					&Identifier{Name: "I"},
					&Identifier{Name: "J"},
				},
			},
			expected: "M(I,J) ",
		},
		{
			name: "IndexExpr with mixed indices",
			expr: &IndexExpr{
				Name: "GRID",
				Indices: []Expression{
					&Identifier{Name: "X"},
					&InfixExpr{
						Left:  &Identifier{Name: "Y"},
						Op:    "*",
						Right: &NumberLiteral{Value: 2},
					},
					&NumberLiteral{Value: 3},
				},
			},
			expected: "GRID(X,Y * 2,3) ",
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
				fmt.Sprintf("StmtExprValue(%T): got %q, expected %q", tt.expr, got, tt.expected),
				got == tt.expected,
			)
		})
	}
}

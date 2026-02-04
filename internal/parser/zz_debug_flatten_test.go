package parser

import (
	"basics/testutils"
	"testing"
)

func TestFlattenIndices(t *testing.T) {
	tests := []struct {
		name     string
		indices  []Expression
		expected string
	}{
		{
			name: "single identifier",
			indices: []Expression{
				&Identifier{Name: "I"},
			},
			expected: "I",
		},
		{
			name: "single number literal",
			indices: []Expression{
				&NumberLiteral{Value: 3},
			},
			expected: "3",
		},
		{
			name: "multiple identifiers",
			indices: []Expression{
				&Identifier{Name: "I"},
				&Identifier{Name: "J"},
			},
			expected: "I,J",
		},
		{
			name: "identifier and number",
			indices: []Expression{
				&Identifier{Name: "I"},
				&NumberLiteral{Value: 10},
			},
			expected: "I,10",
		},
		{
			name: "single infix expression",
			indices: []Expression{
				&InfixExpr{
					Left:  &Identifier{Name: "I"},
					Op:    "+",
					Right: &NumberLiteral{Value: 1},
				},
			},
			expected: "I + 1",
		},
		{
			name: "mixed identifier and infix",
			indices: []Expression{
				&Identifier{Name: "X"},
				&InfixExpr{
					Left:  &Identifier{Name: "I"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2},
				},
			},
			expected: "X,I * 2",
		},
		{
			name: "multiple infix expressions",
			indices: []Expression{
				&InfixExpr{
					Left:  &Identifier{Name: "I"},
					Op:    "+",
					Right: &NumberLiteral{Value: 1},
				},
				&InfixExpr{
					Left:  &Identifier{Name: "J"},
					Op:    "-",
					Right: &NumberLiteral{Value: 2},
				},
			},
			expected: "I + 1,J - 2",
		},
		{
			name:     "empty indices",
			indices:  []Expression{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlattenIndices(tt.indices)
			testutils.Equal(t, tt.name, result, tt.expected)
		})
	}
}

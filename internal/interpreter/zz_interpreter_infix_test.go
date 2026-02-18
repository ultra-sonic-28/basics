package interpreter

import (
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_InfixExpr(t *testing.T) {

	tests := []struct {
		name        string
		expr        *parser.InfixExpr
		expectError bool
		expectType  runtime.ValueType
		expectInt   int
		expectNum   float64
		expectStr   string
	}{
		// =========================
		// STRING
		// =========================
		{
			name: "String concatenation",
			expr: &parser.InfixExpr{
				Left:  &parser.StringLiteral{Value: "APPLE"},
				Op:    "+",
				Right: &parser.StringLiteral{Value: "SOFT"},
			},
			expectType: runtime.STRING,
			expectStr:  "APPLESOFT",
		},
		{
			name: "String equality true",
			expr: &parser.InfixExpr{
				Left:  &parser.StringLiteral{Value: "A"},
				Op:    "=",
				Right: &parser.StringLiteral{Value: "A"},
			},
			expectType: runtime.NUMBER,
			expectNum:  1,
		},
		{
			name: "String equality false",
			expr: &parser.InfixExpr{
				Left:  &parser.StringLiteral{Value: "A"},
				Op:    "=",
				Right: &parser.StringLiteral{Value: "B"},
			},
			expectType: runtime.NUMBER,
			expectNum:  0,
		},
		{
			name: "String invalid operator",
			expr: &parser.InfixExpr{
				Left:  &parser.StringLiteral{Value: "A"},
				Op:    "-",
				Right: &parser.StringLiteral{Value: "B"},
			},
			expectError: true,
		},

		// =========================
		// INTEGER
		// =========================
		{
			name: "Integer addition",
			expr: &parser.InfixExpr{
				Left:  &parser.NumberLiteral{Value: 5},
				Op:    "+",
				Right: &parser.NumberLiteral{Value: 3},
			},
			expectType: runtime.NUMBER,
			expectNum:  8,
		},
		{
			name: "Integer division produces float",
			expr: &parser.InfixExpr{
				Left:  &parser.NumberLiteral{Value: 5},
				Op:    "/",
				Right: &parser.NumberLiteral{Value: 2},
			},
			expectType: runtime.NUMBER,
			expectNum:  2.5,
		},
		{
			name: "Integer division by zero",
			expr: &parser.InfixExpr{
				Left:  &parser.NumberLiteral{Value: 5},
				Op:    "/",
				Right: &parser.NumberLiteral{Value: 0},
			},
			expectError: true,
		},
		{
			name: "Integer comparison true",
			expr: &parser.InfixExpr{
				Left:  &parser.NumberLiteral{Value: 3},
				Op:    "<",
				Right: &parser.NumberLiteral{Value: 5},
			},
			expectType: runtime.NUMBER,
			expectNum:  1,
		},

		// =========================
		// FLOAT
		// =========================
		{
			name: "Float multiplication",
			expr: &parser.InfixExpr{
				Left:  &parser.NumberLiteral{Value: 2.5},
				Op:    "*",
				Right: &parser.NumberLiteral{Value: 4},
			},
			expectType: runtime.NUMBER,
			expectNum:  10,
		},
		{
			name: "Float power",
			expr: &parser.InfixExpr{
				Left:  &parser.NumberLiteral{Value: 2},
				Op:    "^",
				Right: &parser.NumberLiteral{Value: 3},
			},
			expectType: runtime.NUMBER, // dans ton implémentation int^int
			expectNum:  8,
		},

		// =========================
		// MIXED
		// =========================
		{
			name: "Mixed int + float",
			expr: &parser.InfixExpr{
				Left:  &parser.NumberLiteral{Value: 3},
				Op:    "+",
				Right: &parser.NumberLiteral{Value: 2.5},
			},
			expectType: runtime.NUMBER,
			expectNum:  5.5,
		},
		{
			name: "Mixed comparison",
			expr: &parser.InfixExpr{
				Left:  &parser.NumberLiteral{Value: 3},
				Op:    ">",
				Right: &parser.NumberLiteral{Value: 2.5},
			},
			expectType: runtime.NUMBER,
			expectNum:  1,
		},

		// =========================
		// UNKNOWN OPERATOR
		// =========================
		{
			name: "Unknown operator",
			expr: &parser.InfixExpr{
				Left:  &parser.NumberLiteral{Value: 1},
				Op:    "??",
				Right: &parser.NumberLiteral{Value: 2},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rt := runtime.New(nil)

			val, _, err := EvalExpr(tt.expr, rt)

			if tt.expectError {
				testutils.True(t, "expected error", err != nil)
				return
			}

			testutils.True(t, "no error expected", err == nil)
			testutils.Equal(t, "type", val.Type, tt.expectType)

			switch tt.expectType {
			case runtime.INTEGER:
				testutils.Equal(t, "int value", val.Int, tt.expectInt)

			case runtime.NUMBER:
				testutils.Equal(t,
					"float value",
					val.Num,
					tt.expectNum,
				)

			case runtime.STRING:
				testutils.Equal(t, "string value", val.Str, tt.expectStr)
			}
		})
	}
}

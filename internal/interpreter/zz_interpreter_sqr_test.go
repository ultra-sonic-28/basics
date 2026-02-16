package interpreter

import (
	"math"
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_SqrExpr(t *testing.T) {

	tests := []struct {
		name        string
		expr        *parser.SqrExpr
		setupEnv    func(rt *runtime.Runtime)
		expectError bool
		expectNum   float64
	}{
		{
			name: "Positive integer",
			expr: &parser.SqrExpr{
				Expr: &parser.NumberLiteral{Value: 16},
			},
			expectNum: 4,
		},
		{
			name: "Zero",
			expr: &parser.SqrExpr{
				Expr: &parser.NumberLiteral{Value: 0},
			},
			expectNum: 0,
		},
		{
			name: "Positive float",
			expr: &parser.SqrExpr{
				Expr: &parser.NumberLiteral{Value: 5.5},
			},
			expectNum: math.Sqrt(5.5),
		},
		{
			name: "Integer variable",
			expr: &parser.SqrExpr{
				Expr: &parser.Identifier{Name: "I%"},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("I%", runtime.Value{
					Type: runtime.INTEGER,
					Int:  25,
				})
			},
			expectNum: 5,
		},
		{
			name: "Float variable",
			expr: &parser.SqrExpr{
				Expr: &parser.Identifier{Name: "A"},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  2.25,
				})
			},
			expectNum: 1.5,
		},
		{
			name: "Nested expression A * 4",
			expr: &parser.SqrExpr{
				Expr: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "A"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 4},
				},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  9,
				})
			},
			expectNum: 6, // sqrt(9 * 4) = sqrt(36)
		},
		{
			name: "Negative integer",
			expr: &parser.SqrExpr{
				Expr:  &parser.NumberLiteral{Value: -4},
				Line:  10,
				Token: "SQR",
			},
			expectError: true,
		},
		{
			name: "Negative float",
			expr: &parser.SqrExpr{
				Expr:  &parser.NumberLiteral{Value: -5.5},
				Line:  10,
				Token: "SQR",
			},
			expectError: true,
		},
		{
			name: "Type mismatch string",
			expr: &parser.SqrExpr{
				Expr:  &parser.StringLiteral{Value: "HELLO"},
				Line:  10,
				Token: "SQR",
			},
			expectError: true,
		},
		{
			name: "Error propagation from sub-expression",
			expr: &parser.SqrExpr{
				Expr: &parser.InfixExpr{
					Left:  &parser.StringLiteral{Value: "A"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2},
				},
				Line:  10,
				Token: "SQR",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rt := runtime.New(nil)

			if tt.setupEnv != nil {
				tt.setupEnv(rt)
			}

			val, _, err := EvalExpr(tt.expr, rt)

			if tt.expectError {
				testutils.True(t, "expected error", err != nil)
				return
			}

			testutils.True(t, "no error expected", err == nil)
			testutils.Equal(t, "result type", val.Type, runtime.NUMBER)

			testutils.True(t,
				"numeric result correct",
				math.Abs(val.Num-tt.expectNum) < 1e-9,
			)
		})
	}
}

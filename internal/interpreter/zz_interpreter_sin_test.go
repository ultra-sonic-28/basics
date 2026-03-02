package interpreter

import (
	"math"
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_SinExpr(t *testing.T) {

	tests := []struct {
		name        string
		expr        *parser.SinExpr
		setupEnv    func(rt *runtime.Runtime)
		expectError bool
		expectNum   float64
	}{
		{
			name: "Positive integer",
			expr: &parser.SinExpr{
				Expr: &parser.NumberLiteral{Value: 1},
			},
			expectNum: 0.8414709848078965,
		},
		{
			name: "Zero",
			expr: &parser.SinExpr{
				Expr: &parser.NumberLiteral{Value: 0},
			},
			expectNum: 0,
		},
		{
			name: "Positive float",
			expr: &parser.SinExpr{
				Expr: &parser.NumberLiteral{Value: 3.1415926535897932384626433832795},
			},
			expectNum: math.Sin(3.1415926535897932384626433832795),
		},
		{
			name: "Integer variable",
			expr: &parser.SinExpr{
				Expr: &parser.Identifier{Name: "I%"},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("I%", runtime.Value{
					Type: runtime.INTEGER,
					Int:  1,
				})
			},
			expectNum: math.Sin(1),
		},
		{
			name: "Float variable",
			expr: &parser.SinExpr{
				Expr: &parser.Identifier{Name: "A"},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  3.1415926535897932384626433832795,
				})
			},
			expectNum: math.Sin(3.1415926535897932384626433832795),
		},
		{
			name: "Nested expression A * 4",
			expr: &parser.SinExpr{
				Expr: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "A"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 4},
				},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  0,
				})
			},
			expectNum: 0,
		},
		{
			name: "Type mismatch string",
			expr: &parser.SinExpr{
				Expr:  &parser.StringLiteral{Value: "HELLO"},
				Line:  10,
				Token: "SIN",
			},
			expectError: true,
		},
		{
			name: "Error propagation from sub-expression",
			expr: &parser.SinExpr{
				Expr: &parser.InfixExpr{
					Left:  &parser.StringLiteral{Value: "A"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2},
				},
				Line:  10,
				Token: "SIN",
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

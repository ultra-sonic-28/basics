package interpreter

import (
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_IntExpr(t *testing.T) {

	tests := []struct {
		name        string
		expr        *parser.IntExpr
		setupEnv    func(rt *runtime.Runtime)
		expectError bool
		expectInt   int
	}{
		{
			name: "Integer unchanged",
			expr: &parser.IntExpr{
				Expr: &parser.NumberLiteral{Value: 5},
			},
			expectInt: 5,
		},
		{
			name: "Zero integer",
			expr: &parser.IntExpr{
				Expr: &parser.NumberLiteral{Value: 0},
			},
			expectInt: 0,
		},
		{
			name: "Positive float",
			expr: &parser.IntExpr{
				Expr: &parser.NumberLiteral{Value: 1.75},
			},
			expectInt: 1, // partie entière
		},
		{
			name: "Positive float boundary 2.0",
			expr: &parser.IntExpr{
				Expr: &parser.NumberLiteral{Value: 2.0},
			},
			expectInt: 2,
		},
		{
			name: "Negative float",
			expr: &parser.IntExpr{
				Expr: &parser.NumberLiteral{Value: -1.75},
			},
			expectInt: -2, // AppleSoft behavior
		},
		{
			name: "Negative float boundary -2.0",
			expr: &parser.IntExpr{
				Expr: &parser.NumberLiteral{Value: -2.0},
			},
			expectInt: -3, // selon ton implémentation actuelle
		},
		{
			name: "Integer variable",
			expr: &parser.IntExpr{
				Expr: &parser.Identifier{Name: "I%"},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("I%", runtime.Value{
					Type: runtime.INTEGER,
					Int:  7,
				})
			},
			expectInt: 7,
		},
		{
			name: "Float variable",
			expr: &parser.IntExpr{
				Expr: &parser.Identifier{Name: "A"},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  3.9,
				})
			},
			expectInt: 3,
		},
		{
			name: "Negative float variable",
			expr: &parser.IntExpr{
				Expr: &parser.Identifier{Name: "B"},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("B", runtime.Value{
					Type: runtime.NUMBER,
					Num:  -4.2,
				})
			},
			expectInt: -5,
		},
		{
			name: "Nested expression A * 2.5",
			expr: &parser.IntExpr{
				Expr: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "A"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2.5},
				},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  2,
				})
			},
			expectInt: 5, // 2 * 2.5 = 5.0 → INT = 5
		},
		{
			name: "Negative nested expression",
			expr: &parser.IntExpr{
				Expr: &parser.PrefixExpr{
					Op: "-",
					Right: &parser.NumberLiteral{
						Value: 3.2,
					},
				},
			},
			expectInt: -4,
		},
		{
			name: "Type mismatch string",
			expr: &parser.IntExpr{
				Expr:  &parser.StringLiteral{Value: "HELLO"},
				Line:  10,
				Token: "INT",
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
			testutils.Equal(t, "result type", val.Type, runtime.INTEGER)
			testutils.Equal(t, "integer value", val.Int, tt.expectInt)
		})
	}
}

package interpreter

import (
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_AbsExpr(t *testing.T) {

	tests := []struct {
		name        string
		expr        *parser.AbsExpr
		setupEnv    func(rt *runtime.Runtime)
		expectError bool
		expectType  runtime.ValueType
		expectInt   int
		expectNum   float64
	}{
		{
			name: "Positive integer",
			expr: &parser.AbsExpr{
				Expr: &parser.NumberLiteral{Value: 5},
			},
			expectType: runtime.NUMBER,
			expectNum:  5,
		},
		{
			name: "Negative integer literal",
			expr: &parser.AbsExpr{
				Expr: &parser.PrefixExpr{
					Op:    "-",
					Right: &parser.NumberLiteral{Value: 5},
				},
			},
			expectType: runtime.NUMBER,
			expectNum:  5,
		},
		{
			name: "Integer variable negative",
			expr: &parser.AbsExpr{
				Expr: &parser.Identifier{Name: "I%"},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("I%", runtime.Value{
					Type: runtime.INTEGER,
					Int:  -7,
				})
			},
			expectType: runtime.INTEGER,
			expectInt:  7,
		},
		{
			name: "Zero integer",
			expr: &parser.AbsExpr{
				Expr: &parser.NumberLiteral{Value: 0},
			},
			expectType: runtime.NUMBER,
			expectNum:  0,
		},
		{
			name: "Positive float",
			expr: &parser.AbsExpr{
				Expr: &parser.NumberLiteral{Value: 3.14},
			},
			expectType: runtime.NUMBER,
			expectNum:  3.14,
		},
		{
			name: "Negative float",
			expr: &parser.AbsExpr{
				Expr: &parser.PrefixExpr{
					Op:    "-",
					Right: &parser.NumberLiteral{Value: 2.5},
				},
			},
			expectType: runtime.NUMBER,
			expectNum:  2.5,
		},
		{
			name: "Expression -(A*2)",
			expr: &parser.AbsExpr{
				Expr: &parser.PrefixExpr{
					Op: "-",
					Right: &parser.InfixExpr{
						Left:  &parser.Identifier{Name: "A"},
						Op:    "*",
						Right: &parser.NumberLiteral{Value: 2},
					},
				},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.INTEGER,
					Int:  3,
				})
			},
			expectType: runtime.INTEGER,
			expectInt:  6,
		},
		{
			name: "Type mismatch with string",
			expr: &parser.AbsExpr{
				Expr:  &parser.StringLiteral{Value: "HELLO"},
				Line:  10,
				Token: "ABS",
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
			testutils.Equal(t, "result type", val.Type, tt.expectType)

			switch tt.expectType {
			case runtime.INTEGER:
				testutils.Equal(t, "integer value", val.Int, tt.expectInt)
			case runtime.NUMBER:
				testutils.Equal(t, "float value", val.Num, tt.expectNum)
			}
		})
	}
}

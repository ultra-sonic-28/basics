package interpreter

import (
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_SgnExpr(t *testing.T) {

	tests := []struct {
		name        string
		expr        *parser.SgnExpr
		setupEnv    func(rt *runtime.Runtime)
		expectError bool
		expectInt   int
	}{
		{
			name: "Positive integer",
			expr: &parser.SgnExpr{
				Expr: &parser.NumberLiteral{Value: 5},
			},
			expectInt: 1,
		},
		{
			name: "Negative integer",
			expr: &parser.SgnExpr{
				Expr: &parser.PrefixExpr{
					Op:    "-",
					Right: &parser.NumberLiteral{Value: 5},
				},
			},
			expectInt: -1,
		},
		{
			name: "Zero integer",
			expr: &parser.SgnExpr{
				Expr: &parser.NumberLiteral{Value: 0},
			},
			expectInt: 0,
		},
		{
			name: "Positive float",
			expr: &parser.SgnExpr{
				Expr: &parser.NumberLiteral{Value: 3.14},
			},
			expectInt: 1,
		},
		{
			name: "Negative float",
			expr: &parser.SgnExpr{
				Expr: &parser.PrefixExpr{
					Op:    "-",
					Right: &parser.NumberLiteral{Value: 2.5},
				},
			},
			expectInt: -1,
		},
		{
			name: "Zero float",
			expr: &parser.SgnExpr{
				Expr: &parser.NumberLiteral{Value: 0.0},
			},
			expectInt: 0,
		},
		{
			name: "Integer variable negative",
			expr: &parser.SgnExpr{
				Expr: &parser.Identifier{Name: "I%"},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("I%", runtime.Value{
					Type: runtime.INTEGER,
					Int:  -8,
				})
			},
			expectInt: -1,
		},
		{
			name: "Float variable positive",
			expr: &parser.SgnExpr{
				Expr: &parser.Identifier{Name: "A"},
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  7.2,
				})
			},
			expectInt: 1,
		},
		{
			name: "Nested expression -(A*2)",
			expr: &parser.SgnExpr{
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
					Type: runtime.NUMBER,
					Num:  3,
				})
			},
			expectInt: -1,
		},
		{
			name: "Type mismatch with string",
			expr: &parser.SgnExpr{
				Expr:  &parser.StringLiteral{Value: "HELLO"},
				Line:  10,
				Token: "SGN",
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

			// SGN retourne TOUJOURS un INTEGER
			testutils.Equal(t, "result type", val.Type, runtime.INTEGER)
			testutils.Equal(t, "integer value", val.Int, tt.expectInt)
		})
	}
}

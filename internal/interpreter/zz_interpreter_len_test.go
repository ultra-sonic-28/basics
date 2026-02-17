package interpreter

import (
	"fmt"
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_LenExpr(t *testing.T) {

	tests := []struct {
		name        string
		expr        *parser.LenExpr
		setupEnv    func(rt *runtime.Runtime)
		expectError bool
		expectValue int
	}{
		{
			name: "LEN simple",
			expr: &parser.LenExpr{
				Expr:   &parser.StringLiteral{Value: "APPLESOFT"},
				Line:   10,
				Column: 5,
				Token:  "LEN",
			},
			expectValue: 9,
		},
		{
			name: "Length from identifier A$",
			expr: &parser.LenExpr{
				Expr:   &parser.Identifier{Name: "A$"},
				Line:   10,
				Column: 5,
				Token:  "LEN",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A$", runtime.Value{
					Type: runtime.STRING,
					Str:  "APPLESOFT",
				})
			},
			expectValue: 9,
		},
		{
			name: "Length from expression A$+B$",
			expr: &parser.LenExpr{
				Expr: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "A$"},
					Op:    "+",
					Right: &parser.Identifier{Name: "B$"},
				},
				Line:   10,
				Column: 5,
				Token:  "LEN",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A$", runtime.Value{
					Type: runtime.STRING,
					Str:  "APPLESOFT",
				})
				rt.Env.Set("B$", runtime.Value{
					Type: runtime.STRING,
					Str:  "WARE",
				})
			},
			expectValue: 13,
		},
		{
			name: "Error EXPECTED STRING",
			expr: &parser.LenExpr{
				Expr:   &parser.NumberLiteral{Value: 5},
				Line:   10,
				Column: 5,
				Token:  "LEN",
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
				testutils.Equal(t, "result value", val.Int, tt.expectValue)
				testutils.True(t, "expected error", err != nil)
				return
			}

			testutils.True(t, "no error expected", err == nil)
			testutils.Equal(t, "result type", val.Type, runtime.INTEGER)
			testutils.Equal(t, fmt.Sprintf("result value = %d, want %d", val.Int, tt.expectValue), val.Int, tt.expectValue)
		})
	}
}

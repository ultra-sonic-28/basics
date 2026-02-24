package interpreter

import (
	"fmt"
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_AscExpr(t *testing.T) {

	tests := []struct {
		name        string
		expr        *parser.AscExpr
		setupEnv    func(rt *runtime.Runtime)
		expectError bool
		expectValue int
	}{
		{
			name: "ASC('A')",
			expr: &parser.AscExpr{
				Expr:   &parser.StringLiteral{Value: "A"},
				Line:   10,
				Column: 5,
				Token:  "ASC",
			},
			expectValue: 65,
		},
		{
			name: "ASC('ABC')",
			expr: &parser.AscExpr{
				Expr:   &parser.StringLiteral{Value: "ABC"},
				Line:   10,
				Column: 5,
				Token:  "ASC",
			},
			expectValue: 65,
		},
		{
			name: "ASC(A$)",
			expr: &parser.AscExpr{
				Expr:   &parser.Identifier{Name: "A$"},
				Line:   10,
				Column: 5,
				Token:  "ASC",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A$", runtime.Value{
					Type: runtime.STRING,
					Str:  "B",
				})
			},
			expectValue: 66,
		},
		{
			name: "Error EXPECTED STRING",
			expr: &parser.AscExpr{
				Expr:   &parser.NumberLiteral{Value: 5},
				Line:   10,
				Column: 5,
				Token:  "ASC",
			},
			expectError: true,
		},
		{
			name: "Error ILLEGAL QUANTITY (empty string)",
			expr: &parser.AscExpr{
				Expr:   &parser.StringLiteral{Value: ""},
				Line:   10,
				Column: 5,
				Token:  "ASC",
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
			testutils.Equal(t, fmt.Sprintf("result value = %d, want %d", val.Int, tt.expectValue), val.Int, tt.expectValue)
		})
	}
}

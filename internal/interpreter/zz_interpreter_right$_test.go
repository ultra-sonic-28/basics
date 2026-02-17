package interpreter

import (
	"fmt"
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_RightExpr(t *testing.T) {

	tests := []struct {
		name        string
		expr        *parser.RightExpr
		setupEnv    func(rt *runtime.Runtime)
		expectError bool
		expectStr   string
	}{
		{
			name: "RIGHT$ simple",
			expr: &parser.RightExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				LenExpr: &parser.NumberLiteral{Value: 5},
				Line:    10,
				Column:  5,
				Token:   "RIGHT$",
			},
			expectStr: "ESOFT",
		},
		{
			name: "Length greater than string",
			expr: &parser.RightExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLE"},
				LenExpr: &parser.NumberLiteral{Value: 20},
				Line:    10,
				Column:  5,
				Token:   "RIGHT$",
			},
			expectStr: "APPLE",
		},
		{
			name: "Float length truncated",
			expr: &parser.RightExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				LenExpr: &parser.NumberLiteral{Value: 5.8},
				Line:    10,
				Column:  5,
				Token:   "RIGHT$",
			},
			expectStr: "ESOFT",
		},
		{
			name: "Length from expression A*2",
			expr: &parser.RightExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				LenExpr: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "A"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2},
				},
				Line:   10,
				Column: 5,
				Token:  "RIGHT$",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  2,
				})
			},
			expectStr: "SOFT",
		},
		{
			name: "Length from expression A%*2",
			expr: &parser.RightExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				LenExpr: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "A%"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2},
				},
				Line:   10,
				Column: 5,
				Token:  "RIGHT$",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A%", runtime.Value{
					Type: runtime.INTEGER,
					Int:  2,
				})
			},
			expectStr: "SOFT",
		},
		{
			name: "Length from expression A*2+1",
			expr: &parser.RightExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				LenExpr: &parser.InfixExpr{
					Left: &parser.InfixExpr{
						Left:  &parser.Identifier{Name: "A"},
						Op:    "*",
						Right: &parser.NumberLiteral{Value: 2},
					},
					Op:    "+",
					Right: &parser.NumberLiteral{Value: 1},
				},
				Line:   10,
				Column: 5,
				Token:  "RIGHT$",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  2,
				})
			},
			expectStr: "ESOFT",
		},
		{
			name: "Error EXPECTED STRING",
			expr: &parser.RightExpr{
				StrExpr: &parser.NumberLiteral{Value: 5},
				LenExpr: &parser.NumberLiteral{Value: 2},
				Line:    10,
				Column:  5,
				Token:   "RIGHT$",
			},
			expectError: true,
		},
		{
			name: "Error EXPECTED NUMBER",
			expr: &parser.RightExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLE"},
				LenExpr: &parser.StringLiteral{Value: "TEST"},
				Line:    10,
				Column:  5,
				Token:   "RIGHT$",
			},
			expectError: true,
		},
		{
			name: "Error ILLEGAL QUANTITY",
			expr: &parser.RightExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLE"},
				LenExpr: &parser.NumberLiteral{Value: 0},
				Line:    10,
				Column:  5,
				Token:   "RIGHT$",
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
				testutils.Equal(t, "result value", val.Str, tt.expectStr)
				testutils.True(t, "expected error", err != nil)
				return
			}

			testutils.True(t, "no error expected", err == nil)
			testutils.Equal(t, "result type", val.Type, runtime.STRING)
			testutils.Equal(t, fmt.Sprintf("result value = %q, want %q", val.Str, tt.expectStr), val.Str, tt.expectStr)
		})
	}
}

package interpreter

import (
	"fmt"
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_MidExpr(t *testing.T) {

	tests := []struct {
		name        string
		expr        *parser.MidExpr
		setupEnv    func(rt *runtime.Runtime)
		expectError bool
		expectStr   string
	}{
		{
			name: "MID$ simple 2 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start:   &parser.NumberLiteral{Value: 6},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectStr: "SOFT",
		},
		{
			name: "MID$ simple 3 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start:   &parser.NumberLiteral{Value: 6},
				Len:     &parser.NumberLiteral{Value: 2},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectStr: "SO",
		},
		{
			name: "Start greater than string",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start:   &parser.NumberLiteral{Value: 20},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectStr: "",
		},
		{
			name: "Start + Len greater than string",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start:   &parser.NumberLiteral{Value: 6},
				Len:     &parser.NumberLiteral{Value: 10},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectStr: "SOFT",
		},
		{
			name: "Float start truncated",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start:   &parser.NumberLiteral{Value: 5.8},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectStr: "ESOFT",
		},
		{
			name: "Start + Len (both float) greater than string",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start:   &parser.NumberLiteral{Value: 6.8},
				Len:     &parser.NumberLiteral{Value: 5.6},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectStr: "SOFT",
		},
		{
			name: "Start from expression A*2, 2 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "A"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2},
				},
				Line:   10,
				Column: 5,
				Token:  "MID$",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  2,
				})
			},
			expectStr: "LESOFT",
		},
		{
			name: "Start from expression A*2, 3 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "A"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2},
				},
				Len: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "B"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2},
				},
				Line:   10,
				Column: 5,
				Token:  "MID$",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  2,
				})
				rt.Env.Set("B", runtime.Value{
					Type: runtime.NUMBER,
					Num:  1,
				})
			},
			expectStr: "LE",
		},
		{
			name: "Start from expression A%*2, 2 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "A%"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2},
				},
				Line:   10,
				Column: 5,
				Token:  "MID$",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A%", runtime.Value{
					Type: runtime.INTEGER,
					Int:  2,
				})
			},
			expectStr: "LESOFT",
		},
		{
			name: "Start from expression A%*2, 3 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "A%"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2},
				},
				Len: &parser.InfixExpr{
					Left:  &parser.Identifier{Name: "B%"},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 2},
				},
				Line:   10,
				Column: 5,
				Token:  "MID$",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A%", runtime.Value{
					Type: runtime.INTEGER,
					Int:  2,
				})
				rt.Env.Set("B%", runtime.Value{
					Type: runtime.INTEGER,
					Int:  1,
				})
			},
			expectStr: "LE",
		},
		{
			name: "Start from expression A*2+1, 2 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start: &parser.InfixExpr{
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
				Token:  "MID$",
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
			name: "Start from expression A*2+1, 3 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLESOFT"},
				Start: &parser.InfixExpr{
					Left: &parser.InfixExpr{
						Left:  &parser.Identifier{Name: "A"},
						Op:    "*",
						Right: &parser.NumberLiteral{Value: 2},
					},
					Op:    "+",
					Right: &parser.NumberLiteral{Value: 1},
				},
				Len: &parser.InfixExpr{
					Left: &parser.InfixExpr{
						Left:  &parser.Identifier{Name: "B"},
						Op:    "*",
						Right: &parser.NumberLiteral{Value: 2},
					},
					Op:    "+",
					Right: &parser.NumberLiteral{Value: 1},
				},
				Line:   10,
				Column: 5,
				Token:  "MID$",
			},
			setupEnv: func(rt *runtime.Runtime) {
				rt.Env.Set("A", runtime.Value{
					Type: runtime.NUMBER,
					Num:  2,
				})
				rt.Env.Set("B", runtime.Value{
					Type: runtime.NUMBER,
					Num:  1,
				})
			},
			expectStr: "ESO",
		},
		{
			name: "Error EXPECTED STRING 2 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.NumberLiteral{Value: 5},
				Start:   &parser.NumberLiteral{Value: 2},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectError: true,
		},
		{
			name: "Error EXPECTED STRING 3 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.NumberLiteral{Value: 5},
				Start:   &parser.NumberLiteral{Value: 2},
				Len:     &parser.NumberLiteral{Value: 2},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectError: true,
		},
		{
			name: "Error EXPECTED NUMBER 2 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLE"},
				Start:   &parser.StringLiteral{Value: "TEST"},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectError: true,
		},
		{
			name: "Error EXPECTED NUMBER 3 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLE"},
				Start:   &parser.NumberLiteral{Value: 2},
				Len:     &parser.StringLiteral{Value: "TEST"},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectError: true,
		},
		{
			name: "Error ILLEGAL QUANTITY 2 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLE"},
				Start:   &parser.NumberLiteral{Value: 0},
				Line:    10,
				Column:  5,
				Token:   "MID$",
			},
			expectError: true,
		},
		{
			name: "Error ILLEGAL QUANTITY 3 Args",
			expr: &parser.MidExpr{
				StrExpr: &parser.StringLiteral{Value: "APPLE"},
				Start:   &parser.NumberLiteral{Value: 2},
				Len:     &parser.NumberLiteral{Value: 0},
				Line:    10,
				Column:  5,
				Token:   "MID$",
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

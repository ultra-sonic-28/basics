package interpreter

import (
	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
	"testing"
)

func TestEvalExpr_StrExpr(t *testing.T) {
	tests := []struct {
		name       string
		expr       *parser.StrExpr
		want       string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "STR$(3)",
			expr: &parser.StrExpr{
				Expr: &parser.NumberLiteral{Value: 3},
			},
			want: "3",
		},
		{
			name: "STR$(3.5)",
			expr: &parser.StrExpr{
				Expr: &parser.NumberLiteral{Value: 3.5},
			},
			want: "3.5",
		},
		{
			name: "STR$(3*2+4)",
			expr: &parser.StrExpr{
				Expr: &parser.InfixExpr{
					Left: &parser.InfixExpr{
						Left:  &parser.NumberLiteral{Value: 3},
						Op:    "*",
						Right: &parser.NumberLiteral{Value: 2},
					},
					Op:    "+",
					Right: &parser.NumberLiteral{Value: 4},
				},
			},
			want: "10",
		},
		{
			name: "STR$(10000000000)",
			expr: &parser.StrExpr{
				Expr: &parser.NumberLiteral{Value: 10000000000},
			},
			want: "1e+10",
		},
		{
			name: "STR$(\"APPLESOFT\") error",
			expr: &parser.StrExpr{
				Expr: &parser.StringLiteral{Value: "APPLESOFT"},
			},
			wantErr:    true,
			wantErrMsg: "EXPECTED NUMBER",
		},
		{
			name: "STR$(1234567891011121314151617181920)",
			expr: &parser.StrExpr{
				Expr: &parser.NumberLiteral{Value: 1234567891011121314151617181920},
			},
			want: "1.2345678910111213e+30",
		},
		{
			name: "STR$(0.00000000000000000000000000001)",
			expr: &parser.StrExpr{
				Expr: &parser.NumberLiteral{Value: 0.00000000000000000000000000001},
			},
			want: "1e-29",
		},
		{
			name: "STR$(-1234567891011121314151617181920)",
			expr: &parser.StrExpr{
				Expr: &parser.NumberLiteral{Value: -1234567891011121314151617181920},
			},
			want: "-1.2345678910111213e+30",
		},
		{
			name: "STR$(-0.00000000000000000000000000001)",
			expr: &parser.StrExpr{
				Expr: &parser.NumberLiteral{Value: -0.00000000000000000000000000001},
			},
			want: "-1e-29",
		},
		{
			name: "STR$(1.234567891011121314151617181920)",
			expr: &parser.StrExpr{
				Expr: &parser.NumberLiteral{Value: 1.234567891011121314151617181920},
			},
			want: "1.2345678910111213",
		},
		{
			name: "STR$(1e300 * 1e300) error",
			expr: &parser.StrExpr{
				Expr: &parser.InfixExpr{
					Left:  &parser.NumberLiteral{Value: 1e300},
					Op:    "*",
					Right: &parser.NumberLiteral{Value: 1e300},
				},
			},
			wantErr:    true,
			wantErrMsg: "OVER FLOW ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := runtime.New(nil)
			val, _, err := EvalExpr(tt.expr, rt)

			if tt.wantErr {
				testutils.False(t, "error should not be nil", err == nil)
				testutils.Contains(t, "error message", tt.wantErrMsg, err.Msg)
				return
			}

			testutils.True(t, "error should be nil", err == nil)
			testutils.Equal(t, "type should be STRING", val.Type, runtime.STRING)
			testutils.Equal(t, "value should match", val.Str, tt.want)
		})
	}
}

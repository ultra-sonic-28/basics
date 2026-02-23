package interpreter

import (
	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
	"testing"
)

func TestEvalExpr_ChrExpr(t *testing.T) {
	tests := []struct {
		name       string
		expr       *parser.ChrExpr
		want       string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: `CHR$(33)`,
			expr: &parser.ChrExpr{
				Expr: &parser.NumberLiteral{Value: 33},
			},
			want: "!",
		},
		{
			name: `CHR$(66)`,
			expr: &parser.ChrExpr{
				Expr: &parser.NumberLiteral{Value: 66},
			},
			want: "B",
		},
		{
			name: `CHR$(98.62)`,
			expr: &parser.ChrExpr{
				Expr: &parser.NumberLiteral{Value: 98.62},
			},
			want: "b",
		},
		{
			name: `CHR$(35+78)`,
			expr: &parser.ChrExpr{
				Expr: &parser.InfixExpr{
					Left:  &parser.NumberLiteral{Value: 35},
					Op:    "+",
					Right: &parser.NumberLiteral{Value: 78},
				},
			},
			want: "q",
		},
		{
			name: `CHR$(-66) error`,
			expr: &parser.ChrExpr{
				Expr: &parser.NumberLiteral{Value: -66},
			},
			wantErr:    true,
			wantErrMsg: "ILLEGAL QUANTITY ERROR",
		},
		{
			name: `CHR$(300) error`,
			expr: &parser.ChrExpr{
				Expr: &parser.NumberLiteral{Value: 300},
			},
			wantErr:    true,
			wantErrMsg: "ILLEGAL QUANTITY ERROR",
		},
		{
			name: `CHR$("A") error`,
			expr: &parser.ChrExpr{
				Expr: &parser.StringLiteral{Value: "A"},
			},
			wantErr:    true,
			wantErrMsg: "EXPECTED NUMBER",
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

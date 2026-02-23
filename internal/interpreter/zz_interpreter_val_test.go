package interpreter

import (
	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
	"testing"
)

func TestEvalExpr_ValExpr(t *testing.T) {
	tests := []struct {
		name       string
		expr       *parser.ValExpr
		want       float64
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: `VAL("3")`,
			expr: &parser.ValExpr{
				Expr: &parser.StringLiteral{Value: "3"},
			},
			want: 3,
		},
		{
			name: `VAL("3.5")`,
			expr: &parser.ValExpr{
				Expr: &parser.StringLiteral{Value: "3.5"},
			},
			want: 3.5,
		},
		{
			name: `VAL("3e5")`,
			expr: &parser.ValExpr{
				Expr: &parser.StringLiteral{Value: "3e5"},
			},
			want: 300000,
		},
		{
			name: `VAL("3e-2")`,
			expr: &parser.ValExpr{
				Expr: &parser.StringLiteral{Value: "3e-2"},
			},
			want: 0.03,
		},
		{
			name: `VAL("   074")`,
			expr: &parser.ValExpr{
				Expr: &parser.StringLiteral{Value: "   074"},
			},
			want: 74,
		},
		{
			name: `VAL("   1 7 4 2 3")`,
			expr: &parser.ValExpr{
				Expr: &parser.StringLiteral{Value: "   1 7 4 2 3"},
			},
			want: 1,
		},
		{
			name: `VAL("APPLESOFT")`,
			expr: &parser.ValExpr{
				Expr: &parser.StringLiteral{Value: "APPLESOFT"},
			},
			want: 0,
		},
		{
			name: `VAL("35APPLESOFT")`,
			expr: &parser.ValExpr{
				Expr: &parser.StringLiteral{Value: "35APPLESOFT"},
			},
			want: 35,
		},
		{
			name: `VAL("3*2+4")`,
			expr: &parser.ValExpr{
				Expr: &parser.StringLiteral{Value: "3*2+4"},
			},
			want: 3,
		},
		{
			name: `VAL(3) error`,
			expr: &parser.ValExpr{
				Expr: &parser.NumberLiteral{Value: 3},
			},
			wantErr:    true,
			wantErrMsg: "EXPECTED STRING",
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
			testutils.Equal(t, "type should be NUMBER", val.Type, runtime.NUMBER)
			testutils.Equal(t, "value should match", val.Num, tt.want)
		})
	}
}

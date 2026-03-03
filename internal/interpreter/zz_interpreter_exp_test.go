package interpreter

import (
	"basics/internal/errors"
	"basics/internal/lexer"
	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
	"math"
	"testing"
)

func TestEvalExpr_ExpExpr(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		expectNum float64
	}{
		{
			name:      "EXP(0)",
			source:    "EXP(0)",
			expectNum: 1,
		},
		{
			name:      "EXP(1)",
			source:    "EXP(1)",
			expectNum: math.E,
		},
		{
			name:      "EXP(2)",
			source:    "EXP(2)",
			expectNum: math.Exp(2),
		},
		{
			name:      "EXP(-1)",
			source:    "EXP(-1)",
			expectNum: math.Exp(-1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexer.Lex("10 PRINT " + tt.source)
			p := parser.New(tokens)
			prog, errs := p.ParseProgram()
			testutils.Equal(t, "no parser errors", len(errs), 0)

			printStmt := prog.Lines[0].Stmts[0].(*parser.PrintStmt)
			expr := printStmt.Exprs[0]

			rt := runtime.New(nil)
			val, _, err := EvalExpr(expr, rt)
			testutils.Equal(t, "no eval error", err, (*errors.Error)(nil))
			testutils.Equal(t, "type is NUMBER", val.Type, runtime.NUMBER)
			testutils.Equal(t, "value is correct", val.Num, tt.expectNum)
		})
	}
}

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

func TestEvalExpr_CosExpr(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		expectNum float64
	}{
		{
			name:      "COS(0)",
			source:    "COS(0)",
			expectNum: 1,
		},
		{
			name:      "COS(3.1415926535)",
			source:    "COS(3.1415926535)",
			expectNum: math.Cos(3.1415926535),
		},
		{
			name:      "COS(1.5707963267948966)",
			source:    "COS(1.5707963267948966)",
			expectNum: math.Cos(1.5707963267948966),
		},
		{
			name:      "COS(-1.5707963267948966)",
			source:    "COS(-1.5707963267948966)",
			expectNum: math.Cos(-1.5707963267948966),
		},
		{
			name:      "COS(0.7853981633974483)",
			source:    "COS(0.7853981633974483)",
			expectNum: math.Cos(0.7853981633974483),
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

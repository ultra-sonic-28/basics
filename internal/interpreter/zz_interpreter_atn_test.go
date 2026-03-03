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

func TestEvalExpr_AtnExpr(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		expectNum float64
	}{
		{
			name:      "ATN(0)",
			source:    "ATN(0)",
			expectNum: 0,
		},
		{
			name:      "ATN(3.1415926535897932)",
			source:    "ATN(3.1415926535897932)",
			expectNum: math.Atan(3.1415926535897932),
		},
		{
			name:      "ATN(1.5707963267948966)",
			source:    "ATN(1.5707963267948966)",
			expectNum: math.Atan(1.5707963267948966),
		},
		{
			name:      "ATN(1)",
			source:    "ATN(1)",
			expectNum: math.Atan(1),
		},
		{
			name:      "ATN(0.7853981633974483)",
			source:    "ATN(0.7853981633974483)",
			expectNum: math.Atan(0.7853981633974483),
		},
		{
			name:      "ATN(2 * 1.5707963267948966)",
			source:    "ATN(2 * 1.5707963267948966)",
			expectNum: math.Atan(2 * 1.5707963267948966),
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

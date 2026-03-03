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

func TestEvalExpr_LogExpr(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		expectNum float64
	}{
		{
			name:      "LOG(1)",
			source:    "LOG(1)",
			expectNum: 0,
		},
		{
			name:      "LOG(math.E)",
			source:    "LOG(2.718281828459045)",
			expectNum: math.Log(2.718281828459045),
		},
		{
			name:      "LOG(10)",
			source:    "LOG(10)",
			expectNum: math.Log(10),
		},
		{
			name:      "LOG(0.5)",
			source:    "LOG(0.5)",
			expectNum: math.Log(0.5),
		},
		{
			name:      "LOG(0) is -Inf",
			source:    "LOG(0)",
			expectNum: math.Inf(-1),
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
			
			if math.IsInf(tt.expectNum, -1) {
				testutils.True(t, "value is -Inf", math.IsInf(val.Num, -1))
			} else {
				testutils.Equal(t, "value is correct", val.Num, tt.expectNum)
			}
		})
	}
}

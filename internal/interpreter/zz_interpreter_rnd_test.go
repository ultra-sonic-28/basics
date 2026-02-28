package interpreter

import (
	"testing"

	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestEvalExpr_RndExpr(t *testing.T) {
	rt := runtime.New(nil)

	// Test 1: RND(1) should generate a random number
	expr1 := &parser.RndExpr{Expr: &parser.NumberLiteral{Value: 1}}
	val1, _, err := EvalExpr(expr1, rt)
	testutils.True(t, "no error expected", err == nil)
	testutils.Equal(t, "result type", val1.Type, runtime.NUMBER)
	testutils.True(t, "0 <= RND(1) < 1", val1.Num >= 0 && val1.Num < 1)

	// Test 2: RND(0) should return the same number
	expr0 := &parser.RndExpr{Expr: &parser.NumberLiteral{Value: 0}}
	val0, _, err := EvalExpr(expr0, rt)
	testutils.True(t, "no error expected", err == nil)
	testutils.Equal(t, "RND(0) == LastRnd", val0.Num, val1.Num)

	// Test 3: RND(-1) should seed and return a specific number
	exprNeg1 := &parser.RndExpr{Expr: &parser.NumberLiteral{Value: -1}}
	valNeg1, _, err := EvalExpr(exprNeg1, rt)
	testutils.True(t, "no error expected", err == nil)
	
	valNeg1_bis, _, _ := EvalExpr(exprNeg1, rt)
	testutils.Equal(t, "RND(-1) is deterministic", valNeg1.Num, valNeg1_bis.Num)

	// Test 4: Sequence after RND(-1) should be deterministic
	valSeq1, _, _ := EvalExpr(expr1, rt)
	
	// Reseed with -1
	EvalExpr(exprNeg1, rt)
	valSeq1_bis, _, _ := EvalExpr(expr1, rt)
	testutils.Equal(t, "Sequence after RND(-1) is deterministic", valSeq1.Num, valSeq1_bis.Num)

	// Test 5: Different negative seed should give different number
	exprNeg2 := &parser.RndExpr{Expr: &parser.NumberLiteral{Value: -2}}
	valNeg2, _, _ := EvalExpr(exprNeg2, rt)
	testutils.True(t, "RND(-1) != RND(-2)", valNeg1.Num != valNeg2.Num)
}

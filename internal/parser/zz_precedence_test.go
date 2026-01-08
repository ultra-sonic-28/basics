package parser

import (
	"testing"

	"basics/testutils"
)

func TestPrecedenceConstantsOrder(t *testing.T) {
	testutils.True(t, "LOWEST must be > 0", LOWEST > 0)
	testutils.True(t, "EQUALS > LOWEST", EQUALS > LOWEST)
	testutils.True(t, "LESSGREATER > EQUALS", LESSGREATER > EQUALS)
	testutils.True(t, "SUM > LESSGREATER", SUM > LESSGREATER)
	testutils.True(t, "PRODUCT > SUM", PRODUCT > SUM)
	testutils.True(t, "POWER > PRODUCT", POWER > PRODUCT)
	testutils.True(t, "PREFIX > POWER", PREFIX > POWER)
}

func TestOperatorPrecedences(t *testing.T) {
	tests := []struct {
		name string
		op   string
		want int
	}{
		{"equals =", "=", EQUALS},
		{"equals <>", "<>", EQUALS},
		{"less <", "<", LESSGREATER},
		{"greater >", ">", LESSGREATER},
		{"sum +", "+", SUM},
		{"sum -", "-", SUM},
		{"product *", "*", PRODUCT},
		{"product /", "/", PRODUCT},
		{"power ^", "^", POWER},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := precedences[tt.op]
			testutils.True(t, "operator missing: "+tt.op, ok)
			testutils.Equal(t, "precedence mismatch for "+tt.op, got, tt.want)
		})
	}
}

func TestUnknownOperatorHasNoPrecedence(t *testing.T) {
	_, ok := precedences["%"]
	testutils.False(t, "unknown operator must not exist in precedences", ok)
}

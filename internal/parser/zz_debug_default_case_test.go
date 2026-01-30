package parser

import (
	"basics/testutils"
	"fmt"
	"strings"
	"testing"
)

func captureEmitter(lines *[]string) Emitter {
	return func(line string) {
		*lines = append(*lines, line)
	}
}

type fakeStmt struct{}

func (*fakeStmt) stmtNode() {}

type fakeExpr struct{}

func (*fakeExpr) exprNode() {}

func TestDumpStatement_UnknownStatement(t *testing.T) {
	var output []string
	emit := captureEmitter(&output)

	stmt := &fakeStmt{}

	dumpStatement(stmt, "  ", emit)

	msg := fmt.Sprintf("expected 1 line, got %d", len(output))
	testutils.True(t, msg, len(output) == 1)

	out := strings.Trim(output[0], " ")
	msg = fmt.Sprintf("unexpected output: %q", out)
	testutils.Contains(t, msg, out, "UNKNOWN STATEMENT")
}

func TestDumpExpr_UnknownExpr(t *testing.T) {
	var output []string
	emit := captureEmitter(&output)

	expr := &fakeExpr{}

	dumpExpr(expr, "  ", emit)

	msg := fmt.Sprintf("expected 1 line, got %d", len(output))
	testutils.True(t, msg, len(output) == 1)

	out := strings.Trim(output[0], " ")
	msg = fmt.Sprintf("unexpected output: %q", out)
	testutils.Contains(t, msg, out, "UNKNOWN EXPR")
}

func TestDumpStatement_NilStatement(t *testing.T) {
	var output []string
	emit := captureEmitter(&output)

	dumpStatement(nil, "  ", emit)

	msg := fmt.Sprintf("expected no output, got %v", output)
	testutils.True(t, msg, len(output) == 0)
}

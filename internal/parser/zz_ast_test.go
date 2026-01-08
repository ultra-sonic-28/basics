package parser

import (
	"basics/testutils"
	"testing"
)

func TestNumberLiteral_Pos(t *testing.T) {
	n := &NumberLiteral{
		Value:  42,
		Line:   10,
		Column: 5,
		Token:  "42",
	}

	line, col, tok := n.Pos()
	testutils.Equal(t, "line", line, 10)
	testutils.Equal(t, "column", col, 5)
	testutils.Equal(t, "token", tok, "42")
	testutils.Equal(t, "value", n.Value, 42.0)
}

func TestStringLiteral_Pos(t *testing.T) {
	s := &StringLiteral{
		Value:  "hello",
		Line:   3,
		Column: 1,
		Token:  "\"hello\"",
	}

	line, col, tok := s.Pos()
	testutils.Equal(t, "line", line, 3)
	testutils.Equal(t, "column", col, 1)
	testutils.Equal(t, "token", tok, "\"hello\"")
	testutils.Equal(t, "value", s.Value, "hello")
}

func TestIdentifier_Pos(t *testing.T) {
	id := &Identifier{
		Name:   "x",
		Line:   2,
		Column: 4,
		Token:  "x",
	}

	line, col, tok := id.Pos()
	testutils.Equal(t, "line", line, 2)
	testutils.Equal(t, "column", col, 4)
	testutils.Equal(t, "token", tok, "x")
	testutils.Equal(t, "name", id.Name, "x")
}

func TestPrefixExpr_Pos(t *testing.T) {
	right := &NumberLiteral{Value: 5, Line: 1, Column: 2, Token: "5"}
	expr := &PrefixExpr{
		Op:     "-",
		Right:  right,
		Line:   1,
		Column: 1,
		Token:  "-",
	}

	line, col, tok := expr.Pos()
	testutils.Equal(t, "line", line, 1)
	testutils.Equal(t, "column", col, 1)
	testutils.Equal(t, "token", tok, "-")
	testutils.Equal(t, "right value", expr.Right.(*NumberLiteral).Value, 5.0)
}

func TestInfixExpr_Pos(t *testing.T) {
	left := &NumberLiteral{Value: 2, Line: 1, Column: 1, Token: "2"}
	right := &NumberLiteral{Value: 3, Line: 1, Column: 3, Token: "3"}
	expr := &InfixExpr{
		Left:   left,
		Op:     "+",
		Right:  right,
		Line:   1,
		Column: 2,
		Token:  "+",
	}

	line, col, tok := expr.Pos()
	testutils.Equal(t, "line", line, 1)
	testutils.Equal(t, "column", col, 2)
	testutils.Equal(t, "token", tok, "+")
	testutils.Equal(t, "left value", expr.Left.(*NumberLiteral).Value, 2.0)
	testutils.Equal(t, "right value", expr.Right.(*NumberLiteral).Value, 3.0)
	testutils.Equal(t, "operator", expr.Op, "+")
}

func TestProgramAndLines(t *testing.T) {
	line1 := &Line{
		Number: 10,
		Stmts:  []Statement{},
	}
	line2 := &Line{
		Number: 20,
		Stmts:  []Statement{},
	}

	prog := &Program{
		Lines: []*Line{line1, line2},
	}

	testutils.Equal(t, "line count", len(prog.Lines), 2)
	testutils.Equal(t, "first line number", prog.Lines[0].Number, 10)
	testutils.Equal(t, "second line number", prog.Lines[1].Number, 20)
}

func TestStatementsTypes(t *testing.T) {
	printStmt := &PrintStmt{
		Exprs: []Expression{
			&NumberLiteral{Value: 1},
		},
	}
	letStmt := &LetStmt{
		Name:  "x",
		Value: &NumberLiteral{Value: 2},
	}
	forStmt := &ForStmt{
		Var:     "i",
		Start:   &NumberLiteral{Value: 1},
		End:     &NumberLiteral{Value: 10},
		Step:    nil,
		LineNum: 10,
		Column:  1,
	}
	nextStmt := &NextStmt{
		Var:        "i",
		ForLineNum: 10,
	}

	// Vérification basique des types
	var s Statement

	s = printStmt
	_, ok := s.(*PrintStmt)
	testutils.True(t, "PrintStmt type", ok)

	s = letStmt
	_, ok = s.(*LetStmt)
	testutils.True(t, "LetStmt type", ok)

	s = forStmt
	_, ok = s.(*ForStmt)
	testutils.True(t, "ForStmt type", ok)

	s = nextStmt
	_, ok = s.(*NextStmt)
	testutils.True(t, "NextStmt type", ok)
}

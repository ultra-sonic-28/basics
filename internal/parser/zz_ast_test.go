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

// ////////////////////////////////////
// HOME / HTAB / VTAB
// ////////////////////////////////////
func TestHomeStmt_Pos(t *testing.T) {
	h := &HomeStmt{
		Line:   12,
		Column: 3,
	}

	line, col, tok := h.Pos()
	testutils.Equal(t, "line", line, 12)
	testutils.Equal(t, "column", col, 3)
	testutils.Equal(t, "token", tok, "HOME")
}

func TestHTabStmt(t *testing.T) {
	expr := &NumberLiteral{Value: 5}
	stmt := &HTabStmt{Expr: expr}

	var s Statement = stmt
	_, ok := s.(*HTabStmt)
	testutils.True(t, "HTabStmt type", ok)
	testutils.Equal(t, "expr value", stmt.Expr.(*NumberLiteral).Value, 5.0)
}

func TestVTabStmt(t *testing.T) {
	expr := &NumberLiteral{Value: 10}
	stmt := &VTabStmt{Expr: expr}

	var s Statement = stmt
	_, ok := s.(*VTabStmt)
	testutils.True(t, "VTabStmt type", ok)
	testutils.Equal(t, "expr value", stmt.Expr.(*NumberLiteral).Value, 10.0)
}

// ////////////////////////////////////
// Flow Control
// ////////////////////////////////////
func TestEndStmt(t *testing.T) {
	stmt := &EndStmt{}
	var s Statement = stmt

	_, ok := s.(*EndStmt)
	testutils.True(t, "EndStmt type", ok)
}

func TestGotoStmt(t *testing.T) {
	expr := &NumberLiteral{Value: 100}
	stmt := &GotoStmt{Expr: expr}

	var s Statement = stmt
	_, ok := s.(*GotoStmt)
	testutils.True(t, "GotoStmt type", ok)
	testutils.Equal(t, "expr value", stmt.Expr.(*NumberLiteral).Value, 100.0)
}

func TestGosubStmt(t *testing.T) {
	expr := &Identifier{Name: "A"}
	stmt := &GosubStmt{Expr: expr}

	var s Statement = stmt
	_, ok := s.(*GosubStmt)
	testutils.True(t, "GosubStmt type", ok)
	testutils.Equal(t, "expr name", stmt.Expr.(*Identifier).Name, "A")
}

func TestReturnStmt(t *testing.T) {
	stmt := &ReturnStmt{}
	var s Statement = stmt

	_, ok := s.(*ReturnStmt)
	testutils.True(t, "ReturnStmt type", ok)
}

// ////////////////////////////////////
// IF / IFJUMP
// ////////////////////////////////////
func TestIfStmt(t *testing.T) {
	cond := &Identifier{Name: "X"}
	thenStmt := &PrintStmt{}
	elseStmt := &EndStmt{}

	stmt := &IfStmt{
		Cond: cond,
		Then: []Statement{thenStmt},
		Else: []Statement{elseStmt},
	}

	var s Statement = stmt
	_, ok := s.(*IfStmt)
	testutils.True(t, "IfStmt type", ok)

	testutils.Equal(t, "then count", len(stmt.Then), 1)
	testutils.Equal(t, "else count", len(stmt.Else), 1)
}

func TestIfJumpStmt(t *testing.T) {
	cond := &Identifier{Name: "FLAG"}
	stmt := &IfJumpStmt{
		Cond:   cond,
		Target: 42,
	}

	var s Statement = stmt
	_, ok := s.(*IfJumpStmt)
	testutils.True(t, "IfJumpStmt type", ok)
	testutils.Equal(t, "target", stmt.Target, 42)
}

// ////////////////////////////////////
// Maths
// ////////////////////////////////////
func TestIntExpr_Pos(t *testing.T) {
	expr := &IntExpr{
		Expr:   &NumberLiteral{Value: 3.7},
		Line:   5,
		Column: 2,
		Token:  "INT",
	}

	line, col, tok := expr.Pos()
	testutils.Equal(t, "line", line, 5)
	testutils.Equal(t, "column", col, 2)
	testutils.Equal(t, "token", tok, "INT")
}

func TestAbsExpr_Pos(t *testing.T) {
	expr := &AbsExpr{
		Expr:   &PrefixExpr{Op: "-", Right: &NumberLiteral{Value: 5}},
		Line:   8,
		Column: 1,
		Token:  "ABS",
	}

	line, col, tok := expr.Pos()
	testutils.Equal(t, "line", line, 8)
	testutils.Equal(t, "column", col, 1)
	testutils.Equal(t, "token", tok, "ABS")
}

func TestSgnExpr_Pos(t *testing.T) {
	expr := &SgnExpr{
		Expr:   &NumberLiteral{Value: -10},
		Line:   12,
		Column: 4,
		Token:  "SGN",
	}

	line, col, tok := expr.Pos()
	testutils.Equal(t, "line", line, 12)
	testutils.Equal(t, "column", col, 4)
	testutils.Equal(t, "token", tok, "SGN")
}

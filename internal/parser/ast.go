package parser

type Node interface {
	Pos() (line int, col int, token string)
}

type Program struct {
	Lines []*Line
}

type Line struct {
	Number int
	Stmts  []Statement
}

// =========================
// Statements
// =========================
type Statement interface {
	stmtNode()
}

// PRINT
type PrintStmt struct {
	Exprs      []Expression
	Separators []rune // ';' ou ',' pour chaque expression sauf la première
}

func (*PrintStmt) stmtNode() {}

// INPUT
type InputStmt struct {
	Prompt *StringLiteral // nil si absent
	Vars   []*Identifier
	Line   int
	Column int
}

func (*InputStmt) stmtNode() {}

type GetStmt struct {
	Var *Identifier
}

func (*GetStmt) stmtNode() {}

// LET
type LetStmt struct {
	Name  string
	Value Expression
}

func (*LetStmt) stmtNode() {}

// FOR ... TO ... STEP ... NEXT
type ForStmt struct {
	Var     string
	Start   Expression
	End     Expression
	Step    Expression // facultatif
	LineNum int        // Ligne du FOR
	Column  int
}

func (*ForStmt) stmtNode() {}

type NextStmt struct {
	Var        string
	ForLineNum int // Ligne du FOR correspondant
}

func (*NextStmt) stmtNode() {}

// HTAB
type HTabStmt struct {
	Expr Expression
}

func (*HTabStmt) stmtNode() {}

// VTAB
type VTabStmt struct {
	Expr Expression
}

func (*VTabStmt) stmtNode() {}

// END
type EndStmt struct {
}

func (*EndStmt) stmtNode() {}

// =======================
// HOME
// =======================

type HomeStmt struct {
	Line   int
	Column int
}

func (*HomeStmt) stmtNode() {}

func (s *HomeStmt) Pos() (int, int, string) {
	return s.Line, s.Column, "HOME"
}

// =========================
// Flow control
// =========================
type GotoStmt struct {
	Expr Expression
}

func (*GotoStmt) stmtNode() {}

type GosubStmt struct {
	Expr Expression // ligne cible (expression)
}

func (*GosubStmt) stmtNode() {}

type ReturnStmt struct{}

func (*ReturnStmt) stmtNode() {}

type IfStmt struct {
	Cond Expression
	Then []Statement
	Else []Statement // nil si absent
}

func (*IfStmt) stmtNode() {}

type IfJumpStmt struct {
	Cond   Expression
	Target int // PC cible si FAUX
}

func (*IfJumpStmt) stmtNode() {}

// =========================
// Expressions
// =========================
type Expression interface {
	exprNode()
}

// IDENTIFIER
type Identifier struct {
	Name   string
	Line   int
	Column int
	Token  string
}

func (*Identifier) exprNode() {}

func (i *Identifier) Pos() (int, int, string) {
	return i.Line, i.Column, i.Token
}

// NUMBER
type NumberLiteral struct {
	Value  float64
	Line   int
	Column int
	Token  string
}

func (*NumberLiteral) exprNode() {}

func (n *NumberLiteral) Pos() (int, int, string) {
	return n.Line, n.Column, n.Token
}

// PREFIX
type PrefixExpr struct {
	Op     string
	Right  Expression
	Line   int
	Column int
	Token  string
}

func (*PrefixExpr) exprNode() {}

func (pe *PrefixExpr) Pos() (int, int, string) {
	return pe.Line, pe.Column, pe.Token
}

// INFIX
type InfixExpr struct {
	Left   Expression
	Op     string
	Right  Expression
	Line   int
	Column int
	Token  string
}

func (*InfixExpr) exprNode() {}

func (ie *InfixExpr) Pos() (int, int, string) {
	return ie.Line, ie.Column, ie.Token
}

// STRING
type StringLiteral struct {
	Value  string
	Line   int
	Column int
	Token  string
}

func (*StringLiteral) exprNode() {}

func (s *StringLiteral) Pos() (int, int, string) {
	return s.Line, s.Column, s.Token
}

// =========================
// Maths functions
// =========================
// =========================
// INT(expr)
// =========================
type IntExpr struct {
	Expr   Expression
	Line   int
	Column int
	Token  string
}

func (*IntExpr) exprNode() {}

func (i *IntExpr) Pos() (int, int, string) {
	return i.Line, i.Column, i.Token
}

// =========================
// ABS(expr)
// =========================
type AbsExpr struct {
	Expr   Expression
	Line   int
	Column int
	Token  string
}

func (*AbsExpr) exprNode() {}

func (a *AbsExpr) Pos() (int, int, string) {
	return a.Line, a.Column, a.Token
}

// =========================
// SGN(expr)
// =========================

type SgnExpr struct {
	Expr   Expression
	Line   int
	Column int
	Token  string
}

func (*SgnExpr) exprNode() {}

func (s *SgnExpr) Pos() (int, int, string) {
	return s.Line, s.Column, s.Token
}

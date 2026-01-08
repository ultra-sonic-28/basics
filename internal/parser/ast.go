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

type PrintStmt struct {
	Exprs      []Expression
	Separators []rune // ';' ou ',' pour chaque expression sauf la première
}

func (*PrintStmt) stmtNode() {}

type LetStmt struct {
	Name  string
	Value Expression
}

func (*LetStmt) stmtNode() {}

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

// =========================
// Expressions
// =========================
type Expression interface {
	exprNode()
}

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

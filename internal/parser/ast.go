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

type EndStmt struct {
}

func (*EndStmt) stmtNode() {}

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

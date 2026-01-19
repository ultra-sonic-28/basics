package parser

import "fmt"

func StmtName(s Statement) string {
	switch s.(type) {
	case *HomeStmt:
		return "HOME"
	case *PrintStmt:
		return "PRINT"
	case *LetStmt:
		return "LET"
	case *IfStmt:
		return "IF"
	case *IfJumpStmt:
		return "IFMULTI"
	case *GotoStmt:
		return "GOTO"
	case *GosubStmt:
		return "GOSUB"
	case *ReturnStmt:
		return "RETURN"
	case *ForStmt:
		return "FOR"
	case *NextStmt:
		return "NEXT"
	case *EndStmt:
		return "END"
	case *HTabStmt:
		return "HTAB"
	case *VTabStmt:
		return "VTAB"
	default:
		return "UNKNOWN"
	}
}

func StmtArgs(s Statement) string {
	switch stmt := s.(type) {
	case *PrintStmt:
		return " ->"
	case *IfJumpStmt:
		return " ->"
	case *LetStmt:
		return fmt.Sprintf(" %s ->", stmt.Name)
	case *ForStmt:
		return fmt.Sprintf(" %s", stmt.Var)
	case *NextStmt:
		return fmt.Sprintf(" %s", stmt.Var)
	default:
		return ""
	}
}

func StmtExprValue(e Expression) string {
	switch expr := e.(type) {
	case *NumberLiteral:
		return fmt.Sprintf("%g", expr.Value)
	case *StringLiteral:
		return fmt.Sprintf("\"%s\"", expr.Value)
	case *Identifier:
		return expr.Name
	default:
		return ""
	}
}

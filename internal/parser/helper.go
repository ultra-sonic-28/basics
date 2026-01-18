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
		msg := ""
		for iexpr, expr := range stmt.Exprs {
			if len(msg) > 0 {
				msg = msg + string(stmt.Separators[iexpr-1])
			} else {
				msg = " -> "
			}
			msg = msg + StmtExprValue(expr)
		}
		return msg
	case *LetStmt:
		return fmt.Sprintf(" %s -> %s", stmt.Name, StmtExprValue(stmt.Value))
	case *IfStmt:
		return ""
	case *IfJumpStmt:
		return ""
	case *GotoStmt:
		return StmtExprValue(stmt.Expr)
	case *GosubStmt:
		return StmtExprValue(stmt.Expr)
	case *ForStmt:
		return fmt.Sprintf(" %s", stmt.Var)
	case *NextStmt:
		return fmt.Sprintf(" %s", stmt.Var)
	case *HTabStmt:
		return StmtExprValue(stmt.Expr)
	case *VTabStmt:
		return StmtExprValue(stmt.Expr)
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

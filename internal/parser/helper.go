package parser

import (
	"fmt"
	"strings"
)

func StmtName(s Statement) string {
	switch s.(type) {
	case *RemStmt:
		return "REM"
	case *HomeStmt:
		return "HOME"
	case *InputStmt:
		return "INPUT"
	case *GetStmt:
		return "GET"
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
	case *NormalStmt:
		return "NORMAL"
	case *InverseStmt:
		return "INVERSE"
	case *FlashStmt:
		return "FLASH"
	case *ClearStmt:
		return "CLEAR"
	case *DimStmt:
		return "DIM"
	default:
		return "UNKNOWN"
	}
}

func StmtArgs(s Statement) string {
	switch stmt := s.(type) {
	case *RemStmt:
		return " -> " + stmt.Text
	case *InputStmt:
		var allVars string
		for _, v := range stmt.Vars {
			allVars += v.Name
		}
		if stmt.Prompt != nil {
			return fmt.Sprintf(" -> %s -> %s", stmt.Prompt.Value, allVars)
		}
		return fmt.Sprintf(" -> %s", allVars)
	case *PrintStmt:
		var allVars strings.Builder
		for _, e := range stmt.Exprs {
			allVars.WriteString(StmtExprValue(e))
		}
		return fmt.Sprintf(" %s ->", allVars.String())
	case *IfJumpStmt:
		return " ->"
	case *LetStmt:
		if stmt.Indices != nil {
			sIndices := fmt.Sprint(FlattenIndices(stmt.Indices))
			return fmt.Sprintf(" %s(%s) ->", stmt.Name, sIndices)
		} else {
			return fmt.Sprintf(" %s ->", stmt.Name)
		}
	case *GetStmt:
		return fmt.Sprintf(" %s ->", stmt.Var.Name)
	case *ForStmt:
		return fmt.Sprintf(" %s", stmt.Var)
	case *NextStmt:
		return fmt.Sprintf(" %s", stmt.Var)
	case *DimStmt:
		var allVars strings.Builder
		for _, v := range stmt.Arrays {
			allVars.WriteString(v.Name)
			allVars.WriteRune('(')
			for j, d := range v.Dimensions {
				if j > 0 {
					allVars.WriteString(", ")
				}
				allVars.WriteString(StmtExprValue(d))
			}
			allVars.WriteRune(')')
		}
		return fmt.Sprintf(" -> %s", allVars.String())
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
	case *InfixExpr:
		sLeft := StmtExprValue(expr.Left)
		sRight := StmtExprValue(expr.Right)
		return fmt.Sprintf("%s %s %s", sLeft, expr.Op, sRight)
	case *IndexExpr:
		sIndices := fmt.Sprint(FlattenIndices(expr.Indices))
		return fmt.Sprintf("%s(%s) ", expr.Name, sIndices)
	default:
		return ""
	}
}

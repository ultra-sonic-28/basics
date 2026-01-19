package parser

import (
	"basics/internal/logger"
	"fmt"
)

type Emitter func(line string)

func StdoutEmitter(line string) {
	fmt.Println(line)
}

func LoggerEmitter(line string) {
	logger.Debug(line)
}

func DumpProgram(p *Program, emit Emitter) {
	for _, line := range p.Lines {
		emit(fmt.Sprintf("Line %d", line.Number))
		for _, stmt := range line.Stmts {
			dumpStatement(stmt, "  ", emit)
		}
	}
}

func dumpStatement(s Statement, indent string, emit Emitter) {
	switch stmt := s.(type) {

	case *PrintStmt:
		emit(indent + "PRINT")
		for i, expr := range stmt.Exprs {
			emit(fmt.Sprintf("%s  EXPR %d:", indent, i))
			dumpExpr(expr, indent+"    ", emit)
		}

	case *LetStmt:
		emit(fmt.Sprintf("%sLET %s", indent, stmt.Name))
		dumpExpr(stmt.Value, indent+"  ", emit)

	case *ForStmt:
		emit(fmt.Sprintf("%sFOR %s (Line %d)", indent, stmt.Var, stmt.LineNum))
		emit(indent + "  FROM:")
		dumpExpr(stmt.Start, indent+"    ", emit)
		emit(indent + "  TO:")
		dumpExpr(stmt.End, indent+"    ", emit)

		if stmt.Step != nil {
			emit(indent + "  STEP:")
			dumpExpr(stmt.Step, indent+"    ", emit)
		}

	case *NextStmt:
		if stmt.ForLineNum != 0 {
			emit(fmt.Sprintf("%sNEXT %s (FOR Line %d)", indent, stmt.Var, stmt.ForLineNum))
		} else {
			emit(fmt.Sprintf("%sNEXT %s", indent, stmt.Var))
		}

	case *IfStmt:
		emit(indent + "IF")
		dumpExpr(stmt.Cond, indent+"  ", emit)
		emit(indent + "THEN")

		for _, st := range stmt.Then {
			dumpStatement(st, indent+"  ", emit)
		}

		if stmt.Else != nil {
			emit(indent + "ELSE")
			for _, st := range stmt.Else {
				dumpStatement(st, indent+"  ", emit)
			}
		}

	case *GotoStmt:
		emit(indent + "GOTO")
		dumpExpr(stmt.Expr, indent+"  ", emit)

	case *GosubStmt:
		emit(indent + "GOSUB")
		dumpExpr(stmt.Expr, indent+"  ", emit)

	case *ReturnStmt:
		emit(indent + "RETURN")

	case *HomeStmt:
		emit(indent + "HOME")

	case *HTabStmt:
		emit(indent + "HTAB")
		dumpExpr(stmt.Expr, indent+"  ", emit)

	case *VTabStmt:
		emit(indent + "VTAB")
		dumpExpr(stmt.Expr, indent+"  ", emit)

	case *EndStmt:
		emit(indent + "END")

	case nil:
		// REM / instruction vide

	default:
		emit(indent + "UNKNOWN STATEMENT")
	}
}

func dumpExpr(e Expression, indent string, emit Emitter) {
	switch n := e.(type) {

	case *NumberLiteral:
		emit(fmt.Sprintf("%sNumber %v", indent, n.Value))

	case *StringLiteral:
		emit(fmt.Sprintf("%sString %q", indent, n.Value))

	case *Identifier:
		emit(fmt.Sprintf("%sIdent %s", indent, n.Name))

	case *PrefixExpr:
		emit(fmt.Sprintf("%sPrefix %s", indent, n.Op))
		dumpExpr(n.Right, indent+"  ", emit)

	case *InfixExpr:
		emit(fmt.Sprintf("%sInfix %s", indent, n.Op))
		dumpExpr(n.Left, indent+"  ", emit)
		dumpExpr(n.Right, indent+"  ", emit)

	case *IntExpr:
		emit(indent + "INT")
		dumpExpr(n.Expr, indent+"  ", emit)

	case *AbsExpr:
		emit(indent + "ABS")
		dumpExpr(n.Expr, indent+"  ", emit)

	default:
		emit(indent + "UNKNOWN EXPR")
	}
}

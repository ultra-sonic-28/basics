package parser

import (
	"basics/internal/logger"
	"fmt"
)

// DumpProgram affiche tout le programme BASIC
func DumpProgram(p *Program) {
	for _, line := range p.Lines {
		fmt.Printf("Line %d\n", line.Number) // numéro BASIC
		for _, stmt := range line.Stmts {
			dumpStatement(stmt, "  ")
		}
	}
}

func dumpStatement(s Statement, indent string) {
	switch stmt := s.(type) {

	case *PrintStmt:
		fmt.Printf("%sPRINT\n", indent)
		for i, expr := range stmt.Exprs {
			fmt.Printf("%s  EXPR %d:\n", indent, i)
			dumpExpr(expr, indent+"    ")
		}

	case *LetStmt:
		fmt.Printf("%sLET %s\n", indent, stmt.Name)
		dumpExpr(stmt.Value, indent+"  ")

	case *ForStmt:
		// afficher le numéro de ligne BASIC du FOR
		if stmt.Step != nil {
			fmt.Printf("%sFOR %s (Line %d)\n", indent, stmt.Var, stmt.LineNum)
			fmt.Printf("%s  STEP:\n", indent)
			dumpExpr(stmt.Step, indent+"    ")
		} else {
			fmt.Printf("%sFOR %s (Line %d)\n", indent, stmt.Var, stmt.LineNum)
		}
		fmt.Printf("%s  FROM:\n", indent)
		dumpExpr(stmt.Start, indent+"    ")
		fmt.Printf("%s  TO:\n", indent)
		dumpExpr(stmt.End, indent+"    ")

	case *NextStmt:
		// NEXT lié au FOR → afficher le LineNum du FOR
		if stmt.ForLineNum != 0 {
			fmt.Printf("%sNEXT %s (FOR Line %d)\n", indent, stmt.Var, stmt.ForLineNum)
		} else {
			fmt.Printf("%sNEXT %s\n", indent, stmt.Var)
		}

	case *HomeStmt:
		fmt.Printf("%sHOME\n", indent)

	case *EndStmt:
		fmt.Printf("%sEND\n", indent)

	case *HTabStmt:
		fmt.Printf("%sHTAB\n", indent)
		dumpExpr(stmt.Expr, indent+"  ")

	case *VTabStmt:
		fmt.Printf("%sVTAB\n", indent)
		dumpExpr(stmt.Expr, indent+"  ")

	case *GotoStmt:
		fmt.Printf("%sGOTO\n", indent)
		dumpExpr(stmt.Expr, indent+"  ")

	case *GosubStmt:
		fmt.Printf("%sGOSUB\n", indent)
		dumpExpr(stmt.Expr, indent+"  ")

	case *ReturnStmt:
		fmt.Printf("%sRETURN\n", indent)

	case *IfStmt:
		fmt.Printf("%sIF\n", indent)
		dumpExpr(stmt.Cond, indent+"  ")
		fmt.Printf("%sTHEN\n", indent)
		for _, stmt := range stmt.Then {
			dumpStatement(stmt, indent+"  ")
		}

		if stmt.Else != nil {
			fmt.Printf("%sELSE\n", indent)
			for _, stmt := range stmt.Else {
				dumpStatement(stmt, indent+"  ")
			}
		}

	case nil:
		// REM ou instruction vide → rien à afficher

	default:
		fmt.Println(indent, "UNKNOWN STATEMENT")
	}
}

func dumpExpr(e Expression, indent string) {
	switch n := e.(type) {

	case *NumberLiteral:
		fmt.Printf("%sNumber %v\n", indent, n.Value)

	case *StringLiteral:
		fmt.Printf("%sString \"%s\"\n", indent, n.Value)

	case *Identifier:
		fmt.Printf("%sIdent %s\n", indent, n.Name)

	case *PrefixExpr:
		fmt.Printf("%sPrefix %s\n", indent, n.Op)
		dumpExpr(n.Right, indent+"  ")

	case *InfixExpr:
		fmt.Printf("%sInfix %s\n", indent, n.Op)
		dumpExpr(n.Left, indent+"  ")
		dumpExpr(n.Right, indent+"  ")
	}
}

func logStmt(stmt Statement, indent string) {
	switch s := stmt.(type) {

	case *HomeStmt:
		logger.Debug(indent + "HOME")

	case *LetStmt:
		logger.Debug(indent + "LET " + s.Name)
		logExpr(s.Value, indent+"  ")

	case *PrintStmt:
		logger.Debug(indent + "PRINT")
		for i, e := range s.Exprs {
			logger.Debug(fmt.Sprintf("%s  EXPR %d:", indent, i))
			logExpr(e, indent+"    ")
		}

	case *GotoStmt:
		logger.Debug(indent + "GOTO")
		logExpr(s.Expr, indent+"  ")

	case *IfStmt:
		logger.Debug(indent + "IF")
		logger.Debug(indent + "  COND:")
		logExpr(s.Cond, indent+"    ")

		if len(s.Then) > 0 {
			logger.Debug(indent + "  THEN")
			for _, st := range s.Then {
				logStmt(st, indent+"    ")
			}
		}

		if len(s.Else) > 0 {
			logger.Debug(indent + "  ELSE")
			for _, st := range s.Else {
				logStmt(st, indent+"    ")
			}
		}

	case *EndStmt:
		logger.Debug(indent + "END")

	case *ForStmt:
		logger.Debug(indent + "FOR " + s.Var)
		logger.Debug(indent + "  START:")
		logExpr(s.Start, indent+"    ")
		logger.Debug(indent + "  END:")
		logExpr(s.End, indent+"    ")
		if s.Step != nil {
			logger.Debug(indent + "  STEP:")
			logExpr(s.Step, indent+"    ")
		}

	case *NextStmt:
		logger.Debug(indent + "NEXT " + s.Var)

	default:
		logger.Debug(indent + "UNKNOWN STMT")
	}
}

func logExpr(expr Expression, indent string) {
	switch e := expr.(type) {

	case *NumberLiteral:
		logger.Debug(fmt.Sprintf("%sNumber %g", indent, e.Value))

	case *StringLiteral:
		logger.Debug(fmt.Sprintf("%sString %q", indent, e.Value))

	case *Identifier:
		logger.Debug(fmt.Sprintf("%sIdent %s", indent, e.Name))

	case *InfixExpr:
		logger.Debug(fmt.Sprintf("%sInfix %s", indent, e.Op))
		logExpr(e.Left, indent+"  ")
		logExpr(e.Right, indent+"  ")

	default:
		logger.Debug(indent + "UNKNOWN EXPR")
	}
}

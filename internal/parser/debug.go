package parser

import "fmt"

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

	case *EndStmt:
		fmt.Printf("%sEND\n", indent)

	case *HTabStmt:
		fmt.Printf("%sHTAB\n", indent)
		dumpExpr(stmt.Expr, indent+"  ")

	case *VTabStmt:
		fmt.Printf("%sVTAB\n", indent)
		dumpExpr(stmt.Expr, indent+"  ")

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

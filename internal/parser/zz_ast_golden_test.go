package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type row struct {
	Path  string
	Type  string
	Value string
}

func ASTToMarkdownTable(p *Program) string {
	var rows []row

	for _, line := range p.Lines {
		linePath := fmt.Sprintf("Program/Line[%d]", line.Number)
		rows = append(rows, row{
			Path: linePath,
			Type: "Line",
		})

		for j, stmt := range line.Stmts {
			stmtPath := fmt.Sprintf("%s/Stmt[%d]", linePath, j)
			collectStmt(stmtPath, stmt, &rows)
		}
	}

	var b strings.Builder
	b.WriteString("| Path | Type | Value |\n")
	b.WriteString("|------|------|-------|\n")

	for _, r := range rows {
		b.WriteString(fmt.Sprintf(
			"| %s | %s | %s |\n",
			r.Path, r.Type, r.Value,
		))
	}

	return b.String()
}

func collectStmt(path string, s Statement, rows *[]row) {
	switch stmt := s.(type) {

	case *LetStmt:
		*rows = append(*rows, row{path, "LetStmt", stmt.Name})
		collectExpr(path+"/Value", stmt.Value, rows)

	case *PrintStmt:
		*rows = append(*rows, row{path, "PrintStmt", ""})
		for i, e := range stmt.Exprs {
			collectExpr(fmt.Sprintf("%s/Expr[%d]", path, i), e, rows)
		}

	case *ForStmt:
		*rows = append(*rows, row{path, "ForStmt", stmt.Var})
		collectExpr(path+"/Start", stmt.Start, rows)
		collectExpr(path+"/End", stmt.End, rows)

	case *NextStmt:
		*rows = append(*rows, row{path, "NextStmt", stmt.Var})

	case nil:
		*rows = append(*rows, row{path, "Empty", ""})
	}
}

func collectExpr(path string, e Expression, rows *[]row) {
	switch ex := e.(type) {

	case *NumberLiteral:
		*rows = append(*rows, row{
			path, "NumberLiteral", fmt.Sprintf("%g", ex.Value),
		})

	case *StringLiteral:
		*rows = append(*rows, row{
			path, "StringLiteral", ex.Value,
		})

	case *Identifier:
		*rows = append(*rows, row{
			path, "Identifier", ex.Name,
		})

	case *PrefixExpr:
		*rows = append(*rows, row{
			path, "PrefixExpr", ex.Op,
		})
		collectExpr(path+"/Right", ex.Right, rows)

	case *InfixExpr:
		*rows = append(*rows, row{
			path, "InfixExpr", ex.Op,
		})
		collectExpr(path+"/Left", ex.Left, rows)
		collectExpr(path+"/Right", ex.Right, rows)
	}
}

func programToMarkdown(p *Program) string {
	var b strings.Builder

	b.WriteString("| Line | Statement | Details |\n")
	b.WriteString("|------|-----------|---------|\n")

	for _, line := range p.Lines {
		for _, stmt := range line.Stmts {
			switch s := stmt.(type) {

			case *ForStmt:
				b.WriteString(
					formatRow(
						line.Number,
						"FOR",
						"var="+s.Var+
							", from="+exprString(s.Start)+
							", to="+exprString(s.End)+
							", step="+exprString(s.Step),
					),
				)

			case *NextStmt:
				b.WriteString(
					formatRow(
						line.Number,
						"NEXT",
						"var="+s.Var+
							", for_line="+itoa(s.ForLineNum),
					),
				)

			default:
				b.WriteString(
					formatRow(
						line.Number,
						"UNKNOWN",
						"",
					),
				)
			}
		}
	}

	return b.String()
}

func formatRow(line int, stmt, details string) string {
	return "| " + itoa(line) + " | " + stmt + " | " + details + " |\n"
}

func exprString(e Expression) string {
	switch v := e.(type) {
	case *NumberLiteral:
		return trimFloat(v.Value)
	case *Identifier:
		return v.Name
	default:
		return "?"
	}
}

func trimFloat(f float64) string {
	if f == float64(int64(f)) {
		return itoa(int(f))
	}
	return itoa(int(f)) // BASIC → suffisant ici
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

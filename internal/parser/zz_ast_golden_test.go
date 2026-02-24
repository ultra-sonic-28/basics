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
		for i, idx := range stmt.Indices {
			collectExpr(fmt.Sprintf("%s/Index[%d]", path, i), idx, rows)
		}
		collectExpr(path+"/Value", stmt.Value, rows)

	case *DimStmt:
		*rows = append(*rows, row{path, "DimStmt", ""})
		for i, arr := range stmt.Arrays {
			arrPath := fmt.Sprintf("%s/Array[%d]", path, i)
			*rows = append(*rows, row{arrPath, "DimDecl", arr.Name})
			for j, dim := range arr.Dimensions {
				collectExpr(fmt.Sprintf("%s/Dim[%d]", arrPath, j), dim, rows)
			}
		}

	case *ClearStmt:
		*rows = append(*rows, row{path, "ClearStmt", ""})

	case *PrintStmt:
		*rows = append(*rows, row{path, "PrintStmt", ""})
		for i, e := range stmt.Exprs {
			collectExpr(fmt.Sprintf("%s/Expr[%d]", path, i), e, rows)
		}

	case *InputStmt:
		*rows = append(*rows, row{path, "InputStmt", ""})
		if stmt.Prompt != nil {
			collectExpr(path+"/Prompt", stmt.Prompt, rows)
		}
		for i, v := range stmt.Vars {
			collectExpr(fmt.Sprintf("%s/Var[%d]", path, i), v, rows)
		}

	case *GetStmt:
		*rows = append(*rows, row{path, "GetStmt", ""})
		collectExpr(path+"/Var", stmt.Var, rows)

	case *HTabStmt:
		*rows = append(*rows, row{path, "HTabStmt", ""})
		collectExpr(path+"/Expr", stmt.Expr, rows)

	case *VTabStmt:
		*rows = append(*rows, row{path, "VTabStmt", ""})
		collectExpr(path+"/Expr", stmt.Expr, rows)

	case *NormalStmt:
		*rows = append(*rows, row{path, "NormalStmt", ""})

	case *InverseStmt:
		*rows = append(*rows, row{path, "InverseStmt", ""})

	case *FlashStmt:
		*rows = append(*rows, row{path, "FlashStmt", ""})

	case *ForStmt:
		*rows = append(*rows, row{path, "ForStmt", stmt.Var})
		collectExpr(path+"/Start", stmt.Start, rows)
		collectExpr(path+"/End", stmt.End, rows)
		if stmt.Step != nil {
			collectExpr(path+"/Step", stmt.Step, rows)
		}

	case *NextStmt:
		*rows = append(*rows, row{path, "NextStmt", stmt.Var})

	case *IfStmt:
		*rows = append(*rows, row{path, "IfStmt", ""})
		collectExpr(path+"/Cond", stmt.Cond, rows)
		for i, s := range stmt.Then {
			collectStmt(fmt.Sprintf("%s/Then[%d]", path, i), s, rows)
		}
		for i, s := range stmt.Else {
			collectStmt(fmt.Sprintf("%s/Else[%d]", path, i), s, rows)
		}

	case *GotoStmt:
		*rows = append(*rows, row{path, "GotoStmt", ""})
		collectExpr(path+"/Target", stmt.Expr, rows)

	case *GosubStmt:
		*rows = append(*rows, row{path, "GosubStmt", ""})
		collectExpr(path+"/Target", stmt.Expr, rows)

	case *ReturnStmt:
		*rows = append(*rows, row{path, "ReturnStmt", ""})

	case *RemStmt:
		*rows = append(*rows, row{path, "RemStmt", stmt.Text})

	case *EndStmt:
		*rows = append(*rows, row{path, "EndStmt", ""})

	case *PrStmt:
		*rows = append(*rows, row{path, "PrStmt", ""})
		collectExpr(path+"/Slot", stmt.Slot, rows)

	case nil:
		*rows = append(*rows, row{path, "Empty", ""})

	default:
		*rows = append(*rows, row{path, "UNKNOWN", ""})
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

	case *IndexExpr:
		*rows = append(*rows, row{path, "IndexExpr", ex.Name})
		for i, idx := range ex.Indices {
			collectExpr(fmt.Sprintf("%s/Index[%d]", path, i), idx, rows)
		}

	case *AscExpr:
		*rows = append(*rows, row{path, "AscExpr", ""})
		collectExpr(path+"/Expr", ex.Expr, rows)

	case *LenExpr:
		*rows = append(*rows, row{path, "LenExpr", ""})
		collectExpr(path+"/Expr", ex.Expr, rows)

	case *LeftExpr:
		*rows = append(*rows, row{path, "LeftExpr", ""})
		collectExpr(path+"/Str", ex.StrExpr, rows)
		collectExpr(path+"/Len", ex.LenExpr, rows)

	case *RightExpr:
		*rows = append(*rows, row{path, "RightExpr", ""})
		collectExpr(path+"/Str", ex.StrExpr, rows)
		collectExpr(path+"/Len", ex.LenExpr, rows)

	case *MidExpr:
		*rows = append(*rows, row{path, "MidExpr", ""})
		collectExpr(path+"/Str", ex.StrExpr, rows)
		collectExpr(path+"/Start", ex.Start, rows)
		if ex.Len != nil {
			collectExpr(path+"/Len", ex.Len, rows)
		}

	case *ChrExpr:
		*rows = append(*rows, row{path, "ChrExpr", ""})
		collectExpr(path+"/Expr", ex.Expr, rows)

	case *StrExpr:
		*rows = append(*rows, row{path, "StrExpr", ""})
		collectExpr(path+"/Expr", ex.Expr, rows)

	case *ValExpr:
		*rows = append(*rows, row{path, "ValExpr", ""})
		collectExpr(path+"/Expr", ex.Expr, rows)

	case *IntExpr:
		*rows = append(*rows, row{path, "IntExpr", ""})
		collectExpr(path+"/Expr", ex.Expr, rows)

	case *AbsExpr:
		*rows = append(*rows, row{path, "AbsExpr", ""})
		collectExpr(path+"/Expr", ex.Expr, rows)

	case *SqrExpr:
		*rows = append(*rows, row{path, "SqrExpr", ""})
		collectExpr(path+"/Expr", ex.Expr, rows)

	case *SgnExpr:
		*rows = append(*rows, row{path, "SgnExpr", ""})
		collectExpr(path+"/Expr", ex.Expr, rows)

	case *TabExpr:
		*rows = append(*rows, row{path, "TabExpr", ""})
		collectExpr(path+"/Expr", ex.Expr, rows)
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

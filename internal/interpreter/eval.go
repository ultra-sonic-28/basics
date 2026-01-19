package interpreter

import (
	"basics/internal/errors"
	"basics/internal/parser"
	"basics/internal/runtime"
	"math"
	"strconv"
)

func EvalExpr(expr parser.Expression, rt *runtime.Runtime) (runtime.Value, *errors.Error) {

	node, ok := expr.(parser.Node)
	if !ok {
		return runtime.Value{}, errors.NewSemantic(0, "INTERNAL AST ERROR")
	}
	line, col, tok := node.Pos()

	switch e := expr.(type) {

	case *parser.NumberLiteral:
		// NumberLitteral représente soit un flottant, soit un entier
		// La distinction se fait au niveau de l'interpreteur
		return runtime.Value{Type: runtime.NUMBER, Num: e.Value}, nil

	case *parser.StringLiteral:
		return runtime.Value{Type: runtime.STRING, Str: e.Value}, nil

	case *parser.Identifier:
		val, ok := rt.Env.Get(e.Name)
		if !ok {
			return runtime.Value{}, errors.NewSemantic(
				line,
				"UNDEFINED VARIABLE "+e.Name,
			)
		} else {
			// valeur par défaut Applesoft
			switch VarType(e.Name) {
			case "string":
				return runtime.Value{Type: runtime.STRING, Str: val.Str}, nil
			case "int":
				return runtime.Value{Type: runtime.INTEGER, Int: val.Int}, nil
			case "float":
				return runtime.Value{Type: runtime.NUMBER, Num: val.Num}, nil
			}
		}

	case *parser.PrefixExpr:
		right, err := EvalExpr(e.Right, rt)
		if err != nil {
			return runtime.Value{}, err
		}

		switch right.Type {

		case runtime.STRING:
			return runtime.Value{}, errors.NewSyntax(
				line, col, tok,
				"TYPE MISMATCH",
			)

		case runtime.INTEGER:
			switch e.Op {
			case "+":
				return right, nil
			case "-":
				return runtime.Value{
					Type: runtime.INTEGER,
					Int:  -right.Int,
				}, nil
			default:
				return runtime.Value{}, errors.NewSyntax(
					line, col, e.Op,
					"UNKNOWN PREFIX OPERATOR",
				)
			}

		case runtime.NUMBER:
			switch e.Op {
			case "+":
				return right, nil
			case "-":
				return runtime.Value{
					Type: runtime.NUMBER,
					Num:  -right.Num,
				}, nil
			default:
				return runtime.Value{}, errors.NewSyntax(
					line, col, e.Op,
					"UNKNOWN PREFIX OPERATOR",
				)
			}
		}

		// Sécurité (ne devrait jamais arriver)
		return runtime.Value{}, errors.NewSyntax(
			line, col, tok,
			"INVALID PREFIX EXPRESSION",
		)

	case *parser.InfixExpr:
		left, err := EvalExpr(e.Left, rt)
		if err != nil {
			return runtime.Value{}, err
		}

		right, err := EvalExpr(e.Right, rt)
		if err != nil {
			return runtime.Value{}, err
		}

		op := e.Op

		// =========================
		// STRING operations
		// =========================
		if left.Type == runtime.STRING || right.Type == runtime.STRING {

			// Applesoft : seul "+" est autorisé pour les chaînes
			if op != "+" {
				err = errors.NewSyntax(
					line,
					col,
					tok,
					"TYPE MISMATCH",
				)
				return runtime.Value{}, err
			}

			// conversion implicite nombre → string
			ls := ""
			rs := ""

			switch left.Type {
			case runtime.STRING:
				ls = left.Str
			case runtime.INTEGER:
				ls = strconv.Itoa(left.Int)
			default:
				ls = formatNumber(left.Num)
			}

			switch right.Type {
			case runtime.STRING:
				rs = right.Str
			case runtime.INTEGER:
				rs = strconv.Itoa(right.Int)
			default:
				rs = formatNumber(right.Num)
			}

			return runtime.Value{
				Type: runtime.STRING,
				Str:  ls + rs,
			}, nil
		}

		// =========================
		// INTEGER operations
		// =========================
		if left.Type == runtime.INTEGER && right.Type == runtime.INTEGER {

			switch op {

			case "+":
				return runtime.Value{Type: runtime.INTEGER, Int: left.Int + right.Int}, nil
			case "-":
				return runtime.Value{Type: runtime.INTEGER, Int: left.Int - right.Int}, nil
			case "*":
				return runtime.Value{Type: runtime.INTEGER, Int: left.Int * right.Int}, nil
			case "^":
				return runtime.Value{Type: runtime.INTEGER, Int: int(math.Pow(float64(left.Int), float64(right.Int)))}, nil
			case "/":
				// Applesoft : division entière → float
				if right.Int == 0 {
					err = errors.NewSyntax(
						line,
						col,
						tok,
						"DIVISION BY ZERO",
					)
					return runtime.Value{}, err
				}
				return runtime.Value{
					Type: runtime.NUMBER,
					Num:  float64(left.Int) / float64(right.Int),
				}, nil

			case "<":
				return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Int < right.Int}, nil
			case ">":
				return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Int > right.Int}, nil
			case "<=":
				return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Int <= right.Int}, nil
			case ">=":
				return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Int >= right.Int}, nil
			case "=":
				return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Int == right.Int}, nil
			case "<>":
				return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Int != right.Int}, nil
			}

			err = errors.NewSyntax(
				line,
				col,
				tok,
				"SYNTAX ERROR",
			)
			return runtime.Value{}, err
		}

		// =========================
		// MIXED or FLOAT operations
		// =========================

		// conversion implicite int → float
		lf := left.Num
		rf := right.Num

		if left.Type == runtime.INTEGER {
			lf = float64(left.Int)
		}
		if right.Type == runtime.INTEGER {
			rf = float64(right.Int)
		}

		switch op {

		case "+":
			return runtime.Value{Type: runtime.NUMBER, Num: lf + rf}, nil
		case "-":
			return runtime.Value{Type: runtime.NUMBER, Num: lf - rf}, nil
		case "*":
			return runtime.Value{Type: runtime.NUMBER, Num: lf * rf}, nil
		case "^":
			return runtime.Value{Type: runtime.NUMBER, Num: math.Pow(lf, rf)}, nil
		case "/":
			if rf == 0 {
				err = errors.NewSyntax(
					line,
					col,
					tok,
					"DIVISION BY ZERO",
				)
				return runtime.Value{}, err
			}
			return runtime.Value{Type: runtime.NUMBER, Num: lf / rf}, nil

		case "<":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: lf < rf}, nil
		case ">":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: lf > rf}, nil
		case "<=":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: lf <= rf}, nil
		case ">=":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: lf >= rf}, nil
		case "<>":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: lf != rf}, nil
		case "=":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: lf == rf}, nil
		default:
			return runtime.Value{}, errors.NewSyntax(
				line, col, e.Op,
				"UNKNOWN INFIX OPERATOR",
			)
		}

	case *parser.IntExpr:
		val, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, err
		}

		switch val.Type {

		case runtime.STRING:
			return runtime.Value{}, errors.NewSyntax(
				line, col, tok,
				"TYPE MISMATCH",
			)

		case runtime.INTEGER:
			// INT(entier) → entier inchangé
			return val, nil

		case runtime.NUMBER:
			if val.Num >= 0 {
				// Si positif alors partie entière
				// INT (1,75) -> 1
				return runtime.Value{
					Type: runtime.INTEGER,
					Int:  int(val.Num), // cast Go = floor pour positifs
				}, nil
			} else {
				// Si négatif alors partie entière - 1
				// INT (-1,75) -> -2
				return runtime.Value{
					Type: runtime.INTEGER,
					Int:  int(val.Num) - 1,
				}, nil
			}
		}

		return runtime.Value{}, errors.NewSyntax(
			line, col, tok,
			"INVALID INT OPERAND",
		)

	case *parser.AbsExpr:
		val, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, err
		}

		switch val.Type {

		case runtime.STRING:
			return runtime.Value{}, errors.NewSyntax(
				line, col, tok,
				"TYPE MISMATCH",
			)

		case runtime.INTEGER:
			if val.Int < 0 {
				return runtime.Value{
					Type: runtime.INTEGER,
					Int:  -val.Int,
				}, nil
			}
			return val, nil

		case runtime.NUMBER:
			if val.Num < 0 {
				return runtime.Value{
					Type: runtime.NUMBER,
					Num:  -val.Num,
				}, nil
			}
			return val, nil
		}

		return runtime.Value{}, errors.NewSyntax(
			line, col, tok,
			"INVALID ABS OPERAND",
		)

	}

	// =========================
	// Expression inconnue
	// =========================
	return runtime.Value{}, errors.NewSyntax(
		line,
		col,
		tok,
		"INVALID EXPRESSION",
	)

}

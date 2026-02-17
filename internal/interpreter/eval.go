package interpreter

import (
	"basics/internal/common"
	"basics/internal/errors"
	"basics/internal/parser"
	"basics/internal/runtime"
	"math"
	"strconv"
)

func EvalExpr(expr parser.Expression, rt *runtime.Runtime) (runtime.Value, *[]int, *errors.Error) {

	node, ok := expr.(parser.Node)
	if !ok {
		return runtime.Value{}, nil, errors.NewSemantic(0, "INTERNAL AST ERROR")
	}
	line, col, tok := node.Pos()

	switch e := expr.(type) {

	case *parser.NumberLiteral:
		// NumberLitteral représente soit un flottant, soit un entier
		// La distinction se fait au niveau de l'interpreteur
		return runtime.Value{Type: runtime.NUMBER, Num: e.Value}, nil, nil

	case *parser.StringLiteral:
		return runtime.Value{Type: runtime.STRING, Str: e.Value}, nil, nil

	case *parser.IndexExpr:
		// 1. récupérer la variable
		val, ok := rt.Env.Get(e.Name)
		if !ok || val.Type != runtime.ARRAY {
			return runtime.Value{}, nil, errors.NewSemantic(
				line,
				"BAD SUBSCRIPT",
			)
		}

		arr := val.Array

		// 2. évaluer les indices
		var indices []int
		for _, idxExpr := range e.Indices {
			v, _, err := EvalExpr(idxExpr, rt)
			if err != nil {
				return runtime.Value{}, nil, err
			}

			if v.Type != runtime.NUMBER && v.Type != runtime.INTEGER {
				return runtime.Value{}, nil, errors.NewSemantic(
					line,
					"BAD SUBSCRIPT",
				)
			}

			indices = append(indices, int(v.Num))
		}

		// 3. calculer l'index linéaire
		i, err := arr.Index(indices)
		if err != nil {
			return runtime.Value{}, nil, errors.NewSemantic(
				line,
				err.Error(),
			)
		}

		// 4. retourner la valeur (copie)
		elem := arr.Data[i]

		switch arr.BaseType {
		case runtime.STRING:
			return runtime.Value{Type: runtime.STRING, Str: elem.Str}, &indices, nil
		case runtime.INTEGER:
			return runtime.Value{Type: runtime.INTEGER, Int: elem.Int}, &indices, nil
		default:
			return runtime.Value{Type: runtime.NUMBER, Num: elem.Num}, &indices, nil
		}

	case *parser.Identifier:
		val, ok := rt.Env.Get(e.Name)
		if !ok {
			return runtime.Value{}, nil, errors.NewSemantic(
				line,
				"UNDEFINED VARIABLE "+e.Name,
			)
		}

		if val.Type == runtime.ARRAY {
			return runtime.Value{}, nil, errors.NewSemantic(
				line,
				"BAD SUBSCRIPT",
			)
		}

		// valeur par défaut Applesoft
		switch common.VarType(e.Name) {
		case "string":
			return runtime.Value{Type: runtime.STRING, Str: val.Str}, nil, nil
		case "int":
			return runtime.Value{Type: runtime.INTEGER, Int: val.Int}, nil, nil
		case "float":
			return runtime.Value{Type: runtime.NUMBER, Num: val.Num}, nil, nil
		}

	case *parser.PrefixExpr:
		right, _, err := EvalExpr(e.Right, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		switch right.Type {

		case runtime.STRING:
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"TYPE MISMATCH",
			)

		case runtime.INTEGER:
			switch e.Op {
			case "+":
				return right, nil, nil
			case "-":
				return runtime.Value{
					Type: runtime.INTEGER,
					Int:  -right.Int,
				}, nil, nil
			default:
				return runtime.Value{}, nil, errors.NewSyntax(
					line, col, e.Op,
					"UNKNOWN PREFIX OPERATOR",
				)
			}

		case runtime.NUMBER:
			switch e.Op {
			case "+":
				return right, nil, nil
			case "-":
				return runtime.Value{
					Type: runtime.NUMBER,
					Num:  -right.Num,
				}, nil, nil
			default:
				return runtime.Value{}, nil, errors.NewSyntax(
					line, col, e.Op,
					"UNKNOWN PREFIX OPERATOR",
				)
			}
		}

		// Sécurité (ne devrait jamais arriver)
		return runtime.Value{}, nil, errors.NewSyntax(
			line, col, tok,
			"INVALID PREFIX EXPRESSION",
		)

	case *parser.InfixExpr:
		left, _, err := EvalExpr(e.Left, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		right, _, err := EvalExpr(e.Right, rt)
		if err != nil {
			return runtime.Value{}, nil, err
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
				return runtime.Value{}, nil, err
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
			}, nil, nil
		}

		// =========================
		// INTEGER operations
		// =========================
		if left.Type == runtime.INTEGER && right.Type == runtime.INTEGER {

			switch op {

			case "+":
				return runtime.Value{Type: runtime.INTEGER, Int: left.Int + right.Int}, nil, nil
			case "-":
				return runtime.Value{Type: runtime.INTEGER, Int: left.Int - right.Int}, nil, nil
			case "*":
				return runtime.Value{Type: runtime.INTEGER, Int: left.Int * right.Int}, nil, nil
			case "^":
				return runtime.Value{Type: runtime.INTEGER, Int: int(math.Pow(float64(left.Int), float64(right.Int)))}, nil, nil
			case "/":
				// Applesoft : division entière → float
				if right.Int == 0 {
					err = errors.NewSyntax(
						line,
						col,
						tok,
						"DIVISION BY ZERO",
					)
					return runtime.Value{}, nil, err
				}
				return runtime.Value{
					Type: runtime.NUMBER,
					Num:  float64(left.Int) / float64(right.Int),
				}, nil, nil

			case "<":
				if left.Int < right.Int {
					return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
				}
				return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
			case ">":
				if left.Int > right.Int {
					return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
				}
				return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
			case "<=":
				if left.Int <= right.Int {
					return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
				}
				return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
			case ">=":
				if left.Int >= right.Int {
					return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
				}
				return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
			case "=":
				if left.Int == right.Int {
					return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
				}
				return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
			case "<>":
				if left.Int != right.Int {
					return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
				}
				return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
			}

			err = errors.NewSyntax(
				line,
				col,
				tok,
				"SYNTAX ERROR",
			)
			return runtime.Value{}, nil, err
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
			return runtime.Value{Type: runtime.NUMBER, Num: lf + rf}, nil, nil
		case "-":
			return runtime.Value{Type: runtime.NUMBER, Num: lf - rf}, nil, nil
		case "*":
			return runtime.Value{Type: runtime.NUMBER, Num: lf * rf}, nil, nil
		case "^":
			return runtime.Value{Type: runtime.NUMBER, Num: math.Pow(lf, rf)}, nil, nil
		case "/":
			if rf == 0 {
				err = errors.NewSyntax(
					line,
					col,
					tok,
					"DIVISION BY ZERO",
				)
				return runtime.Value{}, nil, err
			}
			return runtime.Value{Type: runtime.NUMBER, Num: lf / rf}, nil, nil

		case "<":
			if lf < rf {
				return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
			}
			return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
		case ">":
			if lf > rf {
				return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
			}
			return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
		case "<=":
			if lf <= rf {
				return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
			}
			return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
		case ">=":
			if lf >= rf {
				return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
			}
			return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
		case "<>":
			if lf != rf {
				return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
			}
			return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
		case "=":
			if lf == rf {
				return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
			}
			return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
		default:
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, e.Op,
				"UNKNOWN INFIX OPERATOR",
			)
		}

	case *parser.IntExpr:
		val, _, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		switch val.Type {

		case runtime.STRING:
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"TYPE MISMATCH",
			)

		case runtime.INTEGER:
			// INT(entier) → entier inchangé
			return val, nil, nil

		case runtime.NUMBER:
			if val.Num >= 0 {
				// Si positif alors partie entière
				// INT (1,75) -> 1
				return runtime.Value{
					Type: runtime.INTEGER,
					Int:  int(val.Num), // cast Go = floor pour positifs
				}, nil, nil
			} else {
				// Si négatif alors partie entière - 1
				// INT (-1,75) -> -2
				return runtime.Value{
					Type: runtime.INTEGER,
					Int:  int(math.Floor(val.Num)),
				}, nil, nil
			}
		}

		return runtime.Value{}, nil, errors.NewSyntax(
			line, col, tok,
			"INVALID INT OPERAND",
		)

	case *parser.AbsExpr:
		val, _, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		switch val.Type {

		case runtime.STRING:
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"TYPE MISMATCH",
			)

		case runtime.INTEGER:
			if val.Int < 0 {
				return runtime.Value{
					Type: runtime.INTEGER,
					Int:  -val.Int,
				}, nil, nil
			}
			return val, nil, nil

		case runtime.NUMBER:
			if val.Num < 0 {
				return runtime.Value{
					Type: runtime.NUMBER,
					Num:  -val.Num,
				}, nil, nil
			}
			return val, nil, nil
		}

		return runtime.Value{}, nil, errors.NewSyntax(
			line, col, tok,
			"INVALID ABS OPERAND",
		)

	case *parser.SqrExpr:
		val, _, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		switch val.Type {

		case runtime.STRING:
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"TYPE MISMATCH",
			)

		case runtime.INTEGER:
			if val.Int < 0 {
				return runtime.Value{}, nil, errors.NewSyntax(
					line, col, tok,
					"EXPRESSION VALUE MUST BE POSITIVE OR NULL",
				)
			}

			result := math.Sqrt(float64(val.Int))
			return runtime.Value{
				Type: runtime.NUMBER,
				Num:  result,
			}, nil, nil

		case runtime.NUMBER:
			if val.Num < 0 {
				return runtime.Value{}, nil, errors.NewSyntax(
					line, col, tok,
					"EXPRESSION VALUE MUST BE POSITIVE OR NULL",
				)
			}

			result := math.Sqrt(val.Num)
			return runtime.Value{
				Type: runtime.NUMBER,
				Num:  result,
			}, nil, nil
		}

		return runtime.Value{}, nil, errors.NewSyntax(
			line, col, tok,
			"INVALID SQR OPERAND",
		)

	case *parser.SgnExpr:
		val, _, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		var retVal float64
		switch val.Type {

		case runtime.STRING:
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"TYPE MISMATCH",
			)

		case runtime.INTEGER:
			if val.Int < 0 {
				retVal = -1
			} else if val.Int > 0 {
				retVal = 1
			} else {
				retVal = 0
			}

			return runtime.Value{
				Type: runtime.INTEGER,
				Int:  int(retVal),
			}, nil, nil

		case runtime.NUMBER:
			if val.Num < 0 {
				retVal = -1
			} else if val.Num > 0 {
				retVal = 1
			} else {
				retVal = 0
			}

			return runtime.Value{
				Type: runtime.INTEGER,
				Int:  int(retVal),
			}, nil, nil
		}

		return runtime.Value{}, nil, errors.NewSyntax(
			line, col, tok,
			"INVALID ABS OPERAND",
		)

	case *parser.LeftExpr:

		strVal, _, err := EvalExpr(e.StrExpr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}
		if strVal.Type != runtime.STRING {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"EXPECTED STRING",
			)
		}

		lenVal, _, err := EvalExpr(e.LenExpr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}
		if lenVal.Type != runtime.INTEGER && lenVal.Type != runtime.NUMBER {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"EXPECTED NUMBER",
			)
		}

		var length int

		switch lenVal.Type {
		case runtime.INTEGER:
			length = lenVal.Int
		case runtime.NUMBER:
			length = int(lenVal.Num) // partie entière
		}

		if length < 1 {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"ILLEGAL QUANTITY ERROR",
			)
		}

		str := strVal.Str

		if length > len(str) {
			length = len(str)
		}

		return runtime.Value{
			Type: runtime.STRING,
			Str:  str[:length],
		}, nil, nil

	case *parser.RightExpr:

		strVal, _, err := EvalExpr(e.StrExpr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}
		if strVal.Type != runtime.STRING {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"EXPECTED STRING",
			)
		}

		lenVal, _, err := EvalExpr(e.LenExpr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}
		if lenVal.Type != runtime.INTEGER && lenVal.Type != runtime.NUMBER {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"EXPECTED NUMBER",
			)
		}

		var length int
		switch lenVal.Type {
		case runtime.INTEGER:
			length = lenVal.Int
		case runtime.NUMBER:
			length = int(lenVal.Num)
		}

		if length < 1 {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"ILLEGAL QUANTITY ERROR",
			)
		}

		str := strVal.Str

		if length >= len(str) {
			return runtime.Value{
				Type: runtime.STRING,
				Str:  str,
			}, nil, nil
		}

		return runtime.Value{
			Type: runtime.STRING,
			Str:  str[len(str)-length:],
		}, nil, nil

	}

	// =========================
	// Expression inconnue
	// =========================
	return runtime.Value{}, nil, errors.NewSyntax(
		line,
		col,
		tok,
		"INVALID EXPRESSION",
	)

}

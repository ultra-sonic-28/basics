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
			if op != "+" && op != "=" {
				err = errors.NewSyntax(
					line,
					col,
					tok,
					"TYPE MISMATCH",
				)
				return runtime.Value{}, nil, err
			}

			switch op {
			// Concaténation de chaines de caractères
			case "+":
				// conversion implicite nombre → string
				ls := ""
				rs := ""

				switch left.Type {
				case runtime.STRING:
					ls = left.Str
				case runtime.INTEGER:
					ls = strconv.Itoa(left.Int)
				default:
					ls = common.FormatNumber(left.Num)
				}

				switch right.Type {
				case runtime.STRING:
					rs = right.Str
				case runtime.INTEGER:
					rs = strconv.Itoa(right.Int)
				default:
					rs = common.FormatNumber(right.Num)
				}

				return runtime.Value{
					Type: runtime.STRING,
					Str:  ls + rs,
				}, nil, nil

			// Test d'égalité de chaines de caractères
			case "=":
				if left.Str == right.Str {
					return runtime.Value{Type: runtime.NUMBER, Num: 1}, nil, nil
				}
				return runtime.Value{Type: runtime.NUMBER, Num: 0}, nil, nil
			}
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
			case "AND":
				if left.Int < -32768 || left.Int > 32767 || right.Int < -32768 || right.Int > 32767 {
					return runtime.Value{}, nil, errors.NewSyntax(line, col, tok, "ILLEGAL QUANTITY")
				}
				return runtime.Value{Type: runtime.INTEGER, Int: int(int16(left.Int) & int16(right.Int))}, nil, nil
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
		case "AND":
			// Applesoft logic: truncate to integer and bitwise &
			li := int(lf)
			ri := int(rf)
			if li < -32768 || li > 32767 || ri < -32768 || ri > 32767 {
				return runtime.Value{}, nil, errors.NewSyntax(line, col, tok, "ILLEGAL QUANTITY")
			}
			return runtime.Value{Type: runtime.INTEGER, Int: int(int16(li) & int16(ri))}, nil, nil
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

	case *parser.StrExpr:

		val, _, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		if val.Type == runtime.STRING {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"EXPECTED NUMBER",
			)
		}

		var out string
		switch val.Type {
		case runtime.NUMBER:
			out = common.FormatNumber(val.Num)
		case runtime.INTEGER:
			out = strconv.Itoa(val.Int)
		default:
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"EXPECTED NUMBER",
			)
		}

		if out == "OVER FLOW ERROR" {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"OVER FLOW ERROR",
			)
		}

		return runtime.Value{
			Type: runtime.STRING,
			Str:  out,
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

	case *parser.MidExpr:

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

		startVal, _, err := EvalExpr(e.Start, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}
		if startVal.Type != runtime.INTEGER && startVal.Type != runtime.NUMBER {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"EXPECTED NUMBER",
			)
		}

		var start int
		switch startVal.Type {
		case runtime.INTEGER:
			start = startVal.Int
		case runtime.NUMBER:
			start = int(startVal.Num)
		}

		if start < 1 {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"ILLEGAL QUANTITY ERROR",
			)
		}

		str := strVal.Str

		// BASIC est 1-based
		startIndex := start - 1

		if startIndex >= len(str) {
			return runtime.Value{
				Type: runtime.STRING,
				Str:  "",
			}, nil, nil
		}

		// CAS 2 paramètres
		if e.Len == nil {
			return runtime.Value{
				Type: runtime.STRING,
				Str:  str[startIndex:],
			}, nil, nil
		}

		// CAS 3 paramètres
		lenVal, _, err := EvalExpr(e.Len, rt)
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

		endIndex := min(startIndex+length, len(str))

		return runtime.Value{
			Type: runtime.STRING,
			Str:  str[startIndex:endIndex],
		}, nil, nil

	case *parser.LenExpr:

		val, _, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		if val.Type != runtime.STRING {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"TYPE MISMATCH",
			)
		}

		return runtime.Value{
			Type: runtime.INTEGER,
			Int:  len(val.Str),
			Num:  float64(len(val.Str)),
		}, nil, nil

	case *parser.AscExpr:

		val, _, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		if val.Type != runtime.STRING {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"EXPECTED STRING",
			)
		}

		if len(val.Str) == 0 {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"ILLEGAL QUANTITY ERROR",
			)
		}

		return runtime.Value{
			Type: runtime.INTEGER,
			Int:  int(val.Str[0]),
			Num:  float64(val.Str[0]),
		}, nil, nil

	case *parser.ValExpr:

		val, _, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		if val.Type != runtime.STRING {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"EXPECTED STRING",
			)
		}

		f, parseErr := common.ParseFloatApplesoft(val.Str)
		if parseErr != nil {
			if parseErr.Error() == "OVER FLOW ERROR" {
				return runtime.Value{}, nil, errors.NewSyntax(
					line, col, tok,
					"OVER FLOW ERROR",
				)
			}
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				parseErr.Error(),
			)
		}

		return runtime.Value{
			Type: runtime.NUMBER,
			Num:  f,
		}, nil, nil

	case *parser.ChrExpr:

		val, _, err := EvalExpr(e.Expr, rt)
		if err != nil {
			return runtime.Value{}, nil, err
		}

		if val.Type == runtime.STRING {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"EXPECTED NUMBER",
			)
		}

		// Coercition en entier
		intValue := int(val.Num)
		if val.Type == runtime.INTEGER {
			intValue = val.Int
		}

		if intValue < 0 || intValue > 255 {
			return runtime.Value{}, nil, errors.NewSyntax(
				line, col, tok,
				"ILLEGAL QUANTITY ERROR",
			)
		}

		return runtime.Value{
			Type: runtime.STRING,
			Str:  string(rune(intValue)),
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

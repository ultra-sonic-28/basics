package interpreter

import (
	"basics/internal/errors"
	"basics/internal/parser"
	"basics/internal/runtime"
	"math"
)

func EvalExpr(expr parser.Expression, rt *runtime.Runtime) (runtime.Value, *errors.Error) {

	node, ok := expr.(parser.Node)
	if !ok {
		return runtime.Value{}, errors.NewSemantic(0, "INTERNAL AST ERROR")
	}
	line, col, tok := node.Pos()

	switch e := expr.(type) {

	case *parser.NumberLiteral:
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
		}
		return val, nil

	case *parser.PrefixExpr:
		right, err := EvalExpr(e.Right, rt)
		if err != nil {
			return runtime.Value{}, err
		}
		if right.Type != runtime.NUMBER {
			return runtime.Value{}, errors.NewSyntax(
				line, col, tok,
				"PREFIX OPERAND MUST BE NUMBER",
			)
		}

		switch e.Op {
		case "-":
			return runtime.Value{Type: runtime.NUMBER, Num: -right.Num}, nil
		case "+":
			return right, nil
		default:
			return runtime.Value{}, errors.NewSyntax(
				line, col, e.Op,
				"UNKNOWN PREFIX OPERATOR",
			)
		}

	case *parser.InfixExpr:
		left, err := EvalExpr(e.Left, rt)
		if err != nil {
			return runtime.Value{}, err
		}
		right, err := EvalExpr(e.Right, rt)
		if err != nil {
			return runtime.Value{}, err
		}

		if e.Op == "+" && (left.Type == runtime.STRING || right.Type == runtime.STRING) {
			return runtime.Value{
				Type: runtime.STRING,
				Str:  left.String() + right.String(),
			}, nil
		}

		if left.Type != runtime.NUMBER || right.Type != runtime.NUMBER {
			return runtime.Value{}, errors.NewSyntax(
				line, col, e.Op,
				"OPERANDS MUST BE NUMBERS",
			)
		}

		switch e.Op {
		case "+":
			return runtime.Value{Type: runtime.NUMBER, Num: left.Num + right.Num}, nil
		case "-":
			return runtime.Value{Type: runtime.NUMBER, Num: left.Num - right.Num}, nil
		case "*":
			return runtime.Value{Type: runtime.NUMBER, Num: left.Num * right.Num}, nil
		case "^":
			return runtime.Value{Type: runtime.NUMBER, Num: math.Pow(left.Num, right.Num)}, nil
		case "/":
			if right.Num == 0 {
				return runtime.Value{}, errors.NewSemantic(
					line,
					"DIVISION BY ZERO",
				)
			}
			return runtime.Value{Type: runtime.NUMBER, Num: left.Num / right.Num}, nil
		case "<":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Num < right.Num}, nil
		case ">":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Num > right.Num}, nil
		case "<=":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Num <= right.Num}, nil
		case ">=":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Num >= right.Num}, nil
		case "<>":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Num != right.Num}, nil
		case "=":
			return runtime.Value{Type: runtime.BOOLEAN, Flag: left.Num == right.Num}, nil
		default:
			return runtime.Value{}, errors.NewSyntax(
				line, col, e.Op,
				"UNKNOWN INFIX OPERATOR",
			)
		}
	}

	return runtime.Value{}, errors.NewSyntax(
		line, col, tok,
		"UNKNOWN EXPRESSION",
	)
}

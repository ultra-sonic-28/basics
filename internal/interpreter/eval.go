package interpreter

import (
	"basics/internal/errors"
	"basics/internal/parser"
	"fmt"
)

func EvalExpr(expr parser.Expression, env *Env) (Value, *errors.Error) {

	node, ok := expr.(parser.Node)
	if !ok {
		return Value{}, errors.NewSemantic(0, "INTERNAL AST ERROR")
	}
	line, col, tok := node.Pos()

	switch e := expr.(type) {

	case *parser.NumberLiteral:
		return Value{Type: NUMBER, Num: e.Value}, nil

	case *parser.StringLiteral:
		return Value{Type: STRING, Str: e.Value}, nil

	case *parser.Identifier:
		val, ok := env.Get(e.Name)
		if !ok {
			return Value{}, errors.NewSemantic(
				line,
				"UNDEFINED VARIABLE "+e.Name,
			)
		}
		return val, nil

	case *parser.PrefixExpr:
		right, err := EvalExpr(e.Right, env)
		if err != nil {
			return Value{}, err
		}
		if right.Type != NUMBER {
			return Value{}, errors.NewSyntax(
				line, col, tok,
				"PREFIX OPERAND MUST BE NUMBER",
			)
		}

		switch e.Op {
		case "-":
			return Value{Type: NUMBER, Num: -right.Num}, nil
		case "+":
			return right, nil
		default:
			return Value{}, errors.NewSyntax(
				line, col, e.Op,
				"UNKNOWN PREFIX OPERATOR",
			)
		}

	case *parser.InfixExpr:
		left, err := EvalExpr(e.Left, env)
		if err != nil {
			return Value{}, err
		}
		right, err := EvalExpr(e.Right, env)
		if err != nil {
			return Value{}, err
		}

		if e.Op == "+" && (left.Type == STRING || right.Type == STRING) {
			return Value{
				Type: STRING,
				Str:  left.String() + right.String(),
			}, nil
		}

		if left.Type != NUMBER || right.Type != NUMBER {
			return Value{}, errors.NewSyntax(
				line, col, e.Op,
				"OPERANDS MUST BE NUMBERS",
			)
		}

		switch e.Op {
		case "+":
			return Value{Type: NUMBER, Num: left.Num + right.Num}, nil
		case "-":
			return Value{Type: NUMBER, Num: left.Num - right.Num}, nil
		case "*":
			return Value{Type: NUMBER, Num: left.Num * right.Num}, nil
		case "/":
			if right.Num == 0 {
				return Value{}, errors.NewSemantic(
					line,
					"DIVISION BY ZERO",
				)
			}
			return Value{Type: NUMBER, Num: left.Num / right.Num}, nil
		default:
			return Value{}, errors.NewSyntax(
				line, col, e.Op,
				"UNKNOWN INFIX OPERATOR",
			)
		}
	}

	return Value{}, errors.NewSyntax(
		line, col, tok,
		"UNKNOWN EXPRESSION",
	)
}

func (v Value) String() string {
	if v.Type == STRING {
		return v.Str
	}
	return fmt.Sprintf("%f", v.Num)
}

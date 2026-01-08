package parser

import (
	"basics/internal/errors"
	"basics/internal/token"
	"basics/testutils"
	"testing"
)

func TestNewParser(t *testing.T) {
	tokens := []token.Token{
		{Type: token.LINENUM, Literal: "10"},
		{Type: token.KEYWORD, Literal: "PRINT"},
		{Type: token.EOF},
	}

	p := New(tokens)

	// Vérifier que curr est le premier token
	testutils.Equal(t, "curr token", p.curr.Literal, "10")

	// Vérifier que peek est le deuxième token
	testutils.Equal(t, "peek token", p.peek.Literal, "PRINT")

	// Vérifier la position initiale
	testutils.Equal(t, "position", p.pos, 0)
}

func TestNext(t *testing.T) {
	tokens := []token.Token{
		{Type: token.LINENUM, Literal: "10"},
		{Type: token.KEYWORD, Literal: "PRINT"},
		{Type: token.EOF},
	}

	p := New(tokens)

	// ✅ état initial
	testutils.Equal(t, "initial pos", p.pos, 0)
	testutils.Equal(t, "initial curr", p.curr.Literal, "10")
	testutils.Equal(t, "initial peek", p.peek.Literal, "PRINT")

	// Appel à next()
	p.next()

	testutils.Equal(t, "pos after 1 next", p.pos, 1)
	testutils.Equal(t, "curr after 1 next", p.curr.Literal, "PRINT")
	testutils.Equal(t, "peek after 1 next", p.peek.Literal, "")

	// Appel à next() encore une fois
	p.next()

	testutils.Equal(t, "pos after 2 next", p.pos, 2)
	testutils.Equal(t, "curr after 2 next", p.curr.Literal, "")

	// peek doit rester EOF ou vide (fin des tokens)
	if len(p.tokens) > p.pos+1 {
		testutils.Equal(t, "peek after 2 next", p.peek.Literal, p.tokens[p.pos+1].Literal)
	} else {
		testutils.True(t, "peek after 2 next is EOF or empty", p.peek.Type == token.EOF)
	}
}

func TestSyntaxError(t *testing.T) {
	tok := token.Token{
		Type:    token.KEYWORD,
		Literal: "PRINT",
		Line:    10,
		Column:  5,
	}

	p := &Parser{
		tokens: []token.Token{tok},
		curr:   tok,
	}

	// Appel de syntaxError
	p.syntaxError("unexpected token")

	// Il doit y avoir exactement 1 erreur
	testutils.Equal(t, "number of errors", len(p.errors), 1)

	errObj := p.errors[0]

	// Vérifier que c'est bien un *errors.Error
	testutils.True(t, "error type", errObj.Kind == errors.Syntax)

	// Vérifier que les champs Line, Column, Token et Msg sont corrects
	testutils.Equal(t, "error line", errObj.Line, 10)
	testutils.Equal(t, "error column", errObj.Column, 5)
	testutils.Equal(t, "error token", errObj.Token, "PRINT")
	testutils.Equal(t, "error message", errObj.Msg, "unexpected token")
}

func TestCurrPrecedence_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		tok      token.Token
		expected int
	}{
		{"EQUAL operator", token.Token{Type: token.EQUAL, Literal: "="}, EQUALS},
		{"NEQ operator", token.Token{Type: token.NEQ, Literal: "<>"}, EQUALS},
		{"LESS operator", token.Token{Type: token.LT, Literal: "<"}, LESSGREATER},
		{"GREATER operator", token.Token{Type: token.GT, Literal: ">"}, LESSGREATER},
		{"PLUS operator", token.Token{Type: token.PLUS, Literal: "+"}, SUM},
		{"MINUS operator", token.Token{Type: token.MINUS, Literal: "-"}, SUM},
		{"PRODUCT operator", token.Token{Type: token.ASTERISK, Literal: "*"}, PRODUCT},
		{"DIV operator", token.Token{Type: token.SLASH, Literal: "/"}, PRODUCT},
		{"POWER operator", token.Token{Type: token.CARET, Literal: "^"}, POWER},
		{"UNKNOWN token", token.Token{Type: token.IDENT, Literal: "XYZ"}, LOWEST},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				tokens: []token.Token{tt.tok},
				curr:   tt.tok,
			}
			got := p.currPrecedence()
			testutils.Equal(t, tt.name, got, tt.expected)
		})
	}
}

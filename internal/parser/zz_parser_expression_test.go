package parser

import (
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestParser_parseExpression(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []token.Token
		wantErr bool
		check   func(t *testing.T, got Expression)
	}{
		{
			name: "number literal",
			tokens: []token.Token{
				{Type: token.NUMBER, Literal: "42", Line: 1, Column: 1},
			},
			check: func(t *testing.T, got Expression) {
				n, ok := got.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", n.Value, 42.0)
			},
		},
		{
			name: "string literal",
			tokens: []token.Token{
				{Type: token.STRING, Literal: "hello", Line: 1, Column: 1},
			},
			check: func(t *testing.T, got Expression) {
				s, ok := got.(*StringLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", s.Value, "hello")
			},
		},
		{
			name: "identifier",
			tokens: []token.Token{
				{Type: token.IDENT, Literal: "X", Line: 1, Column: 1},
			},
			check: func(t *testing.T, got Expression) {
				id, ok := got.(*Identifier)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", id.Name, "X")
			},
		},
		{
			name: "prefix minus number",
			tokens: []token.Token{
				{Type: token.MINUS, Literal: "-", Line: 1, Column: 1},
				{Type: token.NUMBER, Literal: "5", Line: 1, Column: 2},
			},
			check: func(t *testing.T, got Expression) {
				pre, ok := got.(*PrefixExpr)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", pre.Op, "-")
				n, ok := pre.Right.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", n.Value, 5.0)
			},
		},
		{
			name: "prefix equals identifier",
			tokens: []token.Token{
				{Type: token.EQUAL, Literal: "=", Line: 1, Column: 1},
				{Type: token.IDENT, Literal: "Y", Line: 1, Column: 2},
			},
			check: func(t *testing.T, got Expression) {
				pre, ok := got.(*PrefixExpr)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", pre.Op, "=")
				id, ok := pre.Right.(*Identifier)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", id.Name, "Y")
			},
		},
		{
			name: "simple infix addition",
			tokens: []token.Token{
				{Type: token.NUMBER, Literal: "1", Line: 1, Column: 1},
				{Type: token.PLUS, Literal: "+", Line: 1, Column: 2},
				{Type: token.NUMBER, Literal: "2", Line: 1, Column: 3},
			},
			check: func(t *testing.T, got Expression) {
				inf, ok := got.(*InfixExpr)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", inf.Op, "+")
				l, ok := inf.Left.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", l.Value, 1.0)
				r, ok := inf.Right.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", r.Value, 2.0)
			},
		},
		{
			name: "infix with precedence (1 + 2 * 3)",
			tokens: []token.Token{
				{Type: token.NUMBER, Literal: "1", Line: 1, Column: 1},
				{Type: token.PLUS, Literal: "+", Line: 1, Column: 2},
				{Type: token.NUMBER, Literal: "2", Line: 1, Column: 3},
				{Type: token.ASTERISK, Literal: "*", Line: 1, Column: 4},
				{Type: token.NUMBER, Literal: "3", Line: 1, Column: 5},
			},
			check: func(t *testing.T, got Expression) {
				inf, ok := got.(*InfixExpr)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", inf.Op, "+")
				l, ok := inf.Left.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", l.Value, 1.0)

				r, ok := inf.Right.(*InfixExpr)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", r.Op, "*")
				rL, ok := r.Left.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", rL.Value, 2.0)
				rR, ok := r.Right.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", rR.Value, 3.0)
			},
		},
		{
			name: "nested prefix and infix (-1 + 2)",
			tokens: []token.Token{
				{Type: token.MINUS, Literal: "-", Line: 1, Column: 1},
				{Type: token.NUMBER, Literal: "1", Line: 1, Column: 2},
				{Type: token.PLUS, Literal: "+", Line: 1, Column: 3},
				{Type: token.NUMBER, Literal: "2", Line: 1, Column: 4},
			},
			check: func(t *testing.T, got Expression) {
				inf, ok := got.(*InfixExpr)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", inf.Op, "+")
				l, ok := inf.Left.(*PrefixExpr)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", l.Op, "-")
				num, ok := l.Right.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", num.Value, 1.0)
				r, ok := inf.Right.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", r.Value, 2.0)
			},
		},
		{
			name: "invalid token",
			tokens: []token.Token{
				{Type: token.EOL, Literal: "", Line: 1, Column: 1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.tokens)
			got := p.parseExpression(LOWEST)
			if tt.wantErr {
				testutils.True(t, "expected syntax error", len(p.errors) > 0)
				return
			}
			testutils.False(t, "unexpected syntax error", len(p.errors) > 0)
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

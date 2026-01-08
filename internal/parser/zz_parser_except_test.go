package parser

import (
	"basics/internal/token"
	"basics/testutils"
	"testing"
)

func TestExpectLiteral_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		startTok token.Token
		lit      string
		want     bool
	}{
		{
			name:     "matching literal",
			startTok: token.Token{Type: token.EQUAL, Literal: "="},
			lit:      "=",
			want:     true,
		},
		{
			name:     "non-matching literal",
			startTok: token.Token{Type: token.NEQ, Literal: "<>"},
			lit:      "=",
			want:     false,
		},
		{
			name:     "another non-matching literal",
			startTok: token.Token{Type: token.IDENT, Literal: "X"},
			lit:      "Y",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				tokens: []token.Token{tt.startTok},
				curr:   tt.startTok,
			}

			got := p.expectLiteral(tt.lit)
			testutils.Equal(t, "return value", got, tt.want)

			if got {
				// si true, le token doit avoir été consommé → curr avance (dans ce test il n'y a qu'un token, donc curr sera hors liste)
				testutils.NotEqual(t, "token consumed", p.curr.Literal, tt.startTok.Literal)
			} else {
				// si false, le token courant n'est pas consommé
				testutils.Equal(t, "token not consumed", p.curr.Literal, tt.startTok.Literal)
			}

			// Vérifier que si false, une erreur de syntaxe a été ajoutée
			if !tt.want {
				testutils.True(t, "syntax error recorded", len(p.errors) == 1)
				testutils.Equal(t, "error token", p.errors[0].Token, tt.startTok.Literal)
			}
		})
	}
}

func TestExpectKeyword_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		startTok token.Token
		kw       string
		want     bool
	}{
		{
			name:     "matching keyword",
			startTok: token.Token{Type: token.KEYWORD, Literal: "FOR"},
			kw:       "FOR",
			want:     true,
		},
		{
			name:     "non-matching keyword literal",
			startTok: token.Token{Type: token.KEYWORD, Literal: "NEXT"},
			kw:       "FOR",
			want:     false,
		},
		{
			name:     "non-keyword token",
			startTok: token.Token{Type: token.IDENT, Literal: "X"},
			kw:       "FOR",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				tokens: []token.Token{tt.startTok},
				curr:   tt.startTok,
			}

			got := p.expectKeyword(tt.kw)
			testutils.Equal(t, "return value", got, tt.want)

			if got {
				// si true, token consommé → curr change
				testutils.NotEqual(t, "token consumed", p.curr.Literal, tt.startTok.Literal)
			} else {
				// si false, token non consommé
				testutils.Equal(t, "token not consumed", p.curr.Literal, tt.startTok.Literal)
			}

			// Vérifier qu'une erreur de syntaxe est enregistrée si false
			if !tt.want {
				testutils.True(t, "syntax error recorded", len(p.errors) == 1)
				testutils.Equal(t, "error token", p.errors[0].Token, tt.startTok.Literal)
			}
		})
	}
}

func TestExpect_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		startTok  token.Token
		expectTyp token.TokenType
		want      bool
	}{
		{
			name:      "matching type",
			startTok:  token.Token{Type: token.IDENT, Literal: "X"},
			expectTyp: token.IDENT,
			want:      true,
		},
		{
			name:      "non-matching type",
			startTok:  token.Token{Type: token.NUMBER, Literal: "42"},
			expectTyp: token.IDENT,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				tokens: []token.Token{tt.startTok},
				curr:   tt.startTok,
			}

			got := p.expect(tt.expectTyp)
			testutils.Equal(t, "return value", got, tt.want)

			if got {
				// token consommé → curr change
				testutils.NotEqual(t, "token consumed", p.curr.Literal, tt.startTok.Literal)
			} else {
				// token non consommé
				testutils.Equal(t, "token not consumed", p.curr.Literal, tt.startTok.Literal)
			}

			// si false → syntaxError ajouté
			if !tt.want {
				testutils.True(t, "syntax error recorded", len(p.errors) == 1)
				testutils.Equal(t, "error token", p.errors[0].Token, tt.startTok.Literal)
			}
		})
	}
}

package parser

import (
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestParser_parseLine(t *testing.T) {
	tests := []struct {
		name      string
		tokens    []token.Token
		wantLine  int
		wantStmts int
		wantErr   bool
	}{
		{
			name: "PRINT single number",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "10", Line: 1, Column: 1},
				{Type: token.KEYWORD, Literal: "PRINT", Line: 1, Column: 4},
				{Type: token.NUMBER, Literal: "42", Line: 1, Column: 10},
				{Type: token.EOL, Line: 1, Column: 12},
			},
			wantLine:  10,
			wantStmts: 1,
		},
		{
			name: "LET statement",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "20", Line: 2, Column: 1},
				{Type: token.KEYWORD, Literal: "LET", Line: 2, Column: 4},
				{Type: token.IDENT, Literal: "A", Line: 2, Column: 8},
				{Type: token.EQUAL, Literal: "=", Line: 2, Column: 10},
				{Type: token.NUMBER, Literal: "10", Line: 2, Column: 12},
				{Type: token.EOL, Line: 2, Column: 14},
			},
			wantLine:  20,
			wantStmts: 1,
		},
		{
			name: "Multiple statements on one line",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "30", Line: 3, Column: 1},
				{Type: token.KEYWORD, Literal: "PRINT", Line: 3, Column: 4},
				{Type: token.NUMBER, Literal: "1", Line: 3, Column: 10},
				{Type: token.COLON, Literal: ":", Line: 3, Column: 11},
				{Type: token.KEYWORD, Literal: "LET", Line: 3, Column: 12},
				{Type: token.IDENT, Literal: "B", Line: 3, Column: 16},
				{Type: token.EQUAL, Literal: "=", Line: 3, Column: 18},
				{Type: token.NUMBER, Literal: "2", Line: 3, Column: 20},
				{Type: token.EOL, Line: 3, Column: 22},
			},
			wantLine:  30,
			wantStmts: 2,
		},
		{
			name: "FOR ... NEXT",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "40", Line: 4, Column: 1},
				{Type: token.KEYWORD, Literal: "FOR", Line: 4, Column: 4},
				{Type: token.IDENT, Literal: "I", Line: 4, Column: 8},
				{Type: token.EQUAL, Literal: "=", Line: 4, Column: 10},
				{Type: token.NUMBER, Literal: "1", Line: 4, Column: 12},
				{Type: token.KEYWORD, Literal: "TO", Line: 4, Column: 14},
				{Type: token.NUMBER, Literal: "5", Line: 4, Column: 17},
				{Type: token.COLON, Literal: ":", Line: 4, Column: 19},
				{Type: token.KEYWORD, Literal: "NEXT", Line: 4, Column: 21},
				{Type: token.IDENT, Literal: "I", Line: 4, Column: 26},
				{Type: token.EOL, Line: 4, Column: 28},
			},
			wantLine:  40,
			wantStmts: 2,
		},
		{
			name: "LET implicit",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "50", Line: 5, Column: 1},
				{Type: token.IDENT, Literal: "X", Line: 5, Column: 4},
				{Type: token.EQUAL, Literal: "=", Line: 5, Column: 5},
				{Type: token.NUMBER, Literal: "99", Line: 5, Column: 6},
				{Type: token.EOL, Line: 5, Column: 8},
			},
			wantLine:  50,
			wantStmts: 1,
		},
		{
			name: "Unknown statement token",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "60", Line: 6, Column: 1},
				{Type: token.KEYWORD, Literal: "FOO", Line: 6, Column: 4},
				{Type: token.EOL, Line: 6, Column: 7},
			},
			wantLine:  60,
			wantStmts: 1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.tokens)
			line := p.parseLine()
			testutils.Equal(t, "", line.Number, tt.wantLine)
			testutils.Equal(t, "", len(line.Stmts), tt.wantStmts)
			if tt.wantErr {
				testutils.True(t, "expected syntax errors", len(p.errors) > 0)
			} else {
				testutils.False(t, "unexpected syntax errors", len(p.errors) > 0)
			}
		})
	}
}

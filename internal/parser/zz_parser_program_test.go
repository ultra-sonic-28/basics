package parser

import (
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestParser_ParseProgram(t *testing.T) {
	tests := []struct {
		name       string
		tokens     []token.Token
		wantLines  int
		wantErrors bool
	}{
		{
			name: "Simple program with two lines",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "10"},
				{Type: token.KEYWORD, Literal: "PRINT"},
				{Type: token.NUMBER, Literal: "1"},
				{Type: token.EOL},

				{Type: token.LINENUM, Literal: "20"},
				{Type: token.KEYWORD, Literal: "PRINT"},
				{Type: token.NUMBER, Literal: "2"},
				{Type: token.EOL},

				{Type: token.EOF},
			},
			wantLines: 2,
		},
		{
			name: "FOR NEXT valid",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "10"},
				{Type: token.KEYWORD, Literal: "FOR"},
				{Type: token.IDENT, Literal: "I"},
				{Type: token.EQUAL, Literal: "="},
				{Type: token.NUMBER, Literal: "1"},
				{Type: token.KEYWORD, Literal: "TO"},
				{Type: token.NUMBER, Literal: "3"},
				{Type: token.EOL},

				{Type: token.LINENUM, Literal: "20"},
				{Type: token.KEYWORD, Literal: "NEXT"},
				{Type: token.IDENT, Literal: "I"},
				{Type: token.EOL},

				{Type: token.EOF},
			},
			wantLines: 2,
		},
		{
			name: "NEXT without FOR",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "10"},
				{Type: token.KEYWORD, Literal: "NEXT"},
				{Type: token.IDENT, Literal: "I"},
				{Type: token.EOL},
				{Type: token.EOF},
			},
			wantLines:  1,
			wantErrors: true,
		},
		{
			name: "FOR not closed",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "10"},
				{Type: token.KEYWORD, Literal: "FOR"},
				{Type: token.IDENT, Literal: "I"},
				{Type: token.EQUAL, Literal: "="},
				{Type: token.NUMBER, Literal: "1"},
				{Type: token.KEYWORD, Literal: "TO"},
				{Type: token.NUMBER, Literal: "5"},
				{Type: token.EOL},
				{Type: token.EOF},
			},
			wantLines:  1,
			wantErrors: true,
		},
		{
			name: "Syntax error in middle of program",
			tokens: []token.Token{
				{Type: token.LINENUM, Literal: "10"},
				{Type: token.KEYWORD, Literal: "PRINT"},
				{Type: token.NUMBER, Literal: "1"},
				{Type: token.EOL},

				{Type: token.LINENUM, Literal: "20"},
				{Type: token.KEYWORD, Literal: "FOO"},
				{Type: token.EOL},

				{Type: token.LINENUM, Literal: "30"},
				{Type: token.KEYWORD, Literal: "PRINT"},
				{Type: token.NUMBER, Literal: "3"},
				{Type: token.EOL},

				{Type: token.EOF},
			},
			wantLines:  3,
			wantErrors: true,
		},
		{
			name: "Empty program",
			tokens: []token.Token{
				{Type: token.EOF},
			},
			wantLines: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.tokens)
			prog, _ := p.ParseProgram()

			testutils.Equal(t, "line count", len(prog.Lines), tt.wantLines)

			if tt.wantErrors {
				testutils.True(t, "expected errors", len(p.errors) > 0)
			} else {
				testutils.False(t, "unexpected errors", len(p.errors) > 0)
			}
		})
	}
}

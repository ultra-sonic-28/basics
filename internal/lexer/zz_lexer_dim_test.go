package lexer

import (
	"basics/internal/token"
	"basics/testutils"
	"fmt"
	"testing"
)

type dimTestCase struct {
	name   string
	input  string
	tokens []token.TokenType
	lits   []string
}

func TestLexer_DIM_TableDriven(t *testing.T) {
	tests := []dimTestCase{
		{
			name:  "DIM float array",
			input: "20 DIM A(10)",
			tokens: []token.TokenType{
				token.LINENUM, token.KEYWORD, token.IDENT,
				token.LPAREN, token.NUMBER, token.RPAREN,
				token.EOF,
			},
			lits: []string{
				"20", "DIM", "A", "(", "10", ")", "",
			},
		},
		{
			name:  "DIM int array",
			input: "20 DIM A%(10)",
			tokens: []token.TokenType{
				token.LINENUM, token.KEYWORD, token.IDENT,
				token.LPAREN, token.NUMBER, token.RPAREN,
				token.EOF,
			},
			lits: []string{
				"20", "DIM", "A%", "(", "10", ")", "",
			},
		},
		{
			name:  "DIM string array",
			input: "20 DIM A$(10)",
			tokens: []token.TokenType{
				token.LINENUM, token.KEYWORD, token.IDENT,
				token.LPAREN, token.NUMBER, token.RPAREN,
				token.EOF,
			},
			lits: []string{
				"20", "DIM", "A$", "(", "10", ")", "",
			},
		},
		{
			name:  "DIM multi dimensions",
			input: "20 DIM A(10,20)",
			tokens: []token.TokenType{
				token.LINENUM, token.KEYWORD, token.IDENT,
				token.LPAREN,
				token.NUMBER, token.COMMA, token.NUMBER,
				token.RPAREN,
				token.EOF,
			},
			lits: []string{
				"20", "DIM", "A",
				"(", "10", ",", "20", ")", "",
			},
		},
		{
			name:  "DIM multiple vars",
			input: "20 DIM A(10), B%(50), C$(10,10)",
			tokens: []token.TokenType{
				token.LINENUM, token.KEYWORD,
				token.IDENT, token.LPAREN, token.NUMBER, token.RPAREN,
				token.COMMA,
				token.IDENT, token.LPAREN, token.NUMBER, token.RPAREN,
				token.COMMA,
				token.IDENT, token.LPAREN,
				token.NUMBER, token.COMMA, token.NUMBER,
				token.RPAREN,
				token.EOF,
			},
			lits: []string{
				"20", "DIM",
				"A", "(", "10", ")",
				",",
				"B%", "(", "50", ")",
				",",
				"C$", "(", "10", ",", "10", ")",
				"",
			},
		},
		{
			name:  "DIM chained with colon",
			input: "20 DIM OA$(75):DIM OD$(75)",
			tokens: []token.TokenType{
				token.LINENUM, token.KEYWORD, token.IDENT,
				token.LPAREN, token.NUMBER, token.RPAREN,
				token.COLON,
				token.KEYWORD, token.IDENT,
				token.LPAREN, token.NUMBER, token.RPAREN,
				token.EOF,
			},
			lits: []string{
				"20", "DIM", "OA$", "(", "75", ")",
				":",
				"DIM", "OD$", "(", "75", ")",
				"",
			},
		},
		{
			name:  "DIM with variable size",
			input: "30 DIM A%(Size)",
			tokens: []token.TokenType{
				token.LINENUM, token.KEYWORD, token.IDENT,
				token.LPAREN, token.IDENT, token.RPAREN,
				token.EOF,
			},
			lits: []string{
				"30", "DIM", "A%", "(", "Size", ")", "",
			},
		},
		{
			name:  "DIM with expression",
			input: "30 DIM A%(2*Size)",
			tokens: []token.TokenType{
				token.LINENUM, token.KEYWORD, token.IDENT,
				token.LPAREN,
				token.NUMBER, token.ASTERISK, token.IDENT,
				token.RPAREN,
				token.EOF,
			},
			lits: []string{
				"30", "DIM", "A%",
				"(", "2", "*", "Size", ")", "",
			},
		},
		{
			name:  "DIM with complex dimensions",
			input: "30 DIM A%(Size,10,Size/10,2*Size,5)",
			tokens: []token.TokenType{
				token.LINENUM, token.KEYWORD, token.IDENT,
				token.LPAREN,
				token.IDENT, token.COMMA,
				token.NUMBER, token.COMMA,
				token.IDENT, token.SLASH, token.NUMBER, token.COMMA,
				token.NUMBER, token.ASTERISK, token.IDENT, token.COMMA,
				token.NUMBER,
				token.RPAREN,
				token.EOF,
			},
			lits: []string{
				"30", "DIM", "A%",
				"(",
				"Size", ",",
				"10", ",",
				"Size", "/", "10", ",",
				"2", "*", "Size", ",",
				"5",
				")",
				"",
			},
		},
	}

	for ti, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)

			for i, expectedType := range tt.tokens {
				tok := l.NextToken()

				testutils.True(
					t,
					fmt.Sprintf("test[%d] token[%d] type got=%q want=%q",
						ti, i, tok.Type, expectedType),
					tok.Type == expectedType,
				)

				testutils.True(
					t,
					fmt.Sprintf("test[%d] token[%d] literal got=%q want=%q",
						ti, i, tok.Literal, tt.lits[i]),
					tok.Literal == tt.lits[i],
				)
			}
		})
	}
}

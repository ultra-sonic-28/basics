package lexer

import "basics/internal/token"

// Lex tokenize entièrement la source BASIC et retourne tous les tokens (EOF inclus)
func Lex(input string) []token.Token {
	l := New(input)
	var tokens []token.Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type == token.EOF {
			break
		}
	}

	return tokens
}

package lexer

import (
	"basics/internal/logger"
	"basics/internal/token"
	"fmt"
	"os"
)

// Lex tokenize entièrement la source BASIC et retourne tous les tokens (EOF inclus)
func Lex(input string) []token.Token {
	l := New(input)
	var tokens []token.Token

	for {
		tok := l.NextToken()
		logger.Debug(fmt.Sprintf("New token: %s", DumpTokenToLogFile(tok)))
		tokens = append(tokens, tok)

		// We reach end of file
		if tok.Type == token.EOF {
			break
		}

		// We got an invalid token, stop lexing and exit now
		if tok.Type == token.ILLEGAL {
			fmt.Printf(
				"⚠️ Invalid token found in %d (%s)\n",
				tok.Line,
				tok.Literal,
			)

			os.Exit(1)
		}
	}

	return tokens
}

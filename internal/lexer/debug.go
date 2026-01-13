package lexer

import (
	"basics/internal/token"
	"fmt"
)

func DumpTokens(tokens []token.Token) {
	for _, t := range tokens {
		fmt.Println(DumpTokenToLogFile(t))
	}
}

func DumpTokenToLogFile(token token.Token) string {
	var str string = ""

	if token.HasLiteral() {
		str = fmt.Sprintf(
			"[%d:%d] %-8s \"%s\"",
			token.Line,
			token.Column,
			token.TypeName(),
			token.Literal,
		)
	} else {
		str = fmt.Sprintf(
			"[%d:%d] %-8s",
			token.Line,
			token.Column,
			token.TypeName(),
		)
	}

	return str
}

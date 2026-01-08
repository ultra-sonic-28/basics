package lexer

import (
	"basics/internal/token"
	"fmt"
)

func DumpTokens(tokens []token.Token) {
	for _, t := range tokens {
		if t.HasLiteral() {
			fmt.Printf(
				"[%d:%d] %-8s \"%s\"\n",
				t.Line,
				t.Column,
				t.TypeName(),
				t.Literal,
			)
		} else {
			fmt.Printf(
				"[%d:%d] %-8s\n",
				t.Line,
				t.Column,
				t.TypeName(),
			)
		}

	}
}

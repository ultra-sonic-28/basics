package errors

type Kind int

const (
	Lexical Kind = iota
	Syntax
	Semantic
)

func (k Kind) String() string {
	switch k {
	case Lexical:
		return "LEXICAL ERROR"
	case Syntax:
		return "SYNTAX ERROR"
	case Semantic:
		return "SEMANTIC ERROR"
	default:
		return "ERROR"
	}
}

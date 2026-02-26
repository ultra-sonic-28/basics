package parser

const (
	_ int = iota
	LOWEST
	LOGICAL_AND // AND
	EQUALS      // = <> < >
	LESSGREATER // < >
	SUM         // + -
	PRODUCT     // * /
	POWER       // ^
	PREFIX      // -X
)

var precedences = map[string]int{
	"AND": LOGICAL_AND,
	"=":   EQUALS,
	"<>":  EQUALS,
	"<":  LESSGREATER,
	">":  LESSGREATER,
	"<=": LESSGREATER,
	">=": LESSGREATER,
	"+":  SUM,
	"-":  SUM,
	"*":  PRODUCT,
	"/":  PRODUCT,
	"^":  POWER,
}

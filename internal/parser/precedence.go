package parser

const (
	_ int = iota
	LOWEST
	EQUALS      // = <> < >
	LESSGREATER // < >
	SUM         // + -
	PRODUCT     // * /
	POWER       // ^
	PREFIX      // -X
)

var precedences = map[string]int{
	"=":  EQUALS,
	"<>": EQUALS,
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

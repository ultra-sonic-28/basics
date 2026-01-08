package token

type TokenType int

const (
	// Spéciaux
	ILLEGAL TokenType = iota
	EOF
	EOL

	// Spéciaux BASIC
	LINENUM

	// Littéraux
	NUMBER
	STRING
	IDENT

	// Opérateurs
	PLUS
	MINUS
	ASTERISK
	SLASH
	CARET
	EQUAL
	LT
	GT
	LTE
	GTE
	NEQ

	// Délimiteurs
	LPAREN
	RPAREN
	COMMA
	COLON
	SEMICOLON

	// Keywords
	KEYWORD
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var TokenTypeNames = map[TokenType]string{
	// Spéciaux
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	EOL:     "EOL",

	// Spéciaux BASIC
	LINENUM: "LINENUM",

	// Littéraux
	NUMBER: "NUMBER",
	STRING: "STRING",
	IDENT:  "IDENT",

	// Opérateurs
	PLUS:     "+",
	MINUS:    "-",
	ASTERISK: "*",
	SLASH:    "/",
	CARET:    "^",
	EQUAL:    "=",
	LT:       "<",
	GT:       ">",
	LTE:      "<=",
	GTE:      ">=",
	NEQ:      "<>",

	// Délimiteurs
	LPAREN:    "(",
	RPAREN:    ")",
	COMMA:     ",",
	COLON:     ":",
	SEMICOLON: ";",
}

func (t Token) TypeName() string {
	if name, ok := TokenTypeNames[t.Type]; ok {
		return name
	}
	// Par défaut, c'est un KEYWORD du language
	//return fmt.Sprintf("Token(%d)", t.Type)
	return "KEYWORD"
}

func (t Token) HasLiteral() bool {
	return t.Type != EOF && t.Type != EOL
}

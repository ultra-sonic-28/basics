package token

type TokenType int

const (
	// Spéciaux
	ILLEGAL TokenType = iota
	EOF
	EOL
	HASH

	// Spéciaux BASIC
	LINENUM

	// Littéraux
	NUMBER
	STRING
	IDENT
	COMMENT

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
	AND

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
	HASH:    "#",

	// Spéciaux BASIC
	LINENUM: "LINENUM",

	// Littéraux
	NUMBER:  "NUMBER",
	STRING:  "STRING",
	IDENT:   "IDENT",
	COMMENT: "COMMENT",

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
	AND:      "AND",

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

package lexer

import (
	"unicode"

	"basics/internal/logger"
	"basics/internal/token"
)

type Lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune

	line   int
	column int

	expectLineNumber bool
}

func New(input string) *Lexer {
	logger.Info("Instanciate new lexer")
	l := &Lexer{
		input:            []rune(input),
		line:             1,
		expectLineNumber: true,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	tok := token.Token{
		Line:   l.line,
		Column: l.column,
	}

	if tok.Type == token.KEYWORD && tok.Literal == "REM" {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	}

	switch l.ch {
	case 0:
		tok.Type = token.EOF
	case '\n':
		tok.Type = token.EOL
		tok.Literal = "\n"
		l.expectLineNumber = true
		l.readChar()
	case '+':
		tok = l.simpleToken(token.PLUS, "+")
	case '-':
		tok = l.simpleToken(token.MINUS, "-")
	case '*':
		tok = l.simpleToken(token.ASTERISK, "*")
	case '/':
		tok = l.simpleToken(token.SLASH, "/")
	case '^':
		tok = l.simpleToken(token.CARET, "^")
	case '=':
		tok = l.simpleToken(token.EQUAL, "=")
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.LTE
			tok.Literal = string(ch) + string(l.ch)
			l.readChar()
			return tok
		}
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok.Type = token.NEQ
			tok.Literal = string(ch) + string(l.ch)
			l.readChar()
			return tok
		}
		tok = l.simpleToken(token.LT, "<")
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.GTE
			tok.Literal = string(ch) + string(l.ch)
			l.readChar()
			return tok
		}
		tok = l.simpleToken(token.GT, ">")
	case '(':
		tok = l.simpleToken(token.LPAREN, "(")
	case ')':
		tok = l.simpleToken(token.RPAREN, ")")
	case ',':
		tok = l.simpleToken(token.COMMA, ",")
	case ':':
		tok = l.simpleToken(token.COLON, ":")
		l.expectLineNumber = false
		return tok
	case ';':
		tok = l.simpleToken(token.SEMICOLON, ";")
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		if isDigit(l.ch) {
			lit := l.readNumber()

			if l.expectLineNumber {
				tok.Type = token.LINENUM
				l.expectLineNumber = false
			} else {
				tok.Type = token.NUMBER
			}

			tok.Literal = lit
			return tok
		}
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			tok.Literal = lit

			if Keywords[lit] {
				tok.Type = token.KEYWORD

				// ✅ REM : ignorer le reste de la ligne
				if lit == "REM" {
					for l.ch != '\n' && l.ch != 0 {
						l.readChar()
					}
				}

			} else {
				tok.Type = token.IDENT
			}

			return tok
		}

		tok.Type = token.ILLEGAL
		tok.Literal = string(l.ch)
	}

	//l.readChar()
	return tok
}

func (l *Lexer) simpleToken(t token.TokenType, lit string) token.Token {
	tok := token.Token{
		Type:    t,
		Literal: lit,
		Line:    l.line,
		Column:  l.column,
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '$' || l.ch == '%' {
		l.readChar()
	}
	return string(l.input[start:l.position])
}

func (l *Lexer) readNumber() string {
	start := l.position
	dotSeen := false
	for isDigit(l.ch) || (l.ch == '.' && !dotSeen) {
		if l.ch == '.' {
			dotSeen = true
		}
		l.readChar()
	}
	return string(l.input[start:l.position])
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening "
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	out := string(l.input[start:l.position])
	l.readChar() // skip closing "
	return out
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

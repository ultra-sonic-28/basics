package parser

import (
	"fmt"
	"strconv"

	"basics/internal/errors"
	"basics/internal/token"
)

type Parser struct {
	tokens   []token.Token
	pos      int
	curr     token.Token
	peek     token.Token
	errors   []*errors.Error
	forStack []*ForStmt
}

func New(tokens []token.Token) *Parser {
	p := &Parser{tokens: tokens}
	p.curr = tokens[0]
	if len(tokens) > 1 {
		p.peek = tokens[1]
	}
	return p
}

func (p *Parser) next() {
	p.pos++
	p.curr = p.peek
	if p.pos+1 < len(p.tokens) {
		p.peek = p.tokens[p.pos+1]
	}
}

func (p *Parser) ParseProgram() (*Program, []*errors.Error) {
	prog := &Program{}
	seen := make(map[int]bool)

	for p.curr.Type != token.EOF {

		// ✅ ignorer les fins de ligne
		if p.curr.Type == token.EOL {
			p.next()
			continue
		}

		if p.curr.Type != token.LINENUM {
			p.syntaxError("EXPECTED LINE NUMBER")
			p.next()
			continue
		}

		line := p.parseLine()
		if line == nil {
			continue
		}

		if seen[line.Number] {
			p.errors = append(p.errors,
				errors.NewSemantic(
					line.Number,
					"DUPLICATE LINE NUMBER",
				),
			)
		}
		seen[line.Number] = true

		prog.Lines = append(prog.Lines, line)
	}

	// FOR non fermés → erreur
	for _, f := range p.forStack {
		p.errors = append(p.errors,
			errors.NewSyntax(
				f.LineNum,
				f.Column,
				"FOR",
				"MISSING NEXT",
			),
		)
	}

	return prog, p.errors
}

func (p *Parser) parseLine() *Line {
	if p.curr.Type != token.LINENUM {
		p.syntaxError("EXPECTED LINE NUMBER")
		return nil
	}

	num, _ := strconv.Atoi(p.curr.Literal)
	line := &Line{Number: num}

	p.next() // LINENUM

	for {
		stmt := p.parseStatement(num)
		line.Stmts = append(line.Stmts, stmt)

		if p.curr.Type != token.COLON {
			break
		}
		p.next() // :
	}

	if p.curr.Type == token.EOL {
		p.next()
	}

	return line
}

func (p *Parser) parseStatement(lineNum int) Statement {
	if p.curr.Type == token.KEYWORD {
		switch p.curr.Literal {

		case "PRINT":
			return p.parsePrint()

		case "FOR":
			return p.parseFor(lineNum)

		case "NEXT":
			return p.parseNext(lineNum)

		case "LET":
			// LET est optionnel, on le consomme systématiquement
			p.next()
			return p.parseLet()

		case "REM":
			// REM = instruction vide
			p.next()
			return nil

		case "HTAB":
			p.next() // consommer HTAB

			expr := p.parseExpression(LOWEST)
			if expr == nil {
				p.syntaxError("EXPECTED EXPRESSION AFTER HTAB")
				return nil
			}

			return &HTabStmt{
				Expr: expr,
			}

		case "VTAB":
			p.next() // consommer VTAB

			expr := p.parseExpression(LOWEST)
			if expr == nil {
				p.syntaxError("EXPECTED EXPRESSION AFTER VTAB")
				return nil
			}

			return &VTabStmt{
				Expr: expr,
			}
		default:
			p.syntaxError("UNKNOWN KEYWORD")
			p.next()
			return nil
		}
	}

	if p.curr.Type == token.IDENT {
		// IDENT doit être suivi de '='
		if p.peek.Literal != "=" {
			p.syntaxError("UNKNOWN STATEMENT")
			p.next()
			return nil
		}

		return p.parseLet() // LET implicite (Applesoft)
	}

	// --- Token inattendu ---
	p.syntaxError("SYNTAX ERROR")
	p.next()
	return nil
}

func (p *Parser) parsePrint() Statement {
	p.next() // PRINT

	var exprs []Expression
	var separators []rune

	for {
		expr := p.parseExpression(LOWEST)
		if expr == nil {
			break
		}
		exprs = append(exprs, expr)

		// séparateurs PRINT
		if p.curr.Type == token.SEMICOLON || p.curr.Type == token.COMMA {
			separators = append(separators, rune(p.curr.Literal[0]))
			p.next()
			continue
		}

		break
	}

	return &PrintStmt{
		Exprs:      exprs,
		Separators: separators,
	}
}

func (p *Parser) parseLet() Statement {
	name := p.curr.Literal
	p.expect(token.IDENT)

	if !p.expectLiteral("=") {
		return nil
	}

	value := p.parseExpression(LOWEST)
	if value == nil {
		return nil
	}

	return &LetStmt{
		Name:  name,
		Value: value,
	}
}

func (p *Parser) parseFor(lineNum int) Statement {
	// position du mot-clé FOR
	col := p.curr.Column

	p.next() // FOR

	varName := p.curr.Literal
	p.expect(token.IDENT)

	p.expectLiteral("=")

	start := p.parseExpression(LOWEST)

	if !p.expectKeyword("TO") {
		return nil
	}

	end := p.parseExpression(LOWEST)
	if end == nil {
		return nil
	}

	var step Expression = &NumberLiteral{
		Value:  1,
		Line:   p.curr.Line,
		Column: p.curr.Column,
		Token:  p.curr.Literal,
	} // valeur par défaut STEP = 1

	// STEP facultatif
	if p.curr.Type == token.KEYWORD && p.curr.Literal == "STEP" {
		p.next() // STEP
		step = p.parseExpression(LOWEST)
		if step == nil {
			return nil
		}
	}

	stmt := &ForStmt{
		Var:     varName,
		Start:   start,
		End:     end,
		Step:    step,
		LineNum: lineNum,
		Column:  col,
	}

	// Empilement FOR
	p.forStack = append(p.forStack, stmt)

	return stmt
}

func (p *Parser) parseNext(lineNum int) Statement {
	p.next() // NEXT

	if len(p.forStack) == 0 {
		p.syntaxError("NEXT WITHOUT FOR")
		return nil
	}

	name := p.curr.Literal
	p.expect(token.IDENT)

	// récupérer le FOR courant
	top := p.forStack[len(p.forStack)-1]

	if top.Var != name {
		p.syntaxError(
			fmt.Sprintf("MISMATCHED NEXT VARIABLE, expected '%s'", top.Var),
		)
		return nil
	}

	// dépiler
	p.forStack = p.forStack[:len(p.forStack)-1]

	return &NextStmt{
		Var:        name,
		ForLineNum: top.LineNum, // ligne BASIC du FOR correspondant
	}
}

func (p *Parser) parseExpression(precedence int) Expression {
	var left Expression

	// --- prefix ---
	switch p.curr.Type {

	case token.NUMBER:
		val, err := strconv.ParseFloat(p.curr.Literal, 64)
		if err != nil {
			p.syntaxError("INVALID NUMBER")
			return nil
		}
		left = &NumberLiteral{
			Value:  val,
			Line:   p.curr.Line,
			Column: p.curr.Column,
			Token:  p.curr.Literal,
		}
		p.next()

	case token.IDENT:
		left = &Identifier{
			Name:   p.curr.Literal,
			Line:   p.curr.Line,
			Column: p.curr.Column,
			Token:  p.curr.Literal,
		}
		p.next()

	case token.EQUAL, token.MINUS:
		opTok := p.curr
		p.next()
		right := p.parseExpression(PREFIX)
		left = &PrefixExpr{
			Op:     opTok.Literal,
			Right:  right,
			Line:   opTok.Line,
			Column: opTok.Column,
			Token:  opTok.Literal,
		}

	case token.STRING:
		left = &StringLiteral{
			Value:  p.curr.Literal,
			Line:   p.curr.Line,
			Column: p.curr.Column,
			Token:  p.curr.Literal,
		}
		p.next()

	default:
		p.syntaxError("INVALID EXPRESSION")
		return nil
	}

	// --- infix ---
	for p.curr.Type != token.EOL &&
		p.curr.Type != token.COLON &&
		precedence < p.currPrecedence() {

		opTok := p.curr
		prec := p.currPrecedence()
		p.next()

		right := p.parseExpression(prec)
		left = &InfixExpr{
			Left:   left,
			Op:     opTok.Literal,
			Right:  right,
			Line:   opTok.Line,
			Column: opTok.Column,
			Token:  opTok.Literal,
		}
	}

	return left
}

func (p *Parser) expect(t token.TokenType) bool {
	if p.curr.Type != t {
		p.syntaxError(
			fmt.Sprintf("EXPECTED %s", p.curr.TypeName()),
		)
		return false
	}
	p.next()
	return true
}

func (p *Parser) expectKeyword(kw string) bool {
	if p.curr.Type != token.KEYWORD || p.curr.Literal != kw {
		p.syntaxError(
			fmt.Sprintf("EXPECTED KEYWORD %s", kw),
		)
		return false
	}
	p.next()
	return true
}

func (p *Parser) expectLiteral(lit string) bool {
	if p.curr.Literal != lit {
		p.syntaxError(
			fmt.Sprintf("EXPECTED '%s'", lit),
		)
		return false
	}
	p.next()
	return true
}

func (p *Parser) currPrecedence() int {
	if prec, ok := precedences[p.curr.Literal]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peek.Literal]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) syntaxError(msg string) {
	err := errors.NewSyntax(
		p.curr.Line,
		p.curr.Column,
		p.curr.Literal,
		msg,
	)
	p.errors = append(p.errors, err)
}

func (p *Parser) skipToNextStatement() {
	for p.curr.Type != token.COLON &&
		p.curr.Type != token.EOL &&
		p.curr.Type != token.EOF {
		p.next()
	}
}

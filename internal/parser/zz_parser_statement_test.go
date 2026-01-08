package parser

import (
	"testing"

	"basics/internal/token"
	"basics/testutils"
)

func TestParser_parseStatement(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []token.Token
		check   func(t *testing.T, stmt Statement, p *Parser)
		wantErr bool
	}{
		{
			name: "PRINT single number",
			tokens: []token.Token{
				{Type: token.KEYWORD, Literal: "PRINT", Line: 1, Column: 1},
				{Type: token.NUMBER, Literal: "42", Line: 1, Column: 7},
			},
			check: func(t *testing.T, stmt Statement, p *Parser) {
				printStmt, ok := stmt.(*PrintStmt)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", len(printStmt.Exprs), 1)
				num, ok := printStmt.Exprs[0].(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", num.Value, 42.0)
			},
		},
		{
			name: "PRINT with multiple expr and separators",
			tokens: []token.Token{
				{Type: token.KEYWORD, Literal: "PRINT", Line: 1, Column: 1},
				{Type: token.NUMBER, Literal: "1", Line: 1, Column: 7},
				{Type: token.SEMICOLON, Literal: ";", Line: 1, Column: 8},
				{Type: token.NUMBER, Literal: "2", Line: 1, Column: 9},
				{Type: token.COMMA, Literal: ",", Line: 1, Column: 10},
				{Type: token.IDENT, Literal: "X", Line: 1, Column: 11},
			},
			check: func(t *testing.T, stmt Statement, p *Parser) {
				printStmt, ok := stmt.(*PrintStmt)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", len(printStmt.Exprs), 3)
				testutils.Equal(t, "", len(printStmt.Separators), 2)
				testutils.Equal(t, "", printStmt.Separators[0], ';')
				testutils.Equal(t, "", printStmt.Separators[1], ',')
			},
		},
		{
			name: "LET statement",
			tokens: []token.Token{
				{Type: token.KEYWORD, Literal: "LET", Line: 1, Column: 1},
				{Type: token.IDENT, Literal: "A", Line: 1, Column: 5},
				{Type: token.EQUAL, Literal: "=", Line: 1, Column: 7},
				{Type: token.NUMBER, Literal: "10", Line: 1, Column: 9},
			},
			check: func(t *testing.T, stmt Statement, p *Parser) {
				letStmt, ok := stmt.(*LetStmt)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", letStmt.Name, "A")
				num, ok := letStmt.Value.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", num.Value, 10.0)
			},
		},
		{
			name: "implicit LET (IDENT=...)",
			tokens: []token.Token{
				{Type: token.IDENT, Literal: "B", Line: 1, Column: 1},
				{Type: token.EQUAL, Literal: "=", Line: 1, Column: 2},
				{Type: token.NUMBER, Literal: "20", Line: 1, Column: 3},
			},
			check: func(t *testing.T, stmt Statement, p *Parser) {
				letStmt, ok := stmt.(*LetStmt)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", letStmt.Name, "B")
				num, ok := letStmt.Value.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", num.Value, 20.0)
			},
		},
		{
			name: "FOR ... TO ... STEP ...",
			tokens: []token.Token{
				{Type: token.KEYWORD, Literal: "FOR", Line: 1, Column: 1},
				{Type: token.IDENT, Literal: "I", Line: 1, Column: 5},
				{Type: token.EQUAL, Literal: "=", Line: 1, Column: 7},
				{Type: token.NUMBER, Literal: "1", Line: 1, Column: 9},
				{Type: token.KEYWORD, Literal: "TO", Line: 1, Column: 11},
				{Type: token.NUMBER, Literal: "10", Line: 1, Column: 14},
				{Type: token.KEYWORD, Literal: "STEP", Line: 1, Column: 17},
				{Type: token.NUMBER, Literal: "2", Line: 1, Column: 22},
			},
			check: func(t *testing.T, stmt Statement, p *Parser) {
				forStmt, ok := stmt.(*ForStmt)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", forStmt.Var, "I")
				numStart, ok := forStmt.Start.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", numStart.Value, 1.0)
				numEnd, ok := forStmt.End.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", numEnd.Value, 10.0)
				numStep, ok := forStmt.Step.(*NumberLiteral)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", numStep.Value, 2.0)
				testutils.Equal(t, "", len(p.forStack), 1)
			},
		},
		{
			name: "NEXT matching FOR",
			tokens: []token.Token{
				{Type: token.KEYWORD, Literal: "FOR", Line: 1, Column: 1},
				{Type: token.IDENT, Literal: "I", Line: 1, Column: 5},
				{Type: token.EQUAL, Literal: "=", Line: 1, Column: 7},
				{Type: token.NUMBER, Literal: "1", Line: 1, Column: 9},
				{Type: token.KEYWORD, Literal: "TO", Line: 1, Column: 11},
				{Type: token.NUMBER, Literal: "5", Line: 1, Column: 14},
				{Type: token.KEYWORD, Literal: "NEXT", Line: 1, Column: 16},
				{Type: token.IDENT, Literal: "I", Line: 1, Column: 21},
			},
			check: func(t *testing.T, stmt Statement, p *Parser) {
				forStmt, ok := stmt.(*ForStmt)
				testutils.Equal(t, "", ok, true)
				nextStmt := p.parseStatement(LOWEST)
				nextStmtCast, ok := nextStmt.(*NextStmt)
				testutils.Equal(t, "", ok, true)
				testutils.Equal(t, "", nextStmtCast.Var, "I")
				testutils.Equal(t, "", nextStmtCast.ForLineNum, forStmt.LineNum)
				testutils.Equal(t, "", len(p.forStack), 0)
			},
		},
		{
			name: "REM statement",
			tokens: []token.Token{
				{Type: token.KEYWORD, Literal: "REM", Line: 1, Column: 1},
			},
			check: func(t *testing.T, stmt Statement, p *Parser) {
				testutils.Equal(t, "", stmt, nil)
			},
		},
		{
			name: "unknown keyword",
			tokens: []token.Token{
				{Type: token.KEYWORD, Literal: "FOO", Line: 1, Column: 1},
			},
			wantErr: true,
		},
		{
			name: "NEXT without FOR",
			tokens: []token.Token{
				{Type: token.KEYWORD, Literal: "NEXT", Line: 1, Column: 1},
				{Type: token.IDENT, Literal: "X", Line: 1, Column: 6},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.tokens)
			stmt := p.parseStatement(LOWEST)
			if tt.wantErr {
				testutils.True(t, "expected error", len(p.errors) > 0)
				return
			}
			testutils.False(t, "unexpected error", len(p.errors) > 0)
			if tt.check != nil {
				tt.check(t, stmt, p)
			}
		})
	}
}

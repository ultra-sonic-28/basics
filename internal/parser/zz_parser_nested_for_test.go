package parser

import (
	"testing"

	"basics/internal/lexer"
	"basics/testutils"
)

func TestParseNestedFor(t *testing.T) {
	source := `
10 FOR I = 1 TO 10
20 FOR J = 1 TO 5
30 PRINT I, J
40 NEXT J
50 NEXT I
`

	// 🔹 Lexer
	tokens := lexer.Lex(source)

	// 🔹 Parser
	p := New(tokens)
	prog, errs := p.ParseProgram()

	// 🔹 Aucune erreur attendue
	testutils.Equal(t, "should have no errors", len(errs), 0)

	// 🔹 Vérification du nombre de lignes
	testutils.Equal(t, "number of lines", len(prog.Lines), 5)

	// --- Line 10: FOR I = 1 TO 10 ---
	line10 := prog.Lines[0]
	testutils.Equal(t, "line 10 number", line10.Number, 10)

	stmt10, ok := line10.Stmts[0].(*ForStmt)
	testutils.Equal(t, "line 10 type", ok, true)
	testutils.Equal(t, "FOR I var", stmt10.Var, "I")

	//startNum, _, _ := stmt10.Start.(*NumberLiteral).Pos()
	testutils.Equal(t, "FOR I start value", stmt10.Start.(*NumberLiteral).Value, 1)
	testutils.Equal(t, "FOR I end value", stmt10.End.(*NumberLiteral).Value, 10)
	testutils.Equal(t, "FOR I step value", stmt10.Step.(*NumberLiteral).Value, 1)

	// --- Line 20: FOR J = 1 TO 5 ---
	line20 := prog.Lines[1]
	stmt20, ok := line20.Stmts[0].(*ForStmt)
	testutils.Equal(t, "line 20 type", ok, true)
	testutils.Equal(t, "FOR J var", stmt20.Var, "J")
	testutils.Equal(t, "FOR J start value", stmt20.Start.(*NumberLiteral).Value, 1)
	testutils.Equal(t, "FOR J end value", stmt20.End.(*NumberLiteral).Value, 5)
	testutils.Equal(t, "FOR J step value", stmt20.Step.(*NumberLiteral).Value, 1)

	// --- Line 30: PRINT I, J ---
	line30 := prog.Lines[2]
	stmt30, ok := line30.Stmts[0].(*PrintStmt)
	testutils.Equal(t, "line 30 type", ok, true)
	testutils.Equal(t, "PRINT expr count", len(stmt30.Exprs), 2)
	testutils.Equal(t, "PRINT first expr", stmt30.Exprs[0].(*Identifier).Name, "I")
	testutils.Equal(t, "PRINT second expr", stmt30.Exprs[1].(*Identifier).Name, "J")

	// --- Line 40: NEXT J ---
	line40 := prog.Lines[3]
	stmt40, ok := line40.Stmts[0].(*NextStmt)
	testutils.Equal(t, "line 40 type", ok, true)
	testutils.Equal(t, "NEXT J var", stmt40.Var, "J")
	testutils.Equal(t, "NEXT J matches FOR line", stmt40.ForLineNum, 20)

	// --- Line 50: NEXT I ---
	line50 := prog.Lines[4]
	stmt50, ok := line50.Stmts[0].(*NextStmt)
	testutils.Equal(t, "line 50 type", ok, true)
	testutils.Equal(t, "NEXT I var", stmt50.Var, "I")
	testutils.Equal(t, "NEXT I matches FOR line", stmt50.ForLineNum, 10)
}

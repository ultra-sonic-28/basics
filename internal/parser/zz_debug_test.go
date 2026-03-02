package parser

import (
	"basics/testutils"
	"fmt"
	"testing"
)

// ------------------------
// Tests pour dumpExpr
// ------------------------
func TestDumpExpr(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected string
	}{
		{
			name:     "NumberLiteral",
			expr:     &NumberLiteral{Value: 42, Line: 1, Column: 1, Token: "42"},
			expected: "Number 42\n",
		},
		{
			name:     "StringLiteral",
			expr:     &StringLiteral{Value: "hello", Line: 2, Column: 3, Token: "\"hello\""},
			expected: "String \"hello\"\n",
		},
		{
			name:     "Identifier",
			expr:     &Identifier{Name: "X", Line: 3, Column: 5, Token: "X"},
			expected: "Ident X\n",
		},
		{
			name: "LeftExpr",
			expr: &LeftExpr{
				StrExpr: &StringLiteral{
					Value:  "APPLESOFT",
					Line:   4,
					Column: 2,
					Token:  "APPLESOFT",
				},
				LenExpr: &NumberLiteral{
					Value:  5,
					Line:   4,
					Column: 15,
					Token:  "5",
				},
				Line:   1,
				Column: 1,
				Token:  "LEFT$",
			},
			expected: "LEFT\n  String \"APPLESOFT\"\n  Number 5\n",
		},
		{
			name: "LeftExpr with expression length",
			expr: &LeftExpr{
				StrExpr: &Identifier{Name: "A$"},
				LenExpr: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2},
				},
				Token: "LEFT$",
			},
			expected: "LEFT\n  Ident A$\n  Infix *\n    Ident A\n    Number 2\n",
		},
		{
			name: "RightExpr",
			expr: &RightExpr{
				StrExpr: &StringLiteral{
					Value:  "APPLESOFT",
					Line:   4,
					Column: 2,
					Token:  "APPLESOFT",
				},
				LenExpr: &NumberLiteral{
					Value:  5,
					Line:   4,
					Column: 15,
					Token:  "5",
				},
				Line:   1,
				Column: 1,
				Token:  "RIGHT$",
			},
			expected: "RIGHT\n  String \"APPLESOFT\"\n  Number 5\n",
		},
		{
			name: "RightExpr with expression length",
			expr: &RightExpr{
				StrExpr: &Identifier{Name: "A$"},
				LenExpr: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2},
				},
				Token: "RIGHT$",
			},
			expected: "RIGHT\n  Ident A$\n  Infix *\n    Ident A\n    Number 2\n",
		},
		{
			name: "MidExpr",
			expr: &MidExpr{
				StrExpr: &StringLiteral{
					Value:  "APPLESOFT",
					Line:   4,
					Column: 2,
					Token:  "APPLESOFT",
				},
				Start: &NumberLiteral{
					Value:  5,
					Line:   4,
					Column: 15,
					Token:  "5",
				},
				Line:   1,
				Column: 1,
				Token:  "MID$",
			},
			expected: "MID\n  String \"APPLESOFT\"\n  Number 5\n",
		},
		{
			name: "MidExpr with len",
			expr: &MidExpr{
				StrExpr: &StringLiteral{
					Value:  "APPLESOFT",
					Line:   4,
					Column: 2,
					Token:  "APPLESOFT",
				},
				Start: &NumberLiteral{
					Value:  5,
					Line:   4,
					Column: 15,
					Token:  "5",
				},
				Len: &NumberLiteral{
					Value:  1,
					Line:   4,
					Column: 15,
					Token:  "1",
				},
				Line:   1,
				Column: 1,
				Token:  "MID$",
			},
			expected: "MID\n  String \"APPLESOFT\"\n  Number 5\n  Number 1\n",
		},
		{
			name: "MidExpr with start expression",
			expr: &MidExpr{
				StrExpr: &Identifier{Name: "A$"},
				Start: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2},
				},
				Token: "MID$",
			},
			expected: "MID\n  Ident A$\n  Infix *\n    Ident A\n    Number 2\n",
		},
		{
			name: "MidExpr with start and len expression",
			expr: &MidExpr{
				StrExpr: &Identifier{Name: "A$"},
				Start: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "*",
					Right: &NumberLiteral{Value: 3},
				},
				Len: &InfixExpr{
					Left:  &Identifier{Name: "B"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2},
				},
				Token: "MID$",
			},
			expected: "MID\n  Ident A$\n  Infix *\n    Ident A\n    Number 3\n  Infix *\n    Ident B\n    Number 2\n",
		},
		{
			name: "LenExpr",
			expr: &LenExpr{
				Expr: &StringLiteral{
					Value:  "APPLESOFT",
					Line:   4,
					Column: 2,
					Token:  "APPLESOFT",
				},
				Line:   1,
				Column: 1,
				Token:  "LEN",
			},
			expected: "LEN\n  String \"APPLESOFT\"\n",
		},
		{
			name: "LenExpr with expression",
			expr: &LenExpr{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A$"},
					Op:    "+",
					Right: &Identifier{Name: "B$"},
				},
				Token: "LEN",
			},
			expected: "LEN\n  Infix +\n    Ident A$\n    Ident B$\n",
		},
		{
			name: "StrExpr",
			expr: &StrExpr{
				Expr: &NumberLiteral{
					Value:  5,
					Line:   4,
					Column: 15,
					Token:  "5",
				},
				Line:   1,
				Column: 1,
				Token:  "STR$",
			},
			expected: "STR$\n  Number 5\n",
		},
		{
			name: "StrExpr with expression",
			expr: &StrExpr{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "+",
					Right: &Identifier{Name: "B"},
				},
				Line:   1,
				Column: 1,
				Token:  "STR$",
			},
			expected: "STR$\n  Infix +\n    Ident A\n    Ident B\n",
		},
		{
			name: "ChrExpr",
			expr: &ChrExpr{
				Expr: &NumberLiteral{
					Value:  33,
					Line:   4,
					Column: 15,
					Token:  "33",
				},
				Line:   1,
				Column: 1,
				Token:  "CHR$",
			},
			expected: "CHR$\n  Number 33\n",
		},
		{
			name: "ChrExpr with expression",
			expr: &ChrExpr{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "+",
					Right: &Identifier{Name: "B"},
				},
				Line:   1,
				Column: 1,
				Token:  "CHR$",
			},
			expected: "CHR$\n  Infix +\n    Ident A\n    Ident B\n",
		},
		{
			name: "ValExpr",
			expr: &ValExpr{
				Expr: &StringLiteral{
					Value:  "123456",
					Line:   4,
					Column: 2,
					Token:  "123456",
				},
				Line:   1,
				Column: 1,
				Token:  "VAL",
			},
			expected: "VAL\n  String \"123456\"\n",
		},
		{
			name: "ValExpr with expression",
			expr: &ValExpr{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A$"},
					Op:    "+",
					Right: &Identifier{Name: "B$"},
				},
				Token: "VAL",
			},
			expected: "VAL\n  Infix +\n    Ident A$\n    Ident B$\n",
		},
		{
			name: "AscExpr",
			expr: &AscExpr{
				Expr: &StringLiteral{
					Value:  "A",
					Line:   4,
					Column: 2,
					Token:  "A",
				},
				Line:   1,
				Column: 1,
				Token:  "ASC",
			},
			expected: "ASC\n  String \"A\"\n",
		},
		{
			name: "AscExpr with expression",
			expr: &AscExpr{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A$"},
					Op:    "+",
					Right: &Identifier{Name: "B$"},
				},
				Token: "ASC",
			},
			expected: "ASC\n  Infix +\n    Ident A$\n    Ident B$\n",
		},
		{
			name: "TabExpr",
			expr: &TabExpr{
				Expr: &NumberLiteral{
					Value:  16,
					Line:   4,
					Column: 2,
					Token:  "16",
				},
				Line:   1,
				Column: 1,
				Token:  "16",
			},
			expected: "TAB\n  Number 16\n",
		},
		{
			name: "SpcExpr",
			expr: &SpcExpr{
				Expr: &NumberLiteral{
					Value:  16,
					Line:   4,
					Column: 2,
					Token:  "16",
				},
				Line:   1,
				Column: 1,
				Token:  "16",
			},
			expected: "SPC\n  Number 16\n",
		},
		{
			name: "SpcExpr with expression",
			expr: &SpcExpr{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "+",
					Right: &Identifier{Name: "B"},
				},
				Token: "SPC",
			},
			expected: "SPC\n  Infix +\n    Ident A\n    Ident B\n",
		},
		{
			name: "RndExpr",
			expr: &RndExpr{
				Expr: &NumberLiteral{
					Value:  16,
					Line:   4,
					Column: 2,
					Token:  "16",
				},
				Line:   1,
				Column: 1,
				Token:  "16",
			},
			expected: "RND\n  Number 16\n",
		},
		{
			name: "RndExpr with expression",
			expr: &RndExpr{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "+",
					Right: &Identifier{Name: "B"},
				},
				Token: "RND",
			},
			expected: "RND\n  Infix +\n    Ident A\n    Ident B\n",
		},
		{
			name: "SinExpr",
			expr: &SinExpr{
				Expr: &NumberLiteral{
					Value:  16,
					Line:   4,
					Column: 2,
					Token:  "16",
				},
				Line:   1,
				Column: 1,
				Token:  "16",
			},
			expected: "SIN\n  Number 16\n",
		},
		{
			name: "SinExpr with expression",
			expr: &SinExpr{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "+",
					Right: &Identifier{Name: "B"},
				},
				Token: "SIN",
			},
			expected: "SIN\n  Infix +\n    Ident A\n    Ident B\n",
		},
		{
			name: "CosExpr",
			expr: &CosExpr{
				Expr:  &NumberLiteral{Value: 16},
				Token: "COS",
			},
			expected: "COS\n  Number 16\n",
		},
		{
			name: "CosExpr with expression",
			expr: &CosExpr{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "+",
					Right: &Identifier{Name: "B"},
				},
				Token: "COS",
			},
			expected: "COS\n  Infix +\n    Ident A\n    Ident B\n",
		},
		{
			name: "TanExpr",
			expr: &TanExpr{
				Expr:  &NumberLiteral{Value: 16},
				Token: "TAN",
			},
			expected: "TAN\n  Number 16\n",
		},
		{
			name: "TanExpr with expression",
			expr: &TanExpr{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A"},
					Op:    "+",
					Right: &Identifier{Name: "B"},
				},
				Token: "TAN",
			},
			expected: "TAN\n  Infix +\n    Ident A\n    Ident B\n",
		},
		{
			name: "SqrExpr",
			expr: &SqrExpr{
				Expr: &NumberLiteral{
					Value:  16,
					Line:   4,
					Column: 2,
					Token:  "16",
				},
				Line:   1,
				Column: 1,
				Token:  "16",
			},
			expected: "SQR\n  Number 16\n",
		},
		{
			name: "AbsExpr",
			expr: &AbsExpr{
				Expr: &NumberLiteral{
					Value:  16,
					Line:   4,
					Column: 2,
					Token:  "16",
				},
				Line:   1,
				Column: 1,
				Token:  "16",
			},
			expected: "ABS\n  Number 16\n",
		},
		{
			name: "SgnExpr",
			expr: &SgnExpr{
				Expr: &NumberLiteral{
					Value:  16,
					Line:   4,
					Column: 2,
					Token:  "16",
				},
				Line:   1,
				Column: 1,
				Token:  "4",
			},
			expected: "SGN\n  Number 16\n",
		},
		{
			name: "IntExpr",
			expr: &IntExpr{
				Expr: &NumberLiteral{
					Value:  16,
					Line:   4,
					Column: 2,
					Token:  "16",
				},
				Line:   1,
				Column: 1,
				Token:  "4",
			},
			expected: "INT\n  Number 16\n",
		},
		{
			name: "PrefixExpr",
			expr: &PrefixExpr{
				Op:    "-",
				Right: &NumberLiteral{Value: 5, Line: 4, Column: 2, Token: "5"},
			},
			expected: "Prefix -\n  Number 5\n",
		},
		{
			name: "InfixExpr",
			expr: &InfixExpr{
				Left:  &NumberLiteral{Value: 2, Line: 5, Column: 1, Token: "2"},
				Op:    "+",
				Right: &NumberLiteral{Value: 3, Line: 5, Column: 3, Token: "3"},
			},
			expected: "Infix +\n  Number 2\n  Number 3\n",
		},
		{
			name: "IndexExpr single identifier",
			expr: &IndexExpr{
				Name: "A",
				Indices: []Expression{
					&Identifier{Name: "I"},
				},
			},
			expected: "A(I)\n",
		},
		{
			name: "IndexExpr single number",
			expr: &IndexExpr{
				Name: "A",
				Indices: []Expression{
					&NumberLiteral{Value: 3},
				},
			},
			expected: "A(3)\n",
		},
		{
			name: "IndexExpr infix expression",
			expr: &IndexExpr{
				Name: "A",
				Indices: []Expression{
					&InfixExpr{
						Left:  &Identifier{Name: "I"},
						Op:    "+",
						Right: &NumberLiteral{Value: 1},
					},
				},
			},
			expected: "A(I + 1)\n",
		},
		{
			name: "IndexExpr multiple identifiers",
			expr: &IndexExpr{
				Name: "A",
				Indices: []Expression{
					&Identifier{Name: "I"},
					&Identifier{Name: "J"},
				},
			},
			expected: "A(I,J)\n",
		},
		{
			name: "IndexExpr mixed expressions",
			expr: &IndexExpr{
				Name: "A",
				Indices: []Expression{
					&Identifier{Name: "I"},
					&InfixExpr{
						Left:  &Identifier{Name: "J"},
						Op:    "+",
						Right: &NumberLiteral{Value: 2},
					},
				},
			},
			expected: "A(I,J + 2)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := testutils.CaptureStdout(t, func() {
				dumpExpr(tt.expr, "", StdoutEmitter)
			})
			msg := fmt.Sprintf("tests[%s] - dumpExpr output mismatch.", tt.name)
			testutils.Equal(t, msg, output, tt.expected)
		})
	}
}

// ------------------------
// Tests pour dumpStatement
// ------------------------
func TestDumpStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     Statement
		expected string
	}{
		{
			name:     "RemStmt",
			stmt:     &RemStmt{},
			expected: "REM\n",
		},
		{
			name:     "HomeStmt",
			stmt:     &HomeStmt{},
			expected: "HOME\n",
		},
		{
			name:     "PrintStmt",
			stmt:     &PrintStmt{Exprs: []Expression{&NumberLiteral{Value: 42, Line: 1, Column: 1, Token: "42"}}},
			expected: "PRINT\n  EXPR 0:\n    Number 42\n",
		},
		{
			name:     "LetStmt",
			stmt:     &LetStmt{Name: "X", Value: &NumberLiteral{Value: 10, Line: 2, Column: 1, Token: "10"}},
			expected: "LET X\n  Number 10\n",
		},
		{
			name:     "EndStmt",
			stmt:     &EndStmt{},
			expected: "END\n",
		},
		{
			name:     "ClearStmt",
			stmt:     &ClearStmt{},
			expected: "CLEAR\n",
		},
		{
			name:     "GrStmt",
			stmt:     &GrStmt{},
			expected: "GR\n",
		},
		{
			name:     "TextStmt",
			stmt:     &TextStmt{},
			expected: "TEXT\n",
		},
		{
			name: "WaitStmt",
			stmt: &WaitStmt{
				Expr: &NumberLiteral{Value: 1000},
			},
			expected: "WAIT\n  Number 1000\n",
		},
		{
			name: "ColorStmt with literal",
			stmt: &ColorStmt{
				Expr: &NumberLiteral{Value: 15},
			},
			expected: "COLOR\n  Number 15\n",
		},
		{
			name: "ColorStmt with expression",
			stmt: &ColorStmt{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "C"},
					Op:    "+",
					Right: &NumberLiteral{Value: 1},
				},
			},
			expected: "COLOR\n  Infix +\n    Ident C\n    Number 1\n",
		},
		{
			name: "PlotStmt with literals",
			stmt: &PlotStmt{
				X: &NumberLiteral{Value: 10},
				Y: &NumberLiteral{Value: 20},
			},
			expected: "PLOT\n  Number 10\n  Number 20\n",
		},
		{
			name: "PlotStmt with expressions",
			stmt: &PlotStmt{
				X: &InfixExpr{
					Left:  &Identifier{Name: "X"},
					Op:    "+",
					Right: &NumberLiteral{Value: 1},
				},
				Y: &InfixExpr{
					Left:  &Identifier{Name: "Y"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2},
				},
			},
			expected: "PLOT\n  Infix +\n    Ident X\n    Number 1\n  Infix *\n    Ident Y\n    Number 2\n",
		},
		{
			name: "ForStmt without Step",
			stmt: &ForStmt{
				Var:     "I",
				Start:   &NumberLiteral{Value: 1, Line: 3, Column: 1, Token: "1"},
				End:     &NumberLiteral{Value: 10, Line: 3, Column: 3, Token: "10"},
				LineNum: 100,
				Step:    nil,
			},
			expected: "FOR I (Line 100)\n  FROM:\n    Number 1\n  TO:\n    Number 10\n",
		},
		{
			name: "ForStmt with Step",
			stmt: &ForStmt{
				Var:     "J",
				Start:   &NumberLiteral{Value: 0, Line: 4, Column: 1, Token: "0"},
				End:     &NumberLiteral{Value: 5, Line: 4, Column: 3, Token: "5"},
				LineNum: 200,
				Step:    &NumberLiteral{Value: 2, Line: 4, Column: 2, Token: "2"},
			},
			expected: "FOR J (Line 200)\n  FROM:\n    Number 0\n  TO:\n    Number 5\n  STEP:\n    Number 2\n",
		},
		{
			name:     "NextStmt with ForLineNum",
			stmt:     &NextStmt{Var: "I", ForLineNum: 100},
			expected: "NEXT I (FOR Line 100)\n",
		},
		{
			name:     "NextStmt without ForLineNum",
			stmt:     &NextStmt{Var: "J"},
			expected: "NEXT J\n",
		},
		{
			name: "GotoStmt",
			stmt: &GotoStmt{
				Expr: &NumberLiteral{Value: 40, Line: 3, Column: 1, Token: "40"},
			},
			expected: "GOTO\n  Number 40\n",
		},
		{
			name: "GotoStmt with expression",
			stmt: &GotoStmt{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "JUMP", Token: "JUMP"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2, Token: "2"},
				},
			},
			expected: "" +
				"GOTO\n" +
				"  Infix *\n" +
				"    Ident JUMP\n" +
				"    Number 2\n",
		},
		{
			name: "GotoStmt with parenthesized expression",
			stmt: &GotoStmt{
				Expr: &InfixExpr{
					Left: &InfixExpr{
						Left:  &Identifier{Name: "A", Token: "A"},
						Op:    "+",
						Right: &Identifier{Name: "B", Token: "B"},
					},
					Op: "*",
					Right: &NumberLiteral{
						Value: 2,
						Token: "2",
					},
				},
			},
			expected: "" +
				"GOTO\n" +
				"  Infix *\n" +
				"    Infix +\n" +
				"      Ident A\n" +
				"      Ident B\n" +
				"    Number 2\n",
		},
		{
			name: "GosubStmt",
			stmt: &GosubStmt{
				Expr: &NumberLiteral{Value: 40, Line: 3, Column: 1, Token: "40"},
			},
			expected: "GOSUB\n  Number 40\n",
		},
		{
			name: "GosubStmt with expression",
			stmt: &GosubStmt{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "JUMP", Token: "JUMP"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2, Token: "2"},
				},
			},
			expected: "" +
				"GOSUB\n" +
				"  Infix *\n" +
				"    Ident JUMP\n" +
				"    Number 2\n",
		},
		{
			name: "GosubStmt with parenthesized expression",
			stmt: &GosubStmt{
				Expr: &InfixExpr{
					Left: &InfixExpr{
						Left:  &Identifier{Name: "A", Token: "A"},
						Op:    "+",
						Right: &Identifier{Name: "B", Token: "B"},
					},
					Op: "*",
					Right: &NumberLiteral{
						Value: 2,
						Token: "2",
					},
				},
			},
			expected: "" +
				"GOSUB\n" +
				"  Infix *\n" +
				"    Infix +\n" +
				"      Ident A\n" +
				"      Ident B\n" +
				"    Number 2\n",
		},
		{
			name:     "ReturnStmt",
			stmt:     &ReturnStmt{},
			expected: "RETURN\n",
		},
		{
			name: "HTabStmt",
			stmt: &HTabStmt{
				Expr: &NumberLiteral{Value: 10, Line: 1, Column: 1, Token: "10"},
			},
			expected: "HTAB\n  Number 10\n",
		},
		{
			name: "HTabStmt with expression",
			stmt: &HTabStmt{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "A", Token: "A"},
					Op:    "*",
					Right: &NumberLiteral{Value: 2, Token: "2"},
				},
			},
			expected: "" +
				"HTAB\n" +
				"  Infix *\n" +
				"    Ident A\n" +
				"    Number 2\n",
		},
		{
			name: "HTabStmt with parenthesized expression",
			stmt: &HTabStmt{
				Expr: &InfixExpr{
					Left: &InfixExpr{
						Left:  &Identifier{Name: "A", Token: "A"},
						Op:    "+",
						Right: &Identifier{Name: "B", Token: "B"},
					},
					Op: "*",
					Right: &NumberLiteral{
						Value: 2,
						Token: "2",
					},
				},
			},
			expected: "" +
				"HTAB\n" +
				"  Infix *\n" +
				"    Infix +\n" +
				"      Ident A\n" +
				"      Ident B\n" +
				"    Number 2\n",
		},
		{
			name: "VTabStmt",
			stmt: &VTabStmt{
				Expr: &Identifier{Name: "A", Line: 2, Column: 3, Token: "A"},
			},
			expected: "VTAB\n  Ident A\n",
		},
		{
			name: "VTabStmt with expression",
			stmt: &VTabStmt{
				Expr: &InfixExpr{
					Left:  &Identifier{Name: "V", Token: "V"},
					Op:    "+",
					Right: &NumberLiteral{Value: 5, Token: "5"},
				},
			},
			expected: "" +
				"VTAB\n" +
				"  Infix +\n" +
				"    Ident V\n" +
				"    Number 5\n",
		},
		{
			name: "VTabStmt with parenthesized expression",
			stmt: &VTabStmt{
				Expr: &InfixExpr{
					Left: &InfixExpr{
						Left:  &Identifier{Name: "A", Token: "A"},
						Op:    "+",
						Right: &Identifier{Name: "B", Token: "B"},
					},
					Op: "*",
					Right: &NumberLiteral{
						Value: 2,
						Token: "2",
					},
				},
			},
			expected: "" +
				"VTAB\n" +
				"  Infix *\n" +
				"    Infix +\n" +
				"      Ident A\n" +
				"      Ident B\n" +
				"    Number 2\n",
		},
		{
			name: "IfStmt without ELSE",
			stmt: &IfStmt{
				Cond: &InfixExpr{
					Left:  &Identifier{Name: "A", Token: "A"},
					Op:    "<",
					Right: &NumberLiteral{Value: 10, Token: "10"},
				},
				Then: []Statement{
					&GotoStmt{
						Expr: &NumberLiteral{Value: 20, Token: "20"},
					},
				},
			},
			expected: "" +
				"IF\n" +
				"  Infix <\n" +
				"    Ident A\n" +
				"    Number 10\n" +
				"THEN\n" +
				"  GOTO\n" +
				"    Number 20\n",
		},
		{
			name: "IfStmt with ELSE",
			stmt: &IfStmt{
				Cond: &Identifier{Name: "X", Token: "X"},
				Then: []Statement{
					&PrintStmt{
						Exprs: []Expression{
							&StringLiteral{Value: "YES", Token: "\"YES\""},
						},
					},
				},
				Else: []Statement{
					&PrintStmt{
						Exprs: []Expression{
							&StringLiteral{Value: "NO", Token: "\"NO\""},
						},
					},
				},
			},
			expected: "" +
				"IF\n" +
				"  Ident X\n" +
				"THEN\n" +
				"  PRINT\n" +
				"    EXPR 0:\n" +
				"      String \"YES\"\n" +
				"ELSE\n" +
				"  PRINT\n" +
				"    EXPR 0:\n" +
				"      String \"NO\"\n",
		},
		{
			name: "LetStmt with single index identifier",
			stmt: &LetStmt{
				Name: "A",
				Indices: []Expression{
					&Identifier{Name: "I"},
				},
				Value: &NumberLiteral{Value: 10, Token: "10"},
			},
			expected: "" +
				"LET A(I)\n" +
				"  Number 10\n",
		},
		{
			name: "LetStmt with single numeric index",
			stmt: &LetStmt{
				Name: "A",
				Indices: []Expression{
					&NumberLiteral{Value: 3, Token: "3"},
				},
				Value: &Identifier{Name: "X", Token: "X"},
			},
			expected: "" +
				"LET A(3)\n" +
				"  Ident X\n",
		},
		{
			name: "LetStmt with infix index",
			stmt: &LetStmt{
				Name: "A",
				Indices: []Expression{
					&InfixExpr{
						Left:  &Identifier{Name: "I", Token: "I"},
						Op:    "+",
						Right: &NumberLiteral{Value: 1, Token: "1"},
					},
				},
				Value: &NumberLiteral{Value: 42, Token: "42"},
			},
			expected: "" +
				"LET A(I + 1)\n" +
				"  Number 42\n",
		},
		{
			name: "LetStmt with multiple indices",
			stmt: &LetStmt{
				Name: "A",
				Indices: []Expression{
					&Identifier{Name: "I", Token: "I"},
					&Identifier{Name: "J", Token: "J"},
				},
				Value: &Identifier{Name: "X", Token: "X"},
			},
			expected: "" +
				"LET A(I,J)\n" +
				"  Ident X\n",
		},
		{
			name: "LetStmt with mixed indices expressions",
			stmt: &LetStmt{
				Name: "A",
				Indices: []Expression{
					&Identifier{Name: "I", Token: "I"},
					&InfixExpr{
						Left:  &Identifier{Name: "J", Token: "J"},
						Op:    "+",
						Right: &NumberLiteral{Value: 2, Token: "2"},
					},
				},
				Value: &InfixExpr{
					Left:  &Identifier{Name: "X", Token: "X"},
					Op:    "*",
					Right: &NumberLiteral{Value: 3, Token: "3"},
				},
			},
			expected: "" +
				"LET A(I,J + 2)\n" +
				"  Infix *\n" +
				"    Ident X\n" +
				"    Number 3\n",
		},
		{
			name: "DimStmt single array single dimension",
			stmt: &DimStmt{
				Arrays: []DimDecl{
					{
						Name: "A",
						Dimensions: []Expression{
							&NumberLiteral{Value: 10, Token: "10"},
						},
					},
				},
			},
			expected: "" +
				"DIM\n" +
				"  VAR 0: A\n" +
				"  EXPR 0:\n" +
				"    Number 10\n",
		},
		{
			name: "DimStmt single array multiple numeric dimensions",
			stmt: &DimStmt{
				Arrays: []DimDecl{
					{
						Name: "B",
						Dimensions: []Expression{
							&NumberLiteral{Value: 10, Token: "10"},
							&NumberLiteral{Value: 20, Token: "20"},
						},
					},
				},
			},
			expected: "" +
				"DIM\n" +
				"  VAR 0: B\n" +
				"  EXPR 0:\n" +
				"    Number 10\n" +
				"  EXPR 1:\n" +
				"    Number 20\n",
		},
		{
			name: "DimStmt with infix expression dimensions",
			stmt: &DimStmt{
				Arrays: []DimDecl{
					{
						Name: "C",
						Dimensions: []Expression{
							&InfixExpr{
								Left:  &Identifier{Name: "N", Token: "N"},
								Op:    "*",
								Right: &NumberLiteral{Value: 2, Token: "2"},
							},
						},
					},
				},
			},
			expected: "" +
				"DIM\n" +
				"  VAR 0: C\n" +
				"  EXPR 0:\n" +
				"    Infix *\n" +
				"      Ident N\n" +
				"      Number 2\n",
		},
		{
			name: "DimStmt multiple arrays",
			stmt: &DimStmt{
				Arrays: []DimDecl{
					{
						Name: "A",
						Dimensions: []Expression{
							&NumberLiteral{Value: 10, Token: "10"},
						},
					},
					{
						Name: "B$",
						Dimensions: []Expression{
							&Identifier{Name: "SIZE", Token: "SIZE"},
							&NumberLiteral{Value: 5, Token: "5"},
						},
					},
				},
			},
			expected: "" +
				"DIM\n" +
				"  VAR 0: A\n" +
				"  EXPR 0:\n" +
				"    Number 10\n" +
				"  VAR 1: B$\n" +
				"  EXPR 0:\n" +
				"    Ident SIZE\n" +
				"  EXPR 1:\n" +
				"    Number 5\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := testutils.CaptureStdout(t, func() {
				dumpStatement(tt.stmt, "", StdoutEmitter)
			})
			msg := fmt.Sprintf("tests[%s] - dumpStatement output mismatch.", tt.name)
			testutils.Equal(t, msg, output, tt.expected)
		})
	}
}

// ------------------------
// Tests pour DumpProgram
// ------------------------
func TestDumpProgram(t *testing.T) {
	p := &Program{
		Lines: []*Line{
			{
				Number: 10,
				Stmts: []Statement{
					&LetStmt{Name: "X", Value: &NumberLiteral{Value: 5, Line: 1, Column: 1, Token: "5"}},
					&PrintStmt{Exprs: []Expression{&Identifier{Name: "X", Line: 1, Column: 2, Token: "X"}}},
				},
			},
			{
				Number: 20,
				Stmts: []Statement{
					&ForStmt{
						Var:     "I",
						Start:   &NumberLiteral{Value: 1, Line: 2, Column: 1, Token: "1"},
						End:     &NumberLiteral{Value: 3, Line: 2, Column: 3, Token: "3"},
						LineNum: 20,
					},
					&NextStmt{Var: "I", ForLineNum: 20},
				},
			},
		},
	}

	expected := "Line 10\n  LET X\n    Number 5\n  PRINT\n    EXPR 0:\n      Ident X\nLine 20\n  FOR I (Line 20)\n    FROM:\n      Number 1\n    TO:\n      Number 3\n  NEXT I (FOR Line 20)\n"

	output := testutils.CaptureStdout(t, func() {
		DumpProgram(p, StdoutEmitter)
	})

	testutils.Equal(t, "DumpProgram output mismatch", output, expected)
}

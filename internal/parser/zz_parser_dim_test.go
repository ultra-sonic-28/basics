package parser

import (
	"basics/internal/lexer"
	"basics/internal/runtime"
	"basics/testutils"
	"fmt"
	"testing"
)

func TestParse_DIM_Statements(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		lineCount int
		assertFn  func(t *testing.T, prog *Program)
	}{
		{
			name:      "DIM A(10)",
			source:    `20 DIM A(10)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				dim := prog.Lines[0].Stmts[0].(*DimStmt)

				testutils.Equal(t, "array count", len(dim.Arrays), 1)

				arr := dim.Arrays[0]
				testutils.Equal(t, "name", arr.Name, "A")
				testutils.Equal(t, "type", arr.BaseType, runtime.NUMBER)
				testutils.Equal(t, "dims", len(arr.Dimensions), 1)

				_, ok := arr.Dimensions[0].(*NumberLiteral)
				testutils.True(t, "dimension is NumberLiteral", ok)
			},
		},
		{
			name:      "DIM A%(10)",
			source:    `20 DIM A%(10)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				arr := prog.Lines[0].Stmts[0].(*DimStmt).Arrays[0]

				testutils.Equal(t, "name", arr.Name, "A%")
				testutils.Equal(t, "type", arr.BaseType, runtime.INTEGER)
			},
		},
		{
			name:      "DIM A$(10)",
			source:    `20 DIM A$(10)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				arr := prog.Lines[0].Stmts[0].(*DimStmt).Arrays[0]

				testutils.Equal(t, "name", arr.Name, "A$")
				testutils.Equal(t, "type", arr.BaseType, runtime.STRING)
			},
		},
		{
			name:      "DIM A(10,20)",
			source:    `20 DIM A(10,20)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				arr := prog.Lines[0].Stmts[0].(*DimStmt).Arrays[0]
				testutils.Equal(t, "dimension count", len(arr.Dimensions), 2)
			},
		},
		{
			name:      "DIM multiple arrays",
			source:    `20 DIM A(10), B%(50), C$(10,10)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				arrays := prog.Lines[0].Stmts[0].(*DimStmt).Arrays
				testutils.Equal(t, "array count", len(arrays), 3)

				testutils.Equal(t, "A dims", len(arrays[0].Dimensions), 1)
				testutils.Equal(t, "B% dims", len(arrays[1].Dimensions), 1)
				testutils.Equal(t, "C$ dims", len(arrays[2].Dimensions), 2)
			},
		},
		{
			name:      "DIM chained with colon",
			source:    `20 DIM OA$(75):DIM OM(75)`,
			lineCount: 1,
			assertFn: func(t *testing.T, prog *Program) {
				testutils.Equal(t, "stmt count", len(prog.Lines[0].Stmts), 2)

				for _, stmt := range prog.Lines[0].Stmts {
					dim := stmt.(*DimStmt)
					testutils.Equal(t, "one array per DIM", len(dim.Arrays), 1)
				}
			},
		},
		{
			name: "DIM with variable size",
			source: `20 Size = 50
30 DIM A%(Size)`,
			lineCount: 2,
			assertFn: func(t *testing.T, prog *Program) {
				arr := prog.Lines[1].Stmts[0].(*DimStmt).Arrays[0]

				_, ok := arr.Dimensions[0].(*Identifier)
				testutils.True(t, "dimension is Identifier", ok)
			},
		},
		{
			name: "DIM with expression 2*Size",
			source: `20 Size = 50
30 DIM A%(2*Size)`,
			lineCount: 2,
			assertFn: func(t *testing.T, prog *Program) {
				arr := prog.Lines[1].Stmts[0].(*DimStmt).Arrays[0]

				_, ok := arr.Dimensions[0].(*InfixExpr)
				testutils.True(t, "dimension is InfixExpr", ok)
			},
		},
		{
			name: "DIM complex dimensions",
			source: `20 Size = 50
30 DIM A%(Size,10,Size/10,2*Size,5)`,
			lineCount: 2,
			assertFn: func(t *testing.T, prog *Program) {
				arr := prog.Lines[1].Stmts[0].(*DimStmt).Arrays[0]
				testutils.Equal(t, "dimension count", len(arr.Dimensions), 5)

				_, ok := arr.Dimensions[0].(*Identifier)
				testutils.True(t, "dim[0] Identifier", ok)

				_, ok = arr.Dimensions[2].(*InfixExpr)
				testutils.True(t, "dim[2] InfixExpr", ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexer.Lex(tt.source)
			p := New(tokens)
			prog, errs := p.ParseProgram()

			testutils.Equal(t, fmt.Sprintf("test[%s] - no parser errors", tt.name), len(errs), 0)
			testutils.Equal(t, fmt.Sprintf("test[%s] - line count", tt.name), len(prog.Lines), tt.lineCount)

			tt.assertFn(t, prog)
		})
	}
}

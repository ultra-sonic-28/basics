package parser

import (
	"strings"
	"testing"

	"basics/internal/errors"
	"basics/internal/lexer"
	"basics/testutils"
)

func TestParse_LineNumberRequired(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectError bool
	}{
		{
			name: "valid program with line number",
			source: `
10 PRINT "HELLO"
`,
			expectError: false,
		},
		{
			name: "line with only line number is valid (1)",
			source: `
10
`,
			expectError: false,
		},
		{
			name: "line with only line number is valid (2)",
			source: `
10 REM Only linenumber is present
20
30 REM This is valid
40

`,
			expectError: false,
		},
		{
			name: "missing line number after first line",
			source: `
10 REM Test
PRINT "HELLO"
30 END
`,
			expectError: true,
		},
		{
			name: "program without any line number",
			source: `
PRINT "HELLO"
`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexer.Lex(tt.source)
			p := New(tokens)
			prog, errs := p.ParseProgram()

			if tt.expectError {
				testutils.True(t,
					"parser should report errors",
					len(errs) > 0,
				)

				found := false
				for _, err := range errs {
					if strings.Contains(err.Error(), "EXPECTED LINE NUMBER") {
						found = true
						break
					}
				}

				testutils.True(t,
					"error EXPECTED LINE NUMBER should be reported",
					found,
				)
			} else {
				testutils.Equal(t,
					"no parser errors expected",
					len(errs),
					0,
				)

				testutils.True(t,
					"program should have at least one line",
					len(prog.Lines) > 0,
				)
			}
		})
	}
}

func TestParse_DuplicateLineNumbers(t *testing.T) {
	source := `
10 PRINT "HELLO"
10 PRINT "WORLD"
`

	// Lexing
	tokens := lexer.Lex(source)

	// Parsing
	p := New(tokens)
	prog, errs := p.ParseProgram()

	// --- Assertions ---
	testutils.Equal(t, "program parsed", len(prog.Lines), 2)
	testutils.Equal(t, "one error reported", len(errs), 1)

	err := errs[0]

	// Type d’erreur
	testutils.Equal(t, "error is semantic", err.Kind, errors.Semantic)

	// Ligne concernée
	testutils.Equal(t, "duplicate line number", err.Line, 10)

	// Message
	testutils.Equal(
		t,
		"duplicate line number message",
		err.Msg,
		"DUPLICATE LINE NUMBER",
	)
}

package lexer

import (
	"testing"

	"basics/testutils"
)

func TestKeywords_ContainsExpectedKeywords(t *testing.T) {
	expected := []string{
		// Contrôle
		"FOR", "TO", "STEP", "NEXT",
		"IF", "THEN",
		"GOTO", "GOSUB", "RETURN",
		"END", "STOP",

		// Variables & logique
		"LET", "DIM",
		"REM",

		// I/O
		"PRINT", "INPUT", "GET",

		// Math
		"SIN", "COS", "TAN",
		"INT", "ABS", "RND",

		// Graphique / écran
		"GR", "HGR", "TEXT",
		"PLOT", "HPLOT",
		"COLOR", "HCOLOR",
		"HOME",

		// DATA
		"DATA", "READ", "RESTORE",

		// Autres
		"POKE", "PEEK", "CALL",
		"TAB", "VTAB", "HTAB",
		"INVERSE", "NORMAL", "FLASH",

		// Extension
		"SLEEP",
	}

	for _, kw := range expected {
		t.Run(kw, func(t *testing.T) {
			val, ok := Keywords[kw]
			testutils.True(t, "keyword should exist: "+kw, ok)
			testutils.True(t, "keyword value should be true: "+kw, val)
		})
	}
}

func TestKeywords_DoesNotContainInvalidKeywords(t *testing.T) {
	invalid := []string{
		"for",      // lowercase
		"Print",    // mixed case
		"WHILE",    // non supporté
		"ELSE",     // non supporté
		"FUNCTION", // non supporté
		"",         // vide
		"123",      // numérique
	}

	for _, kw := range invalid {
		t.Run(kw, func(t *testing.T) {
			_, ok := Keywords[kw]
			testutils.False(t, "invalid keyword should not exist: "+kw, ok)
		})
	}
}

func TestKeywords_AllValuesAreTrue(t *testing.T) {
	for kw, val := range Keywords {
		t.Run(kw, func(t *testing.T) {
			testutils.True(t, "keyword value must be true: "+kw, val)
		})
	}
}

func TestKeywords_ContainAllExpected(t *testing.T) {
	expected := []string{
		"FOR", "TO", "STEP", "NEXT",
		"IF", "THEN",
		"PRINT", "INPUT",
		"SLEEP",
	}

	testutils.AssertMapContainsKeys(
		t,
		"missing BASIC keywords",
		Keywords,
		expected,
	)
}
